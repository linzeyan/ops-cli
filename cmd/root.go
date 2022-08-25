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
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var (
	ErrArgNotFound   = errors.New("argument not found")
	ErrFileNotFound  = errors.New("file not found")
	ErrFileType      = errors.New("file type not correct")
	ErrInitialFailed = errors.New("initial failed")
	ErrInvalidIP     = errors.New("invalid IP")
	ErrInvalidVar    = errors.New("invalid variable")
	ErrParseCert     = errors.New("can not correctly parse certificate")
	ErrTokenNotFound = errors.New("token not found")
)

var rootCmd = &cobra.Command{
	Use:   common.CommandRoot,
	Short: "OPS useful tools",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

/* Flags. */
var (
	rootConfig     string
	rootOutputJSON bool
	rootOutputYAML bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootOutputJSON, "json", "j", false, common.Usage("Output JSON format"))
	rootCmd.PersistentFlags().BoolVarP(&rootOutputYAML, "yaml", "y", false, common.Usage("Output YAML format"))
	rootCmd.PersistentFlags().StringVar(&rootConfig, "config", "", common.Usage("Specify config path (toml)"))
}

type OutputFormat interface {
	String()
}

func OutputInterfaceString(r OutputFormat) {
	switch {
	case rootOutputJSON:
		PrintJSON(r)
	case rootOutputYAML:
		PrintYAML(r)
	default:
		r.String()
	}
}

func OutputDefaultJSON(i any) {
	if rootOutputYAML {
		PrintYAML(i)
	} else {
		PrintJSON(i)
	}
}

func OutputDefaultNone(i any) {
	if rootOutputJSON {
		PrintJSON(i)
	} else if rootOutputYAML {
		PrintYAML(i)
	}
}

func OutputDefaultString(i any) {
	switch {
	case rootOutputJSON:
		PrintJSON(i)
	case rootOutputYAML:
		PrintYAML(i)
	default:
		PrintString(i)
	}
}

func OutputDefaultYAML(i any) {
	if rootOutputJSON {
		PrintJSON(i)
	} else {
		PrintYAML(i)
	}
}

func PrintJSON(i any) {
	out, err := Encoder.JSONEncode(i)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintYAML(i any) {
	out, err := Encoder.YamlEncode(i)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintString(i any) {
	switch data := i.(type) {
	case string, []byte:
		fmt.Printf("%s\n", data)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		fmt.Printf("%d\n", data)
	default:
		fmt.Printf("%v\n", data)
	}
}
