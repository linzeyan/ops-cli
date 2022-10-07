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
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func initDoc(command *cobra.Command) *cobra.Command {
	var flags struct {
		dir string
	}
	var docCmd = &cobra.Command{
		Use:   CommandDoc,
		Short: "Generate documentation",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
	}

	runE := func(cmd *cobra.Command, _ []string) error {
		_, err := os.Stat(flags.dir)
		if err != nil {
			/* Create directory if not exist. */
			if err = os.Mkdir(flags.dir, os.ModePerm); err != nil {
				return err
			}
		}
		switch cmd.Name() {
		case CommandMan:
			header := &doc.GenManHeader{
				Title:   "MINE",
				Section: "3",
			}
			err = doc.GenManTree(command, header, flags.dir)
		case CommandMarkdown:
			err = doc.GenMarkdownTree(command, flags.dir)
		case CommandReST:
			err = doc.GenReSTTree(command, flags.dir)
		case CommandYaml:
			err = doc.GenYamlTree(command, flags.dir)
		}
		return err
	}

	var docSubCmdMan = &cobra.Command{
		Use:   CommandMan,
		Short: "Generate man page documentation",
		RunE:  runE,
	}

	var docSubCmdMarkdown = &cobra.Command{
		Use:   CommandMarkdown,
		Short: "Generate markdown documentation",
		RunE:  runE,
	}

	var docSubCmdRest = &cobra.Command{
		Use:   CommandReST,
		Short: "Generate rest documentation",
		RunE:  runE,
	}

	var docSubCmdYaml = &cobra.Command{
		Use:   CommandYaml,
		Short: "Generate yaml documentation",
		RunE:  runE,
	}

	docCmd.PersistentFlags().StringVarP(&flags.dir, "dir", "d", "doc", common.Usage("Specify the path to generate documentation"))
	docCmd.AddCommand(docSubCmdMan)
	docCmd.AddCommand(docSubCmdMarkdown)
	docCmd.AddCommand(docSubCmdRest)
	docCmd.AddCommand(docSubCmdYaml)
	return docCmd
}
