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
	tcpingCmd.Flags().BoolVarP(&tcpingFlag.continues, "continues", "c", false, common.Usage("Specify continues connect"))
	tcpingCmd.Flags().StringVarP(&tcpingFlag.protocol, "protocol", "p", "tcp", common.Usage("Specify protocol"))
	tcpingCmd.Flags().DurationVarP(&tcpingFlag.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type TcpingFlag struct {
	continues bool
	protocol  string
	timeout   time.Duration
}

func (t *TcpingFlag) Run(cmd *cobra.Command, args []string) {
	var counter int
	if !t.continues {
		t.Connect(counter, args)
		return
	}
	for {
		t.Connect(counter, args)
		time.Sleep(time.Second)
		counter++
	}
}

func (t *TcpingFlag) Connect(counter int, args []string) {
	conn, err := net.DialTimeout(t.protocol, net.JoinHostPort(args[0], args[1]), t.timeout)
	if err != nil {
		PrintString(err)
		return
	}
	if conn != nil {
		defer conn.Close()
		var p string
		if t.continues {
			p = fmt.Sprintf("seq=%d %s port %s open.", counter, args[0], args[1])
		} else {
			p = fmt.Sprintf("%s port %s open.", args[0], args[1])
		}
		PrintString(p)
	}
}
