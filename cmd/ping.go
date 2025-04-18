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
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func initPing() *cobra.Command {
	var flags struct {
		ipv6             bool
		count, size, ttl int
		interval         time.Duration
		timeout          time.Duration
	}
	var pingCmd = &cobra.Command{
		GroupID: getGroupID(CommandPing),
		Use:     CommandPing + " [host]",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Send ICMP ECHO_REQUEST packets to network hosts",
		Run: func(_ *cobra.Command, args []string) {
			if flags.count == 0 || flags.ttl <= 0 || flags.size <= 0 {
				logger.Error(common.ErrInvalidArg.Error())
				return
			}
			if flags.interval < 50*time.Millisecond {
				flags.interval = 50 * time.Millisecond
			}

			var p = Ping{
				IPv6:     flags.ipv6,
				Count:    flags.count,
				Size:     flags.size,
				TTL:      flags.ttl,
				Interval: flags.interval,
				Timeout:  flags.timeout,
				Data: icmp.Message{
					Type: ipv4.ICMPTypeEcho,
					Code: 0,
					Body: &icmp.Echo{Data: Randoms.GenerateString(flags.size, LowercaseLetters)},
				},
			}

			conn, err := p.Listen()
			if err != nil {
				logger.Warn(err.Error(), common.DefaultField(p))
				printer.Error(err)
				return
			}
			if conn != nil {
				defer conn.Close()
			}
			p.Conn = conn

			ctx, cancel := signal.NotifyContext(common.Context, os.Interrupt)
			defer func() {
				cancel()
			}()
			go func() {
				select {
				case <-ctx.Done():
				default:
				}
			}()

			p.Connect(ctx, args[0])
		},
	}
	pingCmd.Flags().IntVarP(&flags.count, "count", "c", -1, common.Usage("Specify ping counts"))
	pingCmd.Flags().BoolVarP(&flags.ipv6, "ipv6", "6", false, common.Usage("Use ICMPv6"))
	pingCmd.Flags().IntVarP(&flags.size, "size", "s", 56, common.Usage("Specify packet size"))
	pingCmd.Flags().IntVarP(&flags.ttl, "ttl", "", 64, common.Usage("Specify packet ttl"))
	pingCmd.Flags().DurationVarP(&flags.interval, "interval", "i", time.Second, common.Usage("Specify interval"))
	pingCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
	return pingCmd
}

type Ping struct {
	IPv6              bool
	Count, Size, TTL  int
	Interval, Timeout time.Duration
	Data              icmp.Message
	Conn              *icmp.PacketConn

	stat ICMPStat
}

type ICMPStat struct {
	Lost bool

	Hop   int
	DstIP string

	Send, Loss, Receive int
	Min, Avg, Max       time.Duration
	Rtts                []time.Duration
}

func (p *Ping) Listen() (*icmp.PacketConn, error) {
	if p.IPv6 {
		conn, err := icmp.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil {
			logger.Debug(err.Error())
			return nil, err
		}
		if err = conn.IPv6PacketConn().SetHopLimit(p.TTL); err != nil {
			logger.Debug(err.Error(), common.DefaultField(p.TTL))
			return nil, err
		}
		if err = conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true); err != nil {
			logger.Debug(err.Error(), common.DefaultField(ipv6.FlagHopLimit))
			return nil, err
		}
		return conn, err
	}
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	if err = conn.IPv4PacketConn().SetTTL(p.TTL); err != nil {
		logger.Debug(err.Error(), common.DefaultField(p.TTL))
		return nil, err
	}
	if err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true); err != nil {
		logger.Debug(err.Error(), common.DefaultField(ipv4.FlagTTL))
		return nil, err
	}
	return conn, err
}

func (p *Ping) statistics(duration time.Duration) {
	if p.stat.Min == 0 {
		p.stat.Min = duration
	}
	if p.stat.Max == 0 {
		p.stat.Max = duration
	}
	p.stat.Avg += duration
	p.stat.Send++
	p.stat.Rtts = append(p.stat.Rtts, duration)
	if p.stat.Lost {
		p.stat.Loss++
		p.stat.Lost = false
	} else {
		p.stat.Receive++
	}

	if duration < p.stat.Min {
		p.stat.Min = duration
	}
	if duration > p.stat.Max {
		p.stat.Max = duration
	}
}

func (p *Ping) readReply(reply []byte, counter int) (int, any, net.Addr, error) {
	var n int
	var cm any
	var peer net.Addr
	var err error
	if p.IPv6 {
		n, cm, peer, err = p.Conn.IPv6PacketConn().ReadFrom(reply)
	} else {
		n, cm, peer, err = p.Conn.IPv4PacketConn().ReadFrom(reply)
	}
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(reply))
		p.stat.Lost = true
		p.statistics(0)
		e := fmt.Sprintf("Request timeout for icmp_seq %d", counter)
		return 0, nil, nil, errors.New(e)
	}
	return n, cm, peer, err
}

func (p *Ping) printMsg(result *icmp.Message, duration time.Duration, peer net.Addr, cm any) {
	var ttl int
	switch c := cm.(type) {
	case *ipv4.ControlMessage:
		ttl = c.TTL
	case *ipv6.ControlMessage:
		ttl = c.HopLimit
	}
	b, _ := result.Marshal(nil)
	var out string
	switch result.Type {
	case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
		out = fmt.Sprintf("%v bytes from %v: icmp_seq=%d ttl=%d time=%v\n",
			len(b), peer, result.Body.(*icmp.Echo).Seq, ttl, duration)
	case ipv4.ICMPTypeDestinationUnreachable, ipv6.ICMPTypeDestinationUnreachable:
		out = "Destination Unreachable\n"
	default:
		logger.Debug("icmp.Type", common.DefaultField(result.Type))
	}
	p.statistics(duration)
	printer.Printf(out)
}

func (p *Ping) Connect(c context.Context, host string) {
	network := "ip4"
	if p.IPv6 {
		network = "ip6"
	}
	if network == "ip4" {
		addr, err := ParseAnyIPv4Netip(host)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		host = addr.String()
	}
	addr, err := net.ResolveIPAddr(network, host)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if p.IPv6 {
		p.Data.Type = ipv6.ICMPTypeEchoRequest
	}

	reply := make([]byte, 1500)
	allTime := time.Now()

	for i := 0; ; i++ {
		if i == 0 {
			printer.Printf("PING %s (%v): %d data bytes\n", host, addr, p.Size)
		}
		p.Data.Body.(*icmp.Echo).ID = i & 0xffff
		p.Data.Body.(*icmp.Echo).Seq = i & 0xffff

		b, err := p.Data.Marshal(nil)
		if err != nil {
			logger.Debug(err.Error())
			return
		}

		/* Send packet. */
		startTime := time.Now()
		_, err = p.Conn.WriteTo(b, addr)
		if err != nil {
			logger.Debug(err.Error())
		}

		/* Wait receiving. */
		if err = p.Conn.SetReadDeadline(time.Now().Add(p.Timeout)); err != nil {
			logger.Debug(err.Error())
		}
		n, cm, peer, err := p.readReply(reply, i)
		if err != nil {
			logger.Debug(err.Error())
			if i == p.Count-1 {
				p.summary(host, time.Since(allTime))
				return
			}
			select {
			default:
				continue
			case <-c.Done():
				p.summary(host, time.Since(allTime))
				return
			}
		}
		duration := time.Since(startTime)

		proto := 1
		if p.IPv6 {
			proto = 58
		}
		result, err := icmp.ParseMessage(proto, reply[:n])
		if err != nil {
			logger.Debug(err.Error())
		}
		if peer.String() == host {
			p.printMsg(result, duration, peer, cm)
		}
		if i == p.Count-1 {
			p.summary(host, time.Since(allTime))
			return
		}
		time.Sleep(p.Interval)
		select {
		default:
		case <-c.Done():
			p.summary(host, time.Since(allTime))
			return
		}
	}
}

func (p *Ping) summary(host string, t time.Duration) {
	if p.stat.Send == 0 {
		return
	}

	out := "\n"
	out += fmt.Sprintf("--- %s ping statistics ---\n", host)
	out += fmt.Sprintf("%d packets transmitted, %d received, %.1f%% packet loss, time %vms\n",
		p.stat.Send, p.stat.Receive, float64(p.stat.Loss*100)/float64(p.stat.Send), t.Milliseconds())
	if p.stat.Send == p.stat.Loss {
		printer.Printf(out)
		return
	}

	avg := p.stat.Avg / time.Duration(p.stat.Receive)
	var temp float64
	for _, v := range p.stat.Rtts {
		temp += math.Pow(float64(v-avg), 2)
	}
	variance := temp / float64(len(p.stat.Rtts))
	out += fmt.Sprintf("round-trip min/avg/max/mdev = %v/%v/%v/%v\n",
		p.stat.Min, avg, p.stat.Max, time.Duration(math.Sqrt(variance)))
	printer.Printf(out)
}

func ParseAnyIPv4Netip(input string) (netip.Addr, error) {
	parts := strings.Split(input, ".")
	switch len(parts) {
	case 1:
		val, err := strconv.ParseUint(parts[0], 10, 32)
		if err != nil {
			return netip.Addr{}, err
		}
		ip := netip.AddrFrom4([4]byte{
			byte(val >> 24),
			byte(val >> 16),
			byte(val >> 8),
			byte(val),
		})

		return ip, nil
	case 2:
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return netip.Addr{}, err
		}
		b, err := strconv.Atoi(parts[1])
		if err != nil {
			return netip.Addr{}, err
		}
		ip := netip.AddrFrom4([4]byte{
			byte(a),
			byte(b >> 16),
			byte((b >> 8) & 0xFF),
			byte(b & 0xFF),
		})
		return ip, nil
	case 3:
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return netip.Addr{}, err
		}
		b, err := strconv.Atoi(parts[1])
		if err != nil {
			return netip.Addr{}, err
		}
		c, err := strconv.Atoi(parts[2])
		if err != nil {
			return netip.Addr{}, err
		}
		ip := netip.AddrFrom4([4]byte{
			byte(a),
			byte(b),
			byte(c >> 8),
			byte(c & 0xFF),
		})
		return ip, nil
	case 4:
		ip, err := netip.ParseAddr(input)
		if err != nil || !ip.Is4() {
			return netip.Addr{}, errors.New("invalid IPv4 address")
		}
		return ip, nil
	default:
		return netip.Addr{}, errors.New("invalid format")
	}
}

func IPv4NetipToAllForms(ip netip.Addr) (string, string, string, string) {
	if !ip.Is4() {
		panic("not an IPv4 address")
	}
	b := ip.As4()

	f3 := uint32(b[2])<<8 | uint32(b[3])
	f2 := uint32(b[1])<<16 | f3
	f1 := uint32(b[0])<<24 | f2
	form2 := fmt.Sprintf("%d.%d", b[0], f2)
	form3 := fmt.Sprintf("%d.%d.%d", b[0], b[1], f3)
	form4 := ip.String()

	return strconv.FormatUint(uint64(f1), 10), form2, form3, form4
}
