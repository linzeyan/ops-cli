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
	var flags struct {
		size, maxTTL int
		interval     time.Duration
		timeout      time.Duration
	}
	var tracerouteCmd = &cobra.Command{
		Use:  CommandTraceroute + " [host]",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Print the route packets trace to network host",
		Run: func(_ *cobra.Command, args []string) {
			if flags.size <= 0 || flags.maxTTL <= 0 {
				return
			}
			if flags.interval < 50*time.Millisecond {
				flags.interval = 50 * time.Millisecond
			}

			t := Traceroute{
				Size:     flags.size,
				TTL:      flags.maxTTL,
				Retry:    3,
				Interval: flags.interval,
				Timeout:  flags.timeout,
			}
			host := args[0]
			data := Randoms.GenerateString(t.Size, LowercaseLetters)

			ip, err := net.ResolveIPAddr("ip4", host)
			if err != nil {
				PrintString(err)
				return
			}

			conn, err := t.listen()
			if err != nil {
				PrintString(err)
				return
			}
			if conn != nil {
				defer conn.Close()
			}

			icmpData := icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Data: data},
			}
			header := fmt.Sprintf("traceroute to %s (%v), %d hops max, %d byte packets", host, ip, t.TTL, len(data))
			PrintString(header)
			if err = t.Connect(conn, ip, icmpData); err != nil {
				PrintString(err)
				return
			}
		},
	}
	rootCmd.AddCommand(tracerouteCmd)
	tracerouteCmd.Flags().IntVarP(&flags.size, "size", "s", 24, common.Usage("Specify packet size"))
	tracerouteCmd.Flags().IntVarP(&flags.maxTTL, "max-ttl", "m", 64, common.Usage("Specify max hop"))
	tracerouteCmd.Flags().DurationVarP(&flags.interval, "interval", "i", 500*time.Millisecond, common.Usage("Specify interval"))
	tracerouteCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type Traceroute struct {
	Size, TTL, Retry int
	Interval         time.Duration
	Timeout          time.Duration
}

func (*Traceroute) listen() (*icmp.PacketConn, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	return conn, conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL|ipv4.FlagDst|ipv4.FlagInterface|ipv4.FlagSrc, true)
}

func (t *Traceroute) Connect(conn *icmp.PacketConn, addr *net.IPAddr, data icmp.Message) error {
	var err error
	reply := make([]byte, 1500)
	for i := 1; i <= t.TTL; i++ {
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
		time.Sleep(t.Interval)
	}
	return err
}

func (t *Traceroute) sendPacket(hop int, conn *icmp.PacketConn, addr *net.IPAddr, b, reply []byte) (string, error) {
	var err error
	var ip string
	var rtt []string
	for i := 1; i <= t.Retry; i++ {
		/* Send packet. */
		startTime := time.Now()
		_, err = conn.IPv4PacketConn().WriteTo(b, nil, addr)
		if err != nil {
			return "", err
		}
		/* Wait receiving. */
		if err = conn.SetReadDeadline(time.Now().Add(t.Timeout)); err != nil {
			return "", err
		}
		n, _, peer, err := conn.IPv4PacketConn().ReadFrom(reply)
		if err != nil {
			rtt = append(rtt, "*")
			continue
		}
		duration := time.Since(startTime)
		result, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			return peer.String(), err
		}
		switch result.Type {
		case ipv4.ICMPTypeEchoReply, ipv4.ICMPTypeTimeExceeded:
			rtt = append(rtt, duration.String())
		case ipv4.ICMPTypeDestinationUnreachable:
			rtt = append(rtt, "*")
		}
		if peer.String() != "" {
			ip = peer.String()
		}
		if i == t.Retry {
			break
		}
		time.Sleep(t.Interval)
	}
	out := fmt.Sprintf("%-4d  %-16v\t%-10s\t%-10s\t%-10s", hop, ip, rtt[0], rtt[1], rtt[2])
	PrintString(out)
	return ip, err
}
