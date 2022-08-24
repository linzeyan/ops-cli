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
	"os"
	"path/filepath"

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
	RunE: func(_ *cobra.Command, args []string) error {
		if !validator.ValidURL(args[0]) {
			return common.ErrInvalidURL
		}
		result, err := common.HTTPRequestRedirectURL(args[0])
		if err != nil {
			return err
		}
		PrintString(result)
		return err
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
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		if !validator.ValidURL(args[0]) {
			return common.ErrInvalidURL
		}
		filename := filepath.Base(args[0])
		if len(args) > 1 {
			filename = args[1]
		}
		result, err := common.HTTPRequestContent(args[0], nil)
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, result, common.FileModeRAll)
		if err != nil {
			return err
		}
		return err
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
