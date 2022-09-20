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
}

func (p *PingFlag) Run(cmd *cobra.Command, args []string) {
	host := args[0]
	var data RandomString
	data = data.GenerateString(p.size, LowercaseLetters)
	header := fmt.Sprintf("PING %s (%s): %d data bytes", host, host, len(data))
	PrintString(header)

	for i := 0; i != p.count; i++ {
		if err := p.Connect(i, host, data); err != nil {
			PrintString(err)
		}
		time.Sleep(p.interval)
	}
}

func (p *PingFlag) Connect(counter int, host string, icmpData []byte) error {
	if p.ipv6 {
		return p.ping6(counter, host, icmpData)
	}
	return p.ping4(counter, host, icmpData)
}

func (p *PingFlag) ping4(counter int, host string, icmpData []byte) error {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}

	if err = conn.IPv4PacketConn().SetTTL(p.ttl); err != nil {
		return err
	}
	if err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true); err != nil {
		return err
	}

	ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return err
	}

	data := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   counter,
			Seq:  counter,
			Data: icmpData,
		},
	}
	b, err := data.Marshal(nil)
	if err != nil {
		return err
	}

	/* Send packet. */
	startTime := time.Now()
	_, err = conn.WriteTo(b, ip)
	if err != nil {
		return err
	}

	/* Wait receiving. */
	reply := make([]byte, 1500)
	if err = conn.SetReadDeadline(time.Now().Add(p.timeout)); err != nil {
		return err
	}
	n, cm, peer, err := conn.IPv4PacketConn().ReadFrom(reply)
	if err != nil {
		e := fmt.Sprintf("Request timeout for icmp_seq %d", counter)
		return errors.New(e)
	}

	result, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return err
	}
	var out string
	switch result.Type {
	case ipv4.ICMPTypeEchoReply:
		out = fmt.Sprintf("%v bytes from %v: icmp_seq=%d ttl=%d time=%v", len(b), peer, counter, cm.TTL, time.Since(startTime))
	case ipv4.ICMPTypeDestinationUnreachable:
		out = fmt.Sprintf("%v Destination Unreachable", peer)
	}
	PrintString(out)

	return err
}

func (p *PingFlag) ping6(counter int, host string, icmpData []byte) error {
	conn, err := icmp.ListenPacket("ip6:ipv6-icmp", "::")
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}

	if err = conn.IPv6PacketConn().SetHopLimit(p.ttl); err != nil {
		return err
	}
	if err = conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true); err != nil {
		return err
	}
	if err = conn.IPv6PacketConn().SetControlMessage(ipv6.FlagInterface, true); err != nil {
		return err
	}

	ip, err := net.ResolveIPAddr("ip6", host)
	if err != nil {
		return err
	}

	data := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   counter,
			Seq:  counter,
			Data: icmpData,
		},
	}
	b, err := data.Marshal(nil)
	if err != nil {
		return err
	}

	/* Send packet. */
	startTime := time.Now()
	_, err = conn.WriteTo(b, ip)
	if err != nil {
		return err
	}

	/* Wait receiving. */
	reply := make([]byte, 1500)
	if err = conn.SetReadDeadline(time.Now().Add(p.timeout)); err != nil {
		return err
	}
	n, cm, peer, err := conn.IPv6PacketConn().ReadFrom(reply)
	if err != nil {
		e := fmt.Sprintf("Request timeout for icmp_seq %d", counter)
		return errors.New(e)
	}

	result, err := icmp.ParseMessage(58, reply[:n])
	if err != nil {
		return err
	}
	var out string
	switch result.Type {
	case ipv6.ICMPTypeEchoReply:
		out = fmt.Sprintf("%v bytes from %v: icmp_seq=%d ttl=%d time=%v", len(b), peer, counter, cm.HopLimit, time.Since(startTime))
	case ipv6.ICMPTypeDestinationUnreachable:
		out = fmt.Sprintf("%v Destination Unreachable", peer)
	}
	PrintString(out)

	return err
}
