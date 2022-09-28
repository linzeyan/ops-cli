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
	"fmt"
	"net"

	"github.com/j-keck/arping"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var arpingFlag struct {
		check bool
		mac   bool
	}
	var arpingCmd = &cobra.Command{
		Use:   CommandArping,
		Args:  cobra.ExactArgs(1),
		Short: "Discover and probe hosts in a network using the ARP protocol",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if !validator.ValidIPv4(args[0]) {
				return common.ErrInvalidIP
			}
			ip := net.ParseIP(args[0])
			hwAddr, duration, err := arping.Ping(ip)
			if err != nil && !errors.Is(err, arping.ErrTimeout) {
				return err
			}
			switch {
			case arpingFlag.check:
				if errors.Is(err, arping.ErrTimeout) {
					PrintString("offline")
				} else {
					PrintString("online")
				}
			case arpingFlag.mac:
				PrintString(hwAddr)
			default:
				fmt.Printf("response from %s (%s): time=%s\n", ip, hwAddr, duration)
			}
			return err
		},
	}
	rootCmd.AddCommand(arpingCmd)
	arpingCmd.Flags().BoolVarP(&arpingFlag.check, "check", "c", false, "Check if host is online")
	arpingCmd.Flags().BoolVarP(&arpingFlag.mac, "mac", "m", false, "Resolve mac address")
}
