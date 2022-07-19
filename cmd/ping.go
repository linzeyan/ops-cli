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
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Send ICMP echo packets to host",
	Args:  cobra.OnlyValidArgs,
	Run: func(_ *cobra.Command, _ []string) {
		pingExec()
	},
	Example: Examples(`# Ping www.google.com
ops-cli ping -d www.google.com

# Ping www.google.com and specify packet send interval
ops-cli ping -d www.google.com -t 2`),
}

var pingDomain string
var pingCount, pingInterval, pingSize int

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringVarP(&pingDomain, "domain", "d", "", "Specify host")
	pingCmd.Flags().IntVarP(&pingCount, "count", "c", 5, "Specify counts")
	pingCmd.Flags().IntVarP(&pingInterval, "interval", "i", 1, "Specify the packet send time interval")
	pingCmd.Flags().IntVarP(&pingSize, "size", "s", 24, "Specify the packet size")
}

func pingExec() {
	pinger, err := ping.NewPinger(pingDomain)
	if err != nil {
		log.Println(err)
		return
	}
	pinger.Count = pingCount
	pinger.Interval = time.Second * time.Duration(pingInterval)
	pinger.Size = pingSize

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
