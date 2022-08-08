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
	Use:   "doc",
	Short: "Generate documentation",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var docSubCmdMan = &cobra.Command{
	Use:   "man",
	Short: "Generate man page documentation",
	Run:   df.Run,
}

var docSubCmdMarkdown = &cobra.Command{
	Use:   "markdown",
	Short: "Generate markdown documentation",
	Run:   df.Run,
}

var docSubCmdRest = &cobra.Command{
	Use:   "rest",
	Short: "Generate rest documentation",
	Run:   df.Run,
}

var docSubCmdYaml = &cobra.Command{
	Use:   "yaml",
	Short: "Generate yaml documentation",
	Run:   df.Run,
}

var df docFlag

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.PersistentFlags().StringVarP(&df.dir, "dir", "d", "doc", "Specify the path to generate documentation")
	docCmd.AddCommand(docSubCmdMan)
	docCmd.AddCommand(docSubCmdMarkdown)
	docCmd.AddCommand(docSubCmdRest)
	docCmd.AddCommand(docSubCmdYaml)
}

type docFlag struct {
	dir string
}

func (d *docFlag) createDir() error {
	_, err := os.Stat(d.dir)
	if err != nil {
		/* Create directory if not exist. */
		mkErr := os.Mkdir(d.dir, os.ModePerm)
		if mkErr != nil {
			return mkErr
		}
	}
	return nil
}

func (d *docFlag) Run(cmd *cobra.Command, _ []string) {
	var err error
	if err = d.createDir(); err != nil {
		log.Println(err)
		return
	}
	switch cmd.Name() {
	case DocTypeMan:
		header := &doc.GenManHeader{
			Title:   "MINE",
			Section: "3",
		}
		err = doc.GenManTree(rootCmd, header, d.dir)
	case DocTypeMarkdown:
		err = doc.GenMarkdownTree(rootCmd, d.dir)
	case DocTypeReST:
		err = doc.GenReSTTree(rootCmd, d.dir)
	case DocTypeYaml:
		err = doc.GenYamlTree(rootCmd, d.dir)
	default:
		_ = cmd.Help()
		return
	}
	if err != nil {
		log.Println(err)
	}
}
