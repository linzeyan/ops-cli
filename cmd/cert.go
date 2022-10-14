/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"math/big"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func initCert() *cobra.Command {
	var flags struct {
		ip, expiry, days, dns, issuer bool

		port string

		ca, key string
	}
	var certCmd = &cobra.Command{
		Use:   CommandCert + " [host|file]",
		Args:  cobra.ExactArgs(1),
		Short: "Check tls cert expiry time",
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			input := args[0]
			resp := new(Cert)
			switch {
			case validator.ValidFile(input):
				resp, err = resp.CheckFile(input)
			case validator.ValidDomain(input) || validator.ValidIPv4(input):
				resp, err = resp.CheckHost(net.JoinHostPort(input, flags.port))
			default:
				return common.ErrInvalidArg
			}
			if err != nil {
				return err
			}
			if resp == nil {
				return common.ErrResponse
			}
			switch {
			default:
				OutputDefaultJSON(resp)
			case flags.ip:
				PrintString(resp.ServerIP)
			case flags.dns:
				PrintJSON(resp.DNS)
			case flags.expiry:
				PrintString(resp.ExpiryTime)
			case flags.issuer:
				PrintString(resp.Issuer)
			case flags.days:
				PrintString(resp.Days)
			}
			return err
		},
		Example: common.Examples(`# Print certificate expiration time, DNS, IP and issuer
www.google.com

# Only print certificate expiration time
1.1.1.1 --expiry

# Only print certificate DNS
www.google.com --dns

# Print certificate expiration time, DNS and issuer
example.com.crt`, CommandCert),
	}

	certCmd.Flags().StringVarP(&flags.port, "port", "p", "443", common.Usage("Specify host port"))
	certCmd.Flags().BoolVar(&flags.ip, "ip", false, common.Usage("Only print IP"))
	certCmd.Flags().BoolVar(&flags.expiry, "expiry", false, common.Usage("Only print expiry time"))
	certCmd.Flags().BoolVar(&flags.dns, "dns", false, common.Usage("Only print DNS names"))
	certCmd.Flags().BoolVar(&flags.issuer, "issuer", false, common.Usage("Only print issuer"))
	certCmd.Flags().BoolVar(&flags.days, "days", false, common.Usage("Only print the remaining days"))

	var certSubCmdGenerate = &cobra.Command{
		Use:   CommandGenerate,
		Short: "Generate certificates",
		RunE: func(_ *cobra.Command, _ []string) error {
			var c Cert
			return c.Generate()
		}}

	var certSubCmdSign = &cobra.Command{
		Use:   CommandSign,
		Short: "Sign certificates from giving ca files",
		RunE: func(_ *cobra.Command, _ []string) error {
			if flags.ca == "" || flags.key == "" {
				return nil
			}
			var c Cert
			return c.Sign(flags.ca, flags.key)
		}}
	certSubCmdSign.Flags().StringVarP(&flags.ca, "ca", "c", "", common.Usage("Specify CA file"))
	certSubCmdSign.Flags().StringVarP(&flags.key, "key", "k", "", common.Usage("Specify private key file"))
	certCmd.AddCommand(certSubCmdGenerate)
	certCmd.AddCommand(certSubCmdSign)
	return certCmd
}

type Cert struct {
	ExpiryTime string   `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty"`
	Days       int      `json:"days,omitempty" yaml:"days,omitempty"`
	Issuer     string   `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	ServerIP   string   `json:"serverIp,omitempty" yaml:"serverIp,omitempty"`
	DNS        []string `json:"dns,omitempty" yaml:"dns,omitempty"`
}

func (c *Cert) CheckHost(host string) (*Cert, error) {
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	dayRemain := cert.NotAfter.Local().Sub(common.TimeNow)
	var out = Cert{
		ExpiryTime: cert.NotAfter.Local().Format(time.RFC3339),
		DNS:        cert.DNSNames,
		Issuer:     cert.Issuer.String(),
		ServerIP:   conn.RemoteAddr().String(),
		Days:       int(dayRemain.Hours() / 24),
	}
	return &out, err
}

func (c *Cert) CheckFile(fileName string) (*Cert, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	buf := make([]byte, 4096*3)
	var t int
	for {
		n, err := reader.Read(buf)
		if n == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
		}
		t = n
	}
	buf = buf[0:t]
	crtPem, err := Encoder.PemDecode(buf)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificates(crtPem)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, common.ErrInvalidFile
	}

	dayRemain := cert[0].NotAfter.Local().Sub(common.TimeNow)
	var out = Cert{
		ExpiryTime: cert[0].NotAfter.Local().Format(time.RFC3339),
		DNS:        cert[0].DNSNames,
		Issuer:     cert[0].Issuer.String(),
		Days:       int(dayRemain.Hours() / 24),
	}
	return &out, err
}

func (*Cert) defaultSubject() pkix.Name {
	return pkix.Name{
		Organization:       []string{"OPS-CLI"},
		OrganizationalUnit: []string{"Root CA"},
		Country:            []string{"TW"},
		Locality:           []string{"Taipei"},
		CommonName:         "Self-Sign Root CA",
	}
}

func (*Cert) serverSubject() *x509.Certificate {
	type cnf struct {
		Country        string   `json:"C"`
		CommonName     string   `json:"CN"`
		Location       string   `json:"L"`
		Org            string   `json:"O"`
		OrgUnit        string   `json:"OU"`
		State          string   `json:"ST"`
		DNSNames       []string `json:"dnsNames"`
		EmailAddresses []string `json:"emailAddresses"`
		IPAddresses    []string `json:"ipAddresses"`
		URIs           []string `json:"uris"`
		Year           int      `json:"year"`
	}

	var info cnf
	if rootConfig != "" {
		if err := ReadConfig(CommandCert, &info); err != nil {
			return nil
		}
	}
	var ip []net.IP
	for _, v := range info.IPAddresses {
		ip = append(ip, net.ParseIP(v))
	}
	var uri []*url.URL
	for _, v := range info.URIs {
		u, err := url.ParseRequestURI(v)
		if err == nil {
			uri = append(uri, u)
		}
	}
	return &x509.Certificate{
		SerialNumber: big.NewInt(common.TimeNow.Unix()),
		Subject: pkix.Name{
			Country:            []string{info.Country},
			Organization:       []string{info.Org},
			OrganizationalUnit: []string{info.OrgUnit},
			Province:           []string{info.State},
			Locality:           []string{info.Location},
			CommonName:         info.CommonName,
		},
		NotBefore: common.TimeNow.UTC(),
		NotAfter:  common.TimeNow.UTC().AddDate(info.Year, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,

		BasicConstraintsValid: true,
		IsCA:                  true,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},

		DNSNames:       info.DNSNames,
		EmailAddresses: info.EmailAddresses,
		IPAddresses:    ip,
		URIs:           uri,

		SubjectKeyId: []byte{1, 1, 1, 1, 1, 1},
	}
}

func (c *Cert) Generate() error {
	bits := 4096
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(common.TimeNow.Unix()),
		Subject:      c.defaultSubject(),
		NotBefore:    common.TimeNow.UTC(),
		NotAfter:     common.TimeNow.UTC().AddDate(15, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,

		BasicConstraintsValid: true,
		IsCA:                  true,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	caKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	caCert, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	var caCertPem, caKeyPem string
	caKeyPem, err = Encoder.PemEncode(x509.MarshalPKCS1PrivateKey(caKey), "RSA PRIVATE KEY")
	if err != nil {
		return err
	}
	caCertPem, err = Encoder.PemEncode(caCert, "CERTIFICATE")
	if err != nil {
		return err
	}

	subCa := &x509.Certificate{
		SerialNumber: big.NewInt(common.TimeNow.Unix()),
		Subject:      c.defaultSubject(),
		NotBefore:    common.TimeNow.UTC(),
		NotAfter:     common.TimeNow.UTC().AddDate(10, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,

		BasicConstraintsValid: true,
		IsCA:                  true,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	crt, err := x509.CreateCertificate(rand.Reader, subCa, ca, &privKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	var crtPem, keyPem string
	keyPem, err = Encoder.PemEncode(x509.MarshalPKCS1PrivateKey(privKey), "RSA PRIVATE KEY")
	if err != nil {
		return err
	}
	crtPem, err = Encoder.PemEncode(crt, "CERTIFICATE")
	if err != nil {
		return err
	}

	server := c.serverSubject()
	serverKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	serverCert, err := x509.CreateCertificate(rand.Reader, server, subCa, &serverKey.PublicKey, privKey)
	if err != nil {
		return err
	}
	var serverCertPem, serverKeyPem string
	serverKeyPem, err = Encoder.PemEncode(x509.MarshalPKCS1PrivateKey(serverKey), "RSA PRIVATE KEY")
	if err != nil {
		return err
	}
	serverCertPem, err = Encoder.PemEncode(serverCert, "CERTIFICATE")
	if err != nil {
		return err
	}

	_ = os.WriteFile("root.key", []byte(caKeyPem), FileModeROwner)
	_ = os.WriteFile("root.crt", []byte(caCertPem), FileModeRAll)
	_ = os.WriteFile("ca.key", []byte(keyPem), FileModeROwner)
	_ = os.WriteFile("ca.crt", []byte(crtPem), FileModeRAll)
	_ = os.WriteFile("server.key", []byte(serverKeyPem), FileModeROwner)
	_ = os.WriteFile("server.crt", []byte(serverCertPem), FileModeRAll)
	return err
}

func (c *Cert) Sign(caCert, caKey string) error {
	bits := 4096
	caCertFile, err := os.ReadFile(caCert)
	if err != nil {
		return err
	}
	caKeyFile, err := os.ReadFile(caKey)
	if err != nil {
		return err
	}
	caCertDecode, err := Encoder.PemDecode(caCertFile)
	if err != nil {
		return err
	}
	caKeyDecode, err := Encoder.PemDecode(caKeyFile)
	if err != nil {
		return err
	}
	key, err := x509.ParsePKCS1PrivateKey(caKeyDecode)
	if err != nil {
		return err
	}
	ca, err := x509.ParseCertificate(caCertDecode)
	if err != nil {
		return err
	}

	server := c.serverSubject()
	serverKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	serverCert, err := x509.CreateCertificate(rand.Reader, server, ca, &serverKey.PublicKey, key)
	if err != nil {
		return err
	}
	var serverCertPem, serverKeyPem string
	serverKeyPem, err = Encoder.PemEncode(x509.MarshalPKCS1PrivateKey(serverKey), "RSA PRIVATE KEY")
	if err != nil {
		return err
	}
	serverCertPem, err = Encoder.PemEncode(serverCert, "CERTIFICATE")
	if err != nil {
		return err
	}

	_ = os.WriteFile("server.key", []byte(serverKeyPem), FileModeROwner)
	err = os.WriteFile("server.crt", []byte(serverCertPem), FileModeRAll)
	return err
}
