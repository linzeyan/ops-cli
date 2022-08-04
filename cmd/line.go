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

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/cobra"
)

var lineCmd = &cobra.Command{
	Use:   "line [function]",
	Short: "Send message to LINE",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		if args[0] == "" {
			_ = cmd.Help()
			return
		}
		Config(configLINE)
		if line.secret == "" || line.token == "" {
			_ = cmd.Help()
			return
		}
		lineAPI := line.Init()
		if lineAPI == nil {
			_ = cmd.Help()
			return
		}
		switch strings.ToLower(args[0]) {
		case imTypeMsg, imTypeMessage:
			line.Msg(lineAPI)
		case imTypePhoto:
			line.Photo(lineAPI)
		case imTypeVideo:
			line.Video(lineAPI)
		}
	},
	Example: Examples(`# Send message to LINE chat
ops-cli line msg -s secret -t token --id GroupID -a 'Hello LINE'`),
}

var line lineFlag

func init() {
	rootCmd.AddCommand(lineCmd)

	lineCmd.Flags().StringVarP(&line.secret, "secret", "s", "", "Channel Secret")
	lineCmd.Flags().StringVarP(&line.token, "token", "t", "", "Channel Access Token")
	lineCmd.Flags().StringVarP(&line.arg, "arg", "a", "", "Function Argument")
	lineCmd.Flags().StringVar(&line.id, "id", "", "UserID/GroupID/RoomID")
}

type lineFlag struct {
	secret string
	token  string
	id     string
	arg    string
}

func (l lineFlag) Init() *linebot.Client {
	api, err := linebot.New(l.secret, l.token)
	if err != nil {
		log.Println(err)
		return nil
	}
	return api
}

func (l lineFlag) Msg(api *linebot.Client) {
	input := linebot.NewTextMessage(l.arg)
	_, err := api.PushMessage(l.id, input).Do()
	if err != nil {
		log.Println(err)
		return
	}
}

func (l lineFlag) Photo(api *linebot.Client) {
	input := linebot.NewImageMessage(l.arg, l.arg)
	_, err := api.PushMessage(l.id, input).Do()
	if err != nil {
		log.Println(err)
		return
	}
}

func (l lineFlag) Video(api *linebot.Client) {
	input := linebot.NewVideoMessage(l.arg, l.arg)
	_, err := api.PushMessage(l.id, input).Do()
	if err != nil {
		log.Println(err)
		return
	}
}
