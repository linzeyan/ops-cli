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

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docCmd = &cobra.Command{
	Use:   "doc [type]",
	Short: "Generate documentation",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var subCmdMan = &cobra.Command{
	Use:   "man",
	Short: "Generate man page documentation",
	Run:   func(_ *cobra.Command, _ []string) { df.Man() },
}

var subCmdMarkdown = &cobra.Command{
	Use:   "markdown",
	Short: "Generate markdown documentation",
	Run:   func(_ *cobra.Command, _ []string) { df.Markdown() },
}

var subCmdRest = &cobra.Command{
	Use:   "rest",
	Short: "Generate rest documentation",
	Run:   func(_ *cobra.Command, _ []string) { df.Rest() },
}

var subCmdYaml = &cobra.Command{
	Use:   "yaml",
	Short: "Generate yaml documentation",
	Run:   func(_ *cobra.Command, _ []string) { df.Yaml() },
}

var df docFlag

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.PersistentFlags().StringVarP(&df.dir, "dir", "d", "doc", "Specify the path to generate documentation")
	docCmd.AddCommand(subCmdMan)
	docCmd.AddCommand(subCmdMarkdown)
	docCmd.AddCommand(subCmdRest)
	docCmd.AddCommand(subCmdYaml)
}

type docFlag struct {
	dir string
}

func (d docFlag) createDir() {
	_, err := os.Stat(df.dir)
	if err != nil {
		/* Create directory if not exist. */
		mkErr := os.Mkdir(df.dir, 0755)
		if mkErr != nil {
			log.Println(mkErr)
			return
		}
	}
}

func (d docFlag) Man() {
	d.createDir()
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	err := doc.GenManTree(rootCmd, header, d.dir)
	if err != nil {
		log.Println(err)
	}
}

func (d docFlag) Markdown() {
	d.createDir()
	err := doc.GenMarkdownTree(rootCmd, d.dir)
	if err != nil {
		log.Println(err)
	}
}

func (d docFlag) Rest() {
	d.createDir()
	err := doc.GenReSTTree(rootCmd, d.dir)
	if err != nil {
		log.Println(err)
	}
}

func (d docFlag) Yaml() {
	d.createDir()
	err := doc.GenYamlTree(rootCmd, d.dir)
	if err != nil {
		log.Println(err)
	}
}
