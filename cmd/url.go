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

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func initURL() *cobra.Command {
	var urlFlag struct {
		expand  bool
		verbose bool
		output  string
		method  string
		data    string
		headers string
	}
	var urlCmd = &cobra.Command{
		Use:   CommandURL,
		Args:  cobra.ExactArgs(1),
		Short: "Get url content or expand shorten url or download",
		RunE: func(_ *cobra.Command, args []string) error {
			url := args[0]
			if !validator.IsURL(url) {
				return common.ErrInvalidURL
			}
			var err error
			var result any
			switch {
			case urlFlag.expand:
				result, err = common.HTTPRequestRedirectURL(url)
				if err != nil {
					return err
				}
			default:
				body := common.HTTPConfig{
					Body:    urlFlag.data,
					Method:  urlFlag.method,
					Verbose: urlFlag.verbose,
					Headers: urlFlag.headers,
				}
				result, err = common.HTTPRequestContent(url, body)
				if err != nil || urlFlag.verbose {
					return err
				}
				if urlFlag.output != "" {
					return os.WriteFile(urlFlag.output, result.([]byte), FileModeRAll)
				}
			}
			printer.Printf(rootOutputFormat, result)
			return err
		},
		Example: common.Examples(`# Get the file from URL
https://raw.githubusercontent.com/golangci/golangci-lint/master/.golangci.reference.yml -o config.yaml

# Get the real URL from the shortened URL
https://goo.gl/maps/b37Aq3Anc7taXQDd9 -e`,
			CommandURL),
	}
	urlCmd.Flags().BoolVarP(&urlFlag.expand, "expand", "e", false, "Expand shorten url")
	urlCmd.Flags().BoolVarP(&urlFlag.verbose, "verbose", "v", false, "Verbose output")
	urlCmd.Flags().StringVarP(&urlFlag.output, "output", "o", "", "Write to file")
	urlCmd.Flags().StringVarP(&urlFlag.method, "method", "m", "GET", "Request method")
	urlCmd.Flags().StringVarP(&urlFlag.data, "data", "d", "", "Request method")
	urlCmd.Flags().StringVarP(&urlFlag.headers, "headers", "h", "", "Headers")
	return urlCmd
}
