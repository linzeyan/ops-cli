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
	"os"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func init() {
	var tracerouteFlag TracerouteFlag
	var tracerouteCmd = &cobra.Command{
		Use:   CommandTraceroute + " [host]",
		Args:  cobra.ExactArgs(1),
		Short: "Print the route packets trace to network host",
		Run:   tracerouteFlag.Run,
	}
	rootCmd.AddCommand(tracerouteCmd)
	tracerouteCmd.Flags().IntVarP(&tracerouteFlag.size, "size", "s", 24, common.Usage("Specify packet size"))
	tracerouteCmd.Flags().IntVarP(&tracerouteFlag.maxTTL, "max-ttl", "m", 64, common.Usage("Specify max hop"))
	tracerouteCmd.Flags().DurationVarP(&tracerouteFlag.interval, "interval", "i", 500*time.Millisecond, common.Usage("Specify interval"))
	tracerouteCmd.Flags().DurationVarP(&tracerouteFlag.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type TracerouteFlag struct {
	size, maxTTL, retry int
	interval            time.Duration
	timeout             time.Duration
}

func (t *TracerouteFlag) Run(cmd *cobra.Command, args []string) {
	host := args[0]

	var data RandomString
	data = data.GenerateString(t.size, LowercaseLetters)

	ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		PrintString(err)
		return
	}
	t.retry = 3
	if t.interval < 50*time.Millisecond {
		t.interval = 50 * time.Millisecond
	}

	header := fmt.Sprintf("traceroute to %s (%v), %d hops max, %d byte packets", host, ip, t.maxTTL, len(data))
	PrintString(header)
	if err = t.Connect(ip, data); err != nil {
		PrintString(err)
		return
	}
}

func (t *TracerouteFlag) listen() (*icmp.PacketConn, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	return conn, conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL|ipv4.FlagDst|ipv4.FlagInterface|ipv4.FlagSrc, true)
}

func (t *TracerouteFlag) Connect(addr *net.IPAddr, icmpData []byte) error {
	conn, err := t.listen()
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}

	data := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Data: icmpData},
	}

	reply := make([]byte, 1500)
	for i := 1; i <= t.maxTTL; i++ {
		data.Body.(*icmp.Echo).Seq = i
		b, err := data.Marshal(nil)
		if err != nil {
			return err
		}
		if err = conn.IPv4PacketConn().SetTTL(i); err != nil {
			return err
		}
		peer, err := t.sendPacket(i, conn, addr, b, reply)
		if err != nil {
			return err
		}
		if peer == addr.String() {
			break
		}
		time.Sleep(t.interval)
	}
	return err
}

func (t *TracerouteFlag) sendPacket(hop int, conn *icmp.PacketConn, addr *net.IPAddr, b, reply []byte) (string, error) {
	var err error
	var ip string
	for i := 1; i <= t.retry; i++ {
		/* Send packet. */
		startTime := time.Now()
		_, err = conn.IPv4PacketConn().WriteTo(b, nil, addr)
		if err != nil {
			return "", err
		}
		/* Wait receiving. */
		if err = conn.SetReadDeadline(time.Now().Add(t.timeout)); err != nil {
			return "", err
		}
		n, cm, peer, err := conn.IPv4PacketConn().ReadFrom(reply)
		if err != nil {
			if i == 1 {
				fmt.Printf("%-4d%s*", hop, common.IndentTwoSpaces)
			} else {
				fmt.Printf("%s*", common.IndentTwoSpaces)
			}
			if i == t.retry {
				fmt.Print("\n")
			}
			continue
		}
		duration := time.Since(startTime)
		result, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			return peer.String(), err
		}
		var out string
		switch result.Type {
		case ipv4.ICMPTypeEchoReply, ipv4.ICMPTypeTimeExceeded:
			if i == 1 {
				out = fmt.Sprintf("%-4d  %-16v\t%v", hop, cm.Src, duration)
			} else {
				out = fmt.Sprintf("\t%v", duration)
			}
		case ipv4.ICMPTypeDestinationUnreachable:
			if i == 1 {
				out = fmt.Sprintf("%-4d%s*", hop, common.IndentTwoSpaces)
			} else {
				out = fmt.Sprintf("%s*", common.IndentTwoSpaces)
			}
		}
		fmt.Print(out)
		if i == t.retry {
			fmt.Print("\n")
		}
		ip = peer.String()
		time.Sleep(t.interval)
	}
	return ip, err
}
