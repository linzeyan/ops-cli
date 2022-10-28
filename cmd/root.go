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
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var RootCmd = root()
var (
	printer = common.NewPrinter()
	logger  = common.NewLogger()
)

/* Flags. */
var (
	rootConfig       string
	rootOutputFormat string
	rootVerbose      string
)

func root() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   common.RepoName,
		Short: "OPS useful tools",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}
	rootCmd.PersistentFlags().StringVar(&rootOutputFormat, "output", "", common.Usage("Output format, can be json/yaml"))
	rootCmd.PersistentFlags().StringVar(&rootConfig, "config", "", common.Usage("Specify config path"))
	rootCmd.PersistentFlags().StringVar(&rootVerbose, "verbose", "warn", common.Usage("Specify log level (debug/info/warn/error/panic/fatal"))
	rootCmd.PersistentFlags().BoolP("help", "", false, common.Usage("Help for this command"))

	rootCmd.AddCommand(initArping())
	rootCmd.AddCommand(initCert(), initConvert())
	rootCmd.AddCommand(initDate(), initDf(), initDig(), initDiscord(), initDoc(rootCmd), initDos2Unix())
	rootCmd.AddCommand(initEncode(), initEncrypt())
	rootCmd.AddCommand(initFree())
	rootCmd.AddCommand(initGeoip())
	rootCmd.AddCommand(initHash())
	rootCmd.AddCommand(initICP(), initIP())
	rootCmd.AddCommand(initLINE())
	rootCmd.AddCommand(initMTR())
	rootCmd.AddCommand(initNetmask())
	rootCmd.AddCommand(initOTP())
	rootCmd.AddCommand(initPing(), initPs())
	rootCmd.AddCommand(initQrcode())
	rootCmd.AddCommand(initRandom(), initReadlink(), initRedis())
	rootCmd.AddCommand(initSlack(), initSs(), initSSHKeyGen(), initSSL(), initStat(), initSystem())
	rootCmd.AddCommand(initTCPing(), initTelegram(), initTraceroute(), initTree())
	rootCmd.AddCommand(initUpdate(), initURL())
	rootCmd.AddCommand(initVersion())
	rootCmd.AddCommand(initWhois(), initWsping())
	initalize := func() {
		common.SetLoggerLevel(rootVerbose)
	}
	cobra.OnInitialize(initalize)
	addGroup(rootCmd)
	return rootCmd
}

var groupings = map[string]string{
	CommandDiscord:  groupIM,
	CommandLINE:     groupIM,
	CommandSlack:    groupIM,
	CommandTelegram: groupIM,

	CommandArping:     groupNetwork,
	CommandIP:         groupNetwork,
	CommandMTR:        groupNetwork,
	CommandNetmask:    groupNetwork,
	CommandPing:       groupNetwork,
	CommandSs:         groupNetwork,
	CommandTCPing:     groupNetwork,
	CommandTraceroute: groupNetwork,
	CommandWsping:     groupNetwork,
}

func addGroup(cmd *cobra.Command) {
	var groups []*cobra.Group
	im := &cobra.Group{ID: groupings[CommandDiscord], Title: groupings[CommandDiscord]}
	network := &cobra.Group{ID: groupings[CommandArping], Title: groupings[CommandArping]}

	groups = append(groups, im, network)
	cmd.AddGroup(groups...)
}

func ReadConfig(block string, flag any) error {
	v := common.Config(rootConfig, strings.ToLower(block))
	return Encoder.JSONMarshaler(v, flag)
}
