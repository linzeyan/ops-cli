//go:build linux || windows

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
	"github.com/cakturk/go-netstat/netstat"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initSs())
}

func initSs() *cobra.Command {
	var flags struct {
		listen     bool
		ipv4, ipv6 bool
		tcp, udp   bool
	}
	var ssCmd = &cobra.Command{
		Use:   CommandSs,
		Short: "Displays sockets informations",
		Run: func(_ *cobra.Command, _ []string) {
			var fn netstat.AcceptFn
			switch {
			default:
				fn = netstat.NoopFilter
			case flags.listen:
				fn = func(s *netstat.SockTabEntry) bool {
					return s.State == netstat.Listen
				}
			}
			ip := "all"
			switch {
			case flags.ipv4:
				ip = "ipv4"
			case flags.ipv6:
				ip = "ipv6"
			}
			proto := "all"
			switch {
			case flags.tcp:
				proto = "tcp"
			case flags.udp:
				proto = "udp"
			}
			var s SS
			s.String(s.GetData(ip, proto, fn))
		},
	}
	ssCmd.Flags().BoolVarP(&flags.listen, "listen", "l", false, common.Usage("Only listening sockets"))
	ssCmd.Flags().BoolVarP(&flags.ipv4, "ipv4", "4", false, common.Usage("Only IPv4 sockets"))
	ssCmd.Flags().BoolVarP(&flags.ipv6, "ipv6", "6", false, common.Usage("Only IPv6 sockets"))
	ssCmd.Flags().BoolVarP(&flags.tcp, "tcp", "t", false, common.Usage("Only TCP sockets"))
	ssCmd.Flags().BoolVarP(&flags.udp, "udp", "u", false, common.Usage("Only UDP sockets"))
	return ssCmd
}

type SS struct{}

func (*SS) GetData(ip, proto string, fn netstat.AcceptFn) [][]string {
	var data [][]string
	var socks []netstat.SockTabEntry
	var err error
	if ip == "ipv4" || ip == "all" {
		if proto == "tcp" || proto == "all" {
			socks, err = netstat.TCPSocks(fn)
			if err == nil {
				for _, v := range socks {
					data = append(data, []string{"tcp", v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
				}
			}
		}
		if proto == "udp" || proto == "all" {
			socks, err = netstat.UDPSocks(fn)
			if err == nil {
				for _, v := range socks {
					data = append(data, []string{"udp", v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
				}
			}
		}
	}
	if ip == "ipv6" || ip == "all" {
		if proto == "tcp" || proto == "all" {
			socks, err = netstat.TCP6Socks(fn)
			if err == nil {
				for _, v := range socks {
					data = append(data, []string{"tcp6", v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
				}
			}
		}
		if proto == "udp" || proto == "all" {
			socks, err = netstat.UDP6Socks(fn)
			if err == nil {
				for _, v := range socks {
					data = append(data, []string{"udp6", v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
				}
			}
		}
	}
	return data
}

func (*SS) String(data [][]string) {
	header := []string{"Proto", "Local Address", "Foreign Address", "State", "PID/Program name"}
	PrintTable(header, data, tablewriter.ALIGN_RIGHT, "\t", false)
}
