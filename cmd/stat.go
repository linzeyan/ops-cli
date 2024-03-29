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
	"github.com/spf13/cobra"
)

func initStat() *cobra.Command {
	var statCmd = &cobra.Command{
		Use:   CommandStat + " path...",
		Short: "Display file informations",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			for _, v := range args {
				var s FileStat
				if err := s.String(v); err != nil {
					logger.Info(err.Error())
					return
				}
			}
		},
		DisableFlagsInUseLine: true,
	}
	return statCmd
}
