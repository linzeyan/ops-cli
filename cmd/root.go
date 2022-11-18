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

var (
	rootCmd = root()
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
	var cmd = &cobra.Command{
		Use:   common.RepoName,
		Short: "OPS useful tools",
		RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },

		DisableFlagsInUseLine: true,
	}
	cmd.PersistentFlags().StringVar(&rootOutputFormat, "output", "", common.Usage("Output format, can be json/yaml"))
	cmd.PersistentFlags().StringVar(&rootConfig, "config", "", common.Usage("Specify config path"))
	cmd.PersistentFlags().StringVar(&rootVerbose, "verbose", "error", common.Usage("Specify log level (debug/info/warn/error/panic/fatal"))
	cmd.PersistentFlags().BoolP("help", "", false, common.Usage("Help for this command"))

	cmd.AddCommand(initArping())
	cmd.AddCommand(initCert(), initConvert())
	cmd.AddCommand(initDate(), initDf(), initDig(), initDiscord(), initDoc(cmd), initDos2Unix())
	cmd.AddCommand(initEncode(), initEncrypt())
	cmd.AddCommand(initFree())
	cmd.AddCommand(initGeoip())
	cmd.AddCommand(initHash())
	cmd.AddCommand(initICP(), initIP())
	cmd.AddCommand(initLINE())
	cmd.AddCommand(initMTR())
	cmd.AddCommand(initNetmask())
	cmd.AddCommand(initOTP())
	cmd.AddCommand(initPing(), initPs())
	cmd.AddCommand(initQrcode())
	cmd.AddCommand(initRandom(), initReadlink(), initRedis())
	cmd.AddCommand(initSlack(), initSs(), initSSHKeyGen(), initSSL(), initStat(), initSystem())
	cmd.AddCommand(initTCPing(), initTelegram(), initTraceroute(), initTree())
	cmd.AddCommand(initUpdate(), initURL())
	cmd.AddCommand(initVersion())
	cmd.AddCommand(initWhois(), initWsping())
	initalize := func() {
		common.SetLoggerLevel(rootVerbose)
	}
	cobra.OnInitialize(initalize)
	addGroup(cmd)
	return cmd
}

func Run() *cobra.Command {
	return rootCmd
}

var groupings = map[string]string{
	CommandDiscord:  groupIM,
	CommandLINE:     groupIM,
	CommandSlack:    groupIM,
	CommandTelegram: groupIM,

	CommandArping:     groupNetwork,
	CommandDig:        groupNetwork,
	CommandGeoip:      groupNetwork,
	CommandIP:         groupNetwork,
	CommandMTR:        groupNetwork,
	CommandNetmask:    groupNetwork,
	CommandPing:       groupNetwork,
	CommandSs:         groupNetwork,
	CommandTCPing:     groupNetwork,
	CommandTraceroute: groupNetwork,
	CommandURL:        groupNetwork,
	CommandWhois:      groupNetwork,
	CommandWsping:     groupNetwork,
}

func addGroup(cmd *cobra.Command) {
	var groups []*cobra.Group
	im := &cobra.Group{ID: groupings[CommandDiscord], Title: groupings[CommandDiscord]}
	network := &cobra.Group{ID: groupings[CommandArping], Title: groupings[CommandArping]}

	groups = append(groups, im, network)
	cmd.AddGroup(groups...)
}

func getGroupID(cmd string) string {
	group, ok := groupings[cmd]
	if ok {
		return group
	}
	return ""
}

func ReadConfig(block string, flag any) error {
	v := common.Config(rootConfig, strings.ToLower(block))
	return Encoder.JSONMarshaler(v, flag)
}
