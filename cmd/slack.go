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
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var slackCmd = &cobra.Command{
	Use:   "Slack",
	Short: "Send message to Slack",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var slackSubCmdFile = &cobra.Command{
	Use:   "file",
	Short: "Send file to Slack",
	Run:   sf.Run,
	Example: Examples(`# Send file
ops-cli Slack file -a "/tmp/a.txt" --config ~/.config.toml`),
}

var slackSubCmdText = &cobra.Command{
	Use:   "text",
	Short: "Send text to Slack",
	Run:   sf.Run,
	Example: Examples(`# Send text
ops-cli Slack text -a "Hello World!"`),
}

var slackSubCmdPhoto = &cobra.Command{
	Use:   "photo",
	Short: "Send photo to Slack",
	Run:   sf.Run,
	Example: Examples(`# Send photo
ops-cli Slack photo -a "~/robot.png"`),
}

var sf slackFlag

func init() {
	rootCmd.AddCommand(slackCmd)

	slackCmd.PersistentFlags().StringVarP(&sf.token, "token", "t", "", "Bot token (required)")
	slackCmd.PersistentFlags().StringVarP(&sf.channel, "channel", "c", "", "Channel ID")
	slackCmd.PersistentFlags().StringVarP(&sf.arg, "arg", "a", "", "Input argument")

	slackCmd.AddCommand(slackSubCmdFile)
	slackCmd.AddCommand(slackSubCmdText)
	slackCmd.AddCommand(slackSubCmdPhoto)
}

type slackFlag struct {
	token   string
	channel string
	arg     string

	fileName string
	api      *slack.Client
}

func (s *slackFlag) Run(cmd *cobra.Command, _ []string) {
	if s.arg == "" {
		log.Println(ErrArgNotFound)
		os.Exit(1)
	}
	var err error
	if err = s.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	switch cmd.Name() {
	case ImTypeFile:
		s.fileName = cmd.Name()
		err = s.Photo()
	case ImTypePhoto:
		err = s.Photo()
	case ImTypeText:
		err = s.Text()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (s *slackFlag) Init() error {
	var err error
	if s.token == "" && rootConfig != "" {
		if err = Config(ConfigBlockSlack); err != nil {
			return err
		}
	}
	if s.token == "" {
		return ErrTokenNotFound
	}
	s.api = slack.New(s.token)
	if s.api == nil {
		return ErrInitialFailed
	}
	return err
}

func (s *slackFlag) Text() error {
	input := slack.MsgOptionText(s.arg, false)
	_, _, _, err := s.api.SendMessageContext(rootContext, s.channel, input)
	return err
}

func (s *slackFlag) Photo() error {
	var base64Image string
	var err error
	switch {
	case ValidFile(s.arg):
		base64Image, err = s.localFile()
	case ValidURL(s.arg):
		base64Image, err = s.remoteFile()
	}
	if err != nil {
		return err
	}

	uploadFileKey := fmt.Sprintf("upload-f-to-slack-%d", rootNow.UnixNano())
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
	if s.fileName == "" {
		s.fileName = "upload.png"
	}
	_, err = s.api.UploadFileContext(rootContext, slack.FileUploadParameters{
		Filetype: "image/png",
		Filename: s.fileName,
		Channels: []string{s.channel},
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

func (s *slackFlag) localFile() (string, error) {
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

func (s *slackFlag) remoteFile() (string, error) {
	content, err := HTTPRequestContent(s.arg, nil)
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
