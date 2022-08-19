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

	"github.com/linzeyan/icp"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var icpCmd = &cobra.Command{
	Use:   "icp [domain]",
	Args:  cobra.ExactArgs(1),
	Short: "Check ICP status",
	Run: func(_ *cobra.Command, args []string) {
		if (icpCmdGlobalVar.Account == "" || icpCmdGlobalVar.Key == "") && rootConfig != "" {
			v, err := common.Config(rootConfig, common.ICP)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			err = Encoder.JSONMarshaler(v, &icpCmdGlobalVar)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}
		if icpCmdGlobalVar.Account == "" || icpCmdGlobalVar.Key == "" {
			log.Println(ErrTokenNotFound)
			os.Exit(1)
		}
		icp.Domain = args[0]
		icp.WestAccount = icpCmdGlobalVar.Account
		icp.WestApiKey = icpCmdGlobalVar.Key
		OutputDefaultYAML(map[string]string{icp.Domain: icp.Check()})
	},
	Example: common.Examples(`# Print the ICP status
ops-cli icp -a account -k api_key google.com`),
}

var icpCmdGlobalVar IcpFlags

func init() {
	rootCmd.AddCommand(icpCmd)

	icpCmd.Flags().StringVarP(&icpCmdGlobalVar.Account, "account", "a", "", "Enter the WEST account")
	icpCmd.Flags().StringVarP(&icpCmdGlobalVar.Key, "key", "k", "", "Enter the WEST api key")
	icpCmd.MarkFlagsRequiredTogether("account", "key")
}

type IcpFlags struct {
	Account string `json:"account"`
	Key     string `json:"api_key"`
}
