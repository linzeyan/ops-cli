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
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var rootNow = time.Now().Local()
var rootContext = context.Background()

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
		if i.account == "" {
			i.account = viper.GetString("west.account")
		}
		if i.key == "" {
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
		if line.secret == "" {
			line.secret = viper.GetString("line.secret")
		}
		if line.token == "" {
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
	String()
}

func OutputInterfaceString(r rootOutput) {
	switch {
	case rootOutputJSON:
		PrintJSON(r)
	case rootOutputYAML:
		PrintYAML(r)
	default:
		r.String()
	}
}

func OutputDefaultJSON(i interface{}) {
	if rootOutputYAML {
		PrintYAML(i)
	} else {
		PrintJSON(i)
	}
}

func OutputDefaultNone(i interface{}) {
	if rootOutputJSON {
		PrintJSON(i)
	} else if rootOutputYAML {
		PrintYAML(i)
	}
}

func OutputDefaultString(i interface{}) {
	switch {
	case rootOutputJSON:
		PrintJSON(i)
	case rootOutputYAML:
		PrintYAML(i)
	default:
		PrintString(i)
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
	out, err := Encoder.JSONEncode(i)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintYAML(i interface{}) {
	out, err := Encoder.YamlEncode(i)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintString(i interface{}) {
	fmt.Printf("%+v\n", i)
}

/* HttpRequestContent make a simple request to url, and return response body, default request method is get. */
func HTTPRequestContent(url string, body io.Reader, methods ...string) ([]byte, error) {
	var method string
	if len(methods) == 0 {
		method = http.MethodGet
	} else {
		method = methods[0]
	}
	var client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequestWithContext(rootContext, method, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return content, err
	}
	return nil, ErrResponseStatus
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
