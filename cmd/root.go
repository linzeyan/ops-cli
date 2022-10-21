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
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var RootCmd = root()

/* Flags. */
var (
	rootConfig       string
	rootOutputFormat string
)

func root() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   common.RepoName,
		Short: "OPS useful tools",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}
	rootCmd.PersistentFlags().StringVar(&rootOutputFormat, "output", "", common.Usage("Output format, can be json/yaml"))
	rootCmd.PersistentFlags().StringVar(&rootConfig, "config", "", common.Usage("Specify config path"))
	rootCmd.PersistentFlags().BoolP("help", "", false, common.Usage("Help for this command"))

	rootCmd.AddCommand(initArping())
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
	rootCmd.AddCommand(initPing(), initPs())
	rootCmd.AddCommand(initQrcode())
	rootCmd.AddCommand(initRandom(), initReadlink(), initRedis())
	rootCmd.AddCommand(initSlack(), initSs(), initSSHKeyGen(), initStat(), initSystem())
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
	switch rootOutputFormat {
	case CommandJSON:
		PrintJSON(r)
	case CommandYaml:
		PrintYAML(r)
	default:
		r.String()
	}
}

func OutputDefaultJSON(i any) {
	if rootOutputFormat == CommandYaml {
		PrintYAML(i)
	} else {
		PrintJSON(i)
	}
}

func OutputDefaultNone(i any) {
	if rootOutputFormat == CommandJSON {
		PrintJSON(i)
	} else if rootOutputFormat == CommandYaml {
		PrintYAML(i)
	}
}

func OutputDefaultString(i any) {
	switch rootOutputFormat {
	case CommandJSON:
		PrintJSON(i)
	case CommandYaml:
		PrintYAML(i)
	default:
		PrintString(i)
	}
}

func OutputDefaultYAML(i any) {
	if rootOutputFormat == CommandJSON {
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
