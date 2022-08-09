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
	"reflect"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Display system informations",
	Run:   sysf.Run,
}

var sysf systemFlag

func init() {
	rootCmd.AddCommand(systemCmd)
}

type systemFlag struct {
	cpuResp     systemCPUInfoResponse
	cpuTimeResp systemCPUTimesResponse
	diskResp    systemDiskUsageResponse
	hostResp    systemHostInfoResponse
	memResp     systemMemInfoResponse
}

type systemCPUInfoResponse struct {
	/* cpu.Info() */
	VendorID  string `json:"vendorId,omitempty" yaml:"vendorId,omitempty"`
	Cores     int    `json:"cores,omitempty" yaml:"cores,omitempty"`
	ModelName string `json:"modelName,omitempty" yaml:"modelName,omitempty"`
	Mhz       int    `json:"mhz,omitempty" yaml:"mhz,omitempty"`
	CacheSize int    `json:"cacheSize,omitempty" yaml:"cacheSize,omitempty"`
}

func (s *systemFlag) CPUInfo() error {
	info, err := cpu.Info()
	if err != nil {
		return err
	}
	data, err := json.Marshal(info[0])
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.cpuResp); err != nil {
		return err
	}
	return nil
}

type systemCPUTimesResponse struct {
	/* cpu.Times(false) */
	User      string `json:"user,omitempty" yaml:"user,omitempty"`
	System    string `json:"system,omitempty" yaml:"system,omitempty"`
	Idle      string `json:"idle,omitempty" yaml:"idle,omitempty"`
	Nice      string `json:"nice,omitempty" yaml:"nice,omitempty"`
	Iowait    string `json:"iowait,omitempty" yaml:"iowait,omitempty"`
	Irq       string `json:"irq,omitempty" yaml:"irq,omitempty"`
	Softirq   string `json:"softirq,omitempty" yaml:"softirq,omitempty"`
	Steal     string `json:"steal,omitempty" yaml:"steal,omitempty"`
	Guest     string `json:"guest,omitempty" yaml:"guest,omitempty"`
	GuestNice string `json:"guestNice,omitempty" yaml:"guestNice,omitempty"`
}

func (s *systemFlag) CPUTimes() error {
	info, err := cpu.Times(false)
	if err != nil {
		return err
	}
	s.cpuTimeResp = systemCPUTimesResponse{
		User:   (time.Second * time.Duration(info[0].User)).String(),
		System: (time.Second * time.Duration(info[0].System)).String(),
		Idle:   (time.Second * time.Duration(info[0].Idle)).String(),
	}
	switch {
	case info[0].Nice != 0:
		s.cpuTimeResp.Nice = (time.Second * time.Duration(info[0].Nice)).String()
		fallthrough
	case info[0].Iowait != 0:
		s.cpuTimeResp.Iowait = (time.Second * time.Duration(info[0].Iowait)).String()
		fallthrough
	case info[0].Irq != 0:
		s.cpuTimeResp.Irq = (time.Second * time.Duration(info[0].Irq)).String()
		fallthrough
	case info[0].Softirq != 0:
		s.cpuTimeResp.Softirq = (time.Second * time.Duration(info[0].Softirq)).String()
		fallthrough
	case info[0].Steal != 0:
		s.cpuTimeResp.Steal = (time.Second * time.Duration(info[0].Steal)).String()
		fallthrough
	case info[0].Guest != 0:
		s.cpuTimeResp.Guest = (time.Second * time.Duration(info[0].Guest)).String()
		fallthrough
	case info[0].GuestNice != 0:
		s.cpuTimeResp.GuestNice = (time.Second * time.Duration(info[0].GuestNice)).String()
	}
	return nil
}

type systemDiskUsageResponse struct {
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	Fstype      string `json:"fstype,omitempty" yaml:"fstype,omitempty"`
	Total       string `json:"total,omitempty" yaml:"total,omitempty"`
	Free        string `json:"free,omitempty" yaml:"free,omitempty"`
	Used        string `json:"used,omitempty" yaml:"used,omitempty"`
	UsedPercent string `json:"usedPercent,omitempty" yaml:"usedPercent,omitempty"`
}

func (s *systemFlag) DiskUsage() error {
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
	return nil
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

	Temperature []systemHostTemperatureResponse `json:"temperature,omitempty" yaml:"temperature,omitempty"`
}

type systemHostTemperatureResponse struct {
	SensorKey   string  `json:"sensorKey,omitempty" yaml:"sensorKey,omitempty"`
	Temperature float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	High        float64 `json:"sensorHigh,omitempty" yaml:"sensorHigh,omitempty"`
	Critical    float64 `json:"sensorCritical,omitempty" yaml:"sensorCritical,omitempty"`
}

func (s *systemFlag) HostInfo() error {
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

	temp, err := host.SensorsTemperatures()
	if err != nil {
		return err
	}
	for i := range temp {
		if temp[i].Temperature != 0 {
			s.hostResp.Temperature = append(s.hostResp.Temperature, systemHostTemperatureResponse{
				SensorKey:   temp[i].SensorKey,
				Temperature: temp[i].Temperature,
				High:        temp[i].High,
				Critical:    temp[i].Critical,
			})
		}
	}
	return nil
}

type systemMemInfoResponse struct {
	/* mem.VirtualMemory() */
	Total       string `json:"total,omitempty" yaml:"total,omitempty"`
	Available   string `json:"available,omitempty" yaml:"available,omitempty"`
	Free        string `json:"free,omitempty" yaml:"free,omitempty"`
	Used        string `json:"used,omitempty" yaml:"used,omitempty"`
	UsedPercent string `json:"usedPercent,omitempty" yaml:"usedPercent,omitempty"`
}

func (s *systemFlag) MemUsage() error {
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
	return nil
}

func (s systemFlag) Run(cmd *cobra.Command, args []string) {

}

func (s systemFlag) JSON() { PrintJSON(s) }

func (s systemFlag) YAML() { PrintYAML(s) }

func (s systemFlag) String() {
	f := reflect.ValueOf(&s).Elem()
	t := f.Type()
	for i := 0; i < f.NumField(); i++ {
		fmt.Printf("%s\t%v\n", t.Field(i).Name, f.Field(i).Interface())
	}
}
