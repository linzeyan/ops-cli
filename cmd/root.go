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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ByteSize float64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

type ConfigBlock string

const (
	ConfigBlockICP      ConfigBlock = "icp"
	ConfigBlockLINE     ConfigBlock = "line"
	ConfigBlockSlack    ConfigBlock = "slack"
	ConfigBlockTelegram ConfigBlock = "telegram"
)

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

type RandomCharacter string

const (
	LowercaseLetters RandomCharacter = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters RandomCharacter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Symbols          RandomCharacter = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	Numbers          RandomCharacter = "0123456789"
	AllSet           RandomCharacter = LowercaseLetters + UppercaseLetters + Symbols + Numbers
)

var (
	ErrArgNotFound   = errors.New("argument not found")
	ErrConfNotFound  = errors.New("config not found")
	ErrFileNotFound  = errors.New("file not found")
	ErrTokenNotFound = errors.New("token not found")
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

var rootNow = time.Now().Local()
var rootContext = context.Background()

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
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
		if i.account == "" && i.key == "" {
			i.account = viper.GetString("west.account")
			i.key = viper.GetString("west.api_key")
		}
	case ConfigBlockTelegram:
		if tg.token == "" {
			tg.token = viper.GetString("telegram.token")
		}
		if tg.chat == 0 {
			tg.chat = viper.GetInt64("telegram.chat_id")
		}
	case ConfigBlockSlack:
		if sf.token == "" {
			sf.token = viper.GetString("slack.token")
		}
		if sf.channel == "" {
			sf.channel = viper.GetString("slack.channel_id")
		}
	case ConfigBlockLINE:
		if line.secret == "" && line.token == "" {
			line.secret = viper.GetString("line.secret")
			line.token = viper.GetString("line.access_token")
		}
		if line.id == "" {
			line.id = viper.GetString("line.id")
		}
	default:
		return ErrConfNotFound
	}
	return err
}

/* Print examples with color. */
func Examples(s string) string {
	c := color.New(color.FgYellow)
	return c.Sprintf(`%s`, s)
}

type rootOutput interface {
	JSON()
	YAML()
	String()
}

func OutputDefaultString(r rootOutput) {
	switch {
	case rootOutputJSON:
		r.JSON()
	case rootOutputYAML:
		r.YAML()
	default:
		r.String()
	}
}

func OutputDefaultNone(i interface{}) {
	if rootOutputJSON {
		PrintJSON(i)
	} else if rootOutputYAML {
		PrintYAML(i)
	}
}

func OutputDefaultJSON(i interface{}) {
	if rootOutputYAML {
		PrintYAML(i)
	} else {
		PrintJSON(i)
	}
}

func OutputDefaultYAML(i interface{}) {
	if rootOutputJSON {
		PrintJSON(i)
	} else {
		PrintYAML(i)
	}
}

func PrintJSON(i interface{}) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(i); err != nil {
		log.Println(err)
		return
	}
	PrintString(buf.String())
}

func PrintYAML(i interface{}) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(i); err != nil {
		log.Println(err)
		return
	}
	PrintString(buf.String())
}

func PrintString(i interface{}) {
	fmt.Printf("%+v\n", i)
}

/* If i is a domain return true. */
func ValidDomain(i interface{}) bool {
	const elements = "~!@#$%^&*()_+`={}|[]\\:\"<>?,/"
	if val, ok := i.(string); ok {
		if strings.ContainsAny(val, elements) {
			return false
		}
		slice := strings.Split(val, ".")
		l := len(slice)
		if l > 1 {
			n, err := strconv.Atoi(slice[l-1])
			if err != nil {
				return true
			}
			s := strconv.Itoa(n)
			return slice[l-1] != s
		}
	}
	return false
}

/* If f is a valid path return true. */
func ValidFile(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

/* If i is a ipv address return true. */
func ValidIP(i string) bool {
	return net.ParseIP(i) != nil
}

/* If i is a ipv4 address return true. */
func ValidIPv4(i string) bool {
	return net.ParseIP(i).To4() != nil
}

/* If i is a ipv6 address return true. */
func ValidIPv6(i string) bool {
	return net.ParseIP(i).To4() == nil && net.ParseIP(i).To16() != nil
}

/* If u is a valid url return true. */
func ValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
