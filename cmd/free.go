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
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

func init() {
	var flags struct {
		count  uint
		second uint
	}
	var freeCmd = &cobra.Command{
		Use:   CommandFree,
		Short: "Display free memory spaces",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			var err error
			var f Free
			if flags.count == 0 && flags.second == 0 {
				return f.Output()
			}
			var counter uint
			for {
				if flags.second == 0 {
					flags.second = 2
				}
				err = f.Output()
				if err != nil {
					return err
				}
				counter++
				if flags.count > 0 && flags.count == counter {
					return err
				}
				PrintString("")
				time.Sleep(time.Second * time.Duration(flags.second))
			}
		},
	}
	RootCmd.AddCommand(freeCmd)
	freeCmd.Flags().UintVarP(&flags.count, "count", "c", 0, common.Usage("Repeat printing times"))
	freeCmd.Flags().UintVarP(&flags.second, "seconds", "s", 0, common.Usage("Seconds between each repeat printing"))
}

type Free struct{}

func (f *Free) Output() error {
	var err error
	swap, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	var header = []string{"", "total", "used", "free", "available", "use%"}
	var data [][]string
	data = append(data, []string{
		"Mem: ",
		common.ByteSize(memory.Total).String(),
		common.ByteSize(memory.Used).String(),
		common.ByteSize(memory.Free).String(),
		common.ByteSize(memory.Available).String(),
		fmt.Sprintf("%.2f%%", memory.UsedPercent),
	})
	data = append(data, []string{
		"Swap:",
		common.ByteSize(swap.Total).String(),
		common.ByteSize(swap.Used).String(),
		common.ByteSize(swap.Free).String(),
		"",
		fmt.Sprintf("%.2f%%", swap.UsedPercent),
	})
	f.String(header, data)
	return err
}

func (Free) String(header []string, data [][]string) {
	PrintTable(header, data, tablewriter.ALIGN_RIGHT, "\t ", false)
}
