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
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigBlock string

const (
	ConfigBlockDiscord  ConfigBlock = "discord."
	ConfigBlockICP      ConfigBlock = "west."
	ConfigBlockLINE     ConfigBlock = "line."
	ConfigBlockSlack    ConfigBlock = "slack."
	ConfigBlockTelegram ConfigBlock = "telegram."
)

func (c ConfigBlock) String() string {
	return string(c)
}

const (
	DocTypeMan      = "man"
	DocTypeMarkdown = "markdown"
	DocTypeReST     = "rest"
	DocTypeYaml     = "yaml"
)

const (
	FileTypeHTML     = "html"
	FileTypeJSON     = "json"
	FileTypeMarkdown = "markdown"
	FileTypePDF      = "pdf"
	FileTypeTOML     = "toml"
	FileTypeYAML     = "yaml"
)

const (
	ImTypeAudio = "audio"
	ImTypeFile  = "file"
	ImTypeID    = "id"
	ImTypePhoto = "photo"
	ImTypeText  = "text"
	ImTypeVideo = "video"
	ImTypeVoice = "voice"
)

var (
	ErrArgNotFound    = errors.New("argument not found")
	ErrConfNotFound   = errors.New("config not found")
	ErrEmptyResponse  = errors.New("response is empty")
	ErrFileNotFound   = errors.New("file not found")
	ErrFileType       = errors.New("file type not correct")
	ErrInitialFailed  = errors.New("initial failed")
	ErrInvalidIP      = errors.New("invalid IP")
	ErrInvalidLength  = errors.New("invalid length")
	ErrInvalidVar     = errors.New("invalid variable")
	ErrParseCert      = errors.New("can not correctly parse certificate")
	ErrResponseStatus = errors.New("response status code is not 200")
	ErrTokenNotFound  = errors.New("token not found")
)

var rootCmd = &cobra.Command{
	Use:   "ops-cli",
	Short: "OPS useful tools",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

/* Flags. */
var (
	rootConfig     string
	rootOutputJSON bool
	rootOutputYAML bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootOutputJSON, "json", "j", false, "Output JSON format")
	rootCmd.PersistentFlags().BoolVarP(&rootOutputYAML, "yaml", "y", false, "Output YAML format")
	rootCmd.PersistentFlags().StringVar(&rootConfig, "config", "", "Specify config path (toml)")
}

/* Get secret token from config. */
func Config(subFn ConfigBlock) error {
	var err error
	if rootConfig == "" {
		return ErrConfNotFound
	}
	viper.SetConfigFile(rootConfig)
	viper.SetConfigType(FileTypeTOML)
	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	switch subFn {
	case ConfigBlockICP:
		if icpCmdGlobalVar.account == "" {
			icpCmdGlobalVar.account = viper.GetString(ConfigBlockICP.String() + "account")
		}
		if icpCmdGlobalVar.key == "" {
			icpCmdGlobalVar.key = viper.GetString(ConfigBlockICP.String() + "api_key")
		}
	case ConfigBlockTelegram:
		if telegramCmdGlobalVar.token == "" {
			telegramCmdGlobalVar.token = viper.GetString(ConfigBlockTelegram.String() + "token")
		}
		if telegramCmdGlobalVar.chat == 0 {
			telegramCmdGlobalVar.chat = viper.GetInt64(ConfigBlockTelegram.String() + "chat_id")
		}
	case ConfigBlockSlack:
		if slackCmdGlobalVar.token == "" {
			slackCmdGlobalVar.token = viper.GetString(ConfigBlockSlack.String() + "token")
		}
		if slackCmdGlobalVar.channel == "" {
			slackCmdGlobalVar.channel = viper.GetString(ConfigBlockSlack.String() + "channel_id")
		}
	case ConfigBlockLINE:
		if lineCmdGlobalVar.secret == "" {
			lineCmdGlobalVar.secret = viper.GetString(ConfigBlockLINE.String() + "secret")
		}
		if lineCmdGlobalVar.token == "" {
			lineCmdGlobalVar.token = viper.GetString(ConfigBlockLINE.String() + "access_token")
		}
		if lineCmdGlobalVar.id == "" {
			lineCmdGlobalVar.id = viper.GetString(ConfigBlockLINE.String() + "id")
		}
	case ConfigBlockDiscord:
		if discordCmdGlobalVar.token == "" {
			discordCmdGlobalVar.token = viper.GetString(ConfigBlockDiscord.String() + "token")
		}
		if discordCmdGlobalVar.channel == "" {
			discordCmdGlobalVar.channel = viper.GetString(ConfigBlockDiscord.String() + "channel_id")
		}
	default:
		return ErrConfNotFound
	}
	return err
}

type OutputFormat interface {
	String()
}

func OutputInterfaceString(r OutputFormat) {
	switch {
	case rootOutputJSON:
		PrintJSON(r)
	case rootOutputYAML:
		PrintYAML(r)
	default:
		r.String()
	}
}

func OutputDefaultJSON(i any) {
	if rootOutputYAML {
		PrintYAML(i)
	} else {
		PrintJSON(i)
	}
}

func OutputDefaultNone(i any) {
	if rootOutputJSON {
		PrintJSON(i)
	} else if rootOutputYAML {
		PrintYAML(i)
	}
}

func OutputDefaultString(i any) {
	switch {
	case rootOutputJSON:
		PrintJSON(i)
	case rootOutputYAML:
		PrintYAML(i)
	default:
		PrintString(i)
	}
}

func OutputDefaultYAML(i any) {
	if rootOutputJSON {
		PrintJSON(i)
	} else {
		PrintYAML(i)
	}
}

func PrintJSON(i any) {
	out, err := Encoder.JSONEncode(i)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintYAML(i any) {
	out, err := Encoder.YamlEncode(i)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintString(i any) {
	fmt.Printf("%s\n", i)
}
