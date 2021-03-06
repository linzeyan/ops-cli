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
	"encoding/json"
	"fmt"
	"log"
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

var rootCmd = &cobra.Command{
	Use:   "ops-cli",
	Short: "OPS useful tools",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
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
	rootCmd.PersistentFlags().StringVar(&rootConfig, "config", "", "Specify config")
}

func Config(subFn string) {
	if rootConfig == "" {
		return
	}
	viper.SetConfigFile(rootConfig)
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}
	switch subFn {
	case "icp":
		if i.account == "" && i.key == "" {
			i.account = viper.GetString("west.account")
			i.key = viper.GetString("west.api_key")
		}
	case "telegram":
		if tg.token == "" {
			tg.token = viper.GetString("telegram.token")
		}
		if tg.chat == 0 {
			tg.chat = viper.GetInt64("telegram.chat_id")
		}
	case "slack":
		if slk.token == "" {
			slk.token = viper.GetString("slack.token")
		}
		if slk.channel == "" {
			slk.channel = viper.GetString("slack.channel_id")
		}
	}
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

func OutputDefaultJSON(r rootOutput) {
	if rootOutputYAML {
		r.YAML()
	} else {
		r.JSON()
	}
}

func OutputDefaultYAML(r rootOutput) {
	if rootOutputJSON {
		r.JSON()
	} else {
		r.YAML()
	}
}

func PrintJSON(i interface{}) {
	out, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

func PrintYAML(i interface{}) {
	out, err := yaml.Marshal(i)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(out))
}

/* If i is a domain return true. */
func ValidDomain(i interface{}) bool {
	if val, ok := i.(string); ok {
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

/* If u is a valid url return true. */
func ValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
