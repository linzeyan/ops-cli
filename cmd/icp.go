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
	Run: func(_ *cobra.Command, _ []string) {
		icp.ReadConf()
		result := icp.Check()
		fmt.Println(icp.Domain, result)
	},
	Example: Examples(`# Print the ICP status of the domain
ops-cli icp -d apple.com

# Print the ICP status of the domain and specify the configuration path
ops-cli icp -c ~/.env -d baidu.com`),
}

func init() {
	rootCmd.AddCommand(icpCmd)

	icpCmd.Flags().StringVarP(&icp.ConfigFile, "config", "c", "", "Specify config file")
	icpCmd.Flags().StringVarP(&icp.Domain, "domain", "d", "", "Specify domain name")
	icpCmd.MarkFlagRequired("domain")
}
