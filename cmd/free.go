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
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

func initFree() *cobra.Command {
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
		Run: func(_ *cobra.Command, _ []string) {
			var err error
			var f Free
			if flags.count == 0 && flags.second == 0 {
				if err = f.Output(); err != nil {
					logger.Info(err.Error())
				}
				return
			}
			var counter uint
			for {
				if flags.second <= 0 {
					flags.second = 2
				}
				err = f.Output()
				if err != nil {
					logger.Info(err.Error())
					return
				}
				counter++
				if flags.count > 0 && flags.count == counter {
					return
				}
				printer.Printf("\n")
				time.Sleep(time.Second * time.Duration(flags.second))
			}
		},
	}
	freeCmd.Flags().UintVarP(&flags.count, "count", "c", 0, common.Usage("Repeat printing times"))
	freeCmd.Flags().UintVarP(&flags.second, "seconds", "s", 0, common.Usage("Seconds between each repeat printing"))
	return freeCmd
}

type Free struct{}

func (f *Free) Output() error {
	var err error
	swap, err := mem.SwapMemory()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	memory, err := mem.VirtualMemory()
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	var header = []string{"", "total", "used", "free", "available", "use%"}
	var data [][]string
	data = append(data, []string{
		"Mem: ",
		common.ByteSize(memory.Total),
		common.ByteSize(memory.Used),
		common.ByteSize(memory.Free),
		common.ByteSize(memory.Available),
		fmt.Sprintf("%.2f%%", memory.UsedPercent),
	})
	data = append(data, []string{
		"Swap:",
		common.ByteSize(swap.Total),
		common.ByteSize(swap.Used),
		common.ByteSize(swap.Free),
		"",
		fmt.Sprintf("%.2f%%", swap.UsedPercent),
	})
	f.String(header, data)
	return err
}

func (Free) String(header []string, data [][]string) {
	/* tablewriter.ALIGN_RIGHT */
	printer.SetTableAlign(2)
	printer.SetTablePadding("\t ")
	if rootOutputFormat != "" && rootOutputFormat != common.TableFormat {
		printer.Printf(rootOutputFormat, data)
		return
	}
	printer.Printf(printer.SetTableAsDefaultFormat(rootOutputFormat), header, data)
}
