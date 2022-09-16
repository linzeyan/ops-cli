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
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func init() {
	var lineFlag LineFlag
	var lineCmd = &cobra.Command{
		Use:   CommandLINE,
		Short: "Send message to LINE",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	var lineSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send message to LINE",
		Example: common.Examples(`# Send text to LINE chat
-s secret -t token --id GroupID -a 'Hello LINE'`, CommandLINE, CommandText),
		RunE: lineFlag.RunE,
	}

	var lineSubCmdID = &cobra.Command{
		Use:   CommandID,
		Args:  cobra.ExactArgs(1),
		Short: "Get chat ID from LINE",
		Example: common.Examples(`# Get chat ID from LINE,
# execute this command will listen on 80 port,
# type and sent 'id' in the chat,
# then console will print ID.
https://callback_url`, CommandLINE, CommandID),
		RunE: lineFlag.RunE,
	}

	var lineSubCmdPhoto = &cobra.Command{
		Use:   CommandPhoto,
		Short: "Send photo to LINE",
		Example: common.Examples(`# Send photo to LINE chat
-s secret -t token --id GroupID -a https://img.url`, CommandLINE, CommandPhoto),
		RunE: lineFlag.RunE,
	}

	var lineSubCmdVideo = &cobra.Command{
		Use:   CommandVideo,
		Short: "Send video to LINE",
		Example: common.Examples(`# Send video to LINE chat
-s secret -t token --id GroupID -a https://video.url`, CommandLINE, CommandVideo),
		RunE: lineFlag.RunE,
	}
	rootCmd.AddCommand(lineCmd)

	lineCmd.PersistentFlags().StringVarP(&lineFlag.Secret, "secret", "s", "", common.Usage("Channel Secret"))
	lineCmd.PersistentFlags().StringVarP(&lineFlag.Token, "token", "t", "", common.Usage("Channel Access Token"))
	lineCmd.PersistentFlags().StringVarP(&lineFlag.arg, "arg", "a", "", common.Usage("Function Argument"))
	lineCmd.PersistentFlags().StringVar(&lineFlag.ID, "id", "", common.Usage("UserID/GroupID/RoomID"))

	lineCmd.AddCommand(lineSubCmdID, lineSubCmdPhoto, lineSubCmdText, lineSubCmdVideo)
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
	if err = ReadConfig(CommandLINE, l); err != nil {
		return err
	}
	if l.Secret == "" || l.Token == "" {
		return common.ErrInvalidToken
	}
	l.api, err = linebot.New(l.Secret, l.Token)
	if l.api == nil {
		return common.ErrFailedInitial
	}
	return err
}

func (l *LineFlag) GetID() {
	var err error
	http.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		events, err := l.api.ParseRequest(r)
		if err != nil {
			PrintString(err)
			os.Exit(1)
		}
		for i := range events {
			if events[i].Type == linebot.EventTypeMessage && events[i].Message.(*linebot.TextMessage).Text == CommandID {
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

	server := &http.Server{
		Addr: ":80",

		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       300 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		PrintString("error starting server: " + err.Error())
		os.Exit(1)
	}
}

func (l *LineFlag) RunE(cmd *cobra.Command, args []string) error {
	if l.arg == "" {
		return common.ErrInvalidFlag
	}
	var err error
	if err = l.Init(); err != nil {
		return err
	}
	switch cmd.Name() {
	case CommandID:
		l.resp, err = l.api.SetWebhookEndpointURL(args[0]).Do()
		if err != nil {
			return err
		}
		l.GetID()
		return err
	case CommandText:
		input := linebot.NewTextMessage(l.arg)
		l.resp, err = l.api.PushMessage(l.ID, input).Do()
	case CommandPhoto:
		input := linebot.NewImageMessage(l.arg, l.arg)
		l.resp, err = l.api.PushMessage(l.ID, input).Do()
	case CommandVideo:
		input := linebot.NewVideoMessage(l.arg, l.arg)
		l.resp, err = l.api.PushMessage(l.ID, input).Do()
	}
	if err != nil {
		return err
	}
	OutputDefaultNone(l.resp)
	return err
}
