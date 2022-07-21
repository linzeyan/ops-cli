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
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns [host]",
	Short: "Resolve domain name",
	Run: func(cmd *cobra.Command, args []string) {
		var lens = len(args)
		if lens == 0 {
			cmd.Help()
			return
		}
		dnsDomain = args[0]
		var r dnsResolver
		r.dnsResolver()
		if lens == 1 {
			r.dnsTitle()
			r.dnsLookupA("a", dnsDomain)
			return
		}
		if lens > 1 {
			r.dnsTitle()
			for i := range args[1:] {
				switch t := strings.ToLower(args[1:][i]); t {
				case "a":
					r.dnsLookupA(t, dnsDomain)
				case "aaaa":
					r.dnsLookupA(t, dnsDomain)
				case "cname":
					r.dnsLookupCNAME(dnsDomain)
				case "mx":
					r.dnsLookupMX(dnsDomain)
				case "ns":
					r.dnsLookupNS(dnsDomain)
				case "ptr":
					r.dnsLookupAddr(dnsDomain)
				case "txt":
					r.dnsLookupTXT(dnsDomain)
				case "all":
					r.dnsLookupA(t, dnsDomain)
					r.dnsLookupCNAME(dnsDomain)
					r.dnsLookupMX(dnsDomain)
					r.dnsLookupNS(dnsDomain)
					r.dnsLookupTXT(dnsDomain)
				}
			}
			return
		}
	},
	Example: Examples(`# Query A record
ops-cli dns google.com

# Query CNAME record
ops-cli dns google.com CNAME

# Query PTR record
ops-cli dns 8.8.8.8 PTR`),
}

var dnsResolverServer, dnsDomain string

func init() {
	rootCmd.AddCommand(dnsCmd)

	dnsCmd.Flags().StringVarP(&dnsResolverServer, "resolver", "r", "", "Specify DNS server for lookup")
}

type dnsResolver struct {
	resolve *net.Resolver
}

func (d *dnsResolver) dnsResolver() {
	if dnsResolverServer != "" {
		d.resolve = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(2000),
				}
				return d.DialContext(ctx, network, dnsResolverServer+":53")
			},
		}
		return
	}
	d.resolve = net.DefaultResolver

}

func (d *dnsResolver) dnsLookupA(t, domain string) {
	res, err := d.resolve.LookupIP(context.Background(), "ip", domain)
	if err != nil {
		log.Println(err)
		return
	}
	if len(res) == 0 {
		return
	}
	ips := make(map[string]string)
	for i := range res {
		ip := strings.TrimSpace(fmt.Sprintf(` %s `, res[i]))
		if p := net.ParseIP(ip); p.To4() != nil {
			ips[ip] = "A"
		} else {
			ips[ip] = "AAAA"
		}
	}
	if t == "all" {
		for i := range ips {
			d.dnsPrintln(ips[i], i)
		}
		return
	}
	for i := range ips {
		if ips[i] == strings.ToUpper(t) {
			d.dnsPrintln(ips[i], i)
		}
	}
}

func (d *dnsResolver) dnsLookupAddr(ip string) {
	/* Valid IP */
	if net.ParseIP(ip) == nil {
		log.Println("not a valid IP")
		return
	}
	res, err := d.resolve.LookupAddr(context.Background(), ip)
	if err != nil {
		log.Println(err)
		return
	}
	if len(res) == 0 {
		return
	}
	for i := range res {
		d.dnsPrintln("PTR", res[i])
	}
}

func (d *dnsResolver) dnsLookupCNAME(domain string) {
	res, err := d.resolve.LookupCNAME(context.Background(), domain)
	if err != nil {
		log.Println(err)
		return
	}
	if res == "" || strings.TrimRight(res, ".") == domain {
		return
	}
	d.dnsPrintln("CNAME", res)
}

func (d *dnsResolver) dnsLookupMX(domain string) {
	res, err := d.resolve.LookupMX(context.Background(), domain)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return
		}
		log.Println(err)
		return
	}
	if len(res) == 0 {
		return
	}
	for i := range res {

		d.dnsPrintln("MX", fmt.Sprintf(`%d %s`, res[i].Pref, res[i].Host))
	}
}

func (d *dnsResolver) dnsLookupNS(domain string) {
	res, err := d.resolve.LookupNS(context.Background(), domain)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return
		}
		log.Println(err)
		return
	}
	if len(res) == 0 {
		return
	}
	for i := range res {
		d.dnsPrintln("NS", res[i].Host)
	}
}

func (d *dnsResolver) dnsLookupTXT(domain string) {
	res, err := d.resolve.LookupTXT(context.Background(), domain)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return
		}
		log.Println(err)
		return
	}
	if len(res) == 0 {
		return
	}
	for i := range res {
		d.dnsPrintln("TXT", res[i])
	}
}

func (d *dnsResolver) dnsTitle() {
	fmt.Println("Name", "\t\t\t", "Type", "\t\t", "Address")
}

func (d *dnsResolver) dnsPrintln(t, addr string) {
	fmt.Println(dnsDomain, "\t\t", strings.ToUpper(t), "\t\t", addr)
}
