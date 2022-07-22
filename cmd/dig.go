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

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
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
		digDomain = args[0]
		if lens == 1 {
			digNewClient(dns.TypeA)
			if digOutput != nil {
				digPrint()
			}
			return
		}

		if lens > 1 {
			argWithoutDomain := args[1:]
			var argType []string
			for i := range argWithoutDomain {
				if strings.Contains(argWithoutDomain[i], "@") {
					digServer = strings.Replace(argWithoutDomain[i], "@", "", 1)
					/* Copy args and remove @x.x.x.x */
					argType = append(argWithoutDomain[:i], argWithoutDomain[i+1:]...)
					break
				} else {
					argType = append(argType, argWithoutDomain[0])
				}
			}
			switch strings.ToLower(argType[0]) {
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
			default:
				digNewClient(dns.TypeA)
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
	NAME   string `json:"Name"`
	TTL    string `json:"TTL"`
	CLASS  string `json:"Class"`
	TYPE   string `json:"Type"`
	RECORD string `json:"Record"`
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

func digPrint() {
	if rootOutputJson {
		digOutput.Json()
		return
	}
	digOutput.String()
}
