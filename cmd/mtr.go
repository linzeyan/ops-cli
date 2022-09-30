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
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/term"
)

func init() {
	var mtrCmd = &cobra.Command{
		Use:   "mtr",
		Short: "A brief description of your command",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				args = append(args, "1.1.1.1")
			}
			var m MTR
			err := m.getInfo(args[0])
			if err != nil {
				PrintString(err)
				return
			}

			if err := termui.Init(); err != nil {
				PrintString(err)
				return
			}
			defer termui.Close()

			header := widgets.NewParagraph()
			header.Border = false
			header.WrapText = false
			header.Text = "My Traceroute"
			header.TextStyle.Modifier = termui.ModifierBold

			infoL := widgets.NewParagraph()
			infoL.Border = false
			infoL.WrapText = false
			infoL.Text = fmt.Sprintf("%s -> %s", m.LocalHostname, m.RemoteHostname)
			infoR := widgets.NewParagraph()
			infoR.Border = false
			infoR.Text = fmt.Sprintf("%v", time.Now().Local().Format(time.RFC3339))

			keys := widgets.NewParagraph()
			keys.Border = false
			keys.Text = "Keys: (q)uit"

			title1 := widgets.NewParagraph()
			title1.Border = false
			title1.Text = "Packets               Pings"
			title1.TextStyle.Modifier = termui.ModifierBold

			title2L := widgets.NewParagraph()
			title2L.Border = false
			title2L.WrapText = false
			title2L.Text = "Host"
			title2L.TextStyle.Modifier = termui.ModifierBold
			title2R := widgets.NewParagraph()
			title2R.Border = false
			title2R.Text = "Loss%   Snt   Last   Avg  Best  Wrst StDev"
			title2R.TextStyle.Modifier = termui.ModifierBold

			setRect := func() int {
				w, _, _ := term.GetSize(int(os.Stdin.Fd()))
				header.SetRect(w/2-len(header.Text)/2, 1, w/2+len(header.Text), 0)
				infoL.SetRect(0, 3, w/2, 2)
				infoR.SetRect(w-len(time.RFC3339)-2, 3, w, 2)
				keys.SetRect(0, 4, len(keys.Text)+2, 3)
				title1.SetRect(w-len(title1.Text)-15, 5, w, 4)
				title2L.SetRect(0, 6, len(title2L.Text)+2, 5)
				title2R.SetRect(w-len(title2R.Text)-2, 6, w, 5)
				return w
			}
			setRect()

			table := widgets.NewTable()
			table.Border = false
			table.RowSeparator = false
			table.Rows = [][]string{{""}}

			setTable := func(w int) {
				table.SetRect(0, 7, w, 6)
				tbRows := make([][]string, len(m.result))

				for i, v := range m.result {
					spaces := strings.Repeat(" ", w-len(v.Host)-len(title2R.Text))
					r := fmt.Sprintf("%s%s%s  %s  %s  %s  %s  %s  %s",
						v.Host, spaces, v.Loss, v.Snt, v.Last, v.Avg, v.Best, v.Wrst, v.StDev)
					tbRows[i] = []string{r}
				}
				table.Rows = tbRows
				// log.Println(tbRows)
				table.ColumnWidths = []int{w}
			}

			update := func() {
				termui.Clear()
				w := setRect()
				infoR.Text = fmt.Sprintf("%v", time.Now().Local().Format(time.RFC3339))
				m.Run()
				m.Summary()
				setTable(w)
				termui.Render(header, infoL, infoR, keys, title1, title2L, title2R, table)
			}

			tickerCount := 1
			tickerCount++
			uiEvents := termui.PollEvents()
			ticker := time.NewTicker(50 * time.Millisecond).C
			for {
				select {
				case e := <-uiEvents:
					switch e.ID {
					case "q", "<C-c>":
						return
					}
				case <-ticker:
					update()
					tickerCount++
				}
			}
		},
	}
	rootCmd.AddCommand(mtrCmd)
}

type MTR struct {
	IPv6           bool
	LocalHostname  string
	RemoteHostname string
	TargetIP       string

	trace  Traceroute
	result []ICMPStatOutput
}

func (m *MTR) getInfo(host string) error {
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
	addr, err := net.ResolveIPAddr(network, host)
	if err != nil {
		return err
	}
	m.RemoteHostname = fmt.Sprintf("%s (%s)", host, addr.IP)
	m.TargetIP = addr.IP.String()
	ip, _, err := net.SplitHostPort(conn.LocalAddr().String())
	m.LocalHostname = fmt.Sprintf("%s (%s)", hostname, ip)
	return err
}

func (m *MTR) Run() {
	m.trace = Traceroute{
		Size:     24,
		TTL:      64,
		Retry:    1,
		Interval: 100 * time.Millisecond,
		Timeout:  2 * time.Second,
		Record:   true,
	}
	data := Randoms.GenerateString(m.trace.Size, LowercaseLetters)
	m.trace.Data = icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Data: data},
	}

	conn, err := m.trace.Listen()
	if err != nil {
		PrintString(err)
		return
	}
	if conn != nil {
		defer conn.Close()
	}
	m.trace.Connetion = conn
	// m.result = make([]ICMPStatOutput, m.trace.TTL)

	if err = m.trace.Connect(m.TargetIP); err != nil {
		PrintString(err)
		return
	}
}

func (m *MTR) Summary() {
	m.result = make([]ICMPStatOutput, len(m.trace.Stat))
	for i, v := range m.trace.Stat {
		// "Loss%   Snt   Last   Avg  Best  Wrst StDev"
		loss := fmt.Sprintf("%.1f%%", float64(v.Loss*100)/float64(v.Send))
		avg := v.Avg / time.Duration(v.Receive)
		var temp float64
		for _, vv := range v.Rtts {
			temp += math.Pow(float64(vv-avg), 2)
		}
		variance := temp / float64(len(v.Rtts))
		mdev := time.Duration(math.Sqrt(variance))
		m.result[i] = ICMPStatOutput{
			Host:  fmt.Sprintf("%d. %s", v.Hop, v.DstIP),
			Loss:  loss,
			Snt:   strconv.Itoa(v.Send),
			Last:  fmt.Sprintf("%3s", v.Rtts[len(v.Rtts)-1]),
			Avg:   fmt.Sprintf("%3d", avg.Microseconds()),
			Best:  fmt.Sprintf("%3d", v.Min.Microseconds()),
			Wrst:  fmt.Sprintf("%3d", v.Max.Microseconds()),
			StDev: fmt.Sprintf("%3d", mdev.Microseconds()),
		}
	}
}

type ICMPStatOutput struct {
	Host  string
	Loss  string
	Snt   string
	Last  string
	Avg   string
	Best  string
	Wrst  string
	StDev string
}
