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
	"strings"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var slackCmd = &cobra.Command{
	Use:   "slack [function]",
	Short: "Send message to slack",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		if args[0] == "" {
			_ = cmd.Help()
			return
		}
		Config("slack")
		if slk.token == "" {
			_ = cmd.Help()
			return
		}
		slackAPI := slk.Init()
		if slackAPI == nil {
			_ = cmd.Help()
			return
		}
		switch strings.ToLower(args[0]) {
		case "msg", "message":
			slk.Msg(slackAPI)
		}
	}, Example: Examples(`# Send message
ops-cli slack msg -a "Hello World!"`),
}

var slk slackFlag

func init() {
	rootCmd.AddCommand(slackCmd)

	slackCmd.Flags().StringVarP(&slk.token, "token", "t", "", "Bot token (required)")
	slackCmd.Flags().StringVarP(&slk.channel, "channel", "c", "", "Channel ID")
	slackCmd.Flags().StringVarP(&slk.arg, "arg", "a", "", "Input argument")
}

type slackFlag struct {
	token   string
	channel string
	arg     string
}

func (s slackFlag) Init() *slack.Client {
	return slack.New(s.token)
}

func (s slackFlag) Msg(api *slack.Client) {
	input := slack.MsgOptionText(s.arg, false)
	_, _, _, err := api.SendMessageContext(rootContext, s.channel, input)
	if err != nil {
		log.Println(err)
	}
}
