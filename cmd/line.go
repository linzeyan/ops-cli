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
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var lineCmd = &cobra.Command{
	Use:   "LINE",
	Short: "Send message to LINE",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var lineSubCmdText = &cobra.Command{
	Use:   "text",
	Short: "Send message to LINE",
	Example: common.Examples(`# Send text to LINE chat
ops-cli LINE text -s secret -t token --id GroupID -a 'Hello LINE'`),
	Run: lineCmdGlobalVar.Run,
}

var lineSubCmdID = &cobra.Command{
	Use:   "id",
	Args:  cobra.ExactArgs(1),
	Short: "Get chat ID from LINE",
	Example: common.Examples(`# Get chat ID from LINE,
# execute this command will listen on 80 port,
# type and sent 'id' in the chat,
# then console will print ID.
ops-cli LINE id https://callback_url`),
	Run: func(_ *cobra.Command, args []string) {
		var err error
		if err = lineCmdGlobalVar.Init(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		lineCmdGlobalVar.resp, err = lineCmdGlobalVar.api.SetWebhookEndpointURL(args[0]).Do()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		lineCmdGlobalVar.GetID()
	},
}

var lineSubCmdPhoto = &cobra.Command{
	Use:   "photo",
	Short: "Send photo to LINE",
	Example: common.Examples(`# Send photo to LINE chat
ops-cli LINE photo -s secret -t token --id GroupID -a https://img.url`),
	Run: lineCmdGlobalVar.Run,
}

var lineSubCmdVideo = &cobra.Command{
	Use:   "video",
	Short: "Send video to LINE",
	Example: common.Examples(`# Send video to LINE chat
ops-cli LINE video -s secret -t token --id GroupID -a https://video.url`),
	Run: lineCmdGlobalVar.Run,
}

var lineCmdGlobalVar LineFlag

func init() {
	rootCmd.AddCommand(lineCmd)

	lineCmd.PersistentFlags().StringVarP(&lineCmdGlobalVar.Secret, "secret", "s", "", common.Usage("Channel Secret"))
	lineCmd.PersistentFlags().StringVarP(&lineCmdGlobalVar.Token, "token", "t", "", common.Usage("Channel Access Token"))
	lineCmd.PersistentFlags().StringVarP(&lineCmdGlobalVar.arg, "arg", "a", "", common.Usage("Function Argument"))
	lineCmd.PersistentFlags().StringVar(&lineCmdGlobalVar.ID, "id", "", common.Usage("UserID/GroupID/RoomID"))

	lineCmd.AddCommand(lineSubCmdID)
	lineCmd.AddCommand(lineSubCmdPhoto)
	lineCmd.AddCommand(lineSubCmdText)
	lineCmd.AddCommand(lineSubCmdVideo)
}

type LineFlag struct {
	Secret string `json:"secret"`
	Token  string `json:"access_token"`
	ID     string `json:"id"`
	arg    string

	api  *linebot.Client
	resp *linebot.BasicResponse
}

func (l *LineFlag) Init() error {
	var err error
	if (l.Secret == "" || l.Token == "") && rootConfig != "" {
		v := common.Config(rootConfig, common.LINE)
		err = Encoder.JSONMarshaler(v, l)
		if err != nil {
			return err
		}
	}
	if l.Secret == "" || l.Token == "" {
		return ErrTokenNotFound
	}
	l.api, err = linebot.New(l.Secret, l.Token)
	if l.api == nil {
		return ErrInitialFailed
	}
	return err
}

func (l *LineFlag) GetID() {
	var err error
	http.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		events, err := l.api.ParseRequest(r)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		for i := range events {
			if events[i].Type == linebot.EventTypeMessage && events[i].Message.(*linebot.TextMessage).Text == ImTypeID {
				switch s := events[i].Source; s.Type {
				case linebot.EventSourceTypeGroup:
					PrintString(s.GroupID)
					os.Exit(0)
				case linebot.EventSourceTypeRoom:
					PrintString(s.RoomID)
					os.Exit(0)
				case linebot.EventSourceTypeUser:
					PrintString(s.UserID)
					os.Exit(0)
				}
			}
		}
	})

	err = http.ListenAndServe(":80", nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func (l *LineFlag) Run(cmd *cobra.Command, _ []string) {
	if l.arg == "" {
		log.Println(ErrArgNotFound)
		os.Exit(1)
	}
	var err error
	if err = l.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	switch cmd.Name() {
	case ImTypeText:
		input := linebot.NewTextMessage(l.arg)
		l.resp, err = l.api.PushMessage(l.ID, input).Do()
	case ImTypePhoto:
		input := linebot.NewImageMessage(l.arg, l.arg)
		l.resp, err = l.api.PushMessage(l.ID, input).Do()
	case ImTypeVideo:
		input := linebot.NewVideoMessage(l.arg, l.arg)
		l.resp, err = l.api.PushMessage(l.ID, input).Do()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	OutputDefaultNone(l.resp)
}
