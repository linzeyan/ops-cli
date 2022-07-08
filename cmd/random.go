/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/linzeyan/password"
	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate random string",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {
		var result string
		switch mode {
		case "all":
			result = password.GenAll(length)
		case "number":
			result = password.GenNumber(length)
		case "symbol":
			result = password.GenSymbol(length)
		case "upper", "uppercase":
			result = password.GenUpper(length)
		case "lower", "lowercase":
			result = password.GenLower(length)
		default:
			result = password.GeneratePassword(length, minLower, minUpper, minSymbol, minNumber)
		}
		fmt.Println(result)
	},
}

var length, minLower, minUpper, minSymbol, minNumber uint
var mode string

func init() {
	rootCmd.AddCommand(randomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// randomCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	randomCmd.Flags().UintVarP(&length, "length", "l", 24, "Specify the string length")
	randomCmd.Flags().UintVarP(&minLower, "lower", "o", 4, "Number of lowercase letters to include in the string")
	randomCmd.Flags().UintVarP(&minUpper, "upper", "u", 4, "Number of uppercase letters to include in the string")
	randomCmd.Flags().UintVarP(&minSymbol, "symbol", "s", 4, "Number of symbols to include in the string")
	randomCmd.Flags().UintVarP(&minNumber, "number", "n", 4, "Number of digits to include in the string")
	randomCmd.Flags().StringVarP(&mode, "type", "t", "default", "Specifies the string type, which can be number, symbol, upper(case), lower(case), or all but unspecified number of characters")
}
