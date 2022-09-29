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
	"net"
	"reflect"
	"runtime"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/miekg/dns"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	var flags struct {
		network string
		domain  string
		server  string
	}
	var digCmd = &cobra.Command{
		Use:   CommandDig + " [host] [@server] [type]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Resolve domain name",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			var output DigList
			var argsWithoutDomain []string
			var argsType []string
			switch lens := len(args); {
			case lens == 1:
				flags.domain = args[0]
				if output, err = output.Request(dns.TypeA, flags.domain, flags.network, flags.server); err != nil {
					return err
				}
				if output == nil {
					return err
				}
				OutputInterfaceString(&output)
				return nil
			case lens > 1:
				/* Find which arg is domain. */
				for i := range args {
					if validator.ValidDomain(args[i]) || validator.ValidIP(args[i]) {
						flags.domain = args[i]
						argsWithoutDomain = append(argsWithoutDomain, args[:i]...)
						argsWithoutDomain = append(argsWithoutDomain, args[i+1:]...)
						break
					}
				}
				switch lens {
				/* Distinguish remain arg is NS Server or DNS Type */
				case 2:
					if strings.Contains(argsWithoutDomain[0], "@") {
						flags.server = strings.Replace(argsWithoutDomain[0], "@", "", 1)
						argsType = append(argsType, "A")
					} else {
						argsType = append(argsType, argsWithoutDomain[0])
					}
				default:
					for i := range argsWithoutDomain {
						if strings.Contains(argsWithoutDomain[i], "@") {
							flags.server = strings.Replace(argsWithoutDomain[i], "@", "", 1)
							/* Copy args and remove @x.x.x.x */
							argsType = append(argsType, argsWithoutDomain[:i]...)
							argsType = append(argsType, argsWithoutDomain[i+1:]...)
							break
						}
					}
				}
			}
			typ := dns.StringToType[strings.ToUpper(argsType[0])]
			output, err = output.Request(typ, flags.domain, flags.network, flags.server)
			if err != nil {
				return err
			}
			if output == nil {
				return err
			}
			OutputInterfaceString(output)
			return err
		},
		Example: common.Examples(`# Query A record
google.com
@1.1.1.1 google.com A
@8.8.8.8 google.com AAAA

# Query CNAME record
tw.yahoo.com CNAME

# Query ANY record
google.com ANY

# Query PTR record
1.1.1.1 PTR`, CommandDig),
	}
	rootCmd.AddCommand(digCmd)

	digCmd.Flags().StringVarP(&flags.network, "net", "n", "tcp", common.Usage("udp/tcp"))
}

type DigResponse struct {
	NAME   string `json:"name" yaml:"name"`
	TTL    string `json:"ttl" yaml:"ttl"`
	CLASS  string `json:"class" yaml:"class"`
	TYPE   string `json:"type" yaml:"type"`
	RECORD string `json:"record" yaml:"record"`
}

type DigList []DigResponse

func (d *DigList) GetLocalServer() (string, error) {
	if runtime.GOOS == "windows" {
		return "1.1.1.1", nil
	}
	const resolvConfig = "/etc/resolv.conf"
	s, err := dns.ClientConfigFromFile(resolvConfig)
	if err != nil {
		return "", err
	}
	return s.Servers[0], err
}

func (d *DigList) Request(digType uint16, domain, network, server string) (DigList, error) {
	var err error
	/* If Query type is PTR, need to do reverse. */
	if dns.TypeToString[digType] == "PTR" {
		domain, err = dns.ReverseAddr(domain)
		if err != nil {
			return nil, err
		}
		domain = strings.TrimRight(domain, ".")
	}

	var message = dns.Msg{}
	message.SetQuestion(domain+".", digType)
	var client = &dns.Client{Net: network}
	if server == "" {
		server, err = d.GetLocalServer()
		if err != nil {
			return nil, err
		}
	}
	resp, _, err := client.Exchange(&message, net.JoinHostPort(server, "53"))
	if err != nil {
		return nil, err
	}
	if len(resp.Answer) == 0 {
		return nil, err
	}

	for i := range resp.Answer {
		var out DigResponse
		elements := strings.Fields(fmt.Sprintf("%s ", resp.Answer[i]))
		if len(elements) == 5 {
			out = DigResponse{
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

			out = DigResponse{
				NAME:   elements[0],
				TTL:    elements[1],
				CLASS:  elements[2],
				TYPE:   elements[3],
				RECORD: remain,
			}
		}
		*d = append(*d, out)
	}
	return *d, err
}

func (d DigList) String() {
	var header []string
	var dd DigResponse
	f := reflect.ValueOf(&dd).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		header = append(header, t.Field(i).Name)
	}
	var data [][]string
	for i := range d {
		data = append(data, []string{d[i].NAME, d[i].TTL, d[i].CLASS, d[i].TYPE, d[i].RECORD})
	}
	PrintTable(header, data, tablewriter.ALIGN_LEFT, "\t", true)
}
