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
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/spf13/cobra"
)

func init() {
	iface, err := net.Interfaces()
	if err != nil {
		return
	}
	var validArgs []string
	for _, v := range iface {
		validArgs = append(validArgs, v.Name)
	}
	validArgs = append(validArgs, "all")

	var ipCmd = &cobra.Command{
		Use:       CommandIP + " {all|interface}",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: validArgs,
		Short:     "View interfaces configuration",
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			counters, err := net.IOCounters(true)
			if err != nil {
				return err
			}
			idx, out := ParseInterfaces(iface, counters)
			switch args[0] {
			case "all":
				for i := 0; i <= len(out)+1; i++ {
					v, ok := out[i]
					if ok {
						PrintString(fmt.Sprintf("%d: %s", i, v))
					}
				}
			default:
				for _, value := range args {
					PrintString(fmt.Sprintf("%d: %s", idx[value], out[idx[value]]))
				}
			}
			return err
		},
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}
	rootCmd.AddCommand(ipCmd)
}

func ParseInterfaces(iface net.InterfaceStatList, counters []net.IOCountersStat) (map[string]int, map[int]string) {
	idx := make(map[string]int)
	out := make(map[int]string)
	for _, v := range iface {
		idx[v.Name] = v.Index
		var flag string
		for _, f := range v.Flags {
			flag = flag + strings.ToUpper(f) + ","
		}
		flag = strings.TrimRight(flag, ",")

		var value string
		if flag != "" {
			value = fmt.Sprintf("<%s>", flag)
		}

		if v.HardwareAddr != "" {
			value = fmt.Sprintf("%s mtu %d\n\tether %s", value, v.MTU, v.HardwareAddr)
		}

		var addr string
		for _, a := range v.Addrs {
			if validator.ValidIPv4CIDR(a.Addr) {
				addr += fmt.Sprintf("%s\n\tinet %s", addr, a.Addr)
			}
		}
		for _, a := range v.Addrs {
			if validator.ValidIPv6CIDR(a.Addr) {
				addr += fmt.Sprintf("\n\tinet6 %s", a.Addr)
			}
		}
		value = fmt.Sprintf("%s: %s%s", v.Name, value, addr)

		for _, vv := range counters {
			if v.Name == vv.Name && len(v.Addrs) != 0 {
				value = fmt.Sprintf("%s\n\tRX packets %d  bytes %d (%s)\n\tRX errors %d  dropped %d",
					value, vv.PacketsRecv, vv.BytesRecv, common.ByteSize(vv.BytesRecv).String(), vv.Errin, vv.Dropin)
				value = fmt.Sprintf("%s\n\tTX packets %d  bytes %d (%s)\n\tTX errors %d  dropped %d",
					value, vv.PacketsSent, vv.BytesSent, common.ByteSize(vv.BytesSent).String(), vv.Errout, vv.Dropout)
				break
			}
		}
		out[v.Index] = value
	}
	return idx, out
}
