//go:build linux || darwin

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
	"errors"
	"net"
	"time"

	"github.com/j-keck/arping"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func initArping() *cobra.Command {
	var arpingFlag struct {
		check bool
		mac   bool
		iface string
	}
	var arpingCmd = &cobra.Command{
		Use:   CommandArping,
		Args:  cobra.ExactArgs(1),
		Short: "Discover and probe hosts in a network using the ARP protocol",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, args []string) {
			if !validator.IsIPv4(args[0]) {
				printer.Error(common.ErrInvalidIP)
				return
			}
			ip := net.ParseIP(args[0])
			var hwAddr net.HardwareAddr
			var err error
			if arpingFlag.iface != "" {
				hwAddr, _, err = arping.PingOverIfaceByName(ip, arpingFlag.iface)
			} else {
				hwAddr, _, err = arping.Ping(ip)
			}
			if err != nil && !errors.Is(err, arping.ErrTimeout) {
				printer.Error(err)
				return
			}

			switch {
			case arpingFlag.check:
				if errors.Is(err, arping.ErrTimeout) {
					printer.Printf("offline")
				} else {
					printer.Printf("online")
				}
			case arpingFlag.mac:
				printer.Printf(rootOutputFormat, hwAddr)
			default:
				var duration time.Duration
				var j int
				for i := 0; ; i++ {
					if arpingFlag.iface != "" {
						hwAddr, duration, err = arping.PingOverIfaceByName(ip, arpingFlag.iface)
					} else {
						hwAddr, duration, err = arping.Ping(ip)
					}
					if errors.Is(err, arping.ErrTimeout) {
						printer.Printf("seq=%d timeout", i)
						j++
						if j >= 5 {
							break
						}
						continue
					} else if err != nil {
						printer.Error(err)
						return
					}
					printer.Printf("response from %s (%s): index=%d time=%s\n", ip, hwAddr, i, duration)
					time.Sleep(time.Second)
				}
			}
		},
	}
	arpingCmd.Flags().BoolVarP(&arpingFlag.check, "check", "c", false, common.Usage("Check if host is online"))
	arpingCmd.Flags().BoolVarP(&arpingFlag.mac, "mac", "m", false, common.Usage("Resolve mac address"))
	arpingCmd.Flags().StringVarP(&arpingFlag.iface, "interface", "i", "", common.Usage("Specify interface name"))
	return arpingCmd
}
