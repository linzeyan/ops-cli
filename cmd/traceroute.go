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
			}
			data := Randoms.GenerateString(t.Size, LowercaseLetters)
			t.Data = icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Data: data},
			}

			conn, err := t.Listen()
			if err != nil {
				PrintString(err)
				return
			}
			if conn != nil {
				defer conn.Close()
			}
			t.Connetion = conn
			if t.Record {
				t.stat = make([]ICMPStat, t.TTL)
			}
			if err = t.Connect(args[0]); err != nil {
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
	Size, TTL, Retry  int
	Interval, Timeout time.Duration
	Connetion         *icmp.PacketConn
	Data              icmp.Message

	Record bool
	stat   []ICMPStat
}

func (*Traceroute) Listen() (*icmp.PacketConn, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL|ipv4.FlagDst|ipv4.FlagInterface|ipv4.FlagSrc, true)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func (t *Traceroute) Connect(host string) error {
	var err error
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return err
	}
	reply := make([]byte, 1500)
	for i := 1; i <= t.TTL; i++ {
		if i == 1 && !t.Record {
			header := fmt.Sprintf("traceroute to %s (%v), %d hops max, %d byte packets", host, addr, t.TTL, t.Size)
			PrintString(header)
		}
		t.Data.Body.(*icmp.Echo).Seq = i
		b, err := t.Data.Marshal(nil)
		if err != nil {
			return err
		}

		if err = t.Connetion.IPv4PacketConn().SetTTL(i); err != nil {
			return err
		}
		peer, err := t.sendPacket(i, addr, b, reply)
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

func (t *Traceroute) sendPacket(hop int, addr *net.IPAddr, b, reply []byte) (string, error) {
	var err error
	var ip string
	var rtt []string
	for i := 1; i <= t.Retry; i++ {
		/* Send packet. */
		startTime := time.Now()
		_, err = t.Connetion.IPv4PacketConn().WriteTo(b, nil, addr)
		if err != nil {
			return "", err
		}
		/* Wait receiving. */
		if err = t.Connetion.SetReadDeadline(time.Now().Add(t.Timeout)); err != nil {
			return "", err
		}
		n, _, peer, err := t.Connetion.IPv4PacketConn().ReadFrom(reply)
		if err != nil {
			t.stat[hop].Lost = true
			t.statistics(hop, 0)
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
		t.statistics(hop, duration)
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
	out := fmt.Sprintf("%-4d  %-16v\t%-10s\t%-10s\t%-10s", hop, ip, rtt[0], rtt[1], rtt[2])
	PrintString(out)
	return ip, err
}

func (t *Traceroute) statistics(hop int, duration time.Duration) {
	if !t.Record {
		return
	}
	if t.stat[hop].Min == 0 {
		t.stat[hop].Min = duration
	}
	if t.stat[hop].Max == 0 {
		t.stat[hop].Max = duration
	}
	t.stat[hop].Avg += duration
	t.stat[hop].Send++
	t.stat[hop].Rtts = append(t.stat[hop].Rtts, duration)
	if t.stat[hop].Lost {
		t.stat[hop].Loss++
		t.stat[hop].Lost = false
	} else {
		t.stat[hop].Receive++
	}

	if duration < t.stat[hop].Min {
		t.stat[hop].Min = duration
	}
	if duration > t.stat[hop].Max {
		t.stat[hop].Max = duration
	}
}

// func (t *Traceroute) summary() {
// 	for _, v := range t.stat {
// 		if v.Send == 0 {
// 			continue
// 		}

// 		loss := fmt.Sprintf("%.1f%%", float64(v.Loss*100)/float64(v.Send))
// 		avg := v.Avg / time.Duration(v.Receive)
// 		var temp float64
// 		for _, vv := range v.Rtts {
// 			temp += math.Pow(float64(vv-avg), 2)
// 		}
// 		variance := temp / float64(len(v.Rtts))
// 		mdev := time.Duration(math.Sqrt(variance))
// 	}
// }
