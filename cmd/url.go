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
	"log"
	"os"

	"github.com/linzeyan/expandUrl"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url [shorten url]",
	Args:  cobra.ExactArgs(1),
	Short: "Expand shorten url",
	Run: func(_ *cobra.Command, args []string) {
		result, err := expandUrl.Expand(args[0])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		PrintString(result)
	},
	Example: common.Examples(`# Get the real URL from the shortened URL
ops-cli url https://goo.gl/maps/b37Aq3Anc7taXQDd9`),
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(urlCmd)
}
