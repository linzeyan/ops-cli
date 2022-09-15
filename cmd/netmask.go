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
		RunE:  netmaskFlag.RunE,
		Example: common.Examples(`# Print IP address and mask in different format
-b 1.2.3.4/32
-o 2.4.6.8/24
-d 192.168.0.0/16
-x 224.0.0.0/24

# Print IP address ranges
-r 100.90.1.9/17

# Print address lists
-c 1.1.1.1-2.2.2.2
-i 10.0.0.0-10.122.10.0`, CommandNetmask),
	}

	rootCmd.AddCommand(netmaskCmd)
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.ranges, "ranges", "r", false, common.Usage("Print IP address ranges"))
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.binary, "binary", "b", false, common.Usage("Print IP address and mask in binary"))
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.octal, "octal", "o", false, common.Usage("Print IP address and mask in octal"))
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.decimal, "decimal", "d", false, common.Usage("Print IP address and mask in decimal"))
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.hex, "hex", "x", false, common.Usage("Print IP address and mask in hex"))
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.cisco, "cisco", "i", false, common.Usage("Print Cisco style address lists"))
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.cidr, "cidr", "c", false, common.Usage("Print CIDR format address lists"))
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

func (n *NetmaskFlag) RunE(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return cmd.Help()
	}

	if n.ranges {
		for _, v := range args {
			if err = n.Range(v); err != nil {
				PrintString(err)
			}
		}
		return err
	}

	if n.binary ||
		n.octal ||
		n.decimal ||
		n.hex ||
		n.cisco ||
		n.cidr {
		for _, v := range args {
			slice := strings.Split(v, "-")
			if len(slice) == 2 {
				if err = n.CIDR(slice[0], slice[1]); err != nil {
					PrintString(err)
				}
				continue
			}
			if err = n.Address(v); err != nil {
				PrintString(err)
			}
		}
	}
	return err
}

/* ipRange parse string and return CIDR, first IP and last IP. */
func (*NetmaskFlag) ipRange(arg string) (*net.IPNet, netip.Addr, netip.Addr) {
	_, ipnet, err := net.ParseCIDR(arg)
	if err != nil {
		return nil, netip.Addr{}, netip.Addr{}
	}
	l := len(ipnet.IP)
	first := make(net.IP, l)
	last := make(net.IP, l)
	for i := 0; i < l; i++ {
		first[i] = ipnet.IP[i] & ipnet.Mask[i]
		last[i] = first[i] | (ipnet.Mask[i] ^ 0xff)
	}
	return ipnet, netip.MustParseAddr(first.String()), netip.MustParseAddr(last.String())
}

func (n *NetmaskFlag) Range(arg string) error {
	ipnet, first, last := n.ipRange(arg)
	if ipnet == nil {
		return common.ErrInvalidArg
	}
	ones, bits := ipnet.Mask.Size()
	out := fmt.Sprintf("%v -> %v (%d)", first, last,
		/* IP counts. */
		big.NewInt(0).Lsh(big.NewInt(1), uint(bits-ones)))
	PrintString(out)
	return nil
}

func (n *NetmaskFlag) Address(arg string) error {
	var (
		ipv4Len = fmt.Sprintf("/%d", net.IPv4len*8)
		ipv6Len = fmt.Sprintf("/%d", net.IPv6len*8)
	)
	var err error
	ipnet := new(net.IPNet)
	switch {
	case !validator.ValidIP(arg) && validator.ValidIPCIDR(arg):
		_, ipnet, err = net.ParseCIDR(arg)
	case validator.ValidIP(arg) && !validator.ValidIPCIDR(arg):
		if validator.ValidIPv4(arg) {
			_, ipnet, err = net.ParseCIDR(arg + ipv4Len)
		} else if validator.ValidIPv6(arg) {
			_, ipnet, err = net.ParseCIDR(arg + ipv6Len)
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
			mask += fmt.Sprintf("%d.", ipnet.Mask[i]^0xff)
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

/* iterate if not iterate over return next IP and prefix, else return "0" and prefix. */
func (n *NetmaskFlag) iterate(ipa, ipb netip.Addr) (string, string) {
	for i := ipa.BitLen(); i >= 0; i-- {
		p := fmt.Sprintf("%s/%d", ipa.String(), i)
		ipnet, first, last := n.ipRange(p)
		if first.Compare(ipa) == -1 || last.Compare(ipb) == 1 {
			p = fmt.Sprintf("%s/%d", ipa.String(), i+1)
			ipnet, _, last = n.ipRange(p)
			return last.Next().String(), ipnet.String()
		}
		if last.Compare(ipb) == 0 {
			return "0", ipnet.String()
		}
	}
	return "", ""
}

func (n *NetmaskFlag) CIDR(a, b string) error {
	var err error
	if (!validator.ValidIPv4(a) && !validator.ValidIPv4(b)) &&
		(!validator.ValidIPv6(a) && !validator.ValidIPv6(b)) {
		return common.ErrInvalidArg
	}

	ipa := netip.MustParseAddr(a)
	ipb := netip.MustParseAddr(b)

	var out []string
	var next, prefix string
	switch ipa.Compare(ipb) {
	case -1:
		for next != "0" {
			next, prefix = n.iterate(ipa, ipb)
			if next == "0" || next == "" {
				break
			}
			ipa = netip.MustParseAddr(next)
			out = append(out, prefix)
		}
		out = append(out, prefix)
	case 1:
		for next != "0" {
			next, prefix = n.iterate(ipb, ipa)
			if next == "0" || next == "" {
				break
			}
			ipb = netip.MustParseAddr(next)
			out = append(out, prefix)
		}
		out = append(out, prefix)
	case 0:
		out = append(out, fmt.Sprintf("%s/%d", ipa.String(), ipa.BitLen()))
	}

	for _, v := range out {
		if n.cisco {
			if err = n.Address(v); err != nil {
				return err
			}
		} else {
			PrintString(v)
		}
	}
	return err
}
