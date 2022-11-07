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
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initDiscord() *cobra.Command {
	var flags struct {
		Token   string `json:"token"`
		Channel string `json:"channel_id"`
		arg     string
	}
	var discordCmd = &cobra.Command{
		GroupID: groupings[CommandDiscord],
		Use:     CommandDiscord,
		Short:   "Send message to Discord",
		RunE:    func(cmd *cobra.Command, _ []string) error { return cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	run := func(cmd *cobra.Command, _ []string) {
		if flags.arg == "" {
			logger.Info(common.ErrInvalidFlag.Error())
			printer.Error(common.ErrInvalidFlag)
			return
		}
		var err error
		if rootConfig != "" {
			if err = ReadConfig(CommandDiscord, &flags); err != nil {
				logger.Info(err.Error())
				printer.Error(err)
				return
			}
		}
		var d Discord
		if err = d.Init(flags.Token); err != nil {
			logger.Info(err.Error())
			printer.Error(err)
			return
		}
		switch cmd.Name() {
		case CommandFile:
			err = d.File(flags.Channel, flags.arg)
		case CommandText:
			err = d.Text(flags.Channel, flags.arg)
		case CommandText + "TS":
			err = d.TextTTS(flags.Channel, flags.arg)
		}
		if err != nil {
			logger.Info(err.Error())
			printer.Error(err)
			return
		}
		printer.Printf(printer.SetNoneAsDefaultFormat(rootOutputFormat), d.Response)
	}

	var discordSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Short: "Send file to Discord",
		Run:   run,
	}

	var discordSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send text to Discord",
		Run:   run,
	}

	var discordSubCmdTextTS = &cobra.Command{
		Use:   CommandText + "TS",
		Short: "Send text to speech to Discord",
		Run:   run,
	}

	discordCmd.PersistentFlags().StringVarP(&flags.Token, "token", "t", "", common.Usage("Token"))
	discordCmd.PersistentFlags().StringVarP(&flags.Channel, "channel-id", "c", "", common.Usage("Channel ID"))
	discordCmd.PersistentFlags().StringVarP(&flags.arg, "arg", "a", "", common.Usage("Input argument"))

	discordCmd.AddCommand(discordSubCmdFile, discordSubCmdText, discordSubCmdTextTS)
	return discordCmd
}

type Discord struct {
	API      *discordgo.Session
	Response *discordgo.Message
}

func (d *Discord) Init(token string) error {
	var err error
	if token == "" {
		logger.Debug(common.ErrInvalidToken.Error(), common.DefaultField(token))
		return common.ErrInvalidToken
	}
	d.API, err = discordgo.New("Bot " + token)
	if d.API == nil {
		logger.Debug(common.ErrFailedInitial.Error())
		return common.ErrFailedInitial
	}
	return err
}

func (d *Discord) File(channel, arg string) error {
	var err error
	f, err := os.Open(arg)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(arg))
		return err
	}
	filename := filepath.Base(arg)
	defer f.Close()
	d.Response, err = d.API.ChannelFileSend(channel, filename, f)
	if err != nil {
		logger.Debug(err.Error(),
			common.NewField("channel", channel),
			common.NewField("filename", filename),
			common.NewField("f", f),
		)
	}
	return err
}

func (d *Discord) Text(channel, arg string) error {
	var err error
	d.Response, err = d.API.ChannelMessageSend(channel, arg)
	if err != nil {
		logger.Debug(err.Error(),
			common.NewField("channel", channel),
			common.DefaultField(arg),
		)
	}
	return err
}

func (d *Discord) TextTTS(channel, arg string) error {
	var err error
	d.Response, err = d.API.ChannelMessageSendTTS(channel, arg)
	if err != nil {
		logger.Debug(err.Error(),
			common.NewField("channel", channel),
			common.DefaultField(arg),
		)
	}
	return err
}
