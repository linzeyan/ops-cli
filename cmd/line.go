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

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/cobra"
)

var lineCmd = &cobra.Command{
	Use:   "line [function]",
	Short: "Send message to LINE",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var lineSubCmdMsg = &cobra.Command{
	Use:     "msg",
	Aliases: []string{imTypeMessage},
	Example: Examples(`# Send message to LINE chat
ops-cli line msg -s secret -t token --id GroupID -a 'Hello LINE'`),
	Run: func(cmd *cobra.Command, _ []string) {
		if err := line.Init(); err != nil {
			log.Println(err)
			_ = cmd.Help()
			return
		}
		line.Msg()
	},
}

var lineSubCmdPhoto = &cobra.Command{
	Use: "photo",
	Example: Examples(`# Send photo to LINE chat
ops-cli line photo -s secret -t token --id GroupID -a https://img.url`),
	Run: func(cmd *cobra.Command, _ []string) {
		if err := line.Init(); err != nil {
			log.Println(err)
			_ = cmd.Help()
			return
		}
		line.Photo()
	},
}

var lineSubCmdVideo = &cobra.Command{
	Use: "video",
	Example: Examples(`# Send video to LINE chat
ops-cli line video -s secret -t token --id GroupID -a https://video.url`),
	Run: func(cmd *cobra.Command, _ []string) {
		if err := line.Init(); err != nil {
			log.Println(err)
			_ = cmd.Help()
			return
		}
		line.Video()
	},
}

var line lineFlag

func init() {
	rootCmd.AddCommand(lineCmd)

	lineCmd.PersistentFlags().StringVarP(&line.secret, "secret", "s", "", "Channel Secret")
	lineCmd.PersistentFlags().StringVarP(&line.token, "token", "t", "", "Channel Access Token")
	lineCmd.PersistentFlags().StringVarP(&line.arg, "arg", "a", "", "Function Argument")
	lineCmd.PersistentFlags().StringVar(&line.id, "id", "", "UserID/GroupID/RoomID")

	lineCmd.AddCommand(lineSubCmdMsg)
	lineCmd.AddCommand(lineSubCmdPhoto)
	lineCmd.AddCommand(lineSubCmdVideo)
}

type lineFlag struct {
	secret string
	token  string
	id     string
	arg    string

	api *linebot.Client
}

func (l *lineFlag) Init() error {
	Config(configLINE)
	var err error
	l.api, err = linebot.New(l.secret, l.token)
	if err != nil {
		return err
	}
	return nil
}

func (l lineFlag) Msg() {
	input := linebot.NewTextMessage(l.arg)
	_, err := l.api.PushMessage(l.id, input).Do()
	if err != nil {
		log.Println(err)
	}
}

func (l lineFlag) Photo() {
	input := linebot.NewImageMessage(l.arg, l.arg)
	_, err := l.api.PushMessage(l.id, input).Do()
	if err != nil {
		log.Println(err)
	}
}

func (l lineFlag) Video() {
	input := linebot.NewVideoMessage(l.arg, l.arg)
	_, err := l.api.PushMessage(l.id, input).Do()
	if err != nil {
		log.Println(err)
	}
}
