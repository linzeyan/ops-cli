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
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

func init() {
	var slackFlag SlackFlag
	var slackCmd = &cobra.Command{
		Use:   CommandSlack,
		Short: "Send message to Slack",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	var slackSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Short: "Send file to Slack",
		RunE:  slackFlag.RunE,
		Example: common.Examples(`# Send file
-a "/tmp/a.txt" --config ~/.config.toml`, CommandSlack, CommandFile),
	}

	var slackSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send text to Slack",
		RunE:  slackFlag.RunE,
		Example: common.Examples(`# Send text
-a "Hello World!"`, CommandSlack, CommandText),
	}

	var slackSubCmdPhoto = &cobra.Command{
		Use:   CommandPhoto,
		Short: "Send photo to Slack",
		RunE:  slackFlag.RunE,
		Example: common.Examples(`# Send photo
-a "~/robot.png"`, CommandSlack, CommandPhoto),
	}
	rootCmd.AddCommand(slackCmd)

	slackCmd.PersistentFlags().StringVarP(&slackFlag.Token, "token", "t", "", common.Usage("Bot token (required)"))
	slackCmd.PersistentFlags().StringVarP(&slackFlag.Channel, "channel", "c", "", common.Usage("Channel ID"))
	slackCmd.PersistentFlags().StringVarP(&slackFlag.arg, "arg", "a", "", common.Usage("Input argument"))

	slackCmd.AddCommand(slackSubCmdFile)
	slackCmd.AddCommand(slackSubCmdText)
	slackCmd.AddCommand(slackSubCmdPhoto)
}

type SlackFlag struct {
	Token   string `json:"token"`
	Channel string `json:"channel_id"`
	arg     string

	api *slack.Client
}

func (s *SlackFlag) RunE(cmd *cobra.Command, _ []string) error {
	if s.arg == "" {
		return common.ErrInvalidFlag
	}
	var err error
	if err = s.Init(); err != nil {
		return err
	}
	switch cmd.Name() {
	case CommandFile:
		err = s.Photo()
	case CommandPhoto:
		err = s.Photo()
	case CommandText:
		err = s.Text()
	}
	return err
}

func (s *SlackFlag) Init() error {
	var err error
	if s.Token == "" && rootConfig != "" {
		v := common.Config(rootConfig, strings.ToLower(CommandSlack))
		err = Encoder.JSONMarshaler(v, s)
		if err != nil {
			return err
		}
	}
	if s.Token == "" {
		return common.ErrInvalidToken
	}
	s.api = slack.New(s.Token)
	if s.api == nil {
		return common.ErrFailedInitial
	}
	return err
}

func (s *SlackFlag) Text() error {
	input := slack.MsgOptionText(s.arg, false)
	_, _, _, err := s.api.SendMessageContext(common.Context, s.Channel, input)
	return err
}

func (s *SlackFlag) Photo() error {
	var base64Image string
	var err error
	switch {
	case validator.ValidFile(s.arg):
		base64Image, err = s.localFile()
	case validator.ValidURL(s.arg):
		base64Image, err = s.remoteFile()
	}
	if err != nil {
		return err
	}

	uploadFileKey := fmt.Sprintf("upload-f-to-slack-%d", common.TimeNow.UnixNano())
	decode, err := Encoder.Base64StdDecode(base64Image)
	if err != nil {
		return err
	}

	f, err := os.Create(uploadFileKey)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(decode); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	_, err = s.api.UploadFileContext(common.Context, slack.FileUploadParameters{
		Filetype: "image/png",
		Filename: filepath.Base(s.arg),
		Channels: []string{s.Channel},
		File:     uploadFileKey,
	})
	if err != nil {
		return err
	}
	if uploadFileKey != "" {
		if err := os.Remove(uploadFileKey); err != nil {
			return err
		}
	}
	return err
}

func (s *SlackFlag) localFile() (string, error) {
	content, err := os.ReadFile(s.arg)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(content)
	if contentType == "image/jpeg" {
		img, err := jpeg.Decode(bytes.NewReader(content))
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return "", err
		}
		content = buf.Bytes()
	}
	return Encoder.Base64StdEncode(content)
}

func (s *SlackFlag) remoteFile() (string, error) {
	content, err := common.HTTPRequestContent(s.arg)
	if err != nil {
		return "", err
	}
	if http.DetectContentType(content) == "image/jpeg" {
		img, err := jpeg.Decode(bytes.NewReader(content))
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return "", err
		}
		content = buf.Bytes()
	}
	base64Image := base64.StdEncoding.EncodeToString(content)
	return base64Image, err
}
