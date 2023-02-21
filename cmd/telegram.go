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
	"os"
	"strconv"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initTelegram() *cobra.Command {
	var flags struct {
		/* Bind flags */
		Token   string `json:"token"`
		ChatID  string `json:"chat_id"`
		Chat    int64
		arg     string
		caption string
	}
	var telegramCmd = &cobra.Command{
		GroupID: getGroupID(CommandTelegram),
		Use:     CommandTelegram,
		Short:   "Send message to Telegram",
		RunE:    func(cmd *cobra.Command, _ []string) error { return cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	run := func(cmd *cobra.Command, _ []string) {
		if flags.arg == "" {
			logger.Error(common.ErrInvalidFlag.Error(), common.DefaultField(flags.arg))
			return
		}
		var err error
		if rootConfig != "" {
			if err = ReadConfig(CommandTelegram, &flags); err != nil {
				logger.Error(err.Error())
				return
			}
			i, err := strconv.ParseInt(flags.ChatID, 10, 64)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			flags.Chat = i
		}
		var t Telegram
		if err = t.Init(flags.Token); err != nil {
			logger.Error(err.Error())
			return
		}
		if cmd.Name() != CommandText {
			if !common.IsFile(flags.arg) && !common.IsURL(flags.arg) {
				logger.Error(common.ErrInvalidFile.Error())
				return
			}
		}
		switch cmd.Name() {
		case CommandAudio:
			err = t.Audio(flags.Chat, flags.arg, flags.caption)
		case CommandFile:
			err = t.File(flags.Chat, flags.arg, flags.caption)
		case CommandPhoto:
			err = t.Photo(flags.Chat, flags.arg, flags.caption)
		case CommandText:
			err = t.Text(flags.Chat, flags.arg)
		case CommandVideo:
			err = t.Video(flags.Chat, flags.arg, flags.caption)
		case CommandVoice:
			err = t.Voice(flags.Chat, flags.arg, flags.caption)
		}
		if err != nil {
			logger.Error(err.Error())
			return
		}
		printer.Printf(printer.SetNoneAsDefaultFormat(rootOutputFormat), t.Response)
	}

	var telegramSubCmdAudio = &cobra.Command{
		Use:   CommandAudio,
		Short: "Send audio file to Telegram",
		Run:   run,
	}

	var telegramSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Short: "Send file to Telegram",
		Run:   run,
		Example: common.Examples(`# Send file
-t bot_token -c chat_id -a '~/readme.md'`, CommandTelegram, CommandFile),
	}

	var telegramSubCmdID = &cobra.Command{
		Use:   CommandID,
		Short: "Get chat ID",
		Run: func(_ *cobra.Command, _ []string) {
			var err error
			if rootConfig != "" {
				if err = ReadConfig(CommandTelegram, &flags); err != nil {
					logger.Error(err.Error())
					return
				}
			}
			var t Telegram
			if err = t.Init(flags.Token); err != nil {
				logger.Error(err.Error())
				return
			}
			t.GetUpdate()
		},
		Example: common.Examples(`# Execute the command and enter 'id' in the chat to get the chat id.
--config ~/.config.toml`, CommandTelegram, CommandID),
	}

	var telegramSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send text to Telegram",
		Run:   run,
		Example: common.Examples(`# Send message
-t bot_token -c chat_id -a 'Hello word'`, CommandTelegram, CommandText),
	}

	var telegramSubCmdPhoto = &cobra.Command{
		Use:   CommandPhoto,
		Short: "Send photo to Telegram",
		Run:   run,
		Example: common.Examples(`# Send photo
-t bot_token -c chat_id -a 'https://zh.wikipedia.org/wiki/File:Google_Chrome_icon_(February_2022).svg'
-t bot_token -c chat_id -a '~/photo/cat.png'`, CommandTelegram, CommandPhoto),
	}

	var telegramSubCmdVideo = &cobra.Command{
		Use:   CommandVideo,
		Short: "Send video file to Telegram",
		Run:   run,
	}

	var telegramSubCmdVoice = &cobra.Command{
		Use:   CommandVoice,
		Short: "Send voice file to Telegram",
		Run:   run,
	}
	telegramCmd.PersistentFlags().StringVarP(&flags.Token, "token", "t", "", common.Usage("Bot token (required)"))
	telegramCmd.PersistentFlags().Int64VarP(&flags.Chat, "chat-id", "c", 0, common.Usage("Chat ID"))
	telegramCmd.PersistentFlags().StringVarP(&flags.arg, "arg", "a", "", common.Usage("Input argument"))
	telegramCmd.PersistentFlags().StringVarP(&flags.caption, "caption", "", "", common.Usage("Add caption for file"))

	telegramCmd.AddCommand(telegramSubCmdAudio)
	telegramCmd.AddCommand(telegramSubCmdFile)
	telegramCmd.AddCommand(telegramSubCmdID)
	telegramCmd.AddCommand(telegramSubCmdText)
	telegramCmd.AddCommand(telegramSubCmdPhoto)
	telegramCmd.AddCommand(telegramSubCmdVideo)
	telegramCmd.AddCommand(telegramSubCmdVoice)
	return telegramCmd
}

type Telegram struct {
	API      *tgBot.BotAPI
	Response tgBot.Message
}

func (t *Telegram) Init(token string) error {
	var err error
	if token == "" {
		logger.Debug(common.ErrInvalidToken.Error())
		return common.ErrInvalidToken
	}
	t.API, err = tgBot.NewBotAPI(token)
	if t.API == nil {
		logger.Debug(common.ErrFailedInitial.Error())
		return common.ErrFailedInitial
	}
	return err
}

func (t Telegram) Animation(chat int64, arg, caption string) error {
	input := tgBot.NewAnimation(chat, t.parseFile(arg))
	input.Caption = caption
	return t.send(input)
}

func (t *Telegram) Audio(chat int64, arg, caption string) error {
	input := tgBot.NewAudio(chat, t.parseFile(arg))
	input.Caption = caption
	return t.send(input)
}

func (t Telegram) ChatDescription(chat int64, arg string) error {
	input := tgBot.NewChatDescription(chat, arg)
	return t.send(input)
}

func (t Telegram) ChatPhoto(chat int64, arg string) error {
	input := tgBot.NewChatPhoto(chat, t.parseFile(arg))
	return t.send(input)
}

func (t Telegram) ChatTitle(chat int64, arg string) error {
	input := tgBot.NewChatTitle(chat, arg)
	return t.send(input)
}

func (t Telegram) Dice(chat int64) error {
	input := tgBot.NewDice(chat)
	return t.send(input)
}

func (t *Telegram) File(chat int64, arg, caption string) error {
	input := tgBot.NewDocument(chat, t.parseFile(arg))
	input.Caption = caption
	return t.send(input)
}

func (t *Telegram) Text(chat int64, arg string) error {
	input := tgBot.NewMessage(chat, arg)
	input.ParseMode = tgBot.ModeMarkdownV2
	input.DisableWebPagePreview = true
	return t.send(input)
}

func (t *Telegram) Photo(chat int64, arg, caption string) error {
	input := tgBot.NewPhoto(chat, t.parseFile(arg))
	input.Caption = caption
	return t.send(input)
}

func (t *Telegram) Video(chat int64, arg, caption string) error {
	input := tgBot.NewVideo(chat, t.parseFile(arg))
	input.Caption = caption
	return t.send(input)
}

func (t *Telegram) Voice(chat int64, arg, caption string) error {
	input := tgBot.NewVoice(chat, t.parseFile(arg))
	input.Caption = caption
	return t.send(input)
}

func (t *Telegram) GetUpdate() {
	u := tgBot.NewUpdate(0)
	u.Timeout = 60
	updates := t.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == CommandID {
				printer.Printf("%d", update.Message.Chat.ID)
				os.Exit(0)
			}
		}
	}
}

func (t *Telegram) parseFile(s string) tgBot.RequestFileData {
	switch {
	case common.IsFile(s):
		return tgBot.FilePath(s)
	case common.IsURL(s):
		return tgBot.FileURL(s)
	}
	return nil
}

func (t *Telegram) send(c tgBot.Chattable) error {
	var err error
	t.Response, err = t.API.Send(c)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(c))
	}
	return err
}
