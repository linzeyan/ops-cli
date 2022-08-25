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
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var discordCmd = &cobra.Command{
	Use:   common.CommandDiscord,
	Short: "Send message to Discord",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var discordSubCmdFile = &cobra.Command{
	Use:   common.SubCommandFile,
	Short: "Send file to Discord",
	Run:   discordCmdGlobalVar.Run,
}

var discordSubCmdText = &cobra.Command{
	Use:   common.SubCommandText,
	Short: "Send text to Discord",
	Run:   discordCmdGlobalVar.Run,
}

var discordSubCmdTextTS = &cobra.Command{
	Use:   common.SubCommandText + "TS",
	Short: "Send text to speech to Discord",
	Run:   discordCmdGlobalVar.Run,
}

var discordCmdGlobalVar DiscordFlag

func init() {
	rootCmd.AddCommand(discordCmd)

	discordCmd.PersistentFlags().StringVarP(&discordCmdGlobalVar.Token, "token", "t", "", common.Usage("Token"))
	discordCmd.PersistentFlags().StringVarP(&discordCmdGlobalVar.Channel, "channel-id", "c", "", common.Usage("Channel ID"))
	discordCmd.PersistentFlags().StringVarP(&discordCmdGlobalVar.arg, "arg", "a", "", common.Usage("Input argument"))

	discordCmd.AddCommand(discordSubCmdFile, discordSubCmdText, discordSubCmdTextTS)
}

type DiscordFlag struct {
	Token   string `json:"token"`
	Channel string `json:"channel_id"`
	arg     string

	api  *discordgo.Session
	resp *discordgo.Message
}

func (d *DiscordFlag) Run(cmd *cobra.Command, args []string) {
	if d.arg == "" {
		log.Println(ErrArgNotFound)
		os.Exit(1)
	}
	var err error
	if err = d.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	switch cmd.Name() {
	case common.SubCommandFile:
		err = d.File()
	case common.SubCommandText:
		err = d.Text()
	case common.SubCommandText + "TS":
		err = d.TextTTS()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	OutputDefaultNone(d.resp)
}

func (d *DiscordFlag) Init() error {
	var err error
	if d.Token == "" && rootConfig != "" {
		v := common.Config(rootConfig, common.Discord)
		err = Encoder.JSONMarshaler(v, d)
		if err != nil {
			return err
		}
	}
	if d.Token == "" {
		return ErrTokenNotFound
	}

	d.api, err = discordgo.New("Bot " + d.Token)
	if d.api == nil {
		return ErrInitialFailed
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
