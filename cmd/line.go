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
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/cobra"
)

var lineCmd = &cobra.Command{
	Use:   "line",
	Short: "Send message to LINE",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var lineSubCmdText = &cobra.Command{
	Use:   "text",
	Short: "Send message to LINE",
	Example: Examples(`# Send text to LINE chat
ops-cli line text -s secret -t token --id GroupID -a 'Hello LINE'`),
	Run: line.Run,
}

var lineSubCmdPhoto = &cobra.Command{
	Use:   "photo",
	Short: "Send photo to LINE",
	Example: Examples(`# Send photo to LINE chat
ops-cli line photo -s secret -t token --id GroupID -a https://img.url`),
	Run: line.Run,
}

var lineSubCmdVideo = &cobra.Command{
	Use:   "video",
	Short: "Send video to LINE",
	Example: Examples(`# Send video to LINE chat
ops-cli line video -s secret -t token --id GroupID -a https://video.url`),
	Run: line.Run,
}

var line lineFlag

func init() {
	rootCmd.AddCommand(lineCmd)

	lineCmd.PersistentFlags().StringVarP(&line.secret, "secret", "s", "", "Channel Secret")
	lineCmd.PersistentFlags().StringVarP(&line.token, "token", "t", "", "Channel Access Token")
	lineCmd.PersistentFlags().StringVarP(&line.arg, "arg", "a", "", "Function Argument")
	lineCmd.PersistentFlags().StringVar(&line.id, "id", "", "UserID/GroupID/RoomID")

	lineCmd.AddCommand(lineSubCmdText)
	lineCmd.AddCommand(lineSubCmdPhoto)
	lineCmd.AddCommand(lineSubCmdVideo)
}

type lineFlag struct {
	secret string
	token  string
	id     string
	arg    string

	api  *linebot.Client
	resp *linebot.BasicResponse
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

func (l *lineFlag) Run(cmd *cobra.Command, _ []string) {
	if err := l.Init(); err != nil {
		log.Println(err)
		return
	}
	var err error
	switch cmd.Name() {
	case ImTypeText:
		input := linebot.NewTextMessage(l.arg)
		l.resp, err = l.api.PushMessage(l.id, input).Do()
	case ImTypePhoto:
		input := linebot.NewImageMessage(l.arg, l.arg)
		l.resp, err = l.api.PushMessage(l.id, input).Do()
	case ImTypeVideo:
		input := linebot.NewVideoMessage(l.arg, l.arg)
		l.resp, err = l.api.PushMessage(l.id, input).Do()
	}
	if err != nil {
		log.Println(err)
	}
}

func (l lineFlag) JSON() { PrintJSON(l.resp) }

func (l lineFlag) YAML() { PrintYAML(l.resp) }

func (l lineFlag) String() {
	fmt.Printf(`%v`, l.resp)
}
