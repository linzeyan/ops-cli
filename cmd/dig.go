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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var digCmd = &cobra.Command{
	Use:   "dig [host] [@server] [type]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Resolve domain name",
	Run: func(_ *cobra.Command, args []string) {
		var err error
		var argsWithoutDomain []string
		var argsType []string
		switch lens := len(args); {
		case lens == 1:
			digDomain = args[0]
			if err = digOutput.Request(dns.TypeA); err != nil {
				log.Println(err)
				os.Exit(1)
			}
			if digOutput == nil {
				return
			}
			OutputDefaultString(&digOutput)
			return
		case lens > 1:
			/* Find which arg is domain. */
			for i := range args {
				if ValidDomain(args[i]) || ValidIP(args[i]) {
					digDomain = args[i]
					argsWithoutDomain = append(args[:i], args[i+1:]...)
					break
				}
			}
			switch lens {
			/* Distinguish remain arg is NS Server or DNS Type */
			case 2:
				if strings.Contains(argsWithoutDomain[0], "@") {
					digServer = strings.Replace(argsWithoutDomain[0], "@", "", 1)
					argsType = append(argsType, "A")
				} else {
					argsType = append(argsType, argsWithoutDomain[0])
				}
			default:
				for i := range argsWithoutDomain {
					if strings.Contains(argsWithoutDomain[i], "@") {
						digServer = strings.Replace(argsWithoutDomain[i], "@", "", 1)
						/* Copy args and remove @x.x.x.x */
						argsType = append(argsWithoutDomain[:i], argsWithoutDomain[i+1:]...)
						break
					}
				}
			}
		}
		switch strings.ToLower(argsType[0]) {
		case "a":
			err = digOutput.Request(dns.TypeA)
		case "aaaa":
			err = digOutput.Request(dns.TypeAAAA)
		case "caa":
			err = digOutput.Request(dns.TypeCAA)
		case "cname":
			err = digOutput.Request(dns.TypeCNAME)
		case "mx":
			err = digOutput.Request(dns.TypeMX)
		case "ns":
			err = digOutput.Request(dns.TypeNS)
		case "ptr":
			err = digOutput.Request(dns.TypePTR)
		case "soa":
			err = digOutput.Request(dns.TypeSOA)
		case "srv":
			err = digOutput.Request(dns.TypeSRV)
		case "txt":
			err = digOutput.Request(dns.TypeTXT)
		case "any":
			err = digOutput.Request(dns.TypeANY)
		}
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if digOutput == nil {
			return
		}
		OutputDefaultString(&digOutput)
	},
	Example: Examples(`# Query A record
ops-cli dig google.com
ops-cli dig @1.1.1.1 google.com A
ops-cli dig @8.8.8.8 google.com AAAA

# Query CNAME record
ops-cli dig tw.yahoo.com CNAME

# Query ANY record
ops-cli dig google.com ANY

# Query PTR record
ops-cli dig 1.1.1.1 PTR`),
}

var digNetwork, digDomain, digServer string
var digOutput digResponse

func init() {
	rootCmd.AddCommand(digCmd)

	digCmd.Flags().StringVarP(&digNetwork, "net", "n", "tcp", "udp/tcp")
}

type digResponseFormat struct {
	NAME   string `json:"name" yaml:"name"`
	TTL    string `json:"ttl" yaml:"ttl"`
	CLASS  string `json:"class" yaml:"class"`
	TYPE   string `json:"type" yaml:"type"`
	RECORD string `json:"record" yaml:"record"`
}

type digResponse []digResponseFormat

func (d digResponse) Request(digType uint16) error {
	var err error
	/* If Query type is PTR, need to do reverse. */
	if dns.TypeToString[digType] == "PTR" {
		digDomain, err = dns.ReverseAddr(digDomain)
		if err != nil {
			return err
		}
		digDomain = strings.TrimRight(digDomain, ".")
	}

	var message = dns.Msg{}
	message.SetQuestion(digDomain+".", digType)
	var client = &dns.Client{Net: digNetwork}
	if digServer == "" {
		digServer, err = d.GetLocalServer()
		if err != nil {
			return err
		}
	}
	resp, _, err := client.Exchange(&message, digServer+":53")
	if err != nil {
		return err
	}
	if len(resp.Answer) == 0 {
		os.Exit(0)
	}

	for i := range resp.Answer {
		var d digResponseFormat
		elements := strings.Fields(fmt.Sprintf("%s ", resp.Answer[i]))
		if len(elements) == 5 {
			d = digResponseFormat{
				NAME:   elements[0],
				TTL:    elements[1],
				CLASS:  elements[2],
				TYPE:   elements[3],
				RECORD: elements[4],
			}
		} else if len(elements) > 5 {
			var remain string
			slice := elements[4:]
			for i := range slice {
				remain = remain + " " + slice[i]
			}

			d = digResponseFormat{
				NAME:   elements[0],
				TTL:    elements[1],
				CLASS:  elements[2],
				TYPE:   elements[3],
				RECORD: remain,
			}
		}
		digOutput = append(digOutput, d)
	}
	return err
}

func (d digResponse) GetLocalServer() (string, error) {
	const resolvConfig = "/etc/resolv.conf"
	s, err := dns.ClientConfigFromFile(resolvConfig)
	if err != nil {
		return "", err
	}
	return s.Servers[0], err
}

func (d digResponse) String() {
	for i := range d {
		fmt.Printf("%-20s\t%s\t%s\t%s\t%s\n", d[i].NAME, d[i].TTL, d[i].CLASS, d[i].TYPE, d[i].RECORD)
	}
}

func (d *digResponse) JSON() { PrintJSON(d) }

func (d *digResponse) YAML() { PrintYAML(d) }
