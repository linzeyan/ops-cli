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
	"path/filepath"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

var convertCmd = &cobra.Command{
	Use:   CommandConvert,
	Short: "Convert data format",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

/* CSV. */
var convertSubCmdCSV2JSON = &cobra.Command{
	Use:   CommandCsv2JSON,
	Short: "Convert csv to json format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdCSV2TOML = &cobra.Command{
	Use:   CommandCsv2Toml,
	Short: "Convert csv to toml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdCSV2XML = &cobra.Command{
	Use:   CommandCsv2XML,
	Short: "Convert csv to xml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdCSV2YAML = &cobra.Command{
	Use:   CommandCsv2Yaml,
	Short: "Convert csv to yaml format",
	RunE:  convertCmdGlobalVar.RunE,
}

/* DOS. */
var convertSubCmdDOS2Unix = &cobra.Command{
	Use:   CommandDos2Unix + " file...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Convert DOS to Unix format",
	Run: func(_ *cobra.Command, args []string) {
		for _, f := range args {
			if err := common.Dos2Unix(f); err != nil {
				log.Printf("%s: %v\n", f, err)
			}
		}
	},
	DisableFlagsInUseLine: true,
}

/* JSON. */
var convertSubCmdJSON2CSV = &cobra.Command{
	Use:   CommandJSON2Csv,
	Short: "Convert json to csv format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdJSON2TOML = &cobra.Command{
	Use:   CommandJSON2Toml,
	Short: "Convert json to toml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdJSON2XML = &cobra.Command{
	Use:   CommandJSON2XML,
	Short: "Convert json to xml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdJSON2YAML = &cobra.Command{
	Use:   CommandJSON2Yaml,
	Short: "Convert json to yaml format",
	RunE:  convertCmdGlobalVar.RunE,
}

/* TOML. */
var convertSubCmdTOML2CSV = &cobra.Command{
	Use:   CommandToml2Csv,
	Short: "Convert toml to csv format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdTOML2JSON = &cobra.Command{
	Use:   CommandToml2JSON,
	Short: "Convert toml to json format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdTOML2XML = &cobra.Command{
	Use:   CommandToml2XML,
	Short: "Convert toml to xml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdTOML2YAML = &cobra.Command{
	Use:   CommandToml2Yaml,
	Short: "Convert toml to yaml format",
	RunE:  convertCmdGlobalVar.RunE,
}

/* XML. */
var convertSubCmdXML2CSV = &cobra.Command{
	Use:   CommandXML2Csv,
	Short: "Convert xml to csv format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdXML2JSON = &cobra.Command{
	Use:   CommandXML2JSON,
	Short: "Convert xml to json format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdXML2TOML = &cobra.Command{
	Use:   CommandXML2Toml,
	Short: "Convert xml to toml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdXML2YAML = &cobra.Command{
	Use:   CommandXML2Yaml,
	Short: "Convert xml to yaml format",
	RunE:  convertCmdGlobalVar.RunE,
}

/* YAML. */
var convertSubCmdYAML2CSV = &cobra.Command{
	Use:   CommandYaml2Csv,
	Short: "Convert yaml to csv format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdYAML2JSON = &cobra.Command{
	Use:   CommandYaml2JSON,
	Short: "Convert yaml to json format",
	RunE:  convertCmdGlobalVar.RunE,
	Example: common.Examples(`# Convert yaml to json
-i input.yaml -o output.json`, CommandConvert, CommandYaml2JSON),
}
var convertSubCmdYAML2TOML = &cobra.Command{
	Use:   CommandYaml2Toml,
	Short: "Convert yaml to toml format",
	RunE:  convertCmdGlobalVar.RunE,
}
var convertSubCmdYAML2XML = &cobra.Command{
	Use:   CommandYaml2XML,
	Short: "Convert yaml to xml format",
	RunE:  convertCmdGlobalVar.RunE,
}

var convertCmdGlobalVar ConvertFlag

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringVarP(&convertCmdGlobalVar.inFile, "in", "i", "", common.Usage("Input file (required)"))
	convertCmd.PersistentFlags().StringVarP(&convertCmdGlobalVar.outFile, "out", "o", "", common.Usage("Output file"))
	/* dos2unix */
	convertCmd.AddCommand(convertSubCmdDOS2Unix)
	/* CSV */
	convertCmd.AddCommand(convertSubCmdCSV2JSON, convertSubCmdCSV2TOML, convertSubCmdCSV2XML, convertSubCmdCSV2YAML)
	/* JSON */
	convertCmd.AddCommand(convertSubCmdJSON2CSV, convertSubCmdJSON2TOML, convertSubCmdJSON2XML, convertSubCmdJSON2YAML)
	/* TOML */
	convertCmd.AddCommand(convertSubCmdTOML2CSV, convertSubCmdTOML2JSON, convertSubCmdTOML2XML, convertSubCmdTOML2YAML)
	/* XML */
	convertCmd.AddCommand(convertSubCmdXML2CSV, convertSubCmdXML2JSON, convertSubCmdXML2TOML, convertSubCmdXML2YAML)
	/* YAML */
	convertCmd.AddCommand(convertSubCmdYAML2CSV, convertSubCmdYAML2JSON, convertSubCmdYAML2TOML, convertSubCmdYAML2XML)
}

type ConvertFlag struct {
	inFile  string
	inType  string
	outFile string
	outType string
}

func (c *ConvertFlag) RunE(cmd *cobra.Command, _ []string) error {
	if !validator.ValidFile(c.inFile) {
		return common.ErrInvalidFlag
	}
	slice := strings.Split(cmd.Name(), "2")
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
