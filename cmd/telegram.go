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

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

var telegramCmd = &cobra.Command{
	Use:   "telegram [function]",
	Short: "Send message to telegram",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		if args[0] == "" {
			_ = cmd.Help()
			return
		}
		Config("telegram")
		if tg.token == "" {
			_ = cmd.Help()
			return
		}
		telegramAPI := tg.Init()
		if telegramAPI == nil {
			_ = cmd.Help()
			return
		}
		switch strings.ToLower(args[0]) {
		case "audio":
			tg.Audio(telegramAPI)
		case "doc", "file":
			tg.Doc(telegramAPI)
		case "msg", "message":
			tg.Msg(telegramAPI)
		case "photo":
			tg.Photo(telegramAPI)
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

var tg telegramFlag

func init() {
	rootCmd.AddCommand(telegramCmd)

	telegramCmd.Flags().StringVarP(&tg.token, "token", "t", "", "Bot token (required)")
	telegramCmd.Flags().Int64VarP(&tg.chat, "chat-id", "c", 0, "Chat ID")
	telegramCmd.Flags().StringVarP(&tg.arg, "arg", "a", "", "Input argument")
	telegramCmd.Flags().StringVarP(&tg.caption, "caption", "", "", "Add caption for document of photo")
}

type telegramFlag struct {
	/* Bind flags */
	token   string
	chat    int64
	arg     string
	caption string

	resp tgBot.Message
}

func (t telegramFlag) parseFile(s string) tgBot.RequestFileData {
	switch {
	case ValidFile(s):
		return tgBot.FilePath(s)
	case ValidURL(s):
		return tgBot.FileURL(s)
	}
	return nil
}

func (t telegramFlag) Init() *tgBot.BotAPI {
	api, err := tgBot.NewBotAPI(t.token)
	if err != nil {
		log.Println(err)
		return nil
	}
	return api
}

func (t telegramFlag) Animation(api *tgBot.BotAPI) {
	input := tgBot.NewAnimation(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Audio(api *tgBot.BotAPI) {
	input := tgBot.NewAudio(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) ChatDescription(api *tgBot.BotAPI) {
	input := tgBot.NewChatDescription(t.chat, t.arg)
	t.send(api, input)
}

func (t telegramFlag) ChatPhoto(api *tgBot.BotAPI) {
	input := tgBot.NewChatPhoto(t.chat, t.parseFile(t.arg))
	t.send(api, input)
}

func (t telegramFlag) ChatTitle(api *tgBot.BotAPI) {
	input := tgBot.NewChatTitle(t.chat, t.arg)
	t.send(api, input)
}

func (t telegramFlag) Dice(api *tgBot.BotAPI) {
	input := tgBot.NewDice(t.chat)
	t.send(api, input)
}

func (t telegramFlag) Doc(api *tgBot.BotAPI) {
	input := tgBot.NewDocument(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Msg(api *tgBot.BotAPI) {
	input := tgBot.NewMessage(t.chat, t.arg)
	input.ParseMode = tgBot.ModeMarkdownV2
	input.DisableWebPagePreview = true
	t.send(api, input)
}

func (t telegramFlag) Photo(api *tgBot.BotAPI) {
	input := tgBot.NewPhoto(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Video(api *tgBot.BotAPI) {
	input := tgBot.NewVideo(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) Voice(api *tgBot.BotAPI) {
	input := tgBot.NewVoice(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(api, input)
}

func (t telegramFlag) send(api *tgBot.BotAPI, c tgBot.Chattable) {
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
