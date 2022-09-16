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
	"path/filepath"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

func init() {
	var convertFlag ConvertFlag

	var convertCmd = &cobra.Command{
		Use:  CommandConvert,
		Args: cobra.OnlyValidArgs,
		ValidArgs: []string{
			CommandCsv2JSON, CommandCsv2Toml, CommandCsv2XML, CommandCsv2Yaml,
			CommandJSON2Csv, CommandJSON2Toml, CommandJSON2XML, CommandJSON2Yaml,
			CommandToml2Csv, CommandToml2JSON, CommandToml2XML, CommandToml2Yaml,
			CommandYaml2Csv, CommandYaml2JSON, CommandYaml2Toml, CommandYaml2XML,
		},
		Short: "Convert data format, support csv, json, toml, xml, yaml",
		RunE:  convertFlag.RunE,
		Example: common.Examples(`# Convert yaml to json
-i input.yaml -o output.json`, CommandConvert, CommandYaml2JSON),
	}

	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&convertFlag.inFile, "in", "i", "", common.Usage("Input file (required)"))
	convertCmd.Flags().StringVarP(&convertFlag.outFile, "out", "o", "", common.Usage("Output file"))
}

type ConvertFlag struct {
	inFile  string
	inType  string
	outFile string
	outType string
}

func (c *ConvertFlag) RunE(cmd *cobra.Command, args []string) error {
	if !validator.ValidFile(c.inFile) {
		return common.ErrInvalidFlag
	}
	slice := strings.Split(args[0], "2")
	c.inType = slice[0]
	c.outType = slice[1]
	if c.outFile == "" {
		dir, filename := filepath.Split(c.inFile)
		c.outFile = filepath.Join(dir, strings.Replace(filename, filepath.Ext(filename), "."+slice[1], 1))
	}
	return c.Convert()
}

func (c *ConvertFlag) Convert() error {
	node, err := dasel.NewFromFile(c.inFile, c.inType)
	if err != nil {
		return err
	}
	return node.WriteToFile(c.outFile, c.outType, []storage.ReadWriteOption{
		storage.PrettyPrintOption(true),
	})
}
