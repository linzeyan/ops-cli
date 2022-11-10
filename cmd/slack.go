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

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

func initSlack() *cobra.Command {
	var flags struct {
		Token   string `json:"token"`
		Channel string `json:"channel_id"`
		arg     string
	}
	var slackCmd = &cobra.Command{
		GroupID: groupings[CommandSlack],
		Use:     CommandSlack,
		Short:   "Send message to Slack",
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
			if err = ReadConfig(CommandSlack, &flags); err != nil {
				logger.Error(err.Error())
				return
			}
		}
		var s Slack
		if err = s.Init(flags.Token); err != nil {
			logger.Error(err.Error())
			return
		}
		switch cmd.Name() {
		case CommandFile:
			err = s.Photo(flags.Channel, flags.arg)
		case CommandPhoto:
			err = s.Photo(flags.Channel, flags.arg)
		case CommandText:
			err = s.Text(flags.Channel, flags.arg)
		}
		if err != nil {
			logger.Error(err.Error())
		}
	}

	var slackSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Short: "Send file to Slack",
		Run:   run,
		Example: common.Examples(`# Send file
-a "/tmp/a.txt" --config ~/.config.toml`, CommandSlack, CommandFile),
	}

	var slackSubCmdText = &cobra.Command{
		Use:   CommandText,
		Short: "Send text to Slack",
		Run:   run,
		Example: common.Examples(`# Send text
-a "Hello World!"`, CommandSlack, CommandText),
	}

	var slackSubCmdPhoto = &cobra.Command{
		Use:   CommandPhoto,
		Short: "Send photo to Slack",
		Run:   run,
		Example: common.Examples(`# Send photo
-a "~/robot.png"`, CommandSlack, CommandPhoto),
	}

	slackCmd.PersistentFlags().StringVarP(&flags.Token, "token", "t", "", common.Usage("Bot token (required)"))
	slackCmd.PersistentFlags().StringVarP(&flags.Channel, "channel", "c", "", common.Usage("Channel ID"))
	slackCmd.PersistentFlags().StringVarP(&flags.arg, "arg", "a", "", common.Usage("Input argument"))

	slackCmd.AddCommand(slackSubCmdFile, slackSubCmdText, slackSubCmdPhoto)
	return slackCmd
}

type Slack struct {
	API *slack.Client
}

func (s *Slack) Init(token string) error {
	if token == "" {
		logger.Debug(common.ErrInvalidToken.Error(), common.DefaultField(token))
		return common.ErrInvalidToken
	}
	s.API = slack.New(token)
	if s.API == nil {
		logger.Debug(common.ErrFailedInitial.Error())
		return common.ErrFailedInitial
	}
	return nil
}

func (s *Slack) Text(channel, arg string) error {
	input := slack.MsgOptionText(arg, false)
	_, _, _, err := s.API.SendMessageContext(common.Context, channel, input)
	if err != nil {
		logger.Debug(err.Error(),
			common.DefaultField(channel),
			common.NewField("slack.MsgOption", input),
		)
	}
	return err
}

func (s *Slack) Photo(channel, arg string) error {
	var base64Image string
	var err error
	switch {
	case common.IsFile(arg):
		base64Image, err = s.localFile(arg)
	case common.IsURL(arg):
		base64Image, err = s.remoteFile(arg)
	}
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(arg))
		return err
	}

	uploadFileKey := fmt.Sprintf("upload-f-to-slack-%d", common.TimeNow.UnixNano())
	decode, err := Encoder.Base64StdDecode(base64Image)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(base64Image))
		return err
	}

	f, err := os.Create(uploadFileKey)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(uploadFileKey))
		return err
	}
	defer f.Close()
	if _, err := f.Write(decode); err != nil {
		logger.Debug(err.Error(), common.DefaultField(decode))
		return err
	}
	if err := f.Sync(); err != nil {
		logger.Debug(err.Error())
		return err
	}
	_, err = s.API.UploadFileContext(common.Context, slack.FileUploadParameters{
		Filetype: "image/png",
		Filename: filepath.Base(arg),
		Channels: []string{channel},
		File:     uploadFileKey,
	})
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	if uploadFileKey != "" {
		if err := os.Remove(uploadFileKey); err != nil {
			logger.Debug(err.Error(), common.DefaultField(uploadFileKey))
			return err
		}
	}
	return err
}

func (s *Slack) localFile(arg string) (string, error) {
	content, err := os.ReadFile(arg)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(arg))
		return "", err
	}
	contentType := http.DetectContentType(content)
	if contentType == "image/jpeg" {
		img, err := jpeg.Decode(bytes.NewReader(content))
		if err != nil {
			logger.Debug(err.Error(), common.DefaultField(content))
			return "", err
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			logger.Debug(err.Error())
			return "", err
		}
		content = buf.Bytes()
	}
	return Encoder.Base64StdEncode(content)
}

func (s *Slack) remoteFile(arg string) (string, error) {
	content, err := common.HTTPRequestContent(arg)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(arg))
		return "", err
	}
	if http.DetectContentType(content) == "image/jpeg" {
		img, err := jpeg.Decode(bytes.NewReader(content))
		if err != nil {
			logger.Debug(err.Error(), common.DefaultField(content))
			return "", err
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			logger.Debug(err.Error())
			return "", err
		}
		content = buf.Bytes()
	}
	base64Image := base64.StdEncoding.EncodeToString(content)
	return base64Image, err
}
