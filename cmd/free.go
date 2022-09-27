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
	var freeFlag FreeFlag
	var freeCmd = &cobra.Command{
		Use:   CommandFree,
		Short: "Display free memory spaces",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: freeFlag.RunE,
	}
	rootCmd.AddCommand(freeCmd)
	freeCmd.Flags().UintVarP(&freeFlag.count, "count", "c", 0, common.Usage("Repeat printing times"))
	freeCmd.Flags().UintVarP(&freeFlag.second, "seconds", "s", 0, common.Usage("Seconds between each repeat printing"))
}

type FreeFlag struct {
	count  uint
	second uint
}

func (f *FreeFlag) RunE(_ *cobra.Command, _ []string) error {
	var err error
	if f.count == 0 && f.second == 0 {
		return f.Output()
	}
	var counter uint
	for {
		err = f.Output()
		if err != nil {
			return err
		}
		counter++
		if f.count > 0 && f.count == counter {
			return err
		}
		if f.second == 0 {
			f.second = 2
		}
		time.Sleep(time.Second * time.Duration(f.second))
	}
}

func (f *FreeFlag) Output() error {
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

func (FreeFlag) String(header []string, data [][]string) {
	PrintTable(header, data, tablewriter.ALIGN_RIGHT, "\t ", false)
}
