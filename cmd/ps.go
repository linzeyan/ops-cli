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
	"strings"
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

func initPs() *cobra.Command {
	var flags struct {
		command bool
		parent  bool
		user    bool
	}
	var psCmd = &cobra.Command{
		Use:   CommandPs,
		Short: "Display process status",
		RunE: func(_ *cobra.Command, _ []string) error {
			p, err := process.ProcessesWithContext(common.Context)
			if err != nil {
				return err
			}

			/* Generate header. */
			header := []string{"pid"}
			if flags.parent {
				header = append(header, "ppid")
			}
			if flags.user {
				header = append(header, "user", "cpu%", "mem%", "rss", "status", "started")
			}
			header = append(header, "command")

			var data [][]string
			var ps Ps
			for _, v := range p {
				// fmt.Println(v.Connections())
				ps.Process = v
				/* Append values same as header. */
				out := []string{ps.Pid()}
				if flags.parent {
					out = append(out, ps.Ppid())
				}
				if flags.user {
					out = append(out, ps.Username(), ps.CPUPercent(), ps.MemPercent(), ps.RSS(), ps.Status(), ps.CreateTime())
				}
				out = append(out, ps.Exe(flags.command))

				data = append(data, out)
			}
			ps.String(header, data)
			return nil
		},
	}
	psCmd.Flags().BoolVarP(&flags.command, "command", "c", false, "Print full command")
	psCmd.Flags().BoolVarP(&flags.parent, "parent", "p", false, "Print parent ID")
	psCmd.Flags().BoolVarP(&flags.user, "user", "u", false, "Print username")
	return psCmd
}

type Ps struct {
	Process *process.Process
}

func (p *Ps) CPUGuestTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Guest)
}

func (p *Ps) CPUGuestNiceTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.GuestNice)
}

func (p *Ps) CPUIdleTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Idle)
}

func (p *Ps) CPUIowaitTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Iowait)
}

func (p *Ps) CPUIrqTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Irq)
}

func (p *Ps) CPUNiceTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Nice)
}

func (p *Ps) CPUPercent() string {
	cpuPercent, err := p.Process.CPUPercentWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f%%", cpuPercent)
}

func (p *Ps) CPUSoftirqTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Softirq)
}

func (p *Ps) CPUStealTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.Steal)
}

func (p *Ps) CPUSystemTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.System)
}

func (p *Ps) CPUUserTime() string {
	cpu, err := p.Process.TimesWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f", cpu.User)
}

func (p *Ps) CreateTime() string {
	createTime, err := p.Process.CreateTimeWithContext(common.Context)
	if err != nil {
		return ""
	}
	procTime := time.Unix(createTime/1000, 0)
	if common.TimeNow.Sub(procTime) > 24*time.Hour {
		return procTime.Format("2006-01-02 15:04:05")
	}
	return procTime.Format("15:04:05")
}

func (p *Ps) Data() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.Data)
}

func (p *Ps) Exe(b bool) string {
	const width = 30
	exe, err := p.Process.ExeWithContext(common.Context)
	if err != nil {
		return ""
	}
	if b {
		return exe
	}
	s := strings.Fields(exe)[0]
	if len(s) < width {
		return s
	}
	return s[:width]
}

func (p *Ps) HWM() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.HWM)
}

func (p *Ps) Locked() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.Locked)
}

func (p *Ps) MemPercent() string {
	memoryPercent, err := p.Process.MemoryPercentWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f%%", memoryPercent)
}

func (p *Ps) Nice() string {
	nice, err := p.Process.NiceWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", nice)
}

func (p *Ps) Pid() string {
	return fmt.Sprintf("%d", p.Process.Pid)
}

func (p *Ps) Ppid() string {
	ppid, err := p.Process.PpidWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", ppid)
}

func (p *Ps) ProcessName() string {
	name, err := p.Process.NameWithContext(common.Context)
	if err != nil {
		return ""
	}
	return name
}

func (p *Ps) RSS() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.RSS)
}

func (p *Ps) Stack() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.Stack)
}

func (p *Ps) Status() string {
	status, err := p.Process.StatusWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", status)
}

func (p *Ps) Swap() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.Swap)
}

func (p *Ps) Thread() string {
	threads, err := p.Process.NumThreadsWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", threads)
}

func (p *Ps) Username() string {
	username, err := p.Process.UsernameWithContext(common.Context)
	if err != nil {
		return ""
	}
	return username
}

func (p *Ps) VMS() string {
	mem, err := p.Process.MemoryInfoWithContext(common.Context)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", mem.VMS)
}

func (*Ps) String(header []string, data [][]string) {
	/* tablewriter.ALIGN_LEFT */
	printer.SetTableAlign(3)
	printer.SetTablePadding(IndentTwoSpaces)
	printer.SetTableFormatHeaders(true)
	if rootOutputFormat != "" && rootOutputFormat != common.TableFormat {
		printer.Printf(rootOutputFormat, data)
		return
	}
	printer.Printf(printer.SetTableAsDefaultFormat(rootOutputFormat), header, data)
}
