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
	"bytes"
	"context"
	"errors"
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
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/term"
)

func initMTR() *cobra.Command {
	var flags struct {
		output   string
		count    int
		interval time.Duration
		timeout  time.Duration
	}
	var mtrCmd = &cobra.Command{
		GroupID: groupings[CommandMTR],
		Use:     CommandMTR,
		Short:   "Combined traceroute and ping",
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, args []string) {
			if !common.IsDomain(args[0]) && !common.IsIP(args[0]) {
				logger.Info(common.ErrInvalidArg.Error(), common.DefaultField(args))
				printer.Error(common.ErrInvalidArg)
				return
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
				logger.Info(err.Error())
				printer.Error(err)
				return
			}

			if err := termui.Init(); err != nil {
				logger.Info(err.Error())
				printer.Error(err)
				return
			}
			defer termui.Close()

			/* My Traceroute. */
			header := widgets.NewParagraph()
			header.Border = false
			header.WrapText = false
			headerString := "My Traceroute"
			header.TextStyle.Modifier = termui.ModifierBold

			/* Hostname (Host IP) -> Remote Hostname (Target IP)    time.RFC3339. */
			info := widgets.NewParagraph()
			info.Border = false
			info.WrapText = false
			infoString := fmt.Sprintf("%s -> %s", m.LocalHostname, m.RemoteHostname)

			/* Keys: (q)uit. */
			keys := widgets.NewParagraph()
			keys.Border = false
			keys.Text = "Keys: (q)uit"

			/* Packets               Pings */
			title1 := widgets.NewParagraph()
			title1.Border = false
			t1String := "Packets               Pings"
			title1.TextStyle.Modifier = termui.ModifierBold

			/* Host        Loss%   Snt   Last   Avg  Best  Wrst StDev */
			title2 := widgets.NewParagraph()
			title2.Border = false
			title2.WrapText = false
			titl2String := "Host"
			title2.TextStyle.Modifier = termui.ModifierBold

			setRect := func() {
				w, _, _ := term.GetSize(int(os.Stdin.Fd()))
				n := (w-len(headerString))/2 - 2
				spaces := strings.Repeat(" ", n)
				header.Text = spaces + headerString + spaces
				header.SetRect(0, 1, w, 0)

				n = w - len(infoString) - len(time.RFC3339) - 2
				spaces = strings.Repeat(" ", n)
				info.Text = infoString + spaces + fmt.Sprintf("%v", time.Now().Local().Format(time.RFC3339))
				info.SetRect(0, 2, w, 1)

				keys.SetRect(0, 3, w, 2)

				n = w - len(mtrStatHeader)
				spaces = strings.Repeat(" ", n)
				title1.Text = spaces + t1String
				title1.SetRect(0, 4, w, 3)

				n = w - len(titl2String) - len(mtrStatHeader) - 2
				spaces = strings.Repeat(" ", n)
				title2.Text = titl2String + spaces + mtrStatHeader
				title2.SetRect(0, 5, w, 4)
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
						table.SetRect(0, 5, m.TerminalWidth, 36)
						table.Rows = m.Statistics
						table.ColumnWidths = []int{m.TerminalWidth}
						termui.Render(header, info, keys, title1, title2, table)
					}
				}
			}()
			err = m.Run(ctx)
			if errors.Is(err, common.ErrResponse) {
				if flags.count != -1 && flags.output != "" {
					var buf bytes.Buffer
					buf.WriteString(header.Text + "\n" + info.Text + "\n" + keys.Text + "\n" + title1.Text + "\n" + title2.Text + "\n")
					for _, i := range m.Statistics {
						buf.WriteString(i[0] + "\n")
					}
					wrErr := os.WriteFile(flags.output, buf.Bytes(), FileModeRAll)
					if wrErr != nil {
						logger.Info(wrErr.Error())
						printer.Error(wrErr)
					}
				}
			} else if err != nil {
				logger.Info(err.Error())
				printer.Error(err)
			}
		},
	}
	mtrCmd.Flags().StringVarP(&flags.output, "output", "o", "", common.Usage("Specify output file name"))
	mtrCmd.Flags().IntVarP(&flags.count, "count", "c", -1, common.Usage("Specify ping counts"))
	mtrCmd.Flags().DurationVarP(&flags.interval, "interval", "i", 100*time.Millisecond, common.Usage("Specify interval"))
	mtrCmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 800*time.Millisecond, common.Usage("Specify timeout"))
	return mtrCmd
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
		logger.Debug(err.Error())
		return err
	}
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		logger.Debug(err.Error())
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
		logger.Debug(err.Error())
		return err
	}
	m.RemoteHostname = fmt.Sprintf("%s (%s)", m.trace.Host, m.trace.Target)
	ip, _, err := net.SplitHostPort(conn.LocalAddr().String())
	m.LocalHostname = fmt.Sprintf("%s (%s)", hostname, ip)
	m.Statistics = [][]string{{""}}
	m.trace.Size = 60
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
		logger.Debug(err.Error())
		return err
	}
	if conn != nil {
		defer conn.Close()
	}
	m.trace.Connetion = conn

	reply := make([]byte, 1500)

	go func() {
		for {
			if len(m.trace.Stat) != 0 {
				m.Summary()
			}
		}
	}()
	for round := 0; ; round++ {
		if err = m.trace.Connect(ctx, reply); err != nil {
			return err
		}
		if round == m.trace.Count {
			return common.ErrResponse
		}
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			continue
		}
	}
}

func (m *MTR) Summary() {
	var rows [][]string
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
		if len(v.Rtts)-1 < 0 {
			time.Sleep(time.Second)
		}
		stats := fmt.Sprintf("%5s%% %5s  %5s %5s %5s %5s %5s",
			fmt.Sprintf("%.1f", float64(v.Loss*100)/float64(v.Send)),
			strconv.Itoa(v.Send),
			m.trim(v.Rtts[len(v.Rtts)-1].String()),
			m.trim(avg.String()), m.trim(v.Min.String()),
			m.trim(v.Max.String()), m.trim(mdev.String()))
		spaces := strings.Repeat(" ", m.TerminalWidth-19-len(mtrStatHeader)-3)

		rows = append(rows, []string{fmt.Sprintf("%-19s", host) + spaces + stats})
	}
	m.Statistics = rows
}

func (m *MTR) trim(s string) string {
	i := strings.Index(s, ".")
	s = strings.Replace(s, "ms", "", 1)
	if strings.Contains(s, "µ") {
		s = strings.Replace(s, "µs", "", 1)
		num, err := strconv.ParseFloat(s, 64)
		if err == nil {
			ms := num / 1000
			s = fmt.Sprintf("%f", ms)
		} else {
			logger.Debug(err.Error())
		}
	}

	if len(s[:i+2]) > 5 {
		return s[0:i]
	}
	return s[:i+2]
}
