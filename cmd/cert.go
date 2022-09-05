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
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var certFlag CertFlag

	var certCmd = &cobra.Command{
		Use:   CommandCert + " [host|file]",
		Args:  cobra.ExactArgs(1),
		Short: "Check tls cert expiry time",
		RunE:  certFlag.RunE,
		Example: common.Examples(`# Print certificate expiration time, DNS, IP and issuer
www.google.com

# Only print certificate expiration time
1.1.1.1 --expiry

# Only print certificate DNS
www.google.com --dns

# Print certificate expiration time, DNS and issuer
example.com.crt`, CommandCert),
	}
	rootCmd.AddCommand(certCmd)

	certCmd.Flags().StringVarP(&certFlag.port, "port", "p", "443", common.Usage("Specify host port"))
	certCmd.Flags().BoolVar(&certFlag.ip, "ip", false, common.Usage("Only print IP"))
	certCmd.Flags().BoolVar(&certFlag.expiry, "expiry", false, common.Usage("Only print expiry time"))
	certCmd.Flags().BoolVar(&certFlag.dns, "dns", false, common.Usage("Only print DNS names"))
	certCmd.Flags().BoolVar(&certFlag.issuer, "issuer", false, common.Usage("Only print issuer"))
	certCmd.Flags().BoolVar(&certFlag.days, "days", false, common.Usage("Only print the remaining days"))
}

type CertFlag struct {
	ip, expiry, days, dns, issuer bool

	port string

	resp *certResponse
}

func (cf *CertFlag) RunE(cmd *cobra.Command, args []string) error {
	var err error
	input := args[0]
	switch {
	case validator.ValidFile(input):
		cf.resp, err = cf.resp.CheckFile(input)
	case validator.ValidDomain(input) || validator.ValidIPv4(input):
		cf.resp, err = cf.resp.CheckHost(input + ":" + cf.port)
	default:
		return common.ErrInvalidArg
	}
	if err != nil {
		return err
	}
	if cf.resp == nil {
		return common.ErrResponse
	}
	switch {
	default:
		OutputDefaultJSON(cf.resp)
	case cf.ip:
		PrintString(cf.resp.ServerIP)
	case cf.dns:
		PrintJSON(cf.resp.DNS)
	case cf.expiry:
		PrintString(cf.resp.ExpiryTime)
	case cf.issuer:
		PrintString(cf.resp.Issuer)
	case cf.days:
		PrintString(cf.resp.Days)
	}
	return err
}

type certResponse struct {
	ExpiryTime string   `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty"`
	Days       int      `json:"days,omitempty" yaml:"days,omitempty"`
	Issuer     string   `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	ServerIP   string   `json:"serverIp,omitempty" yaml:"serverIp,omitempty"`
	DNS        []string `json:"dns,omitempty" yaml:"dns,omitempty"`
}

func (c *certResponse) CheckHost(host string) (*certResponse, error) {
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	dayRemain := cert.NotAfter.Local().Sub(common.TimeNow)
	var out = certResponse{
		ExpiryTime: cert.NotAfter.Local().Format(time.RFC3339),
		DNS:        cert.DNSNames,
		Issuer:     cert.Issuer.String(),
		ServerIP:   conn.RemoteAddr().String(),
		Days:       int(dayRemain.Hours() / 24),
	}
	return &out, err
}

func (c *certResponse) CheckFile(fileName string) (*certResponse, error) {
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
	var out = certResponse{
		ExpiryTime: cert[0].NotAfter.Local().Format(time.RFC3339),
		DNS:        cert[0].DNSNames,
		Issuer:     cert[0].Issuer.String(),
		Days:       int(dayRemain.Hours() / 24),
	}
	return &out, err
}
