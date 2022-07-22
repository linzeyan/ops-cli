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

// digCmd represents the dig command
var digCmd = &cobra.Command{
	Use:   "dig [host] [@server] [type]",
	Short: "Resolve domain name",
	Run: func(cmd *cobra.Command, args []string) {
		var lens = len(args)
		if lens == 0 {
			cmd.Help()
			return
		}

		if lens == 1 {
			digDomain = args[0]
			digNewClient(dns.TypeA)
			if digOutput != nil {
				digPrint()
			}
			return
		}

		if lens > 1 {
			var argsWithoutDomain []string
			var argsType []string

			for i := range args {
				if govalidator.IsDNSName(args[i]) {
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

			switch strings.ToLower(argsType[0]) {
			case "a":
				digNewClient(dns.TypeA)
			case "aaaa":
				digNewClient(dns.TypeAAAA)
			case "caa":
				digNewClient(dns.TypeCAA)
			case "cname":
				digNewClient(dns.TypeCNAME)
			case "mx":
				digNewClient(dns.TypeMX)
			case "ns":
				digNewClient(dns.TypeNS)
			case "ptr":
				digNewClient(dns.TypePTR)
			case "soa":
				digNewClient(dns.TypeSOA)
			case "srv":
				digNewClient(dns.TypeSRV)
			case "txt":
				digNewClient(dns.TypeTXT)
			case "any":
				digNewClient(dns.TypeANY)
			}
			if digOutput != nil {
				digPrint()
			}
			return
		}
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
var digOutput digResponseOutput

func init() {
	rootCmd.AddCommand(digCmd)

	digCmd.Flags().StringVarP(&digNetwork, "net", "n", "tcp", "udp/tcp")
}

func digNewClient(digType uint16) {
	// if dns.TypeToString[digType] == "PTR" {
	// 	var err error
	// 	digDomain, err = dns.ReverseAddr(digDomain)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }

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
		elements := strings.Fields(fmt.Sprintf("%s ", resp.Answer[i]))
		var d = digResponseFormat{
			NAME:   elements[0],
			TTL:    elements[1],
			CLASS:  elements[2],
			TYPE:   elements[3],
			RECORD: elements[4],
		}
		digOutput = append(digOutput, d)
	}
}

type digResponseFormat struct {
	NAME   string `json:"Name" yaml:"Name"`
	TTL    string `json:"TTL" yaml:"TTL"`
	CLASS  string `json:"Class" yaml:"Class"`
	TYPE   string `json:"Type" yaml:"Type"`
	RECORD string `json:"Record" yaml:"Record"`
}

type digResponseOutput []digResponseFormat

func (d digResponseOutput) String() {
	for i := range d {
		fmt.Printf("%-20s\t%s\t%s\t%s\t%s\n", d[i].NAME, d[i].TTL, d[i].CLASS, d[i].TYPE, d[i].RECORD)
	}

}

func (d *digResponseOutput) Json() {
	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func (d *digResponseOutput) Yaml() {
	out, err := yaml.Marshal(d)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func digPrint() {
	if rootOutputJson {
		digOutput.Json()
		return
	}
	if rootOutputYaml {
		digOutput.Yaml()
		return
	}
	digOutput.String()
}
