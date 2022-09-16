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
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func init() {
	var discordFlag DiscordFlag
	var discordCmd = &cobra.Command{
		Use:   CommandDiscord,
		Short: "Send message to Discord",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	var discordSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Short: "Send file to Discord",
		RunE:  discordFlag.RunE,
	}

	var discordSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send text to Discord",
		RunE:  discordFlag.RunE,
	}

	var discordSubCmdTextTS = &cobra.Command{
		Use:   CommandText + "TS",
		Short: "Send text to speech to Discord",
		RunE:  discordFlag.RunE,
	}
	rootCmd.AddCommand(discordCmd)

	discordCmd.PersistentFlags().StringVarP(&discordFlag.Token, "token", "t", "", common.Usage("Token"))
	discordCmd.PersistentFlags().StringVarP(&discordFlag.Channel, "channel-id", "c", "", common.Usage("Channel ID"))
	discordCmd.PersistentFlags().StringVarP(&discordFlag.arg, "arg", "a", "", common.Usage("Input argument"))

	discordCmd.AddCommand(discordSubCmdFile, discordSubCmdText, discordSubCmdTextTS)
}

type DiscordFlag struct {
	Token   string `json:"token"`
	Channel string `json:"channel_id"`
	arg     string

	api  *discordgo.Session
	resp *discordgo.Message
}

func (d *DiscordFlag) RunE(cmd *cobra.Command, args []string) error {
	if d.arg == "" {
		return common.ErrInvalidFlag
	}
	var err error
	if err = d.Init(); err != nil {
		return err
	}
	switch cmd.Name() {
	case CommandFile:
		err = d.File()
	case CommandText:
		err = d.Text()
	case CommandText + "TS":
		err = d.TextTTS()
	}
	if err != nil {
		return err
	}
	OutputDefaultNone(d.resp)
	return err
}

func (d *DiscordFlag) Init() error {
	var err error
	if rootConfig != "" {
		v := common.Config(rootConfig, strings.ToLower(CommandDiscord))
		err = Encoder.JSONMarshaler(v, d)
		if err != nil {
			return err
		}
	}
	if d.Token == "" {
		return common.ErrInvalidToken
	}

	d.api, err = discordgo.New("Bot " + d.Token)
	if d.api == nil {
		return common.ErrFailedInitial
	}
	return err
}

func (d *DiscordFlag) File() error {
	var err error
	f, err := os.Open(d.arg)
	if err != nil {
		return err
	}
	filename := filepath.Base(d.arg)
	defer f.Close()
	d.resp, err = d.api.ChannelFileSend(d.Channel, filename, f)
	return err
}

func (d *DiscordFlag) Text() error {
	var err error
	d.resp, err = d.api.ChannelMessageSend(d.Channel, d.arg)
	return err
}

func (d *DiscordFlag) TextTTS() error {
	var err error
	d.resp, err = d.api.ChannelMessageSendTTS(d.Channel, d.arg)
	return err
}
