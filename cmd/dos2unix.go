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
	"regexp"

	"github.com/spf13/cobra"
)

var dos2unixCmd = &cobra.Command{
	Use:   "dos2unix",
	Short: "Convert file eol to unix style",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
		for _, f := range args {
			if err := Dos2Unix(f); err != nil {
				log.Printf("%s: %v\n", f, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dos2unixCmd)
}

func Dos2Unix(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	eol := regexp.MustCompile(`\r\n`)
	f = eol.ReplaceAllLiteral(f, []byte{'\n'})
	return os.WriteFile(filename, f, os.ModePerm)
}
