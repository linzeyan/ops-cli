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
	"net"
	"os"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func initTraceroute() *cobra.Command {
	var flags struct {
		size, maxTTL int
		interval     time.Duration
		timeout      time.Duration
	}
	var tracerouteCmd = &cobra.Command{
		GroupID: getGroupID(CommandTraceroute),
		Use:     CommandTraceroute + " [host]",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Print the route packets trace to network host",
		Run: func(_ *cobra.Command, args []string) {
			if flags.size <= 0 || flags.maxTTL <= 0 {
				logger.Info(common.ErrInvalidArg.Error())
				return
			}
			if flags.maxTTL > 64 {
				flags.maxTTL = 64
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
				Count:    1,
			}
			data := Randoms.GenerateString(t.Size, LowercaseLetters)
			t.Data = icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Data: data},
			}

			conn, err := t.Listen()
			if err != nil {
				logger.Info(err.Error())
				printer.Error(err)
				return
			}
			if conn != nil {
				defer conn.Close()
			}
			t.Connetion = conn
			t.Host = args[0]
			t.Target, err = net.ResolveIPAddr("ip4", t.Host)
			if err != nil {
				logger.Info(err.Error())
				printer.Error(err)
				return
			}
			reply := make([]byte, 1500)
			if err = t.Connect(common.Context, reply); err != nil {
				logger.Info(err.Error())
				printer.Error(err)
				return
			}
		},
	}
	tracerouteCmd.Flags().IntVarP(&flags.size, "size", "s", 60, common.Usage("Specify packet size"))
	tracerouteCmd.Flags().IntVarP(&flags.maxTTL, "max-ttl", "m", 64, common.Usage("Specify max hop"))
	tracerouteCmd.Flags().DurationVarP(&flags.interval, "interval", "i", 500*time.Millisecond, common.Usage("Specify interval"))
	tracerouteCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
	return tracerouteCmd
}

type Traceroute struct {
	Size, TTL, Retry  int
	Interval, Timeout time.Duration
	Connetion         *icmp.PacketConn
	Data              icmp.Message

	Host   string
	Target *net.IPAddr

	Count int

	lost   bool
	Record bool
	Stat   []ICMPStat
}

func (*Traceroute) Listen() (*icmp.PacketConn, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL|ipv4.FlagDst|ipv4.FlagInterface|ipv4.FlagSrc, true)
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return conn, err
}

func (t *Traceroute) Connect(ctx context.Context, reply []byte) error {
	var err error
	for i := 1; i <= t.TTL; i++ {
		if i == 1 && !t.Record {
			printer.Printf("traceroute to %s (%v), %d hops max, %d byte packets\n", t.Host, t.Target, t.TTL, t.Size)
		}
		t.Data.Body.(*icmp.Echo).Seq = i
		b, err := t.Data.Marshal(nil)
		if err != nil {
			logger.Debug(err.Error())
			return err
		}

		if err = t.Connetion.IPv4PacketConn().SetTTL(i); err != nil {
			logger.Debug(err.Error())
			return err
		}
		peer, err := t.sendPacket(i, b, reply)
		if err != nil {
			logger.Debug(err.Error())
			return err
		}
		if peer == t.Target.String() {
			logger.Debug("peer == target")
			break
		}
		time.Sleep(t.Interval)
	}
	return err
}

func (t *Traceroute) sendPacket(hop int, b, reply []byte) (string, error) {
	var err error
	var ip string
	var rtt []string
	for i := 1; i <= t.Retry; i++ {
		/* Send packet. */
		startTime := time.Now()
		_, err = t.Connetion.IPv4PacketConn().WriteTo(b, nil, t.Target)
		if err != nil {
			logger.Debug(err.Error())
			return "", err
		}
		/* Wait receiving. */
		if err = t.Connetion.SetReadDeadline(time.Now().Add(t.Timeout)); err != nil {
			logger.Debug(err.Error())
			return "", err
		}
		n, _, peer, err := t.Connetion.IPv4PacketConn().ReadFrom(reply)
		if err != nil {
			logger.Debug(err.Error())
			t.lost = true
			t.statistics(hop, "*", 0)
			rtt = append(rtt, "*")
			continue
		}
		if peer.String() == inetLocalhostIP && (t.Host == inetLocalhost || t.Host == inetLocalhostIP) {
			logger.Debug("peer == localhost")
			continue
		}
		duration := time.Since(startTime)
		result, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			logger.Debug(err.Error())
			return peer.String(), err
		}
		switch result.Type {
		case ipv4.ICMPTypeEchoReply, ipv4.ICMPTypeTimeExceeded:
			rtt = append(rtt, duration.String())
		case ipv4.ICMPTypeDestinationUnreachable:
			rtt = append(rtt, "*")
		}
		t.statistics(hop, peer.String(), duration)
		if peer.String() != "" {
			ip = peer.String()
		}
		if i == t.Retry {
			break
		}
		time.Sleep(t.Interval)
	}
	if t.Record {
		return ip, err
	}
	printer.Printf("%2d. %-16v\t%-10s\t%-10s\t%-10s\n", hop, ip, rtt[0], rtt[1], rtt[2])
	return ip, err
}

func (t *Traceroute) statistics(hop int, ip string, duration time.Duration) {
	if !t.Record {
		return
	}
	if len(t.Stat) < hop {
		t.Stat = append(t.Stat, ICMPStat{
			Hop: hop,
		})
	}
	if t.Stat[hop-1].Min == 0 {
		t.Stat[hop-1].Min = duration
	}
	if t.Stat[hop-1].Max == 0 {
		t.Stat[hop-1].Max = duration
	}
	t.Stat[hop-1].Avg += duration
	t.Stat[hop-1].DstIP = ip
	t.Stat[hop-1].Send++
	t.Stat[hop-1].Rtts = append(t.Stat[hop-1].Rtts, duration)
	if t.lost {
		t.Stat[hop-1].Loss++
		t.lost = false
	} else {
		t.Stat[hop-1].Receive++
	}

	if duration < t.Stat[hop-1].Min && duration != 0 {
		t.Stat[hop-1].Min = duration
	}
	if duration > t.Stat[hop-1].Max {
		t.Stat[hop-1].Max = duration
	}
}
