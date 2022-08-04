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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var slackCmd = &cobra.Command{
	Use:   "slack [function]",
	Short: "Send message to slack",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		if args[0] == "" {
			_ = cmd.Help()
			return
		}
		Config("slack")
		if slk.token == "" {
			_ = cmd.Help()
			return
		}
		slackAPI := slk.Init()
		if slackAPI == nil {
			_ = cmd.Help()
			return
		}
		switch strings.ToLower(args[0]) {
		case imTypeDoc, imTypeDocument, imTypeFile:
			slk.fileName = args[0]
			slk.Photo(slackAPI)
		case imTypeMsg, imTypeMessage:
			slk.Msg(slackAPI)
		case imTypePhoto:
			slk.Photo(slackAPI)
		}
	}, Example: Examples(`# Send message
ops-cli slack msg -a "Hello World!"`),
}

var slk slackFlag

func init() {
	rootCmd.AddCommand(slackCmd)

	slackCmd.Flags().StringVarP(&slk.token, "token", "t", "", "Bot token (required)")
	slackCmd.Flags().StringVarP(&slk.channel, "channel", "c", "", "Channel ID")
	slackCmd.Flags().StringVarP(&slk.arg, "arg", "a", "", "Input argument")
}

type slackFlag struct {
	token   string
	channel string
	arg     string

	fileName string
}

func (s slackFlag) Init() *slack.Client {
	return slack.New(s.token)
}

func (s slackFlag) Msg(api *slack.Client) {
	input := slack.MsgOptionText(s.arg, false)
	_, _, _, err := api.SendMessageContext(rootContext, s.channel, input)
	if err != nil {
		log.Println(err)
	}
}

func (s slackFlag) Photo(api *slack.Client) {
	var base64Image string
	switch {
	case ValidFile(s.arg):
		base64Image = s.localFile()
	case ValidURL(s.arg):
		base64Image = s.remoteFile()
	}

	uploadFileKey := fmt.Sprintf("upload-f-to-slack-%d", rootNow.UnixNano())
	decode, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.Create(uploadFileKey)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	if _, err := f.Write(decode); err != nil {
		log.Println(err)
	}
	if err := f.Sync(); err != nil {
		log.Println(err)
	}
	if s.fileName == "" {
		s.fileName = "upload.png"
	}
	_, err = api.UploadFileContext(rootContext, slack.FileUploadParameters{
		Filetype: "image/png",
		Filename: s.fileName,
		Channels: []string{s.channel},
		File:     uploadFileKey,
	})
	if err != nil {
		log.Println(err)
		return
	}
	if uploadFileKey != "" {
		if err := os.Remove(uploadFileKey); err != nil {
			log.Println(err)
		}
	}
}

func (s slackFlag) localFile() string {
	content, err := os.ReadFile(s.arg)
	if err != nil {
		log.Println(err)
		return ""
	}
	contentType := http.DetectContentType(content)
	if contentType == "image/jpeg" {
		img, err := jpeg.Decode(bytes.NewReader(content))
		if err != nil {
			log.Println(err)
			return ""
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return ""
		}
		content = buf.Bytes()
	}
	base64Image := base64.StdEncoding.EncodeToString(content)
	return base64Image
}

func (s slackFlag) remoteFile() string {
	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequestWithContext(rootContext, http.MethodGet, s.arg, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println(err)
		return ""
	}
	var content []byte
	if resp.StatusCode == http.StatusOK {
		content, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
	}
	if http.DetectContentType(content) == "image/jpeg" {
		img, err := jpeg.Decode(bytes.NewReader(content))
		if err != nil {
			log.Println(err)
			return ""
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return ""
		}
		content = buf.Bytes()
	}
	base64Image := base64.StdEncoding.EncodeToString(content)
	return base64Image
}
