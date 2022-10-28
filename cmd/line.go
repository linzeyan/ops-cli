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

func initLINE() *cobra.Command {
	var flags struct {
		Secret string `json:"secret"`
		Token  string `json:"access_token"`
		ID     string `json:"id"`
		arg    string
	}
	var lineCmd = &cobra.Command{
		GroupID: groupings[CommandLINE],
		Use:     CommandLINE,
		Short:   "Send message to LINE",
		Run:     func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	runE := func(cmd *cobra.Command, args []string) error {
		if flags.arg == "" {
			return common.ErrInvalidFlag
		}
		var err error
		if rootConfig != "" {
			if err = ReadConfig(CommandLINE, &flags); err != nil {
				return err
			}
		}
		var l LINE
		if err = l.Init(flags.Secret, flags.Token); err != nil {
			return err
		}
		switch cmd.Name() {
		case CommandID:
			l.Response, err = l.API.SetWebhookEndpointURL(args[0]).Do()
			if err != nil {
				return err
			}
			l.GetID()
			return err
		case CommandText:
			input := linebot.NewTextMessage(flags.arg)
			l.Response, err = l.API.PushMessage(flags.ID, input).Do()
		case CommandPhoto:
			input := linebot.NewImageMessage(flags.arg, flags.arg)
			l.Response, err = l.API.PushMessage(flags.ID, input).Do()
		case CommandVideo:
			input := linebot.NewVideoMessage(flags.arg, flags.arg)
			l.Response, err = l.API.PushMessage(flags.ID, input).Do()
		}
		if err != nil {
			return err
		}
		printer.Printf(printer.SetNoneAsDefaultFormat(rootOutputFormat), l.Response)
		return err
	}

	var lineSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send message to LINE",
		Example: common.Examples(`# Send text to LINE chat
-s secret -t token --id GroupID -a 'Hello LINE'`, CommandLINE, CommandText),
		RunE: runE,
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
		RunE: runE,
	}

	var lineSubCmdPhoto = &cobra.Command{
		Use:   CommandPhoto,
		Short: "Send photo to LINE",
		Example: common.Examples(`# Send photo to LINE chat
-s secret -t token --id GroupID -a https://img.url`, CommandLINE, CommandPhoto),
		RunE: runE,
	}

	var lineSubCmdVideo = &cobra.Command{
		Use:   CommandVideo,
		Short: "Send video to LINE",
		Example: common.Examples(`# Send video to LINE chat
-s secret -t token --id GroupID -a https://video.url`, CommandLINE, CommandVideo),
		RunE: runE,
	}

	lineCmd.PersistentFlags().StringVarP(&flags.Secret, "secret", "s", "", common.Usage("Channel Secret"))
	lineCmd.PersistentFlags().StringVarP(&flags.Token, "token", "t", "", common.Usage("Channel Access Token"))
	lineCmd.PersistentFlags().StringVarP(&flags.arg, "arg", "a", "", common.Usage("Function Argument"))
	lineCmd.PersistentFlags().StringVar(&flags.ID, "id", "", common.Usage("UserID/GroupID/RoomID"))

	lineCmd.AddCommand(lineSubCmdID, lineSubCmdPhoto, lineSubCmdText, lineSubCmdVideo)
	return lineCmd
}

type LINE struct {
	API      *linebot.Client
	Response *linebot.BasicResponse
}

func (l *LINE) Init(secret, token string) error {
	var err error
	if secret == "" || token == "" {
		return common.ErrInvalidToken
	}
	l.API, err = linebot.New(secret, token)
	if l.API == nil {
		return common.ErrFailedInitial
	}
	return err
}

func (l *LINE) GetID() {
	var err error
	http.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		events, err := l.API.ParseRequest(r)
		if err != nil {
			printer.Error(err)
			os.Exit(1)
		}
		for i := range events {
			if events[i].Type == linebot.EventTypeMessage && events[i].Message.(*linebot.TextMessage).Text == CommandID {
				switch s := events[i].Source; s.Type {
				case linebot.EventSourceTypeGroup:
					printer.Printf(s.GroupID)
					os.Exit(0)
				case linebot.EventSourceTypeRoom:
					printer.Printf(s.RoomID)
					os.Exit(0)
				case linebot.EventSourceTypeUser:
					printer.Printf(s.UserID)
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
		printer.Printf("error starting server: %s", err)
		os.Exit(1)
	}
}
