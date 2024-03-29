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
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
)

func initSystem() *cobra.Command {
	var systemCmd = &cobra.Command{
		Use:   CommandSystem,
		Short: "Display system informations",
		RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	run := func(cmd *cobra.Command, _ []string) {
		var s System
		var err error
		var resp any
		switch cmd.Name() {
		case CommandCPU:
			resp, err = s.CPUInfo()
		case CommandDisk:
			resp, err = s.DiskUsage()
		case CommandHost:
			resp, err = s.HostInfo()
		case CommandLoad:
			resp, err = s.LoadAvg()
		case CommandMemory:
			resp, err = s.MemUsage()
		case CommandNetwork:
			resp, err = s.NetInfo()
		}
		if err != nil {
			logger.Info(err.Error())
			return
		}
		printer.Printf(printer.SetJSONAsDefaultFormat(rootOutputFormat), resp)
	}

	var systemSubCmdCPU = &cobra.Command{
		Use:   CommandCPU,
		Short: "Display cpu informations",
		Run:   run,

		DisableFlagsInUseLine: true,
	}

	var systemSubCmdDisk = &cobra.Command{
		Use:   CommandDisk,
		Short: "Display disk informations",
		Run:   run,

		DisableFlagsInUseLine: true,
	}

	var systemSubCmdHost = &cobra.Command{
		Use:   CommandHost,
		Short: "Display host informations",
		Run:   run,

		DisableFlagsInUseLine: true,
	}

	var systemSubCmdLoad = &cobra.Command{
		Use:   CommandLoad,
		Short: "Display load informations",
		Run:   run,

		DisableFlagsInUseLine: true,
	}

	var systemSubCmdMemory = &cobra.Command{
		Use:   CommandMemory,
		Short: "Display memory informations",
		Run:   run,

		DisableFlagsInUseLine: true,
	}

	var systemSubCmdNetwork = &cobra.Command{
		Use:   CommandNetwork,
		Short: "Display network informations",
		Run:   run,

		DisableFlagsInUseLine: true,
	}
	systemCmd.AddCommand(systemSubCmdCPU)
	systemCmd.AddCommand(systemSubCmdDisk)
	systemCmd.AddCommand(systemSubCmdHost)
	systemCmd.AddCommand(systemSubCmdLoad)
	systemCmd.AddCommand(systemSubCmdMemory)
	systemCmd.AddCommand(systemSubCmdNetwork)
	return systemCmd
}

type System struct{}

func (s *System) CPUInfo() (any, error) {
	info, err := cpu.InfoWithContext(common.Context)
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	resp := struct {
		VendorID, ModelName string
		Cores, CacheSize    int32
		GHz                 float64
	}{
		VendorID:  info[0].VendorID,
		Cores:     info[0].Cores,
		ModelName: info[0].ModelName,
		GHz:       info[0].Mhz / 1000,
		CacheSize: info[0].CacheSize,
	}
	return &resp, err
}

func (s *System) DiskUsage() (any, error) {
	info, err := disk.UsageWithContext(common.Context, "/")
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	resp := struct {
		Path, FsType, Total, Free, Used, UsedPercent string
	}{
		Path:        info.Path,
		FsType:      info.Fstype,
		Total:       common.ByteSize(info.Total),
		Free:        common.ByteSize(info.Free),
		Used:        common.ByteSize(info.Used),
		UsedPercent: fmt.Sprintf("%0.2f%%", info.UsedPercent),
	}
	return &resp, err
}

func (s *System) HostInfo() (any, error) {
	info, err := host.InfoWithContext(common.Context)
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	resp := struct {
		Hostname, HostID, Uptime, BootTime            string
		OS, Platform, PlatformFamily, PlatformVersion string
		KernelVersion, KernelArch                     string
		VirtualizationSystem, VirtualizationRole      string
		Procs                                         uint64
	}{
		Hostname:             info.Hostname,
		HostID:               info.HostID,
		Uptime:               (time.Second * time.Duration(info.Uptime)).String(),
		BootTime:             (time.Second * time.Duration(info.BootTime)).String(),
		OS:                   info.OS,
		Platform:             info.Platform,
		PlatformFamily:       info.PlatformFamily,
		PlatformVersion:      info.PlatformVersion,
		KernelVersion:        info.KernelVersion,
		KernelArch:           info.KernelArch,
		VirtualizationSystem: info.VirtualizationSystem,
		VirtualizationRole:   info.VirtualizationRole,
		Procs:                info.Procs,
	}
	return &resp, err
}

func (s *System) LoadAvg() (any, error) {
	info, err := load.AvgWithContext(common.Context)
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	resp := struct {
		Load1, Load5, Load15 string
	}{
		Load1:  fmt.Sprintf("%0.2f", info.Load1),
		Load5:  fmt.Sprintf("%0.2f", info.Load5),
		Load15: fmt.Sprintf("%0.2f", info.Load15),
	}
	return &resp, err
}

func (s *System) MemUsage() (any, error) {
	info, err := mem.VirtualMemoryWithContext(common.Context)
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	resp := struct {
		Total, Available, Free, Used, UsedPercent string
	}{
		Total:       common.ByteSize(info.Total),
		Available:   common.ByteSize(info.Available),
		Free:        common.ByteSize(info.Free),
		Used:        common.ByteSize(info.Used),
		UsedPercent: fmt.Sprintf("%0.2f%%", info.UsedPercent),
	}
	return &resp, err
}

func (s *System) NetInfo() (any, error) {
	info, err := net.IOCountersWithContext(common.Context, false)
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	type systemNetIOResponse struct {
		BytesSent   uint64 `json:"bytesSent,omitempty" yaml:"bytesSent,omitempty"`
		BytesRecv   uint64 `json:"bytesRecv,omitempty" yaml:"bytesRecv,omitempty"`
		PacketsSent uint64 `json:"packetsSent,omitempty" yaml:"packetsSent,omitempty"`
		PacketsRecv uint64 `json:"packetsRecv,omitempty" yaml:"packetsRecv,omitempty"`
		Errin       uint64 `json:"errin,omitempty" yaml:"errin,omitempty"`
		Errout      uint64 `json:"errout,omitempty" yaml:"errout,omitempty"`
		Dropin      uint64 `json:"dropin,omitempty" yaml:"dropin,omitempty"`
		Dropout     uint64 `json:"dropout,omitempty" yaml:"dropout,omitempty"`
		Fifoin      uint64 `json:"fifoin,omitempty" yaml:"fifoin,omitempty"`
		Fifoout     uint64 `json:"fifoout,omitempty" yaml:"fifoout,omitempty"`

		Interfaces net.InterfaceStatList `json:"interface,omitempty" yaml:"interface,omitempty"`
	}
	var netResp systemNetIOResponse
	err = Encoder.JSONMarshaler(info[0], &netResp)
	if err != nil {
		logger.Debug(err.Error())
		return &netResp, err
	}
	inet, err := net.InterfacesWithContext(common.Context)
	if err != nil {
		logger.Debug(err.Error())
		return &netResp, err
	}
	if err = Encoder.JSONMarshaler(inet, &netResp.Interfaces); err != nil {
		logger.Debug(err.Error())
	}
	return &netResp, err
}
