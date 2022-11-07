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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"net/url"
	"os"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initSSL() *cobra.Command {
	var flags struct {
		ca, key string
	}

	var sslCmd = &cobra.Command{
		Use:   CommandSSL,
		Short: "Genreate self-sign certificate",
		RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}

	var sslSubCmdGenerate = &cobra.Command{
		Use:   CommandGenerate,
		Short: "Generate certificates",
		RunE: func(_ *cobra.Command, _ []string) error {
			var s SSL
			return s.Generate()
		}}

	var sslSubCmdSign = &cobra.Command{
		Use:   CommandSign,
		Short: "Sign certificates from giving ca files",
		RunE: func(_ *cobra.Command, _ []string) error {
			if flags.ca == "" || flags.key == "" {
				return nil
			}
			var s SSL
			return s.Sign(flags.ca, flags.key)
		}}
	sslSubCmdSign.Flags().StringVarP(&flags.ca, "ca", "c", "", common.Usage("Specify CA file"))
	sslSubCmdSign.Flags().StringVarP(&flags.key, "key", "k", "", common.Usage("Specify private key file"))
	sslCmd.AddCommand(sslSubCmdGenerate, sslSubCmdSign)
	return sslCmd
}

type SSL struct{}

func (*SSL) defaultSubject() pkix.Name {
	const (
		defaultCountry            = "TW"
		defaultState              = "TP"
		defaultLocality           = "Xinyi"
		defaultOrganization       = common.RepoName
		defaultOrganizationalUnit = "Root CA"
		defaultCommonName         = "Self-Sign Root CA"
	)
	return pkix.Name{
		Organization:       []string{defaultOrganization},
		OrganizationalUnit: []string{defaultOrganizationalUnit},
		Country:            []string{defaultCountry},
		Province:           []string{defaultState},
		Locality:           []string{defaultLocality},
		CommonName:         defaultCommonName,
	}
}

func (s *SSL) serverSubject() *x509.Certificate {
	type cnf struct {
		Country        string   `json:"C"`
		CommonName     string   `json:"CN"`
		Locality       string   `json:"L"`
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
				Locality:           []string{info.Locality},
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
	d := s.defaultSubject()
	return &x509.Certificate{
		SerialNumber: big.NewInt(common.TimeNow.Unix()),
		Subject: pkix.Name{
			Country:            d.Country,
			Organization:       d.Organization,
			OrganizationalUnit: d.OrganizationalUnit,
			Province:           d.Province,
			Locality:           d.Locality,
			CommonName:         d.CommonName,
		},
		NotBefore: common.TimeNow.UTC(),
		NotAfter:  common.TimeNow.UTC().AddDate(1, 0, 0),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,

		BasicConstraintsValid: true,
		IsCA:                  true,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},

		DNSNames:       []string{"localhost", "*.localhost"},
		EmailAddresses: nil,
		IPAddresses:    nil,
		URIs:           nil,

		SubjectKeyId: []byte{0, 1, 1, 1, 1, 1},
	}
}

func (s *SSL) Generate() error {
	bits := 4096
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(common.TimeNow.Unix()),
		Subject:      s.defaultSubject(),
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
		Subject:      s.defaultSubject(),
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

	server := s.serverSubject()
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

func (s *SSL) Sign(caCert, caKey string) error {
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

	server := s.serverSubject()
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
