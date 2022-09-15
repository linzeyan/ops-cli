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
	"strconv"
	"strings"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var telegramFlag TelegramFlag
	var telegramCmd = &cobra.Command{
		Use:   CommandTelegram,
		Short: "Send message to Telegram",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	var telegramSubCmdAudio = &cobra.Command{
		Use:   CommandAudio,
		Short: "Send audio file to Telegram",
		RunE:  telegramFlag.RunE,
	}

	var telegramSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Short: "Send file to Telegram",
		RunE:  telegramFlag.RunE,
		Example: common.Examples(`# Send file
-t bot_token -c chat_id -a '~/readme.md'`, CommandTelegram, CommandFile),
	}

	var telegramSubCmdID = &cobra.Command{
		Use:   CommandID,
		Short: "Get chat ID",
		RunE: func(_ *cobra.Command, _ []string) error {
			var err error
			if err = telegramFlag.Init(); err != nil {
				return err
			}
			telegramFlag.GetUpdate()
			return err
		},
		Example: common.Examples(`# Execute the command and enter 'id' in the chat to get the chat id.
--config ~/.config.toml`, CommandTelegram, CommandID),
	}

	var telegramSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send text to Telegram",
		RunE:  telegramFlag.RunE,
		Example: common.Examples(`# Send message
-t bot_token -c chat_id -a 'Hello word'`, CommandTelegram, CommandText),
	}

	var telegramSubCmdPhoto = &cobra.Command{
		Use:   CommandPhoto,
		Short: "Send photo to Telegram",
		RunE:  telegramFlag.RunE,
		Example: common.Examples(`# Send photo
-t bot_token -c chat_id -a 'https://zh.wikipedia.org/wiki/File:Google_Chrome_icon_(February_2022).svg'
-t bot_token -c chat_id -a '~/photo/cat.png'`, CommandTelegram, CommandPhoto),
	}

	var telegramSubCmdVideo = &cobra.Command{
		Use:   CommandVideo,
		Short: "Send video file to Telegram",
		RunE:  telegramFlag.RunE,
	}

	var telegramSubCmdVoice = &cobra.Command{
		Use:   CommandVoice,
		Short: "Send voice file to Telegram",
		RunE:  telegramFlag.RunE,
	}
	rootCmd.AddCommand(telegramCmd)

	telegramCmd.PersistentFlags().StringVarP(&telegramFlag.Token, "token", "t", "", common.Usage("Bot token (required)"))
	telegramCmd.PersistentFlags().Int64VarP(&telegramFlag.Chat, "chat-id", "c", 0, common.Usage("Chat ID"))
	telegramCmd.PersistentFlags().StringVarP(&telegramFlag.arg, "arg", "a", "", common.Usage("Input argument"))
	telegramCmd.PersistentFlags().StringVarP(&telegramFlag.caption, "caption", "", "", common.Usage("Add caption for file"))

	telegramCmd.AddCommand(telegramSubCmdAudio)
	telegramCmd.AddCommand(telegramSubCmdFile)
	telegramCmd.AddCommand(telegramSubCmdID)
	telegramCmd.AddCommand(telegramSubCmdText)
	telegramCmd.AddCommand(telegramSubCmdPhoto)
	telegramCmd.AddCommand(telegramSubCmdVideo)
	telegramCmd.AddCommand(telegramSubCmdVoice)
}

type TelegramFlag struct {
	/* Bind flags */
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
	Chat    int64
	arg     string
	caption string

	api  *tgBot.BotAPI
	resp tgBot.Message
}

func (t *TelegramFlag) Init() error {
	var err error
	if t.Token == "" && rootConfig != "" {
		v := common.Config(rootConfig, strings.ToLower(CommandTelegram))
		err = Encoder.JSONMarshaler(v, t)
		if err != nil {
			return err
		}
		i, err := strconv.ParseInt(t.ChatID, 10, 64)
		if err != nil {
			return err
		}
		t.Chat = i
	}
	if t.Token == "" {
		return common.ErrInvalidToken
	}
	t.api, err = tgBot.NewBotAPI(t.Token)
	if t.api == nil {
		return common.ErrFailedInitial
	}
	return err
}

func (t TelegramFlag) Animation() error {
	input := tgBot.NewAnimation(t.Chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *TelegramFlag) Audio() error {
	input := tgBot.NewAudio(t.Chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t TelegramFlag) ChatDescription() error {
	input := tgBot.NewChatDescription(t.Chat, t.arg)
	return t.send(input)
}

func (t TelegramFlag) ChatPhoto() error {
	input := tgBot.NewChatPhoto(t.Chat, t.parseFile(t.arg))
	return t.send(input)
}

func (t TelegramFlag) ChatTitle() error {
	input := tgBot.NewChatTitle(t.Chat, t.arg)
	return t.send(input)
}

func (t TelegramFlag) Dice() error {
	input := tgBot.NewDice(t.Chat)
	return t.send(input)
}

func (t *TelegramFlag) File() error {
	input := tgBot.NewDocument(t.Chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *TelegramFlag) Text() error {
	input := tgBot.NewMessage(t.Chat, t.arg)
	input.ParseMode = tgBot.ModeMarkdownV2
	input.DisableWebPagePreview = true
	return t.send(input)
}

func (t *TelegramFlag) Photo() error {
	input := tgBot.NewPhoto(t.Chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *TelegramFlag) Video() error {
	input := tgBot.NewVideo(t.Chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *TelegramFlag) Voice() error {
	input := tgBot.NewVoice(t.Chat, t.parseFile(t.arg))
	input.Caption = t.caption
	return t.send(input)
}

func (t *TelegramFlag) GetUpdate() {
	u := tgBot.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == CommandID {
				PrintString(update.Message.Chat.ID)
				break
			}
		}
	}
}

func (t *TelegramFlag) RunE(cmd *cobra.Command, _ []string) error {
	if t.arg == "" {
		return common.ErrInvalidFlag
	}
	var err error
	if err = t.Init(); err != nil {
		return err
	}
	switch cmd.Name() {
	case CommandAudio:
		err = t.Audio()
	case CommandFile:
		err = t.File()
	case CommandPhoto:
		err = t.Photo()
	case CommandText:
		err = t.Text()
	case CommandVideo:
		err = t.Video()
	case CommandVoice:
		err = t.Voice()
	}
	if err != nil {
		return err
	}
	OutputDefaultNone(t.resp)
	return err
}

func (t *TelegramFlag) parseFile(s string) tgBot.RequestFileData {
	switch {
	case validator.ValidFile(s):
		return tgBot.FilePath(s)
	case validator.ValidURL(s):
		return tgBot.FileURL(s)
	}
	return nil
}

func (t *TelegramFlag) send(c tgBot.Chattable) error {
	var err error
	t.resp, err = t.api.Send(c)
	return err
}
