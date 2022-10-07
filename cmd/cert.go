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
	"net"
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
