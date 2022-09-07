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

	"github.com/spf13/cobra"
)

func init() {
	var netmaskFlag NetmaskFlag
	var netmaskCmd = &cobra.Command{
		Use:   CommandNetmask,
		Short: "Netmask",
		Run: func(_ *cobra.Command, args []string) {
			for _, v := range args {
				if netmaskFlag.ranges {
					err := netmaskFlag.Range(v)
					if err != nil {
						log.Println(err)
					}
				}
			}
		},
	}

	rootCmd.AddCommand(netmaskCmd)
	netmaskCmd.Flags().BoolVarP(&netmaskFlag.ranges, "ranges", "r", false, "Print ip address ranges")
}

type NetmaskFlag struct {
	ranges bool
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
