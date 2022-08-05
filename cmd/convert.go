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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}
		if args[0] == "" {
			_ = cmd.Help()
			return
		}
		slice := strings.Split(strings.ToLower(args[0]), "2")
		if len(slice) != 2 {
			_ = cmd.Help()
			return
		}
		cf.Select(slice)
	},
	Example: Examples(`# Convert yaml to json
ops-cli convert yaml2json -i input.yaml -o output.json`),
}

var cf convertFlag

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&cf.inFile, "in", "i", "", "Input file")
	convertCmd.Flags().StringVarP(&cf.outFile, "out", "o", "", "Output file")
}

type convertFlag struct {
	inFile  string
	inType  fileSelector
	outFile string
	outType fileSelector

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

func (c convertFlag) Select(slice []string) {
	cf.inType = fileSelector(slice[0])
	cf.outType = fileSelector(slice[1])

	switch cf.inType {
	case fileJSON:
		cf.ParseJSON()
	case fileTOML:
		cf.ParseTOML()
	case fileYAML:
		cf.ParseYAML()
	default:
		log.Println("Input file format not support")
		os.Exit(0)
	}

	switch cf.outType {
	case fileJSON:
		cf.ToJSON()
	case fileTOML:
		cf.ToTOML()
	case fileYAML:
		cf.ToYAML()
	default:
		log.Println("Output file format not support")
		os.Exit(0)
	}
}

func (c convertFlag) WriteFile(content []byte) {
	if err := os.WriteFile(c.outFile, content, 0600); err != nil {
		log.Println(err)
	}
}
