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
	"os"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

var telegramCmd = &cobra.Command{
	Use:   "Telegram",
	Short: "Send message to Telegram",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var telegramSubCmdAudio = &cobra.Command{
	Use:   "audio",
	Short: "Send audio file to Telegram",
	Run:   tg.Run,
}

var telegramSubCmdFile = &cobra.Command{
	Use:   "file",
	Short: "Send file to Telegram",
	Run:   tg.Run,
	Example: Examples(`# Send file
ops-cli Telegram file -t bot_token -c chat_id -a '~/readme.md'`),
}

var telegramSubCmdID = &cobra.Command{
	Use:   "id",
	Short: "Get chat ID",
	Run:   tg.Run,
	Example: Examples(`# Execute the command and enter 'id' in the chat to get the chat id.
ops-cli Telegram id --config ~/.config.toml`),
}

var telegramSubCmdText = &cobra.Command{
	Use:   "text",
	Short: "Send text to Telegram",
	Run:   tg.Run,
	Example: Examples(`# Send message
ops-cli Telegram text -t bot_token -c chat_id -a 'Hello word'`),
}

var telegramSubCmdPhoto = &cobra.Command{
	Use:   "photo",
	Short: "Send photo to Telegram",
	Run:   tg.Run,
	Example: Examples(`# Send photo
ops-cli Telegram photo -t bot_token -c chat_id -a 'https://zh.wikipedia.org/wiki/File:Google_Chrome_icon_(February_2022).svg'
ops-cli Telegram photo -t bot_token -c chat_id -a '~/photo/cat.png'`),
}

var telegramSubCmdVideo = &cobra.Command{
	Use:   "video",
	Short: "Send video file to Telegram",
	Run:   tg.Run,
}

var telegramSubCmdVoice = &cobra.Command{
	Use:   "voice",
	Short: "Send voice file to Telegram",
	Run:   tg.Run,
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
	telegramCmd.AddCommand(telegramSubCmdText)
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
	var err error
	if tg.token == "" && rootConfig != "" {
		if err = Config(ConfigBlockTelegram); err != nil {
			return err
		}
	}
	if tg.token == "" {
		return ErrTokenNotFound
	}
	t.api, err = tgBot.NewBotAPI(t.token)
	if t.api == nil {
		return ErrInitialFailed
	}
	return err
}

func (t telegramFlag) Animation() error {
	input := tgBot.NewAnimation(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *telegramFlag) Audio() error {
	input := tgBot.NewAudio(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t telegramFlag) ChatDescription() error {
	input := tgBot.NewChatDescription(t.chat, t.arg)
	return t.send(input)
}

func (t telegramFlag) ChatPhoto() error {
	input := tgBot.NewChatPhoto(t.chat, t.parseFile(t.arg))
	return t.send(input)
}

func (t telegramFlag) ChatTitle() error {
	input := tgBot.NewChatTitle(t.chat, t.arg)
	return t.send(input)
}

func (t telegramFlag) Dice() error {
	input := tgBot.NewDice(t.chat)
	return t.send(input)
}

func (t *telegramFlag) File() error {
	input := tgBot.NewDocument(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *telegramFlag) Text() error {
	input := tgBot.NewMessage(t.chat, t.arg)
	input.ParseMode = tgBot.ModeMarkdownV2
	input.DisableWebPagePreview = true
	return t.send(input)
}

func (t *telegramFlag) Photo() error {
	input := tgBot.NewPhoto(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *telegramFlag) Video() error {
	input := tgBot.NewVideo(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *telegramFlag) Voice() error {
	input := tgBot.NewVoice(t.chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *telegramFlag) GetUpdate() {
	u := tgBot.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == ImTypeID {
				fmt.Println(update.Message.Chat.ID)
				break
			}
		}
	}
}

func (t *telegramFlag) Run(cmd *cobra.Command, _ []string) {
	if t.arg == "" {
		log.Println(ErrArgNotFound)
		os.Exit(1)
	}
	var err error
	if err = t.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	switch cmd.Name() {
	case ImTypeAudio:
		err = t.Audio()
	case ImTypeFile:
		err = t.File()
	case ImTypeID:
		t.GetUpdate()
	case ImTypePhoto:
		err = t.Photo()
	case ImTypeText:
		err = t.Text()
	case ImTypeVideo:
		err = t.Video()
	case ImTypeVoice:
		err = t.Voice()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	OutputDefaultNone(t.resp)
}

func (t *telegramFlag) parseFile(s string) tgBot.RequestFileData {
	switch {
	case ValidFile(s):
		return tgBot.FilePath(s)
	case ValidURL(s):
		return tgBot.FileURL(s)
	}
	return nil
}

func (t *telegramFlag) send(c tgBot.Chattable) error {
	var err error
	t.resp, err = t.api.Send(c)
	return err
}
