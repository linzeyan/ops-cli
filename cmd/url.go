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

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "URL expand or download",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
	DisableFlagParsing:    true,
}

var urlSubCmdExpand = &cobra.Command{
	Use:   "expand [url]",
	Args:  cobra.ExactArgs(1),
	Short: "Expand shorten url",
	Run: func(_ *cobra.Command, args []string) {
		if !validator.ValidURL(args[0]) {
			log.Println(common.ErrInvalidURL)
			os.Exit(1)
		}
		result, err := common.HTTPRequestRedirectURL(args[0])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		PrintString(result)
	},
	Example: common.Examples(`# Get the real URL from the shortened URL
ops-cli url expand https://goo.gl/maps/b37Aq3Anc7taXQDd9`),
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    true,
}

var urlSubCmdGet = &cobra.Command{
	Use:   "get [url] [output]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get file from url",
	Run: func(_ *cobra.Command, args []string) {
		if !validator.ValidURL(args[0]) {
			log.Println(common.ErrInvalidURL)
			os.Exit(1)
		}
		filename := path.Base(args[0])
		if len(args) > 1 {
			filename = args[1]
		}
		result, err := common.HTTPRequestContent(args[0], nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		err = os.WriteFile(filename, result, os.FileMode(0644))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
	Example: common.Examples(`# Get the file from URL
ops-cli url get https://raw.githubusercontent.com/golangci/golangci-lint/master/.golangci.reference.yml`),
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    true,
}

func init() {
	rootCmd.AddCommand(urlCmd)

	urlCmd.AddCommand(urlSubCmdExpand)
	urlCmd.AddCommand(urlSubCmdGet)
}
