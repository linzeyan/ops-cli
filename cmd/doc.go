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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docCmd = &cobra.Command{
	Use:   "doc [type]",
	Short: "Generate documentation",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}

		_, err := os.Stat(d.dir)
		if err != nil {
			mkErr := os.Mkdir(d.dir, 0755)
			if mkErr != nil {
				log.Println(mkErr)
				return
			}
		}
		switch strings.ToLower(args[0]) {
		case "man":
			d.Man()
		case "markdown":
			d.Markdown()
		case "rest":
			d.Rest()
		case "yaml":
			d.Yaml()
		default:
			_ = cmd.Help()
		}
	},
	Example: Examples(`# Generate different type documents
ops-cli doc man
ops-cli doc markdown
ops-cli doc rest
ops-cli doc yaml`),
}

var d docGenerate

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.Flags().StringVarP(&d.dir, "dir", "d", "doc", "Specify the path to generate documentation")
}

type docGenerate struct {
	dir string
}

func (docGenerate) Man() {
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	err := doc.GenManTree(rootCmd, header, d.dir)
	if err != nil {
		log.Println(err)
	}
}

func (docGenerate) Markdown() {
	err := doc.GenMarkdownTree(rootCmd, d.dir)
	if err != nil {
		log.Println(err)
	}
}

func (docGenerate) Rest() {
	err := doc.GenReSTTree(rootCmd, d.dir)
	if err != nil {
		log.Println(err)
	}
}

func (docGenerate) Yaml() {
	err := doc.GenYamlTree(rootCmd, d.dir)
	if err != nil {
		log.Println(err)
	}
}
