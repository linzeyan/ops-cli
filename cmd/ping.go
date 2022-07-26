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
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send ICMP echo packets to host",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		var pingDomain = args[0]
		pinger, err := ping.NewPinger(pingDomain)
		if err != nil {
			log.Println(err)
			return
		}
		pinger.Count = pingCount
		pinger.Interval = time.Second * time.Duration(pingInterval)
		pinger.Size = pingSize
		if !pingTimeout {
			pinger.Timeout = time.Second * 600
		}
		pinger.TTL = pingTTL
		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		// Listen for Ctrl-C.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				pinger.Stop()
			}
		}()

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%v time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt)
		}

		pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%v time=%v (DUP!)\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt)
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}

		fmt.Printf("PING %s (%s): %d data bytes\n", pinger.Addr(), pinger.IPAddr(), pinger.Size)
		err = pinger.Run()
		if err != nil {
			log.Println(err)
			return
		}
	},
	Example: Examples(`# Ping www.google.com
ops-cli ping www.google.com

# Ping 1.1.1.1 and specify packet send interval
ops-cli ping 1.1.1.1 -i 2`),
}

var pingCount, pingInterval, pingSize, pingTTL int
var pingTimeout bool

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVarP(&pingCount, "count", "c", 5, "Specify echo counts")
	pingCmd.Flags().IntVarP(&pingInterval, "interval", "i", 1, "Specify the packet send time interval")
	pingCmd.Flags().IntVarP(&pingSize, "size", "s", 24, "Specify the packet size")
	pingCmd.Flags().BoolVarP(&pingTimeout, "timeout", "t", false, "Do not timeout before exiting")
	pingCmd.Flags().IntVarP(&pingTTL, "ttl", "", 64, "Specify the packet ttl")
}
