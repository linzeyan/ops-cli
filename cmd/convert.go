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
	"path/filepath"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

var convertCmd = &cobra.Command{
	Use:   common.CommandConvert,
	Short: "Convert data format",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

/* CSV. */
var convertSubCmdCSV2JSON = &cobra.Command{
	Use:   "csv2json",
	Short: "Convert csv to json format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdCSV2TOML = &cobra.Command{
	Use:   "csv2toml",
	Short: "Convert csv to toml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdCSV2XML = &cobra.Command{
	Use:   "csv2xml",
	Short: "Convert csv to xml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdCSV2YAML = &cobra.Command{
	Use:   "csv2yaml",
	Short: "Convert csv to yaml format",
	Run:   convertCmdGlobalVar.Run,
}

/* DOS. */
var convertSubCmdDOS2Unix = &cobra.Command{
	Use:   "dos2unix [file...]",
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
	DisableFlagParsing:    true,
}

/* JSON. */
var convertSubCmdJSON2CSV = &cobra.Command{
	Use:   "json2csv",
	Short: "Convert json to csv format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdJSON2TOML = &cobra.Command{
	Use:   "json2toml",
	Short: "Convert json to toml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdJSON2XML = &cobra.Command{
	Use:   "json2xml",
	Short: "Convert json to xml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdJSON2YAML = &cobra.Command{
	Use:   "json2yaml",
	Short: "Convert json to yaml format",
	Run:   convertCmdGlobalVar.Run,
}

/* TOML. */
var convertSubCmdTOML2CSV = &cobra.Command{
	Use:   "toml2csv",
	Short: "Convert toml to csv format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdTOML2JSON = &cobra.Command{
	Use:   "toml2json",
	Short: "Convert toml to json format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdTOML2XML = &cobra.Command{
	Use:   "toml2xml",
	Short: "Convert toml to xml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdTOML2YAML = &cobra.Command{
	Use:   "toml2yaml",
	Short: "Convert toml to yaml format",
	Run:   convertCmdGlobalVar.Run,
}

/* XML. */
var convertSubCmdXML2CSV = &cobra.Command{
	Use:   "xml2csv",
	Short: "Convert xml to csv format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdXML2JSON = &cobra.Command{
	Use:   "xml2json",
	Short: "Convert xml to json format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdXML2TOML = &cobra.Command{
	Use:   "xml2toml",
	Short: "Convert xml to toml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdXML2YAML = &cobra.Command{
	Use:   "xml2yaml",
	Short: "Convert xml to yaml format",
	Run:   convertCmdGlobalVar.Run,
}

/* YAML. */
var convertSubCmdYAML2CSV = &cobra.Command{
	Use:   "yaml2csv",
	Short: "Convert yaml to csv format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdYAML2JSON = &cobra.Command{
	Use:   "yaml2json",
	Short: "Convert yaml to json format",
	Run:   convertCmdGlobalVar.Run,
	Example: common.Examples(`# Convert yaml to json
-i input.yaml -o output.json`, common.CommandConvert, common.SubCommandYaml2JSON),
}
var convertSubCmdYAML2TOML = &cobra.Command{
	Use:   "yaml2toml",
	Short: "Convert yaml to toml format",
	Run:   convertCmdGlobalVar.Run,
}
var convertSubCmdYAML2XML = &cobra.Command{
	Use:   "yaml2xml",
	Short: "Convert yaml to xml format",
	Run:   convertCmdGlobalVar.Run,
}

var convertCmdGlobalVar ConvertFlag

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringVarP(&convertCmdGlobalVar.inFile, "in", "i", "", common.Usage("Input file"))
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

func (c *ConvertFlag) Run(cmd *cobra.Command, _ []string) {
	if !validator.ValidFile(c.inFile) {
		log.Println(`Error: required flag(s) "in" not set`)
		os.Exit(1)
	}
	slice := strings.Split(cmd.Name(), "2")
	c.inType = slice[0]
	c.outType = slice[1]
	if c.outFile == "" {
		dir, filename := filepath.Split(c.inFile)
		c.outFile = filepath.Join(dir, strings.Replace(filename, filepath.Ext(filename), "."+slice[1], 1))
	}
	if err := c.Convert(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
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
