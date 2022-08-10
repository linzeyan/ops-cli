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
	"encoding/pem"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Args:  cobra.ExactArgs(1),
	Short: "Check tls cert expiry time",
	Run:   crtf.Run,
	Example: Examples(`# Print certificate expiration time, DNS, IP and issuer
ops-cli cert www.google.com

# Only print certificate expiration time
ops-cli cert 1.1.1.1 --expiry

# Only print certificate DNS
ops-cli cert www.google.com --dns

# Print certificate expiration time, DNS and issuer
ops-cli cert example.com.crt`),
}

var crtf certFlag

func init() {
	rootCmd.AddCommand(certCmd)

	certCmd.Flags().StringVarP(&crtf.port, "port", "p", "443", "Specify host port")
	certCmd.Flags().BoolVar(&crtf.ip, "ip", false, "Only print IP")
	certCmd.Flags().BoolVar(&crtf.expiry, "expiry", false, "Only print expiry time")
	certCmd.Flags().BoolVar(&crtf.dns, "dns", false, "Only print DNS names")
	certCmd.Flags().BoolVar(&crtf.issuer, "issuer", false, "Only print issuer")
	certCmd.Flags().BoolVar(&crtf.days, "days", false, "Only print the remaining days")
}

type certFlag struct {
	ip, expiry, days, dns, issuer bool

	port string

	resp *certResponse
}

func (cf *certFlag) Run(cmd *cobra.Command, args []string) {
	var err error
	input := args[0]
	switch {
	case ValidFile(input):
		cf.resp, err = cf.resp.CheckFile(input)
	case ValidDomain(input) || net.ParseIP(input).To4() != nil:
		cf.resp, err = cf.resp.CheckHost(input + ":" + cf.port)
	default:
		_ = cmd.Help()
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	if cf.resp == nil {
		log.Println("response is empty")
		return
	}
	switch {
	default:
		OutputDefaultJSON(cf.resp)
	case cf.ip:
		PrintString(cf.resp.ServerIP)
	case cf.dns:
		PrintString(cf.resp.DNS)
	case cf.expiry:
		PrintString(cf.resp.ExpiryTime)
	case cf.issuer:
		PrintString(cf.resp.Issuer)
	case cf.days:
		PrintString(cf.resp.Days)
	}
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
