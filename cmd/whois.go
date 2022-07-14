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
	"encoding/json"
	"fmt"
	"log"

	"github.com/linzeyan/whois"
	"github.com/spf13/cobra"
)

// whoisCmd represents the whois command
var whoisCmd = &cobra.Command{
	Use:   "whois",
	Short: "List domain name information",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		if whoisDomain != "" {
			switch whoisServer {
			case "WhoisXML", "whoisxml", "WHOISXML":
				if whoisKey != "" {
					whois.WhoisXMLAPIKey = whoisKey
				} else {
					whois.WhoisXMLAPIKey = whoisWhoisXMLAPIKey
				}
				result, err := whois.RequestWhoisXML(whoisDomain)
				if err != nil {
					log.Println(err)
					return
				}
				out, err := json.MarshalIndent(whois.ParserWhoisXML(result), "", "  ")
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(string(out))
			case "IP2Whois", "ip2whois", "IP2WHOIS":
				if whoisKey != "" {
					whois.IP2WhoisKey = whoisKey
				} else {
					whois.IP2WhoisKey = whoisIP2WhoisKey
				}
				result, err := whois.RequestIp2Whois(whoisDomain)
				if err != nil {
					log.Println(err)
					return
				}
				out, err := json.MarshalIndent(whois.ParserIp2Whois(result), "", "  ")
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(string(out))
			case "WhoApi", "whoapi", "WHOAPI":
				if whoisKey != "" {
					whois.WhoApiKey = whoisKey
				} else {
					whois.WhoApiKey = whoisWhoApiKey
				}
				result, err := whois.RequestWhoApi(whoisDomain)
				if err != nil {
					log.Println(err)
					return
				}
				out, err := json.MarshalIndent(whois.ParserWhoApi(result), "", "  ")
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(string(out))
			case "ApiNinjas", "apininjas", "APININJAS":
				if whoisKey != "" {
					whois.ApiNinjasKey = whoisKey
				} else {
					whois.ApiNinjasKey = whoisApiNinjasKey
				}
				result, err := whois.RequestApiNinjas(whoisDomain)
				if err != nil {
					log.Println(err)
					return
				}
				out, err := json.MarshalIndent(whois.ParserApiNinjas(result), "", "  ")
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(string(out))
			default:
				if whoisKey != "" {
					whois.ApiNinjasKey = whoisKey
				} else {
					whois.ApiNinjasKey = whoisApiNinjasKey
				}
				result, err := whois.RequestApiNinjas(whoisDomain)
				if err != nil {
					log.Println(err)
					return
				}
				out, err := json.MarshalIndent(whois.ParserApiNinjas(result), "", "  ")
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(string(out))
			}
			return
		}
		cmd.Help()
	},
}

//go:embed key_whoisxmlapi
var whoisWhoisXMLAPIKey string

//go:embed key_ip2whois
var whoisIP2WhoisKey string

//go:embed key_whoapi
var whoisWhoApiKey string

//go:embed key_apininjas
var whoisApiNinjasKey string

var whoisDomain, whoisServer, whoisKey string

func init() {
	rootCmd.AddCommand(whoisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	whoisCmd.Flags().StringVarP(&whoisDomain, "domain", "d", "", "Specify domain")
	whoisCmd.Flags().StringVarP(&whoisServer, "server", "s", "ApiNinjas", "Specify request server, can be WhoisXML, IP2Whois, WhoApi, ApiNinjas")
	whoisCmd.Flags().StringVarP(&whoisKey, "key", "k", "", "Specify API Key")
}
