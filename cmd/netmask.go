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
	"math/big"
	"net"
	"net/netip"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var netmaskFlag NetmaskFlag
	var netmaskCmd = &cobra.Command{
		Use:   CommandNetmask,
		Short: "Print IP/Mask pair, list address ranges",
		Run: func(_ *cobra.Command, args []string) {
			if netmaskFlag.ranges {
				for _, v := range args {
					if err := netmaskFlag.Range(v); err != nil {
						log.Println(err)
					}
				}
				return
			}
			if netmaskFlag.binary ||
				netmaskFlag.octal ||
				netmaskFlag.decimal ||
				netmaskFlag.hex ||
				netmaskFlag.cisco ||
				netmaskFlag.cidr {
				for _, v := range args {
					slice := strings.Split(v, "-")
					if len(slice) == 2 {
						if err := netmaskFlag.CIDR(slice[0], slice[1]); err != nil {
							log.Println(err)
						}
						continue
					}
					if err := netmaskFlag.Address(v); err != nil {
						log.Println(err)
					}
				}
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
	l := len(ipnet.IP)
	first := make(net.IP, l)
	last := make(net.IP, l)
	for i := 0; i < l; i++ {
		first[i] = ipnet.IP[i] & ipnet.Mask[i]
		last[i] = first[i] + (1<<8 - 1 - ipnet.Mask[i])
	}
	ones, bits := ipnet.Mask.Size()
	out := fmt.Sprintf("%v -> %v (%d)", first, last,
		big.NewInt(0).Lsh(big.NewInt(1), uint(bits-ones)))
	PrintString(out)
	return err
}

func (n *NetmaskFlag) Address(arg string) error {
	var err error
	ipnet := new(net.IPNet)
	switch {
	case !validator.ValidIP(arg) && validator.ValidIPCIDR(arg):
		_, ipnet, err = net.ParseCIDR(arg)
	case validator.ValidIP(arg) && !validator.ValidIPCIDR(arg):
		if validator.ValidIPv4(arg) {
			_, ipnet, err = net.ParseCIDR(arg + "/32")
		} else if validator.ValidIPv6(arg) {
			_, ipnet, err = net.ParseCIDR(arg + "/64")
		}
	default:
		return common.ErrInvalidArg
	}
	if err != nil {
		return err
	}

	var ip, mask string
	for i := 0; i < len(ipnet.IP); i++ {
		var f string
		switch {
		case n.binary:
			f = "%08b "
		case n.octal:
			f = "%o "
		case n.decimal:
			f = "%d."
		case n.hex:
			f = "%x "
		case n.cisco:
			ip += fmt.Sprintf("%d.", ipnet.IP[i])
			mask += fmt.Sprintf("%d.", 1<<8-1-ipnet.Mask[i])
			continue
		}
		ip += fmt.Sprintf(f, ipnet.IP[i])
		mask += fmt.Sprintf(f, ipnet.Mask[i])
	}
	switch {
	case n.binary, n.octal, n.hex:
		ip = strings.TrimRight(ip, " ")
		mask = strings.TrimRight(mask, " ")
	case n.decimal, n.cisco:
		ip = strings.TrimRight(ip, ".")
		mask = strings.TrimRight(mask, ".")
	}
	PrintString(fmt.Sprintf("%s / %s", ip, mask))
	return err
}

func (n *NetmaskFlag) CIDR(a, b string) error {
	var err error
	if (!validator.ValidIPv4(a) && !validator.ValidIPv4(b)) &&
		(!validator.ValidIPv6(a) && !validator.ValidIPv6(b)) {
		return common.ErrInvalidArg
	}

	ipa, err := netip.ParseAddr(a)
	if err != nil {
		return err
	}
	ipb, err := netip.ParseAddr(b)
	if err != nil {
		return err
	}

	var out []string
	var next, prefix string
	switch ipa.Compare(ipb) {
	case -1:
		for next != "0" {
			log.Printf("next: %s\n", next)
			next, prefix = n.calculate(ipa, ipb)
			if next == "" {
				continue
			}
			if next == "0" {
				break
			}
			ipa = netip.MustParseAddr(next)
			nextPrefix := netip.MustParsePrefix(prefix)
			for nextPrefix.Contains(ipa) && ipa.Compare(ipb) != 0 {
				ipa = ipa.Next()
			}
			if ipa.Compare(ipb) >= 0 && nextPrefix.Contains(ipb) {
				break
			}
			out = append(out, prefix)
		}
		out = append(out, prefix)
	case 1:
		n.calculate(ipb, ipa)
	case 0:
		out = append(out, fmt.Sprintf("%s/%d", ipa.String(), ipa.BitLen()))
	}

	for _, v := range out {
		PrintString(v)
	}
	return err
}

func (n *NetmaskFlag) calculate(ipa, ipb netip.Addr) (string, string) {
	temp := ipa
	for i := temp.BitLen(); i >= 0; i-- {
		if temp.Compare(ipb) == 1 {
			break
		}
		ipnet := netip.MustParsePrefix(fmt.Sprintf("%s/%d", ipa.String(), i))
		if ipnet.Contains(ipb) && !ipnet.Contains(ipb.Next()) {
			return "0", ipnet.String()
		}

		if ipnet.Contains(ipb) && ipnet.Contains(ipb.Next()) {
			return temp.String(), fmt.Sprintf("%s/%d", ipa.String(), i+1)
		}
		// log.Printf("i: %d\n", i)
		// log.Printf("ipnet: %v\n", ipnet)
		// log.Printf("temp: %v\n", temp)
		temp = temp.Next()
	}
	return "", ""
}
