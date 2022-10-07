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
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var RootCmd = root()

/* Flags. */
var (
	rootConfig     string
	rootOutputJSON bool
	rootOutputYAML bool
)

func root() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   common.RepoName,
		Short: "OPS useful tools",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}
	rootCmd.PersistentFlags().BoolVarP(&rootOutputJSON, "json", "j", false, common.Usage("Output JSON format"))
	rootCmd.PersistentFlags().BoolVarP(&rootOutputYAML, "yaml", "y", false, common.Usage("Output YAML format"))
	rootCmd.PersistentFlags().StringVar(&rootConfig, "config", "", common.Usage("Specify config path"))

	if !validator.IsWindows() {
		rootCmd.AddCommand(initArping())
	}
	rootCmd.AddCommand(initCert(), initConvert())
	rootCmd.AddCommand(initDate(), initDF(), initDig(), initDiscord(), initDoc(rootCmd), initDos2Unix())
	rootCmd.AddCommand(initEncode(), initEncrypt())
	rootCmd.AddCommand(initFree())
	rootCmd.AddCommand(initGeoip())
	rootCmd.AddCommand(initHash())
	rootCmd.AddCommand(initIcp(), initIP())
	rootCmd.AddCommand(initLINE())
	rootCmd.AddCommand(initMtr())
	rootCmd.AddCommand(initNetmask())
	rootCmd.AddCommand(initOtp())
	rootCmd.AddCommand(initPing())
	rootCmd.AddCommand(initQrcode())
	rootCmd.AddCommand(initRandom(), initReadlink(), initRedis())
	rootCmd.AddCommand(initSlack(), initSSHKeyGen(), initStat(), initSystem())
	rootCmd.AddCommand(initTcping(), initTelegram(), initTraceroute(), initTree())
	rootCmd.AddCommand(initUpdate(), initURL())
	rootCmd.AddCommand(initVersion())
	rootCmd.AddCommand(initWhois(), initWsping())
	return rootCmd
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
		PrintString(err)
		os.Exit(1)
	}
	PrintString(out)
}

func PrintYAML(i any) {
	out, err := Encoder.YamlEncode(i)
	if err != nil {
		PrintString(err)
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
	case map[int]string:
		b := new(bytes.Buffer)
		for key, value := range data {
			fmt.Fprintf(b, "%d: %s", key, value)
		}
		fmt.Printf("%s\n", b.String())
	case map[string]string:
		b := new(bytes.Buffer)
		for key, value := range data {
			fmt.Fprintf(b, "%s: %s", key, value)
		}
		fmt.Printf("%s\n", b.String())
	default:
		fmt.Printf("%v\n", data)
	}
}

func PrintTable(header []string, data [][]string, align int, padding string, format bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(format)
	table.SetHeaderAlignment(align)
	table.SetAlignment(align)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding(padding)
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
}

func ReadConfig(block string, flag any) error {
	v := common.Config(rootConfig, strings.ToLower(block))
	return Encoder.JSONMarshaler(v, flag)
}
