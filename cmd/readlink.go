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
	"os"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initReadlink() *cobra.Command {
	var flags bool
	var readlinkCmd = &cobra.Command{
		Use:   CommandReadlink,
		Short: "Get symlink information",
		Run: func(_ *cobra.Command, args []string) {
			var out string
			var err error
			for _, v := range args {
				temp := v
				for {
					out, err = os.Readlink(temp)
					if err != nil {
						if !flags {
							break
						}
						out = temp
						break
					}
					if !flags {
						break
					}
					temp = out
				}
				printer.Printf(out)
			}
		},
	}
	readlinkCmd.Flags().BoolVarP(&flags, "follow", "f", false, common.Usage("Follow all symlinks"))
	return readlinkCmd
}
