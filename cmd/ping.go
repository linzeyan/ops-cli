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
	"math"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func init() {
	var flags struct {
		ipv6             bool
		count, size, ttl int
		interval         time.Duration
		timeout          time.Duration
	}
	var pingCmd = &cobra.Command{
		Use:  CommandPing + " [host]",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Send ICMP ECHO_REQUEST packets to network hosts",
		Run: func(_ *cobra.Command, args []string) {
			if flags.count == 0 || flags.ttl <= 0 || flags.size <= 0 {
				return
			}
			if flags.interval < 50*time.Millisecond {
				flags.interval = 50 * time.Millisecond
			}

			host := args[0]
			var p = Ping{
				IPv6:     flags.ipv6,
				TTL:      flags.ttl,
				Interval: flags.interval,
				Timeout:  flags.timeout,
				Data:     Randoms.GenerateString(flags.size, LowercaseLetters),
			}

			network := "ip4"
			if flags.ipv6 {
				network = "ip6"
			}
			ip, err := net.ResolveIPAddr(network, host)
			if err != nil {
				PrintString(err)
				return
			}

			conn, err := p.Listen()
			if err != nil {
				PrintString(err)
				return
			}
			if conn != nil {
				defer conn.Close()
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			startTime := time.Now()
			for i := 0; ; i++ {
				if i == 0 {
					header := fmt.Sprintf("PING %s (%v): %d data bytes", host, ip, len(p.Data))
					PrintString(header)
				}
				if err := p.Connect(conn, i, ip, p.Data); err != nil {
					PrintString(err)
				}
				if i == flags.count-1 {
					p.summary(host, time.Since(startTime))
					return
				}
				time.Sleep(flags.interval)
				select {
				default:
				case <-quit:
					p.summary(host, time.Since(startTime))
					return
				}
			}
		},
	}
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&flags.count, "count", "c", -1, common.Usage("Specify ping counts"))
	pingCmd.Flags().BoolVarP(&flags.ipv6, "ipv6", "6", false, common.Usage("Use ICMPv6"))
	pingCmd.Flags().IntVarP(&flags.size, "size", "s", 56, common.Usage("Specify packet size"))
	pingCmd.Flags().IntVarP(&flags.ttl, "ttl", "", 64, common.Usage("Specify packet ttl"))
	pingCmd.Flags().DurationVarP(&flags.interval, "interval", "i", time.Second, common.Usage("Specify interval"))
	pingCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type Ping struct {
	IPv6              bool
	TTL               int
	Interval, Timeout time.Duration
	Data              []byte

	lost bool
	sta  struct {
		send, loss, receive int
		min, avg, max       time.Duration
		rtts                []time.Duration
	}
}

func (p *Ping) Listen() (*icmp.PacketConn, error) {
	if p.IPv6 {
		conn, err := icmp.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil {
			return nil, err
		}
		if err = conn.IPv6PacketConn().SetHopLimit(p.TTL); err != nil {
			return nil, err
		}
		return conn, conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
	}
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	if err = conn.IPv4PacketConn().SetTTL(p.TTL); err != nil {
		return nil, err
	}
	return conn, conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
}

func (p *Ping) statistics(duration time.Duration) {
	if p.sta.min == 0 {
		p.sta.min = duration
	}
	if p.sta.max == 0 {
		p.sta.max = duration
	}
	p.sta.avg += duration
	p.sta.send++
	p.sta.rtts = append(p.sta.rtts, duration)
	if p.lost {
		p.sta.loss++
		p.lost = false
	} else {
		p.sta.receive++
	}

	if duration < p.sta.min {
		p.sta.min = duration
	}
	if duration > p.sta.max {
		p.sta.max = duration
	}
}

func (p *Ping) readReply(conn *icmp.PacketConn, reply []byte, counter int) (int, any, net.Addr, error) {
	var n int
	var cm any
	var peer net.Addr
	var err error
	if p.IPv6 {
		n, cm, peer, err = conn.IPv6PacketConn().ReadFrom(reply)
	} else {
		n, cm, peer, err = conn.IPv4PacketConn().ReadFrom(reply)
	}
	if err != nil {
		p.lost = true
		p.statistics(0)
		e := fmt.Sprintf("Request timeout for icmp_seq %d", counter)
		return 0, nil, nil, errors.New(e)
	}
	return n, cm, peer, err
}

func (p *Ping) printMsg(result *icmp.Message, duration time.Duration, peer net.Addr, counter, size int, cm any) string {
	var ttl int
	switch c := cm.(type) {
	case *ipv4.ControlMessage:
		ttl = c.TTL
	case *ipv6.ControlMessage:
		ttl = c.HopLimit
	}

	var out string
	switch result.Type {
	case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
		out = fmt.Sprintf("%v bytes from %v: icmp_seq=%d ttl=%d time=%v", size, peer, counter, ttl, duration)
	case ipv4.ICMPTypeDestinationUnreachable, ipv6.ICMPTypeDestinationUnreachable:
		out = fmt.Sprintf("%v Destination Unreachable", peer)
	}
	p.statistics(duration)
	return out
}

func (p *Ping) Connect(conn *icmp.PacketConn, counter int, addr *net.IPAddr, icmpData []byte) error {
	data := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: counter & 0xffff, Seq: counter & 0xffff, Data: icmpData},
	}
	if p.IPv6 {
		data.Type = ipv6.ICMPTypeEchoRequest
	}
	b, err := data.Marshal(nil)
	if err != nil {
		return err
	}

	/* Send packet. */
	startTime := time.Now()
	_, err = conn.WriteTo(b, addr)
	if err != nil {
		return err
	}

	/* Wait receiving. */
	reply := make([]byte, 1500)
	if err = conn.SetReadDeadline(time.Now().Add(p.Timeout)); err != nil {
		return err
	}
	n, cm, peer, err := p.readReply(conn, reply, counter)
	if err != nil {
		return err
	}
	duration := time.Since(startTime)

	proto := 1
	if p.IPv6 {
		proto = 58
	}
	result, err := icmp.ParseMessage(proto, reply[:n])
	if err != nil {
		return err
	}
	out := p.printMsg(result, duration, peer, counter, len(b), cm)
	PrintString(out)

	return err
}

func (p *Ping) summary(host string, t time.Duration) {
	if p.sta.send == 0 {
		return
	}

	out := "\n"
	out += fmt.Sprintf("--- %s ping statistics ---\n", host)
	out += fmt.Sprintf("%d packets transmitted, %d received, %.1f%% packet loss, time %vms",
		p.sta.send, p.sta.receive, float64(p.sta.loss*100)/float64(p.sta.send), t.Milliseconds())
	if p.sta.send == p.sta.loss {
		PrintString(out)
		return
	}

	out += "\n"
	avg := p.sta.avg / time.Duration(p.sta.receive)
	var temp float64
	for _, v := range p.sta.rtts {
		temp += math.Pow(float64(v-avg), 2)
	}
	variance := temp / float64(len(p.sta.rtts))
	out += fmt.Sprintf("round-trip min/avg/max/mdev = %v/%v/%v/%v",
		p.sta.min, avg, p.sta.max, time.Duration(math.Sqrt(variance)))
	PrintString(out)
}
