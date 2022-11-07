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

func initTCPing() *cobra.Command {
	var flags struct {
		count    int
		protocol string
		timeout  time.Duration
	}
	var tcpingCmd = &cobra.Command{
		GroupID: groupings[CommandTCPing],
		Use:     CommandTCPing + " [host] [port]",
		Args:    cobra.ExactArgs(2),
		Short:   "Connect to a port of a host",
		Run: func(_ *cobra.Command, args []string) {
			t := TCPing{
				Protocal: flags.protocol,
				Timeout:  flags.timeout,
			}
			for i := 0; ; {
				addr, duration, err := t.Connect(i, args)
				if err != nil {
					logger.Info(err.Error())
					printer.Error(err)
				}
				ip, port, err := net.SplitHostPort(addr)
				if err != nil {
					logger.Info(err.Error())
					printer.Error(err)
				}
				if flags.count == 1 {
					printer.Printf("%s response from %s (%s) port %s [open] %v", t.Protocal, args[0], ip, port, duration)
					return
				}
				printer.Printf("seq %d: %s response from %s (%s) port %s [open] %v", i, t.Protocal, args[0], ip, port, duration)
				i++
				if i == flags.count {
					break
				}
				time.Sleep(time.Second)
			}
		},
	}
	tcpingCmd.Flags().IntVarP(&flags.count, "count", "c", 1, common.Usage("Specify tcping counts"))
	tcpingCmd.Flags().StringVarP(&flags.protocol, "protocol", "p", "tcp", common.Usage("Specify protocol"))
	tcpingCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
	return tcpingCmd
}

type TCPing struct {
	Protocal string
	Timeout  time.Duration
}

func (t *TCPing) Connect(counter int, args []string) (string, time.Duration, error) {
	var p string
	startTime := time.Now()
	conn, err := net.DialTimeout(t.Protocal, net.JoinHostPort(args[0], args[1]), t.Timeout)
	if err != nil {
		logger.Debug(err.Error())
		p = fmt.Sprintf("Connect error: %s", err.Error())
		return "", 0, errors.New(p)
	}

	if conn != nil {
		duration := time.Since(startTime)
		defer conn.Close()
		return conn.RemoteAddr().String(), duration, err
	}
	p = fmt.Sprintf("Connect error for seq %d", counter)
	return "", 0, errors.New(p)
}
