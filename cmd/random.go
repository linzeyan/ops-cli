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

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate random string",
	Args:  cobra.OnlyValidArgs,
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
	Example: Examples(`# Generate a random string
ops-cli random

# Generate a random string of length 32
ops-cli random -l 32

# Generate a random string of length 32 consisting of 10 symbols, 10 lowercase letters, 10 uppercase letters, 2 numbers
ops-cli random -l 32 -s 10 -o 10 -u 10 -n 2

# Generate a random string of numbers of length 100
ops-cli random -t number -l 100

# Generate a random string of symbols
ops-cli random -t symbol

# Generate a random string of uppercase letters
ops-cli random -t upper

# Generate a random string of lowercase letters
ops-cli random -t lower

# Generate a random string
ops-cli random -t all`),
}

var randomLength, randomLower, randomUpper, randomSymbol, randomNumber uint
var randomMode string

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.Flags().UintVarP(&randomLength, "length", "l", 24, "Specify the string length")
	randomCmd.Flags().UintVarP(&randomLower, "lower", "o", 4, "Number of lowercase letters to include in the string")
	randomCmd.Flags().UintVarP(&randomUpper, "upper", "u", 4, "Number of uppercase letters to include in the string")
	randomCmd.Flags().UintVarP(&randomSymbol, "symbol", "s", 4, "Number of symbols to include in the string")
	randomCmd.Flags().UintVarP(&randomNumber, "number", "n", 4, "Number of digits to include in the string")
	randomCmd.Flags().StringVarP(&randomMode, "type", "t", "default", "Specifies the string type, which can be number, symbol, upper(case), lower(case), or all but unspecified number of characters")
}
