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
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

var telegramCmd = &cobra.Command{
	Use:   "telegram [function]",
	Short: "Send message to telegram",
	Run: func(cmd *cobra.Command, args []string) {
		Config("telegram")
		if t.token == "" {
			_ = cmd.Help()
			return
		}
		telegramAPI := t.Init()
		if telegramAPI == nil || len(args) != 1 {
			_ = cmd.Help()
			return
		}
		if args[0] == "" {
			_ = cmd.Help()
			return
		}
		switch strings.ToLower(args[0]) {
		case "audio":
			t.Audio(telegramAPI)
		case "doc", "file":
			t.Doc(telegramAPI)
		case "msg", "message":
			t.Msg(telegramAPI)
		case "photo":
			t.Photo(telegramAPI)
		}
	},
	Example: Examples(`# Send message
ops-cli telegram msg -t bot_token -c chat_id -a 'Hello word'

# Send file
ops-cli telegram file -t bot_token -c chat_id -a '~/readme.md'

# Send photo
ops-cli telegram photo -t bot_token -c chat_id -a 'https://zh.wikipedia.org/wiki/File:Google_Chrome_icon_(February_2022).svg'
ops-cli telegram photo -t bot_token -c chat_id -a '~/photo/cat.png'`),
}

var t telegramFlag

func init() {
	rootCmd.AddCommand(telegramCmd)

	telegramCmd.Flags().StringVarP(&t.token, "token", "t", "", "Bot token (required)")
	telegramCmd.Flags().Int64VarP(&t.chat, "chat-id", "c", 0, "Chat ID")
	telegramCmd.Flags().StringVarP(&t.arg, "arg", "a", "", "Input argument")
	telegramCmd.Flags().StringVarP(&t.caption, "caption", "", "", "Add caption for document of photo")
}

type telegramFlag struct {
	/* Bind flags */
	token   string
	chat    int64
	arg     string
	caption string

	resp tg.Message
}

func (t telegramFlag) parseFile(s string) tg.RequestFileData {
	switch {
	case ValidFile(s):
		return tg.FilePath(s)
	case ValidURL(s):
		return tg.FileURL(s)
	}
	return nil
}

func (t telegramFlag) Init() *tg.BotAPI {
	api, err := tg.NewBotAPI(t.token)
	if err != nil {
		log.Println(err)
		return nil
	}
	return api
}

func (t telegramFlag) Animation(api *tg.BotAPI) {
	input := tg.NewAnimation(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Audio(api *tg.BotAPI) {
	input := tg.NewAudio(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) ChatDescription(api *tg.BotAPI) {
	input := tg.NewChatDescription(t.chat, t.arg)
	t.send(api, input)
}

func (t telegramFlag) ChatPhoto(api *tg.BotAPI) {
	input := tg.NewChatPhoto(t.chat, t.parseFile(t.arg))
	t.send(api, input)
}

func (t telegramFlag) ChatTitle(api *tg.BotAPI) {
	input := tg.NewChatTitle(t.chat, t.arg)
	t.send(api, input)
}

func (t telegramFlag) Dice(api *tg.BotAPI) {
	input := tg.NewDice(t.chat)
	t.send(api, input)
}

func (t telegramFlag) Doc(api *tg.BotAPI) {
	input := tg.NewDocument(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Msg(api *tg.BotAPI) {
	input := tg.NewMessage(t.chat, t.arg)
	input.ParseMode = tg.ModeMarkdownV2
	input.DisableWebPagePreview = true
	t.send(api, input)
}

func (t telegramFlag) Photo(api *tg.BotAPI) {
	input := tg.NewPhoto(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Video(api *tg.BotAPI) {
	input := tg.NewVideo(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Voice(api *tg.BotAPI) {
	input := tg.NewVoice(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) send(api *tg.BotAPI, c tg.Chattable) {
	var err error
	t.resp, err = api.Send(c)
	if err != nil {
		log.Println(err)
		return
	}
	if rootOutputJSON {
		t.JSON()
	} else if rootOutputYAML {
		t.YAML()
	}
}

func (t telegramFlag) JSON() { PrintJSON(t.resp) }

func (t telegramFlag) YAML() { PrintYAML(t.resp) }

func (t telegramFlag) String() {
	fmt.Printf(`%v`, t.resp)
}
