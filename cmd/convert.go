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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

func initConvert() *cobra.Command {
	var flags struct {
		inFile  string
		outFile string
	}
	validArg := []string{
		CommandCsv2JSON, CommandCsv2Toml, CommandCsv2XML, CommandCsv2Yaml,
		CommandJSON2Csv, CommandJSON2Toml, CommandJSON2XML, CommandJSON2Yaml,
		CommandToml2Csv, CommandToml2JSON, CommandToml2XML, CommandToml2Yaml,
		CommandYaml2Csv, CommandYaml2JSON, CommandYaml2Toml, CommandYaml2XML,
	}

	var convertCmd = &cobra.Command{
		Use:       CommandConvert,
		Args:      cobra.OnlyValidArgs,
		ValidArgs: validArg,
		Short:     "Convert data format, support csv, json, toml, xml, yaml",
		Run: func(_ *cobra.Command, args []string) {
			if !common.IsFile(flags.inFile) {
				logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(flags.inFile))
				printer.Error(common.ErrInvalidFlag)
				return
			}
			slice := strings.Split(args[0], "2")
			inType := slice[0]
			outType := slice[1]
			if flags.outFile == "" {
				dir, filename := filepath.Split(flags.inFile)
				flags.outFile = filepath.Join(dir, strings.Replace(filename, filepath.Ext(filename), "."+slice[1], 1))
			}
			node, err := dasel.NewFromFile(flags.inFile, inType)
			if err != nil {
				logger.Debug(err.Error())
				printer.Error(err)
				return
			}
			err = node.WriteToFile(
				flags.outFile,
				outType,
				[]storage.ReadWriteOption{
					storage.PrettyPrintOption(true),
				})
			if err != nil {
				logger.Debug(err.Error())
				printer.Error(err)
				return
			}
		},
		Example: common.Examples(`# Convert yaml to json
-i input.yaml -o output.json`, CommandConvert, CommandYaml2JSON) + `

Available Commands:
` + fmt.Sprintf("  %-10s %-10s %-10s %-10s\n  %-10s %-10s %-10s %-10s\n  %-10s %-10s %-10s %-10s\n  %-10s %-10s %-10s %-10s",
			common.SliceStringToInterface(validArg)...),
	}

	convertCmd.Flags().StringVarP(&flags.inFile, "in", "i", "", common.Usage("Input file (required)"))
	convertCmd.Flags().StringVarP(&flags.outFile, "out", "o", "", common.Usage("Output file"))
	return convertCmd
}
