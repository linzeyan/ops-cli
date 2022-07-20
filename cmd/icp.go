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

	"github.com/linzeyan/icp"
	"github.com/spf13/cobra"
)

var icpCmd = &cobra.Command{
	Use:   "icp",
	Short: "Check ICP status",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			icp.Domain = args[0]
			icp.ReadConf()
			fmt.Println(icp.Domain, icp.Check())
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Print the ICP status of the domain
ops-cli icp apple.com

# Print the ICP status of the domain and specify the configuration path
ops-cli icp -c ~/.env baidu.com`),
}

func init() {
	rootCmd.AddCommand(icpCmd)

	icpCmd.Flags().StringVarP(&icp.ConfigFile, "config", "c", "", "Specify config file")
}
