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
	RunE:  docCmdGlobalVar.RunE,
}

var docSubCmdMarkdown = &cobra.Command{
	Use:   "markdown",
	Short: "Generate markdown documentation",
	RunE:  docCmdGlobalVar.RunE,
}

var docSubCmdRest = &cobra.Command{
	Use:   "rest",
	Short: "Generate rest documentation",
	RunE:  docCmdGlobalVar.RunE,
}

var docSubCmdYaml = &cobra.Command{
	Use:   "yaml",
	Short: "Generate yaml documentation",
	RunE:  docCmdGlobalVar.RunE,
}

var docCmdGlobalVar DocFlag

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.PersistentFlags().StringVarP(&docCmdGlobalVar.dir, "dir", "d", "doc", "Specify the path to generate documentation")
	docCmd.AddCommand(docSubCmdMan)
	docCmd.AddCommand(docSubCmdMarkdown)
	docCmd.AddCommand(docSubCmdRest)
	docCmd.AddCommand(docSubCmdYaml)
}

type DocFlag struct {
	dir string
}

func (d *DocFlag) createDir() error {
	var err error
	_, err = os.Stat(d.dir)
	if err != nil {
		/* Create directory if not exist. */
		if err = os.Mkdir(d.dir, os.ModePerm); err != nil {
			return err
		}
	}
	return err
}

func (d *DocFlag) RunE(cmd *cobra.Command, _ []string) error {
	var err error
	if err = d.createDir(); err != nil {
		return err
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
	}
	return err
}
