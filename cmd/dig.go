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
	"os"
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
	var digFlag DigFlag
	var digCmd = &cobra.Command{
		Use:   CommandDig + " [host] [@server] [type]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Resolve domain name",
		RunE:  digFlag.RunE,
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

	digCmd.Flags().StringVarP(&digFlag.network, "net", "n", "tcp", common.Usage("udp/tcp"))
}

type DigFlag struct {
	network string
	domain  string
	server  string
	output  digResponse
}

func (d *DigFlag) RunE(_ *cobra.Command, args []string) error {
	var err error
	var argsWithoutDomain []string
	var argsType []string
	switch lens := len(args); {
	case lens == 1:
		d.domain = args[0]
		if err = d.Request(dns.TypeA); err != nil {
			return err
		}
		if d.output == nil {
			return err
		}
		OutputInterfaceString(&d.output)
		return nil
	case lens > 1:
		/* Find which arg is domain. */
		for i := range args {
			if validator.ValidDomain(args[i]) || validator.ValidIP(args[i]) {
				d.domain = args[i]
				argsWithoutDomain = append(argsWithoutDomain, args[:i]...)
				argsWithoutDomain = append(argsWithoutDomain, args[i+1:]...)
				break
			}
		}
		switch lens {
		/* Distinguish remain arg is NS Server or DNS Type */
		case 2:
			if strings.Contains(argsWithoutDomain[0], "@") {
				d.server = strings.Replace(argsWithoutDomain[0], "@", "", 1)
				argsType = append(argsType, "A")
			} else {
				argsType = append(argsType, argsWithoutDomain[0])
			}
		default:
			for i := range argsWithoutDomain {
				if strings.Contains(argsWithoutDomain[i], "@") {
					d.server = strings.Replace(argsWithoutDomain[i], "@", "", 1)
					/* Copy args and remove @x.x.x.x */
					argsType = append(argsType, argsWithoutDomain[:i]...)
					argsType = append(argsType, argsWithoutDomain[i+1:]...)
					break
				}
			}
		}
	}
	err = d.Assertion(argsType[0])
	if err != nil {
		return err
	}
	if d.output == nil {
		return err
	}
	OutputInterfaceString(&d.output)
	return err
}

func (d *DigFlag) Request(digType uint16) error {
	var err error
	/* If Query type is PTR, need to do reverse. */
	if dns.TypeToString[digType] == "PTR" {
		d.domain, err = dns.ReverseAddr(d.domain)
		if err != nil {
			return err
		}
		d.domain = strings.TrimRight(d.domain, ".")
	}

	var message = dns.Msg{}
	message.SetQuestion(d.domain+".", digType)
	var client = &dns.Client{Net: d.network}
	if d.server == "" {
		d.server, err = d.GetLocalServer()
		if err != nil {
			return err
		}
	}
	resp, _, err := client.Exchange(&message, d.server+":53")
	if err != nil {
		return err
	}
	if len(resp.Answer) == 0 {
		return err
	}

	for i := range resp.Answer {
		var out digResponseFormat
		elements := strings.Fields(fmt.Sprintf("%s ", resp.Answer[i]))
		if len(elements) == 5 {
			out = digResponseFormat{
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

			out = digResponseFormat{
				NAME:   elements[0],
				TTL:    elements[1],
				CLASS:  elements[2],
				TYPE:   elements[3],
				RECORD: remain,
			}
		}
		d.output = append(d.output, out)
	}
	return err
}

func (d *DigFlag) GetLocalServer() (string, error) {
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

func (d *DigFlag) Assertion(arg string) error {
	var err error
	switch strings.ToLower(arg) {
	case "a":
		err = d.Request(dns.TypeA)
	case "aaaa":
		err = d.Request(dns.TypeAAAA)
	case "caa":
		err = d.Request(dns.TypeCAA)
	case "cname":
		err = d.Request(dns.TypeCNAME)
	case "mx":
		err = d.Request(dns.TypeMX)
	case "ns":
		err = d.Request(dns.TypeNS)
	case "ptr":
		err = d.Request(dns.TypePTR)
	case "soa":
		err = d.Request(dns.TypeSOA)
	case "srv":
		err = d.Request(dns.TypeSRV)
	case "txt":
		err = d.Request(dns.TypeTXT)
	case "any":
		err = d.Request(dns.TypeANY)
	}
	return err
}

type digResponseFormat struct {
	NAME   string `json:"name" yaml:"name"`
	TTL    string `json:"ttl" yaml:"ttl"`
	CLASS  string `json:"class" yaml:"class"`
	TYPE   string `json:"type" yaml:"type"`
	RECORD string `json:"record" yaml:"record"`
}

type digResponse []digResponseFormat

func (d digResponse) String() {
	var header []string
	var dd digResponseFormat
	f := reflect.ValueOf(&dd).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		header = append(header, t.Field(i).Name)
	}
	var data [][]string
	for i := range d {
		data = append(data, []string{d[i].NAME, d[i].TTL, d[i].CLASS, d[i].TYPE, d[i].RECORD})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
}
