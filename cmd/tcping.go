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
	var tcpingFlag TcpingFlag
	var tcpingCmd = &cobra.Command{
		Use:   CommandTcping + " [host] [port]",
		Args:  cobra.ExactArgs(2),
		Short: "Connect to a port of a host",
		Run:   tcpingFlag.Run,
	}
	rootCmd.AddCommand(tcpingCmd)
	tcpingCmd.Flags().IntVarP(&tcpingFlag.count, "count", "c", 1, common.Usage("Specify tcping counts"))
	tcpingCmd.Flags().StringVarP(&tcpingFlag.protocol, "protocol", "p", "tcp", common.Usage("Specify protocol"))
	tcpingCmd.Flags().DurationVarP(&tcpingFlag.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type TcpingFlag struct {
	count    int
	protocol string
	timeout  time.Duration
}

func (t *TcpingFlag) Run(cmd *cobra.Command, args []string) {
	for i := 0; i != t.count; i++ {
		if err := t.Connect(i, args); err != nil {
			PrintString(err)
		}
		time.Sleep(time.Second)
	}
}

func (t *TcpingFlag) Connect(counter int, args []string) error {
	var p string
	startTime := time.Now()
	conn, err := net.DialTimeout(t.protocol, net.JoinHostPort(args[0], args[1]), t.timeout)
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
		if t.count != 1 {
			p = fmt.Sprintf("seq %d: %s response from %s (%s) port %s [open] %v", counter, t.protocol, args[0], ip, port, duration)
		} else {
			p = fmt.Sprintf("%s response from %s (%s) port %s [open] %v", t.protocol, args[0], ip, port, duration)
		}
		PrintString(p)
		return err
	}
	p = fmt.Sprintf("Connect error for seq %d", counter)
	return errors.New(p)
}
