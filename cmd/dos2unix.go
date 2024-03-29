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
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initDos2Unix() *cobra.Command {
	var dos2unixCmd = &cobra.Command{
		Use:   CommandDos2Unix + " file...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Convert file eol to unix style",
		Run: func(_ *cobra.Command, args []string) {
			for _, f := range args {
				if err := common.Dos2Unix(f); err != nil {
					logger.Warn(err.Error(), common.DefaultField(f))
					printer.Printf("%s: %s\n", f, err)
					continue
				}
				printer.Printf("Converting file %s to Unix format...\n", f)
			}
		},
		DisableFlagsInUseLine: true,
	}
	return dos2unixCmd
}
