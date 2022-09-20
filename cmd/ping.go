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
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func init() {
	var pingFlag PingFlag
	var pingCmd = &cobra.Command{
		Use:   CommandPing + " [host]",
		Args:  cobra.ExactArgs(1),
		Short: "Send ICMP ECHO_REQUEST packets to network hosts.",
		Run:   pingFlag.Run,
	}
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&pingFlag.count, "count", "c", -1, common.Usage("Specify ping counts"))
	pingCmd.Flags().BoolVarP(&pingFlag.ipv6, "ipv6", "6", false, common.Usage("Use ICMPv6"))
	pingCmd.Flags().IntVarP(&pingFlag.size, "size", "s", 24, common.Usage("Specify packet size"))
	pingCmd.Flags().IntVarP(&pingFlag.ttl, "ttl", "", 64, common.Usage("Specify packet ttl"))
	pingCmd.Flags().DurationVarP(&pingFlag.interval, "interval", "i", time.Second, common.Usage("Specify interval"))
	pingCmd.Flags().DurationVarP(&pingFlag.timeout, "timeout", "t", 5*time.Second, common.Usage("Specify timeout"))
}

type PingFlag struct {
	ipv6             bool
	count, size, ttl int
	interval         time.Duration
	timeout          time.Duration

	lost bool
	sta  struct {
		send, loss, receive int
		min, avg, max       time.Duration
		stddevRaw           []time.Duration
	}
}

func (p *PingFlag) Run(cmd *cobra.Command, args []string) {
	host := args[0]

	var data RandomString
	data = data.GenerateString(p.size, LowercaseLetters)

	network := "ip4"
	if p.ipv6 {
		network = "ip6"
	}
	ip, err := net.ResolveIPAddr(network, host)
	if err != nil {
		PrintString(err)
		return
	}

	header := fmt.Sprintf("PING %s (%v): %d data bytes", host, ip, len(data))
	PrintString(header)
	for i := 0; i != p.count; i++ {
		if err := p.Connect(i, ip, data); err != nil {
			PrintString(err)
		}
		if i == p.count {
			break
		}
		time.Sleep(p.interval)
	}
	p.output(host)
}

func (p *PingFlag) listen() (*icmp.PacketConn, error) {
	if p.ipv6 {
		conn, err := icmp.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil {
			return nil, err
		}
		if err = conn.IPv6PacketConn().SetHopLimit(p.ttl); err != nil {
			return nil, err
		}
		return conn, conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
	}
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	if err = conn.IPv4PacketConn().SetTTL(p.ttl); err != nil {
		return nil, err
	}
	return conn, conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
}

func (p *PingFlag) statistics(duration time.Duration) {
	if p.sta.min == 0 {
		p.sta.min = duration
	}
	if p.sta.max == 0 {
		p.sta.max = duration
	}
	p.sta.avg += duration
	p.sta.send++
	p.sta.stddevRaw = append(p.sta.stddevRaw, duration)

	if duration < p.sta.min {
		p.sta.min = duration
	}
	if duration > p.sta.max {
		p.sta.max = duration
	}
	if p.lost {
		p.sta.loss++
		p.lost = false
	} else {
		p.sta.receive++
	}
}

func (p *PingFlag) readReply(conn *icmp.PacketConn, reply []byte, counter int) (int, any, net.Addr, error) {
	var err error
	if p.ipv6 {
		n, cm, peer, err := conn.IPv6PacketConn().ReadFrom(reply)
		if err != nil {
			p.lost = true
			p.statistics(0)
			e := fmt.Sprintf("Request timeout for icmp_seq %d", counter)
			return 0, nil, nil, errors.New(e)
		}
		return n, cm, peer, err
	}
	n, cm, peer, err := conn.IPv4PacketConn().ReadFrom(reply)
	if err != nil {
		p.lost = true
		p.statistics(0)
		e := fmt.Sprintf("Request timeout for icmp_seq %d", counter)
		return 0, nil, nil, errors.New(e)
	}
	return n, cm, peer, err
}

func (p *PingFlag) printMsg(result *icmp.Message, duration time.Duration, peer net.Addr, counter, size int, cm any) string {
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

func (p *PingFlag) Connect(counter int, addr *net.IPAddr, icmpData []byte) error {
	conn, err := p.listen()
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}

	data := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: counter, Seq: counter, Data: icmpData},
	}
	if p.ipv6 {
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
	if err = conn.SetReadDeadline(time.Now().Add(p.timeout)); err != nil {
		return err
	}
	n, cm, peer, err := p.readReply(conn, reply, counter)
	if err != nil {
		return err
	}
	duration := time.Since(startTime)

	proto := 1
	if p.ipv6 {
		proto = 58
	}
	result, err := icmp.ParseMessage(proto, reply[:n])
	if err != nil {
		return err
	}
	size := len(b)
	out := p.printMsg(result, duration, peer, counter, size, cm)
	PrintString(out)

	return err
}

func (p *PingFlag) output(host string) {
	out := "\n"
	out += fmt.Sprintf("--- %s ping statistics ---\n", host)

	if p.sta.send == 0 {
		return
	}

	out += fmt.Sprintf("%d packets transmitted, %d packets received, %.2f%% packet loss", p.sta.send, p.sta.receive, float64(p.sta.loss*100/p.sta.send))
	if p.sta.send == p.sta.loss {
		PrintString(out)
		return
	}

	out += "\n"
	avg := p.sta.avg / time.Duration(p.sta.receive)
	var temp float64
	for _, v := range p.sta.stddevRaw {
		temp += math.Pow(float64(v-avg), 2)
	}
	variance := temp / float64(len(p.sta.stddevRaw))
	out += fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v",
		p.sta.min, avg, p.sta.max, time.Duration(math.Sqrt(variance)))
	PrintString(out)
}
