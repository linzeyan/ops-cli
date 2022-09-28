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
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func init() {
	var flags struct {
		count    int
		protocol string
		timeout  time.Duration
	}
	var tcpingCmd = &cobra.Command{
		Use:   CommandTcping + " [host] [port]",
		Args:  cobra.ExactArgs(2),
		Short: "Connect to a port of a host",
		Run: func(_ *cobra.Command, args []string) {
			var t TCPing
			for i := 0; ; {
				if err := t.Connect(i, args, flags.protocol, flags.count, flags.timeout); err != nil {
					PrintString(err)
				}
				i++
				if i == flags.count {
					break
				}
				time.Sleep(time.Second)
			}
		},
	}
	rootCmd.AddCommand(tcpingCmd)
	tcpingCmd.Flags().IntVarP(&flags.count, "count", "c", 1, common.Usage("Specify tcping counts"))
	tcpingCmd.Flags().StringVarP(&flags.protocol, "protocol", "p", "tcp", common.Usage("Specify protocol"))
	tcpingCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type TCPing struct{}

func (t *TCPing) Connect(counter int, args []string, protocol string, limit int, timeout time.Duration) error {
	var p string
	startTime := time.Now()
	conn, err := net.DialTimeout(protocol, net.JoinHostPort(args[0], args[1]), timeout)
	if err != nil {
		p = fmt.Sprintf("Connect error: %s", err.Error())
		return errors.New(p)
	}

	if conn != nil {
		duration := time.Since(startTime)
		defer conn.Close()
		ip, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			return err
		}
		if limit != 1 {
			p = fmt.Sprintf("seq %d: %s response from %s (%s) port %s [open] %v", counter, protocol, args[0], ip, port, duration)
		} else {
			p = fmt.Sprintf("%s response from %s (%s) port %s [open] %v", protocol, args[0], ip, port, duration)
		}
		PrintString(p)
		return err
	}
	p = fmt.Sprintf("Connect error for seq %d", counter)
	return errors.New(p)
}
