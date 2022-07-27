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
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var digCmd = &cobra.Command{
	Use:   "dig [host] [@server] [type]",
	Short: "Resolve domain name",
	Run: func(cmd *cobra.Command, args []string) {
		var argsWithoutDomain []string
		var argsType []string
		switch lens := len(args); {
		case lens == 0:
			_ = cmd.Help()
			return
		case lens == 1:
			digDomain = args[0]
			digOutput.Request(dns.TypeA)
			if digOutput == nil {
				log.Println("response is empty")
				return
			}
			outputDefaultString(&digOutput)
			return
		case lens > 1:
			for i := range args {
				if govalidator.IsDNSName(args[i]) || govalidator.IsIP(args[i]) {
					digDomain = args[i]
					argsWithoutDomain = append(args[:i], args[i+1:]...)
					break
				}
			}
			switch lens {
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
			digOutput.Request(dns.TypeA)
		case "aaaa":
			digOutput.Request(dns.TypeAAAA)
		case "caa":
			digOutput.Request(dns.TypeCAA)
		case "cname":
			digOutput.Request(dns.TypeCNAME)
		case "mx":
			digOutput.Request(dns.TypeMX)
		case "ns":
			digOutput.Request(dns.TypeNS)
		case "ptr":
			digOutput.Request(dns.TypePTR)
		case "soa":
			digOutput.Request(dns.TypeSOA)
		case "srv":
			digOutput.Request(dns.TypeSRV)
		case "txt":
			digOutput.Request(dns.TypeTXT)
		case "any":
			digOutput.Request(dns.TypeANY)
		}
		if digOutput == nil {
			log.Println("response is empty")
			return
		}
		outputDefaultString(&digOutput)
	},
	Example: Examples(`# Query A record
ops-cli dig google.com
ops-cli dig google.com A
ops-cli dig google.com AAAA

# Query CNAME record
ops-cli dig tw.yahoo.com CNAME

# Query ANY record
ops-cli dig google.com ANY`),
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

func (digResponse) Request(digType uint16) {
	if dns.TypeToString[digType] == "PTR" {
		var err error
		digDomain, err = dns.ReverseAddr(digDomain)
		if err != nil {
			log.Println(err)
			return
		}
		digDomain = strings.TrimRight(digDomain, ".")
	}

	var message = dns.Msg{}
	message.SetQuestion(digDomain+".", digType)
	var client = dns.Client{Net: digNetwork}
	if digServer == "" {
		digServer = "8.8.8.8"
	}
	resp, _, err := client.Exchange(&message, digServer+":53")
	if err != nil {
		log.Println(err)
		return
	}
	if len(resp.Answer) == 0 {
		return
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
}
func (d digResponse) String() {
	for i := range d {
		fmt.Printf("%-20s\t%s\t%s\t%s\t%s\n", d[i].NAME, d[i].TTL, d[i].CLASS, d[i].TYPE, d[i].RECORD)
	}
}

func (d *digResponse) JSON() {
	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (d *digResponse) YAML() {
	out, err := yaml.Marshal(d)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}
