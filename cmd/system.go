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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Display system informations",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var systemSubCmdCPU = &cobra.Command{
	Use:   "cpu",
	Short: "Display cpu informations",
	Run: func(_ *cobra.Command, _ []string) {
		if err := systemCmdGlobalVar.CPUInfo(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(systemCmdGlobalVar.cpuResp)
	},
}

var systemSubCmdDisk = &cobra.Command{
	Use:   "disk",
	Short: "Display disk informations",
	Run: func(_ *cobra.Command, _ []string) {
		if err := systemCmdGlobalVar.DiskUsage(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(systemCmdGlobalVar.diskResp)
	},
}

var systemSubCmdHost = &cobra.Command{
	Use:   "host",
	Short: "Display host informations",
	Run: func(_ *cobra.Command, _ []string) {
		if err := systemCmdGlobalVar.HostInfo(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(systemCmdGlobalVar.hostResp)
	},
}

var systemSubCmdLoad = &cobra.Command{
	Use:   "load",
	Short: "Display load informations",
	Run: func(_ *cobra.Command, _ []string) {
		if err := systemCmdGlobalVar.LoadAvg(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(systemCmdGlobalVar.loadResp)
	},
}

var systemSubCmdMemory = &cobra.Command{
	Use:   "memory",
	Short: "Display memory informations",
	Run: func(_ *cobra.Command, _ []string) {
		if err := systemCmdGlobalVar.MemUsage(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(systemCmdGlobalVar.memResp)
	},
}

var systemSubCmdNetwork = &cobra.Command{
	Use:   "network",
	Short: "Display network informations",
	Run: func(_ *cobra.Command, _ []string) {
		if err := systemCmdGlobalVar.NetInfo(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		OutputDefaultJSON(systemCmdGlobalVar.netResp)
	},
}

var systemCmdGlobalVar SystemFlag

func init() {
	rootCmd.AddCommand(systemCmd)

	systemCmd.AddCommand(systemSubCmdCPU)
	systemCmd.AddCommand(systemSubCmdDisk)
	systemCmd.AddCommand(systemSubCmdHost)
	systemCmd.AddCommand(systemSubCmdLoad)
	systemCmd.AddCommand(systemSubCmdMemory)
	systemCmd.AddCommand(systemSubCmdNetwork)

	systemSubCmdNetwork.Flags().BoolVarP(&systemCmdGlobalVar.aiface, "all-interfaces", "a", false, "Display all interfaces")
	systemSubCmdNetwork.Flags().BoolVarP(&systemCmdGlobalVar.iface, "interface", "i", false, "Display interfaces")
}

type SystemFlag struct {
	aiface bool
	iface  bool

	cpuResp  systemCPUInfoResponse
	diskResp systemDiskUsageResponse
	hostResp systemHostInfoResponse
	loadResp systemLoadAvgResponse
	memResp  systemMemInfoResponse
	netResp  systemNetIOResponse
}

type systemCPUInfoResponse struct {
	/* cpu.Info() */
	VendorID  string `json:"vendorId,omitempty" yaml:"vendorId,omitempty"`
	Cores     string `json:"cores,omitempty" yaml:"cores,omitempty"`
	ModelName string `json:"modelName,omitempty" yaml:"modelName,omitempty"`
	Mhz       string `json:"mhz,omitempty" yaml:"mhz,omitempty"`
	CacheSize string `json:"cacheSize,omitempty" yaml:"cacheSize,omitempty"`
}

func (s *SystemFlag) CPUInfo() error {
	info, err := cpu.Info()
	if err != nil {
		return err
	}
	s.cpuResp = systemCPUInfoResponse{
		VendorID:  info[0].VendorID,
		Cores:     fmt.Sprintf("%d", info[0].Cores),
		ModelName: info[0].ModelName,
		Mhz:       fmt.Sprintf("%d", int(info[0].Mhz)),
		CacheSize: fmt.Sprintf("%d", info[0].CacheSize),
	}
	return err
}

type systemDiskUsageResponse struct {
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	Fstype      string `json:"fstype,omitempty" yaml:"fstype,omitempty"`
	Total       string `json:"total,omitempty" yaml:"total,omitempty"`
	Free        string `json:"free,omitempty" yaml:"free,omitempty"`
	Used        string `json:"used,omitempty" yaml:"used,omitempty"`
	UsedPercent string `json:"usedPercent,omitempty" yaml:"usedPercent,omitempty"`
}

func (s *SystemFlag) DiskUsage() error {
	info, err := disk.Usage("/")
	if err != nil {
		return err
	}
	s.diskResp = systemDiskUsageResponse{
		Path:        info.Path,
		Fstype:      info.Fstype,
		Total:       ByteSize(info.Total).String(),
		Free:        ByteSize(info.Free).String(),
		Used:        ByteSize(info.Used).String(),
		UsedPercent: fmt.Sprintf("%0.2f%%", info.UsedPercent),
	}
	return err
}

type systemHostInfoResponse struct {
	Hostname             string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Uptime               string `json:"uptime,omitempty" yaml:"uptime,omitempty"`
	BootTime             string `json:"bootTime,omitempty" yaml:"bootTime,omitempty"`
	Procs                uint64 `json:"procs,omitempty" yaml:"procs,omitempty"`
	OS                   string `json:"os,omitempty" yaml:"os,omitempty"`
	Platform             string `json:"platform,omitempty" yaml:"platform,omitempty"`
	PlatformFamily       string `json:"platformFamily,omitempty" yaml:"platformFamily,omitempty"`
	PlatformVersion      string `json:"platformVersion,omitempty" yaml:"platformVersion,omitempty"`
	KernelVersion        string `json:"kernelVersion,omitempty" yaml:"kernelVersion,omitempty"`
	KernelArch           string `json:"kernelArch,omitempty" yaml:"kernelArch,omitempty"`
	VirtualizationSystem string `json:"virtualizationSystem,omitempty" yaml:"virtualizationSystem,omitempty"`
	VirtualizationRole   string `json:"virtualizationRole,omitempty" yaml:"virtualizationRole,omitempty"`
	HostID               string `json:"hostId,omitempty" yaml:"hostId,omitempty"`
}

func (s *SystemFlag) HostInfo() error {
	info, err := host.Info()
	if err != nil {
		return err
	}
	s.hostResp = systemHostInfoResponse{
		Hostname:             info.Hostname,
		Uptime:               (time.Second * time.Duration(info.Uptime)).String(),
		BootTime:             (time.Second * time.Duration(info.BootTime)).String(),
		Procs:                info.Procs,
		OS:                   info.OS,
		Platform:             info.Platform,
		PlatformFamily:       info.PlatformFamily,
		PlatformVersion:      info.PlatformVersion,
		KernelVersion:        info.KernelVersion,
		KernelArch:           info.KernelArch,
		VirtualizationSystem: info.VirtualizationSystem,
		VirtualizationRole:   info.VirtualizationRole,
		HostID:               info.HostID,
	}
	return err
}

type systemLoadAvgResponse struct {
	Load1  string `json:"load1,omitempty" yaml:"load1,omitempty"`
	Load5  string `json:"load5,omitempty" yaml:"load5,omitempty"`
	Load15 string `json:"load15,omitempty" yaml:"load15,omitempty"`
}

func (s *SystemFlag) LoadAvg() error {
	info, err := load.Avg()
	if err != nil {
		return err
	}
	s.loadResp = systemLoadAvgResponse{
		Load1:  fmt.Sprintf("%0.2f", info.Load1),
		Load5:  fmt.Sprintf("%0.2f", info.Load5),
		Load15: fmt.Sprintf("%0.2f", info.Load15),
	}
	return err
}

type systemMemInfoResponse struct {
	/* mem.VirtualMemory() */
	Total       string `json:"total,omitempty" yaml:"total,omitempty"`
	Available   string `json:"available,omitempty" yaml:"available,omitempty"`
	Free        string `json:"free,omitempty" yaml:"free,omitempty"`
	Used        string `json:"used,omitempty" yaml:"used,omitempty"`
	UsedPercent string `json:"usedPercent,omitempty" yaml:"usedPercent,omitempty"`
}

func (s *SystemFlag) MemUsage() error {
	info, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	s.memResp = systemMemInfoResponse{
		Total:       ByteSize(info.Total).String(),
		Available:   ByteSize(info.Available).String(),
		Free:        ByteSize(info.Free).String(),
		Used:        ByteSize(info.Used).String(),
		UsedPercent: fmt.Sprintf("%0.2f%%", info.UsedPercent),
	}
	return err
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

	Interfaces []systemNetInterfaceResponse `json:"interface,omitempty" yaml:"interface,omitempty"`
}

type systemNetInterfaceResponse struct {
	Index        int      `json:"index"`
	MTU          int      `json:"mtu"`
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardwareAddr"`
	Flags        []string `json:"flags"`

	Addrs net.InterfaceAddrList `json:"addrs"`
}

func (s *SystemFlag) NetInfo() error {
	info, err := net.IOCounters(false)
	if err != nil {
		return err
	}
	data, err := json.Marshal(info[0])
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.netResp); err != nil {
		return err
	}
	if !s.iface && !s.aiface {
		return err
	}
	inet, err := net.Interfaces()
	if err != nil {
		return err
	}

	for i := range inet {
		if s.aiface {
			s.netResp.Interfaces = append(s.netResp.Interfaces, systemNetInterfaceResponse{
				Index:        inet[i].Index,
				MTU:          inet[i].MTU,
				Name:         inet[i].Name,
				HardwareAddr: inet[i].HardwareAddr,
				Flags:        inet[i].Flags,
				Addrs:        inet[i].Addrs,
			})
		} else if len(inet[i].Addrs) != 0 && inet[i].HardwareAddr != "" {
			s.netResp.Interfaces = append(s.netResp.Interfaces, systemNetInterfaceResponse{
				Index:        inet[i].Index,
				MTU:          inet[i].MTU,
				Name:         inet[i].Name,
				HardwareAddr: inet[i].HardwareAddr,
				Flags:        inet[i].Flags,
				Addrs:        inet[i].Addrs,
			})
		}
	}
	return err
}
