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

func initNetmask() *cobra.Command {
	var flags struct {
		ranges bool

		binary  bool
		octal   bool
		decimal bool
		hex     bool
		cisco   bool
		cidr    bool
	}
	var netmaskCmd = &cobra.Command{
		Use:   CommandNetmask,
		Short: "Print IP/Mask pair, list address ranges",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 0 {
				return cmd.Help()
			}
			var n Netmask
			if flags.ranges {
				for _, v := range args {
					if err = n.Range(v); err != nil {
						PrintString(err)
					}
				}
				return err
			}

			if flags.binary ||
				flags.octal ||
				flags.decimal ||
				flags.hex ||
				flags.cisco ||
				flags.cidr {
				var typ string
				switch {
				case flags.binary:
					typ = TypeBinary
				case flags.octal:
					typ = TypeOctal
				case flags.decimal:
					typ = TypeDecimal
				case flags.hex:
					typ = TypeHex
				case flags.cisco:
					typ = TypeCisco
				}
				for _, v := range args {
					slice := strings.Split(v, "-")
					if len(slice) == 2 {
						if err = n.CIDR(slice[0], slice[1], typ); err != nil {
							PrintString(err)
						}
						continue
					}
					if err = n.Address(v, typ); err != nil {
						PrintString(err)
					}
				}
			}
			return err
		},
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

	netmaskCmd.Flags().BoolVarP(&flags.ranges, "ranges", "r", false, common.Usage("Print IP address ranges"))
	netmaskCmd.Flags().BoolVarP(&flags.binary, TypeBinary, "b", false, common.Usage("Print IP address and mask in binary"))
	netmaskCmd.Flags().BoolVarP(&flags.octal, TypeOctal, "o", false, common.Usage("Print IP address and mask in octal"))
	netmaskCmd.Flags().BoolVarP(&flags.decimal, TypeDecimal, "d", false, common.Usage("Print IP address and mask in decimal"))
	netmaskCmd.Flags().BoolVarP(&flags.hex, TypeHex, "x", false, common.Usage("Print IP address and mask in hex"))
	netmaskCmd.Flags().BoolVarP(&flags.cisco, TypeCisco, "i", false, common.Usage("Print Cisco style address lists"))
	netmaskCmd.Flags().BoolVarP(&flags.cidr, "cidr", "c", false, common.Usage("Print CIDR format address lists"))
	return netmaskCmd
}

type Netmask struct{}

/* ipRange parse string and return CIDR, first IP and last IP. */
func (*Netmask) ipRange(arg string) (*net.IPNet, netip.Addr, netip.Addr) {
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

func (n *Netmask) Range(arg string) error {
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

func (n *Netmask) Address(arg, typ string) error {
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
		switch typ {
		case TypeBinary:
			f = "%08b "
		case TypeOctal:
			f = "%o "
		case TypeDecimal:
			f = "%d."
		case TypeHex:
			f = "%x "
		case TypeCisco:
			ip += fmt.Sprintf("%d.", ipnet.IP[i])
			mask += fmt.Sprintf("%d.", ipnet.Mask[i]^0xff)
			continue
		}
		ip += fmt.Sprintf(f, ipnet.IP[i])
		mask += fmt.Sprintf(f, ipnet.Mask[i])
	}
	switch typ {
	case TypeBinary, TypeOctal, TypeHex:
		ip = strings.TrimRight(ip, " ")
		mask = strings.TrimRight(mask, " ")
	case TypeDecimal, TypeCisco:
		ip = strings.TrimRight(ip, ".")
		mask = strings.TrimRight(mask, ".")
	}
	PrintString(fmt.Sprintf("%s / %s", ip, mask))
	return err
}

/* iterate if not iterate over return next IP and prefix, else return "0" and prefix. */
func (n *Netmask) iterate(ipa, ipb netip.Addr) (string, string) {
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

func (n *Netmask) CIDR(a, b, typ string) error {
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
	default:
		for next != "0" {
			next, prefix = n.iterate(ipa, ipb)
			if next == "0" || next == "" {
				break
			}
			ipa = netip.MustParseAddr(next)
			out = append(out, prefix)
		}
	case 1:
		for next != "0" {
			next, prefix = n.iterate(ipb, ipa)
			if next == "0" || next == "" {
				break
			}
			ipb = netip.MustParseAddr(next)
			out = append(out, prefix)
		}
	}
	out = append(out, prefix)

	for _, v := range out {
		if typ == TypeCisco {
			if err = n.Address(v, typ); err != nil {
				return err
			}
		} else {
			PrintString(v)
		}
	}
	return err
}
