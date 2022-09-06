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
	"time"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func init() {
	var dateFlag DateFlag
	var dateCmd = &cobra.Command{
		Use:       CommandDate,
		ValidArgs: []string{"milli", "micro", "nano"},
		Short: "Print date time" + common.Usage(`

# Specific format use Golang time format
Year: "2006" "06"
Month: "Jan" "January" "01" "1"
Day of the week: "Mon" "Monday"
Day of the month: "2" "_2" "02"
Day of the year: "__2" "002"
Hour: "15" "3" "03" (PM or AM)
Minute: "4" "04"
Second: "5" "05"
Milli Second: ".000" ".999"
Micro Second: ".000000" ".999999"
Nano Second: ".000000000" ".999999999"
AM/PM mark: "PM"
Time zone:
"Z0700" "-0700"         Z or ±hhmm
"Z07:00" "-07:00"       Z or ±hh:mm
"Z07" "-07"             Z or ±hh
"Z070000" "-070000"     Z or ±hhmmss
"Z07:00:00" "-07:00:00" Z or ±hh:mm:ss`),
		Run: func(_ *cobra.Command, args []string) {
			if dateFlag.seconds {
				dateFlag.PrintUnixTime(args)
				return
			}
			if dateFlag.date {
				PrintString(dateFlag.Now().Format("2006-01-02"))
				return
			}
			if dateFlag.time {
				PrintString(dateFlag.Now().Format("15:04:05"))
				return
			}
			if dateFlag.format {
				dateFlag.PrintFormat(args)
				return
			}
			PrintString(dateFlag.Now().Format(time.RFC3339))
		},
	}
	rootCmd.AddCommand(dateCmd)
	dateCmd.Flags().BoolVarP(&dateFlag.format, "format", "f", false, common.Usage("Print date using specific format"))
	dateCmd.Flags().BoolVarP(&dateFlag.seconds, "seconds", "s", false, common.Usage("Print Unix time"))
	dateCmd.Flags().BoolVarP(&dateFlag.utc, "utc", "u", false, common.Usage("Print date using UTC time"))
	dateCmd.Flags().BoolVarP(&dateFlag.date, "date", "D", false, common.Usage(`Print date using '2006-01-02' format`))
	dateCmd.Flags().BoolVarP(&dateFlag.time, "time", "T", false, common.Usage("Print time using '15:04:05' format"))
}

type DateFlag struct {
	utc     bool
	seconds bool
	format  bool
	date    bool
	time    bool
}

func (d *DateFlag) Now() time.Time {
	if d.utc {
		return common.TimeNow.UTC()
	}
	return common.TimeNow
}

func (d *DateFlag) PrintFormat(args []string) {
	if len(args) == 0 {
		PrintString(d.Now().Format(time.RFC3339))
		return
	}
	PrintString(d.Now().Format(args[0]))
}

func (d *DateFlag) PrintUnixTime(args []string) {
	var s = map[string]string{
		"milli": time.StampMilli,
		"micro": time.StampMicro,
		"nano":  time.StampNano,
	}
	if len(args) == 0 {
		PrintString(common.TimeNow.Unix())
		return
	}
	switch s[args[0]] {
	default:
		PrintString(common.TimeNow.Unix())
	case time.StampMilli:
		PrintString(common.TimeNow.UnixMilli())
	case time.StampMicro:
		PrintString(common.TimeNow.UnixMicro())
	case time.StampNano:
		PrintString(common.TimeNow.UnixNano())
	}
}
