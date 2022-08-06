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
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert data format",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var convertSubCmdJSON2TOML = &cobra.Command{
	Use:   "json2toml",
	Short: "Convert json to toml format",
	Run:   cf.Run,
}

var convertSubCmdJSON2YAML = &cobra.Command{
	Use:   "json2yaml",
	Short: "Convert json to yaml format",
	Run:   cf.Run,
}

var convertSubCmdTOML2JSON = &cobra.Command{
	Use:   "toml2json",
	Short: "Convert toml to json format",
	Run:   cf.Run,
}

var convertSubCmdTOML2YAML = &cobra.Command{
	Use:   "toml2yaml",
	Short: "Convert toml to yaml format",
	Run:   cf.Run,
}

var convertSubCmdYAML2JSON = &cobra.Command{
	Use:   "yaml2json",
	Short: "Convert yaml to json format",
	Run:   cf.Run,
	Example: Examples(`# Convert yaml to json
ops-cli convert yaml2json -i input.yaml -o output.json`),
}

var convertSubCmdYAML2TOML = &cobra.Command{
	Use:   "yaml2toml",
	Short: "Convert yaml to toml format",
	Run:   cf.Run,
}

var cf convertFlag

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringVarP(&cf.inFile, "in", "i", "", "Input file")
	convertCmd.PersistentFlags().StringVarP(&cf.outFile, "out", "o", "", "Output file")

	convertCmd.AddCommand(convertSubCmdJSON2TOML)
	convertCmd.AddCommand(convertSubCmdJSON2YAML)
	convertCmd.AddCommand(convertSubCmdTOML2JSON)
	convertCmd.AddCommand(convertSubCmdTOML2YAML)
	convertCmd.AddCommand(convertSubCmdYAML2JSON)
	convertCmd.AddCommand(convertSubCmdYAML2TOML)
}

type convertFlag struct {
	inFile  string
	inType  string
	outFile string
	outType string

	readFile      []byte
	unmarshalData interface{}
}

func (c *convertFlag) Load() {
	var err error
	if c.readFile, err = os.ReadFile(c.inFile); err != nil {
		log.Println(err)
		os.Exit(0)
	}
}

func (c *convertFlag) ParseJSON() {
	c.Load()
	if err := json.Unmarshal(c.readFile, &c.unmarshalData); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.unmarshalData = c.Convert(c.unmarshalData)
}

func (c *convertFlag) ParseTOML() {
	c.Load()
	if err := toml.Unmarshal(c.readFile, &c.unmarshalData); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.unmarshalData = c.Convert(c.unmarshalData)
}

func (c *convertFlag) ParseYAML() {
	c.Load()
	if err := yaml.Unmarshal(c.readFile, &c.unmarshalData); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.unmarshalData = c.Convert(c.unmarshalData)
}

func (c convertFlag) ToJSON() {
	out, err := json.MarshalIndent(c.unmarshalData, "", "  ")
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.WriteFile(out)
}

func (c convertFlag) ToTOML() {
	out, err := toml.Marshal(c.unmarshalData)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.WriteFile(out)
}

func (c convertFlag) ToYAML() {
	out, err := yaml.Marshal(c.unmarshalData)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	c.WriteFile(out)
}

func (c convertFlag) Convert(i interface{}) interface{} {
	switch val := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range val {
			m2[k.(string)] = c.Convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range val {
			val[i] = c.Convert(v)
		}
	}
	return i
}

func (c *convertFlag) WriteFile(content []byte) {
	if err := os.WriteFile(c.outFile, content, os.ModePerm); err != nil {
		log.Println(err)
	}
}

func (c *convertFlag) Run(cmd *cobra.Command, _ []string) {
	if !ValidFile(c.inFile) || c.outFile == "" {
		_ = cmd.Help()
		return
	}
	slice := strings.Split(cmd.Name(), "2")
	if len(slice) != 2 {
		_ = cmd.Help()
		return
	}
	c.inType = slice[0]
	c.outType = slice[1]

	switch c.inType {
	case FileTypeJSON:
		c.ParseJSON()
	case FileTypeTOML:
		c.ParseTOML()
	case FileTypeYAML:
		c.ParseYAML()
	}

	switch c.outType {
	case FileTypeJSON:
		c.ToJSON()
	case FileTypeTOML:
		c.ToTOML()
	case FileTypeYAML:
		c.ToYAML()
	}
}
