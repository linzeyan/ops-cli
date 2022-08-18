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
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var discordCmd = &cobra.Command{
	Use:   "Discord",
	Short: "Send message to Discord",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var discordSubCmdFile = &cobra.Command{
	Use:   "file",
	Short: "Send file to Discord",
	Run:   discordCmdGlobalVar.Run,
}

var discordSubCmdText = &cobra.Command{
	Use:   "text",
	Short: "Send text to Discord",
	Run:   discordCmdGlobalVar.Run,
}

var discordSubCmdTextTS = &cobra.Command{
	Use:   "textTS",
	Short: "Send text to speech to Discord",
	Run:   discordCmdGlobalVar.Run,
}

var discordCmdGlobalVar DiscordFlag

func init() {
	rootCmd.AddCommand(discordCmd)

	discordCmd.PersistentFlags().StringVarP(&discordCmdGlobalVar.token, "token", "t", "", "Token")
	discordCmd.PersistentFlags().StringVarP(&discordCmdGlobalVar.channel, "channel-id", "c", "", "Channel ID")
	discordCmd.PersistentFlags().StringVarP(&discordCmdGlobalVar.arg, "arg", "a", "", "Input argument")

	discordCmd.AddCommand(discordSubCmdFile, discordSubCmdText, discordSubCmdTextTS)
}

type DiscordFlag struct {
	arg     string
	token   string
	channel string

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
	case ImTypeFile:
		err = d.File()
	case ImTypeText:
		err = d.Text()
	case ImTypeText + "TS":
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
	if d.token == "" && rootConfig != "" {
		if err = Config(ConfigBlockDiscord); err != nil {
			return err
		}
	}
	if d.token == "" {
		return ErrTokenNotFound
	}

	d.api, err = discordgo.New("Bot " + d.token)
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
	filename := path.Base(d.arg)
	defer f.Close()
	d.resp, err = d.api.ChannelFileSend(d.channel, filename, f)
	return err
}

func (d *DiscordFlag) Text() error {
	var err error
	d.resp, err = d.api.ChannelMessageSend(d.channel, d.arg)
	return err
}

func (d *DiscordFlag) TextTTS() error {
	var err error
	d.resp, err = d.api.ChannelMessageSendTTS(d.channel, d.arg)
	return err
}
