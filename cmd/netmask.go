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
	"net"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var netmaskFlag NetmaskFlag
	var netmaskCmd = &cobra.Command{
		Use:   CommandNetmask,
		Short: "Netmask",
		Run: func(_ *cobra.Command, args []string) {
			if netmaskFlag.ranges {
				for _, v := range args {
					err := netmaskFlag.Range(v)
					if err != nil {
						log.Println(err)
					}
				}
				return
			}
			err := netmaskFlag.Address(args)
			if err != nil {
				log.Println(err)
				return
			}
		},
	}

	rootCmd.AddCommand(netmaskCmd)
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.ranges, "ranges", "r", false, "Print ip address ranges")
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.binary, "binary", "b", false, "Print ip address in binary")
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.octal, "octal", "o", false, "Print ip address in octal")
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.decimal, "decimal", "d", false, "Print ip address in decimal")
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.hex, "hex", "x", false, "Print ip address in hex")
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.cisco, "cisco", "i", false, "Print Cisco style address lists")
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.cidr, "cidr", "c", false, "Print CIDR format address lists")
}

type NetmaskFlag struct {
	ranges bool

	binary  bool
	octal   bool
	decimal bool
	hex     bool
	cisco   bool
	cidr    bool
}

func (n *NetmaskFlag) Range(arg string) error {
	_, ipnet, err := net.ParseCIDR(arg)
	if err != nil {
		return err
	}
	ip := ipnet.IP
	mask := ipnet.Mask
	first := make(net.IP, len(ip))
	first[0] = ip[0] & mask[0]
	first[1] = ip[1] & mask[1]
	first[2] = ip[2] & mask[2]
	first[3] = ip[3] & mask[3]
	last := make(net.IP, len(ip))
	last[0] = first[0] + (255 - mask[0])
	last[1] = first[1] + (255 - mask[1])
	last[2] = first[2] + (255 - mask[2])
	last[3] = first[3] + (255 - mask[3])
	sum := (256 - int(mask[0])) * (256 - int(mask[1])) * (256 - int(mask[2])) * (256 - int(mask[3]))
	PrintString(fmt.Sprintf("%v -> %v (%d)", first, last, sum))
	return err
}

func (n *NetmaskFlag) Address(args []string) error {
	var err error
	for _, v := range args {
		ipnet := new(net.IPNet)
		switch {
		case !validator.ValidIP(v) && validator.ValidIPCIDR(v):
			_, ipnet, err = net.ParseCIDR(v)
		case validator.ValidIP(v) && !validator.ValidIPCIDR(v):
			if validator.ValidIPv4(v) {
				_, ipnet, err = net.ParseCIDR(v + "/32")
			} else if validator.ValidIPv6(v) {
				_, ipnet, err = net.ParseCIDR(v + "/64")
			}
		default:
			return common.ErrInvalidArg
		}
		if err != nil {
			return err
		}

		var ip, mask string
		for i := 0; i < len(ipnet.IP); i++ {
			switch {
			case n.binary:
				ip += fmt.Sprintf("%08b ", ipnet.IP[i])
				mask += fmt.Sprintf("%08b ", ipnet.Mask[i])
			case n.octal:
				ip += fmt.Sprintf("%o ", ipnet.IP[i])
				mask += fmt.Sprintf("%o ", ipnet.Mask[i])
			case n.decimal:
				ip += fmt.Sprintf("%d.", ipnet.IP[i])
				mask += fmt.Sprintf("%d.", ipnet.Mask[i])
			case n.hex:
				ip += fmt.Sprintf("%x ", ipnet.IP[i])
				mask += fmt.Sprintf("%x ", ipnet.Mask[i])
			case n.cisco:
				ip += fmt.Sprintf("%d.", ipnet.IP[i])
				mask += fmt.Sprintf("%d.", 255-ipnet.Mask[i])
			}
		}
		switch {
		case n.binary, n.octal, n.hex:
			PrintString(fmt.Sprintf("%s / %s", strings.TrimRight(ip, " "), strings.TrimRight(mask, " ")))
		case n.decimal, n.cisco:
			PrintString(fmt.Sprintf("%s / %s", strings.TrimRight(ip, "."), strings.TrimRight(mask, ".")))
		}
	}
	return err
}
