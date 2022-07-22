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
	_ "embed"
	"log"
	"strings"

	"github.com/linzeyan/whois"
	"github.com/spf13/cobra"
)

var whoisCmd = &cobra.Command{
	Use:   "whois",
	Short: "List domain name information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}
		var whoisDomain = args[0]
		var resp *whois.Response
		var err error
		var data whois.Servers
		if whoisDomain != "" {
			switch strings.ToLower(whoisServer) {
			case "whoisxml":
				data = whois.WhoisXML{}
			case "ip2whois":
				data = whois.Ip2Whois{}
			case "whoapi":
				data = whois.WhoApi{}
			case "apininjas":
				data = whois.ApiNinjas{}
			default:
				data = whois.Verisign{}
			}
			resp, err = whois.Request(data, whoisDomain)
			if err != nil {
				log.Println(err)
				return
			}
			if rootOutputJson {
				resp.Json()
			} else if rootOutputYaml {
				resp.Yaml()
			} else {
				resp.String()
			}
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Search domain
ops-cli whois apple.com

# Search domains using the specified whois server that requires an api key
ops-cli whois -s ApiNinjas -k your_api_key google.com`),
}

var whoisServer string

func init() {
	rootCmd.AddCommand(whoisCmd)

	whoisCmd.Flags().StringVarP(&whoisServer, "server", "s", "whois.verisign-grs.com", "Specify request server, can be WhoisXML, IP2Whois, WhoApi, ApiNinjas")
	whoisCmd.Flags().StringVarP(&whois.Key, "key", "k", "", "Specify API Key")
}
