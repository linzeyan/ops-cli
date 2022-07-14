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
		switch randomMode {
		case "all":
			result = password.GenAll(randomLength)
		case "num", "number":
			result = password.GenNumber(randomLength)
		case "sym", "symbol":
			result = password.GenSymbol(randomLength)
		case "upper", "uppercase":
			result = password.GenUpper(randomLength)
		case "lower", "lowercase":
			result = password.GenLower(randomLength)
		default:
			result = password.GeneratePassword(randomLength, randomLower, randomUpper, randomSymbol, randomNumber)
		}
		fmt.Println(result)
	},
	Example: `ops-cli random
ops-cli random -l 32
ops-cli random -l 32 -s 10 -o 10 -u 10 -n 2
ops-cli random -t number -l 100
ops-cli random -t symbol
ops-cli random -t upper
ops-cli random -t lower
ops-cli random -t all`,
}

var randomLength, randomLower, randomUpper, randomSymbol, randomNumber uint
var randomMode string

func init() {
	rootCmd.AddCommand(randomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// randomCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// randomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	randomCmd.Flags().UintVarP(&randomLength, "length", "l", 24, "Specify the string length")
	randomCmd.Flags().UintVarP(&randomLower, "lower", "o", 4, "Number of lowercase letters to include in the string")
	randomCmd.Flags().UintVarP(&randomUpper, "upper", "u", 4, "Number of uppercase letters to include in the string")
	randomCmd.Flags().UintVarP(&randomSymbol, "symbol", "s", 4, "Number of symbols to include in the string")
	randomCmd.Flags().UintVarP(&randomNumber, "number", "n", 4, "Number of digits to include in the string")
	randomCmd.Flags().StringVarP(&randomMode, "type", "t", "default", "Specifies the string type, which can be number, symbol, upper(case), lower(case), or all but unspecified number of characters")
}
