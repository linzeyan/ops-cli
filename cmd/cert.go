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
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/linzeyan/tlsCheck"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Check tls cert",
	Run: func(_ *cobra.Command, _ []string) {
		var out *certOutput
		var err error

		if certHost != "domain:port" {
			out, err = certOutputByHost()
			if err != nil {
				log.Println(err)
				return
			}
		} else if certCrt != "" {
			out, err = certOutputByFile()
			if err != nil {
				log.Println(err)
				return
			}
		}
		if out == nil {
			log.Println("response is empty")
			return
		}
		outputDefaultString(out)
	},
	Example: Examples(`# Print certificate expiration time, DNS, IP and issuer
ops-cli cert -d www.google.com:443

# Print certificate expiration time
ops-cli cert -d www.google.com:443 --expiry

# Print certificate DNS
ops-cli cert -d www.google.com:443 --dns

# Print certificate expiration time, DNS and issuer
ops-cli cert -f example.com.crt`),
}

var certCrt, certHost string
var certIP, certExpiry, certDNS, certIssuer bool

func init() {
	rootCmd.AddCommand(certCmd)

	certCmd.Flags().StringVarP(&certCrt, "file", "f", "", "Specify .crt file path")
	certCmd.Flags().StringVarP(&certHost, "domain", "d", "domain:port", "Specify domain and host port")
	certCmd.MarkFlagsMutuallyExclusive("file", "domain")
	certCmd.Flags().BoolVarP(&certIP, "ip", "", false, "Print IP")
	certCmd.Flags().BoolVarP(&certExpiry, "expiry", "", false, "Print expiry time")
	certCmd.Flags().BoolVarP(&certDNS, "dns", "", false, "Print DNS names")
	certCmd.Flags().BoolVarP(&certIssuer, "issuer", "", false, "Print issuer")
}

type certOutput struct {
	ExpiryTime string   `json:"ExpiryTime,omitempty" yaml:"ExpiryTime,omitempty"`
	Issuer     string   `json:"Issuer,omitempty" yaml:"Issuer,omitempty"`
	ServerIP   string   `json:"ServerIP,omitempty" yaml:"ServerIP,omitempty"`
	DNS        []string `json:"DNS,omitempty" yaml:"DNS,omitempty"`
}

func (c certOutput) Json() {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (c certOutput) Yaml() {
	out, err := yaml.Marshal(c)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (c certOutput) String() {
	if certIP {
		fmt.Println(c.ServerIP)
		return
	}
	if certDNS {
		fmt.Println(c.DNS)
		return
	}
	if certExpiry {
		fmt.Println(c.ExpiryTime)
		return
	}
	if certIssuer {
		fmt.Println(c.Issuer)
		return
	}

	var s strings.Builder
	f := reflect.ValueOf(&c).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		_, err := s.WriteString(fmt.Sprintf("%-10s\t%v\n", t.Field(i).Name, f.Field(i).Interface()))
		//f.Field(i).Type()
		if err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Println(s.String())
}

func certOutputByHost() (*certOutput, error) {
	conn, err := tlsCheck.CheckByHost(certHost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	var out = certOutput{
		ExpiryTime: cert.NotAfter.Local().Format(time.RFC3339),
		DNS:        cert.DNSNames,
		Issuer:     cert.Issuer.String(),
		ServerIP:   conn.RemoteAddr().String(),
	}
	return &out, nil
}

func certOutputByFile() (*certOutput, error) {
	cert, err := tlsCheck.CheckByFile(certCrt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if cert == nil {
		log.Println(cert)
		return nil, err
	}
	var out = certOutput{
		ExpiryTime: cert[0].NotAfter.Local().Format(time.RFC3339),
		DNS:        cert[0].DNSNames,
		Issuer:     cert[0].Issuer.String(),
	}
	return &out, nil
}
