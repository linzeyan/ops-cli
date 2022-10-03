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
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/term"
)

func init() {
	var flags struct {
		count    int
		interval time.Duration
		timeout  time.Duration
	}
	var mtrCmd = &cobra.Command{
		Use:   CommandMtr,
		Short: "Combined traceroute and ping",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if !validator.ValidDomain(args[0]) && !validator.ValidIP(args[0]) {
				return common.ErrInvalidArg
			}

			if flags.interval < 50*time.Millisecond {
				flags.interval = 50 * time.Millisecond
			}

			var m MTR
			m.trace.Host = args[0]
			m.trace.Interval = flags.interval
			m.trace.Timeout = flags.timeout
			m.trace.Count = flags.count
			err := m.init()
			if err != nil {
				return err
			}

			if err := termui.Init(); err != nil {
				return err
			}
			defer termui.Close()

			/* My Traceroute. */
			header := widgets.NewParagraph()
			header.Border = false
			header.WrapText = false
			header.Text = "My Traceroute"
			header.TextStyle.Modifier = termui.ModifierBold

			/* Hostname (Host IP) -> Remote Hostname (Target IP)    time.RFC3339. */
			infoL := widgets.NewParagraph()
			infoL.Border = false
			infoL.WrapText = false
			infoL.Text = fmt.Sprintf("%s -> %s", m.LocalHostname, m.RemoteHostname)
			infoR := widgets.NewParagraph()
			infoR.Border = false
			infoR.Text = fmt.Sprintf("%v", time.Now().Local().Format(time.RFC3339))

			/* Keys: (q)uit. */
			keys := widgets.NewParagraph()
			keys.Border = false
			keys.Text = "Keys: (q)uit"

			/* Packets               Pings */
			title1 := widgets.NewParagraph()
			title1.Border = false
			title1.Text = "Packets               Pings"
			title1.TextStyle.Modifier = termui.ModifierBold

			/* Host        Loss%   Snt   Last   Avg  Best  Wrst StDev */
			title2L := widgets.NewParagraph()
			title2L.Border = false
			title2L.WrapText = false
			title2L.Text = "Host"
			title2L.TextStyle.Modifier = termui.ModifierBold
			title2R := widgets.NewParagraph()
			title2R.Border = false
			title2R.Text = mtrStatHeader
			title2R.TextStyle.Modifier = termui.ModifierBold

			setRect := func() {
				w, _, _ := term.GetSize(int(os.Stdin.Fd()))
				header.SetRect(w/2-len(header.Text)/2, 1, w/2+len(header.Text), 0)
				infoL.SetRect(0, 3, w/2, 2)
				infoR.SetRect(w-len(time.RFC3339)-2, 3, w, 2)
				keys.SetRect(0, 4, len(keys.Text)+2, 3)
				title1.SetRect(w-len(title1.Text)-15, 5, w, 4)
				title2L.SetRect(0, 6, len(title2L.Text)+2, 5)
				title2R.SetRect(w-len(title2R.Text)-2, 6, w, 5)
				m.TerminalWidth = w
			}
			setRect()

			table := widgets.NewTable()
			table.Border = false
			table.RowSeparator = false
			table.Rows = m.Statistics

			uiEvents := termui.PollEvents()
			ticker := time.NewTicker(10 * time.Millisecond).C
			ctx, cancel := signal.NotifyContext(common.Context, os.Interrupt)
			defer func() {
				cancel()
			}()
			go func() {
				for {
					select {
					case e := <-uiEvents:
						switch e.ID {
						case "q", "<C-c>":
							termui.Close()
							os.Exit(0)
						}
					case <-ticker:
						termui.Clear()
						setRect()
						infoR.Text = fmt.Sprintf("%v", time.Now().Local().Format(time.RFC3339))
						table.SetRect(0, 6, m.TerminalWidth, 36)
						table.Rows = m.Statistics
						table.ColumnWidths = []int{m.TerminalWidth}
						termui.Render(header, infoL, infoR, keys, title1, title2L, title2R, table)
					}
				}
			}()

			return m.Run(ctx)
		},
	}
	rootCmd.AddCommand(mtrCmd)
	mtrCmd.Flags().IntVarP(&flags.count, "count", "c", -1, common.Usage("Specify ping counts"))
	mtrCmd.Flags().DurationVarP(&flags.interval, "interval", "i", 100*time.Millisecond, common.Usage("Specify interval"))
	mtrCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 2*time.Second, common.Usage("Specify timeout"))
}

type MTR struct {
	IPv6           bool
	LocalHostname  string
	RemoteHostname string
	TerminalWidth  int
	Statistics     [][]string

	trace Traceroute
}

func (m *MTR) init() error {
	var err error
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}
	network := "ip4"
	if m.IPv6 {
		network = "ip6"
	}
	m.trace.Target, err = net.ResolveIPAddr(network, m.trace.Host)
	if err != nil {
		return err
	}
	m.RemoteHostname = fmt.Sprintf("%s (%s)", m.trace.Host, m.trace.Target)
	ip, _, err := net.SplitHostPort(conn.LocalAddr().String())
	m.LocalHostname = fmt.Sprintf("%s (%s)", hostname, ip)
	m.Statistics = [][]string{{""}}
	m.trace.Size = 24
	m.trace.TTL = 64
	m.trace.Retry = 1
	m.trace.Record = true
	return err
}

func (m *MTR) Run(ctx context.Context) error {
	data := Randoms.GenerateString(m.trace.Size, LowercaseLetters)
	m.trace.Data = icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Data: data},
	}

	conn, err := m.trace.Listen()
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}
	m.trace.Connetion = conn

	reply := make([]byte, 1500)
	for round := 0; ; round++ {
		if err = m.trace.Connect(ctx, reply); err != nil {
			return err
		}
		if round == m.trace.Count {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m.Summary()
		}
	}
}

func (m *MTR) Summary() {
	var rows [][]string
	headers := strings.Fields(mtrStatHeader)
	for _, v := range m.trace.Stat {
		var avg time.Duration
		if v.Receive == 0 {
			avg = 0
		} else {
			avg = v.Avg / time.Duration(v.Receive)
		}
		var temp float64
		for _, vv := range v.Rtts {
			temp += math.Pow(float64(vv-avg), 2)
		}
		variance := temp / float64(len(v.Rtts))
		mdev := time.Duration(math.Sqrt(variance))

		host := fmt.Sprintf("%d. %s", v.Hop, v.DstIP)
		stats := fmt.Sprintf("%4s%%   %3s   %4s   %3s  %4s  %4s %5s",
			m.trim(fmt.Sprintf("%.1f", float64(v.Loss*100)/float64(v.Send)), strings.Replace(headers[0], "%", "", 1)),
			strconv.Itoa(v.Send),
			m.trim(v.Rtts[len(v.Rtts)-1].String(), headers[2]),
			m.trim(avg.String(), headers[3]), m.trim(v.Min.String(), headers[4]),
			m.trim(v.Max.String(), headers[5]), m.trim(mdev.String(), headers[6]))
		spaces := strings.Repeat(" ", m.TerminalWidth-19-len(mtrStatHeader)-2)

		rows = append(rows, []string{fmt.Sprintf("%-19s", host) + spaces + stats})
	}
	m.Statistics = rows
}

func (m *MTR) trim(s, header string) string {
	i := strings.Index(s, ".")
	s = strings.Replace(s, "ms", "", 1)
	if strings.Contains(s, "µ") {
		s = strings.Replace(s, "µs", "", 1)
		num, err := strconv.ParseFloat(s, 64)
		if err == nil {
			ms := num / 1000
			s = fmt.Sprintf("%f", ms)
		}
	}

	if len(s[:i+2]) > len(header) {
		return s[0:i]
	}
	return s[:i+2]
}
