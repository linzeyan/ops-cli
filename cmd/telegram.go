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
	"fmt"
	"log"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

var telegramCmd = &cobra.Command{
	Use:     "telegram",
	Aliases: []string{"tg"},
	Short:   "Send message to telegram",
	Run:     func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var telegramSubCmdAudio = &cobra.Command{
	Use:   "audio",
	Short: "Send audio file to telegram",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.Audio()
	},
}

var telegramSubCmdFile = &cobra.Command{
	Use:     "file",
	Aliases: []string{imTypeDoc, imTypeDocument},
	Short:   "Send file to telegram",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.File()
	},
	Example: Examples(`# Send file
ops-cli telegram file -t bot_token -c chat_id -a '~/readme.md'`),
}

var telegramSubCmdID = &cobra.Command{
	Use:   "id",
	Short: "Get chat ID",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.GetUpdate()
	},
	Example: Examples(`# Execute the command and enter 'id' in the chat to get the chat id.
ops-cli telegram id --config ~/.config.toml`),
}

var telegramSubCmdMsg = &cobra.Command{
	Use:     "msg",
	Aliases: []string{imTypeMessage},
	Short:   "Send message to telegram",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.Msg()
	},
	Example: Examples(`# Send message
ops-cli telegram msg -t bot_token -c chat_id -a 'Hello word'`),
}

var telegramSubCmdPhoto = &cobra.Command{
	Use:   "photo",
	Short: "Send photo to telegram",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.Photo()
	},
	Example: Examples(`# Send photo
ops-cli telegram photo -t bot_token -c chat_id -a 'https://zh.wikipedia.org/wiki/File:Google_Chrome_icon_(February_2022).svg'
ops-cli telegram photo -t bot_token -c chat_id -a '~/photo/cat.png'`),
}

var telegramSubCmdVideo = &cobra.Command{
	Use:   "video",
	Short: "Send video file to telegram",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.Video()
	},
}

var telegramSubCmdVoice = &cobra.Command{
	Use:   "voice",
	Short: "Send voice file to telegram",
	Run: func(_ *cobra.Command, _ []string) {
		if err := tg.Init(); err != nil {
			log.Println(err)
			return
		}
		tg.Voice()
	},
}

var tg telegramFlag

func init() {
	rootCmd.AddCommand(telegramCmd)

	telegramCmd.PersistentFlags().StringVarP(&tg.token, "token", "t", "", "Bot token (required)")
	telegramCmd.PersistentFlags().Int64VarP(&tg.chat, "chat-id", "c", 0, "Chat ID")
	telegramCmd.PersistentFlags().StringVarP(&tg.arg, "arg", "a", "", "Input argument")
	telegramCmd.PersistentFlags().StringVarP(&tg.caption, "caption", "", "", "Add caption for file")

	telegramCmd.AddCommand(telegramSubCmdAudio)
	telegramCmd.AddCommand(telegramSubCmdFile)
	telegramCmd.AddCommand(telegramSubCmdID)
	telegramCmd.AddCommand(telegramSubCmdMsg)
	telegramCmd.AddCommand(telegramSubCmdPhoto)
	telegramCmd.AddCommand(telegramSubCmdVideo)
	telegramCmd.AddCommand(telegramSubCmdVoice)
}

type telegramFlag struct {
	/* Bind flags */
	token   string
	chat    int64
	arg     string
	caption string

	api  *tgBot.BotAPI
	resp tgBot.Message
}

func (t *telegramFlag) Init() error {
	Config(configTelegram)
	if tg.token == "" {
		return errors.New("token is empty")
	}
	var err error
	t.api, err = tgBot.NewBotAPI(t.token)
	if err != nil {
		return err
	}
	return nil
}

func (t telegramFlag) Animation() {
	input := tgBot.NewAnimation(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(input)
}

func (t telegramFlag) Audio() {
	input := tgBot.NewAudio(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(input)
}

func (t telegramFlag) ChatDescription() {
	input := tgBot.NewChatDescription(t.chat, t.arg)
	t.send(input)
}

func (t telegramFlag) ChatPhoto() {
	input := tgBot.NewChatPhoto(t.chat, t.parseFile(t.arg))
	t.send(input)
}

func (t telegramFlag) ChatTitle() {
	input := tgBot.NewChatTitle(t.chat, t.arg)
	t.send(input)
}

func (t telegramFlag) Dice() {
	input := tgBot.NewDice(t.chat)
	t.send(input)
}

func (t telegramFlag) File() {
	input := tgBot.NewDocument(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(input)
}

func (t telegramFlag) Msg() {
	input := tgBot.NewMessage(t.chat, t.arg)
	input.ParseMode = tgBot.ModeMarkdownV2
	input.DisableWebPagePreview = true
	t.send(input)
}

func (t telegramFlag) Photo() {
	input := tgBot.NewPhoto(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(input)
}

func (t telegramFlag) Video() {
	input := tgBot.NewVideo(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(input)
}

func (t telegramFlag) Voice() {
	input := tgBot.NewVoice(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	t.send(input)
}

func (t telegramFlag) GetUpdate() {
	u := tgBot.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == imTypeID {
				fmt.Println(update.Message.Chat.ID)
				break
			}
		}
	}
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

func (t *telegramFlag) send(c tgBot.Chattable) {
	var err error
	t.resp, err = t.api.Send(c)
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
