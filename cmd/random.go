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

	"github.com/linzeyan/password"
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate random string",
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
