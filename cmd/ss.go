//go:build linux || (windows && amd64)

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
			inet := All
			switch {
			case flags.ipv4:
				inet = IPv4
			case flags.ipv6:
				inet = IPv6
			}
			proto := All
			switch {
			case flags.tcp:
				proto = TCP
			case flags.udp:
				proto = UDP
			}
			var s SS
			s.String(s.GetData(inet, proto, fn))
		},
	}
	ssCmd.Flags().BoolVarP(&flags.listen, "listen", "l", false, common.Usage("Only listening sockets"))
	ssCmd.Flags().BoolVarP(&flags.ipv4, IPv4, "4", false, common.Usage("Only IPv4 sockets"))
	ssCmd.Flags().BoolVarP(&flags.ipv6, IPv6, "6", false, common.Usage("Only IPv6 sockets"))
	ssCmd.Flags().BoolVarP(&flags.tcp, TCP, "t", false, common.Usage("Only TCP sockets"))
	ssCmd.Flags().BoolVarP(&flags.udp, UDP, "u", false, common.Usage("Only UDP sockets"))
	return ssCmd
}

type SS struct{}

func (*SS) getTCPSocks(inet string, fn netstat.AcceptFn) [][]string {
	var data [][]string
	if inet == IPv4 || inet == All {
		socks, err := netstat.TCPSocks(fn)
		if err != nil {
			return nil
		}
		for _, v := range socks {
			data = append(data, []string{TCP, v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
		}
	}
	if inet == IPv4 {
		return data
	}
	socks6, err := netstat.TCP6Socks(fn)
	if err != nil {
		return nil
	}
	for _, v := range socks6 {
		data = append(data, []string{TCP6, v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
	}
	return data
}

func (*SS) getUDPSocks(inet string, fn netstat.AcceptFn) [][]string {
	var data [][]string
	if inet == IPv4 || inet == All {
		socks, err := netstat.UDPSocks(fn)
		if err != nil {
			return nil
		}
		for _, v := range socks {
			data = append(data, []string{UDP, v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
		}
		if inet == IPv4 {
			return data
		}
	}
	socks, _ := netstat.UDP6Socks(fn)
	for _, v := range socks {
		data = append(data, []string{UDP6, v.LocalAddr.String(), v.RemoteAddr.String(), v.State.String(), v.Process.String()})
	}
	return data
}

func (s *SS) GetData(inet, proto string, fn netstat.AcceptFn) [][]string {
	var data [][]string
	if proto == TCP || proto == All {
		data = s.getTCPSocks(inet, fn)
	}
	if proto == TCP {
		return data
	}
	return append(data, s.getUDPSocks(inet, fn)...)
}

func (*SS) String(data [][]string) {
	header := []string{"Proto", "Local Address", "Foreign Address", "State", "PID/Program name"}
	PrintTable(header, data, tablewriter.ALIGN_LEFT, "\t", false)
}
