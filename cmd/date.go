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
	var flags struct {
		utc          bool
		seconds      bool
		milliseconds bool
		microseconds bool
		nanoseconds  bool
		format       string
		timezone     string
		date         bool
		time         bool
	}
	var dateCmd = &cobra.Command{
		Use:   CommandDate,
		Short: "Print date time",
		Long: "Print date time" + common.Usage(`

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
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, _ []string) {
			/* Set timezone. */
			t := common.TimeNow
			switch {
			case flags.utc:
				t = common.TimeNow.UTC()
			case flags.timezone != "":
				z, err := time.LoadLocation(flags.timezone)
				if err != nil {
					PrintString(err)
					return
				}
				t = common.TimeNow.In(z)
			}
			/* Print format. */
			switch {
			case flags.date:
				PrintString(t.Format("2006-01-02"))
			case flags.time:
				PrintString(t.Format("15:04:05"))
			case flags.seconds:
				PrintString(t.Unix())
			case flags.milliseconds:
				PrintString(t.UnixMilli())
			case flags.microseconds:
				PrintString(t.UnixMicro())
			case flags.nanoseconds:
				PrintString(t.UnixNano())
			case flags.format == "":
				PrintString(t.Format("2006-01-02T15:04:05-07:00"))
			case flags.format != "":
				PrintString(t.Format(flags.format))
			}
		},
	}
	rootCmd.AddCommand(dateCmd)
	dateCmd.Flags().StringVarP(&flags.format, "format", "f", "", common.Usage("Print date using specific format"))
	dateCmd.Flags().StringVarP(&flags.timezone, "timezone", "z", "", common.Usage("Specify timezone"))
	dateCmd.Flags().BoolVarP(&flags.seconds, "seconds", "s", false, common.Usage("Print Unix time"))
	dateCmd.Flags().BoolVarP(&flags.milliseconds, "milliseconds", "m", false, common.Usage("Print Unix time in milliseconds"))
	dateCmd.Flags().BoolVarP(&flags.microseconds, "microseconds", "M", false, common.Usage("Print Unix time in microseconds"))
	dateCmd.Flags().BoolVarP(&flags.nanoseconds, "nanoseconds", "n", false, common.Usage("Print Unix time in nanoseconds"))
	dateCmd.Flags().BoolVarP(&flags.utc, "utc", "u", false, common.Usage("Print date using UTC time"))
	dateCmd.Flags().BoolVarP(&flags.date, "date", "D", false, common.Usage(`Print date using '2006-01-02' format`))
	dateCmd.Flags().BoolVarP(&flags.time, "time", "T", false, common.Usage("Print time using '15:04:05' format"))
}
