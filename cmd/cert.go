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
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Check tls cert expiry time",
	Run: func(cmd *cobra.Command, args []string) {
		var out *certResponse
		var err error
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		input := args[0]
		switch {
		case ValidFile(input):
			out, err = out.CheckFile(input)
		case ValidDomain(input) || net.ParseIP(input).To4() != nil:
			out, err = out.CheckHost(input + ":" + certPort)
		default:
			_ = cmd.Help()
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		if out == nil {
			log.Println("response is empty")
			return
		}
		OutputDefaultString(out)
	},
	Example: Examples(`# Print certificate expiration time, DNS, IP and issuer
ops-cli cert www.google.com

# Only print certificate expiration time
ops-cli cert 1.1.1.1 --expiry

# Only print certificate DNS
ops-cli cert www.google.com --dns

# Print certificate expiration time, DNS and issuer
ops-cli cert example.com.crt`),
}

var certPort string
var certIP, certExpiry, certRemainDays, certDNS, certIssuer bool

func init() {
	rootCmd.AddCommand(certCmd)

	certCmd.Flags().StringVarP(&certPort, "port", "p", "443", "Specify host port")
	certCmd.Flags().BoolVar(&certIP, "ip", false, "Only print IP")
	certCmd.Flags().BoolVar(&certExpiry, "expiry", false, "Only print expiry time")
	certCmd.Flags().BoolVar(&certDNS, "dns", false, "Only print DNS names")
	certCmd.Flags().BoolVar(&certIssuer, "issuer", false, "Only print issuer")
	certCmd.Flags().BoolVar(&certRemainDays, "days", false, "Only print the remaining days")
}

type certResponse struct {
	ExpiryTime string   `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty"`
	Days       int      `json:"days,omitempty" yaml:"days,omitempty"`
	Issuer     string   `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	ServerIP   string   `json:"serverIp,omitempty" yaml:"serverIp,omitempty"`
	DNS        []string `json:"dns,omitempty" yaml:"dns,omitempty"`
}

func (c certResponse) JSON() {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (c certResponse) YAML() {
	out, err := yaml.Marshal(c)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (c certResponse) String() {
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
	if certRemainDays {
		fmt.Println(c.Days)
		return
	}

	var s strings.Builder
	f := reflect.ValueOf(&c).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		_, err := s.WriteString(fmt.Sprintf("%-10s\t%v\n", t.Field(i).Name, f.Field(i).Interface()))
		// f.Field(i).Type()
		if err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Println(s.String())
}

func (c *certResponse) CheckHost(host string) (*certResponse, error) {
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	dayRemain := cert.NotAfter.Local().Sub(rootNow)
	var out = certResponse{
		ExpiryTime: cert.NotAfter.Local().Format(time.RFC3339),
		DNS:        cert.DNSNames,
		Issuer:     cert.Issuer.String(),
		ServerIP:   conn.RemoteAddr().String(),
		Days:       int(dayRemain.Hours() / 24),
	}
	return &out, nil
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
				log.Println(err)
				break
			}
		}
		t = n
	}
	buf = buf[0:t]
	crtPem, _ := pem.Decode(buf)
	if crtPem == nil {
		return nil, errors.New("file type not correct")
	}
	cert, err := x509.ParseCertificates(crtPem.Bytes)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, errors.New("can not correctly parse")
	}

	dayRemain := cert[0].NotAfter.Local().Sub(rootNow)
	var out = certResponse{
		ExpiryTime: cert[0].NotAfter.Local().Format(time.RFC3339),
		DNS:        cert[0].DNSNames,
		Issuer:     cert[0].Issuer.String(),
		Days:       int(dayRemain.Hours() / 24),
	}
	return &out, nil
}
