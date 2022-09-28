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
	"reflect"
	"strconv"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/cobra"
)

func init() {
	partition, err := disk.Partitions(true)
	if err != nil {
		return
	}
	var validArgs []string
	for _, v := range partition {
		validArgs = append(validArgs, v.Mountpoint)
	}

	var dfCmd = &cobra.Command{
		Use:       CommandDf,
		Short:     "Display free disk spaces",
		ValidArgs: validArgs,
		Args:      cobra.OnlyValidArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			var d DF
			var data [][]string
			switch {
			case len(args) == 0:
				for _, v := range partition {
					usage, err := disk.Usage(v.Mountpoint)
					if err != nil {
						return err
					}
					d.ParseDevices(usage, partition)
					data = append(data, d.OutputData())
				}
			default:
				for _, v := range args {
					usage, err := disk.Usage(v)
					if err != nil {
						return err
					}
					d.ParseDevices(usage, partition)
					data = append(data, d.OutputData())
				}
			}
			d.String(data)
			return err
		},
	}
	rootCmd.AddCommand(dfCmd)
}

type DF struct {
	Filesystem  string `json:"Filesystem"`
	Size        string `json:"Size"`
	Used        string `json:"Used"`
	Avail       string `json:"Avail"`
	UsedPercent string `json:"Use%"`
	MountedOn   string `json:"Mounted on"`
	Fstype      string `json:"FsType"`
	// Opts        []string `json:"Opts"`
	// InodesSize        uint64 `json:"iSize"`
	InodesUsed        string `json:"iUsed"`
	InodesFree        string `json:"iFree"`
	InodesUsedPercent string `json:"iUse%"`
}

func (d *DF) ParseDevices(usage *disk.UsageStat, partition []disk.PartitionStat) {
	for _, v := range partition {
		if usage.Path == v.Mountpoint {
			d.Filesystem = v.Device
			d.Size = common.ByteSize(usage.Total).String()
			d.Used = common.ByteSize(usage.Used).String()
			d.Avail = common.ByteSize(usage.Free).String()
			d.UsedPercent = fmt.Sprintf("%.2f%%", usage.UsedPercent)
			d.MountedOn = v.Mountpoint
			d.Fstype = v.Fstype
			// d.Opts = v.Opts
			// d.InodesSize = usage.InodesTotal
			d.InodesUsed = strconv.Itoa(int(usage.InodesUsed))
			d.InodesFree = strconv.Itoa(int(usage.InodesFree))
			d.InodesUsedPercent = fmt.Sprintf("%.2f%%", usage.InodesUsedPercent)
			break
		}
	}
}

func (d DF) OutputData() []string {
	var value []string
	f := reflect.ValueOf(&d).Elem()
	for i := 0; i < f.NumField(); i++ {
		value = append(value, fmt.Sprintf("%s", f.Field(i).Interface()))
	}
	return value
}

func (d DF) String(value any) {
	var header []string
	f := reflect.ValueOf(&d).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		tag := strings.TrimRight(strings.Replace(string(t.Field(i).Tag), `json:"`, "", 1), `"`)
		header = append(header, tag)
	}

	var data [][]string
	switch i := value.(type) {
	case []string:
		data = append(data, i)
	case [][]string:
		data = i
	}
	PrintTable(header, data, tablewriter.ALIGN_LEFT, IndentTwoSpaces, false)
}
