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
	"strings"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

// digCmd represents the dig command
var digCmd = &cobra.Command{
	Use:   "dig",
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
				digTitle()
				for i := range digOutput {
					fmt.Println(digOutput[i])
				}
			}
			return
		}

		if lens > 1 {
			a := args[1:]
			for i := range a {
				if strings.Contains(a[i], "@") {
					digServer = strings.Replace(a[i], "@", "", 1)
					break
				}
			}
			for i := range a {
				switch strings.ToLower(a[i]) {
				case "a":
					digNewClient(dns.TypeA)
				case "aaaa":
					digNewClient(dns.TypeAAAA)
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
				}
			}
			if digOutput != nil {
				digTitle()
				for i := range digOutput {
					fmt.Println(digOutput[i])
				}
			}
			return
		}
	},
}

var digNetwork, digDomain, digServer string
var digOutput []string

func init() {
	rootCmd.AddCommand(digCmd)

	digCmd.Flags().StringVarP(&digNetwork, "net", "n", "udp", "udp/tcp")
}

func digNewClient(digType uint16) {
	var message = dns.Msg{}
	message.SetQuestion(digDomain+".", digType)
	var client = dns.Client{Net: digNetwork}
	if digServer == "" {
		digServer = "8.8.8.8"
	}
	resp, _, err := client.Exchange(&message, digServer+":53")
	if err != nil {
		log.Println(err)
		log.Println(digType)
		return
	}
	if len(resp.Answer) == 0 {
		return
	}
	for i := range resp.Answer {
		digOutput = append(digOutput, fmt.Sprintf("%s", resp.Answer[i]))
	}
}

func digTitle() {
	fmt.Println("NAME\t\tTTL\tCLASS\tTYPE\tADDRESS")
}
