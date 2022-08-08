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
	Args:  cobra.NoArgs,
	Short: "Generate random string",
	Run:   r.Run,
	Example: Examples(`# Generate a random string
ops-cli random

# Generate a random string of length 32
ops-cli random -l 32

# Generate a random string of length 32 consisting of 10 symbols, 10 lowercase letters, 10 uppercase letters, 2 numbers
ops-cli random -l 32 -s 10 -o 10 -u 10 -n 2`),
}

var randomSubCmdLower = &cobra.Command{
	Use:   "lowercase",
	Short: "Generate a string consisting of lowercase letters",
	Run:   r.Run,
	Example: Examples(` Generate a random string of lowercase letters
ops-cli random lowercase`),
}

var randomSubCmdNumber = &cobra.Command{
	Use:   "number",
	Short: "Generate a string consisting of numbers",
	Run:   r.Run,
	Example: Examples(`Generate a random string of numbers of length 100
ops-cli random number -l 100`),
}

var randomSubCmdSymbol = &cobra.Command{
	Use:   "symbol",
	Short: "Generate a string consisting of symbols",
	Run:   r.Run,
	Example: Examples(`# Generate a random string of symbols
ops-cli random symbol`),
}

var randomSubCmdUpper = &cobra.Command{
	Use:   "uppercase",
	Short: "Generate a string consisting of uppercase letters",
	Run:   r.Run,
	Example: Examples(`# Generate a random string of uppercase letters
ops-cli random uppercase`),
}

var r ran

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().UintVarP(&r.length, "length", "l", 24, "Specify the string length")
	randomCmd.Flags().UintVarP(&r.lower, "lower", "o", 4, "Number of lowercase letters to include in the string")
	randomCmd.Flags().UintVarP(&r.upper, "upper", "u", 4, "Number of uppercase letters to include in the string")
	randomCmd.Flags().UintVarP(&r.symbol, "symbol", "s", 4, "Number of symbols to include in the string")
	randomCmd.Flags().UintVarP(&r.number, "number", "n", 4, "Number of digits to include in the string")

	randomCmd.AddCommand(randomSubCmdLower)
	randomCmd.AddCommand(randomSubCmdNumber)
	randomCmd.AddCommand(randomSubCmdSymbol)
	randomCmd.AddCommand(randomSubCmdUpper)
}

type ran struct {
	/* Bind flags */
	length, lower, upper, symbol, number uint

	/* Output string */
	result string
}

func (r *ran) Run(cmd *cobra.Command, _ []string) {
	switch cmd.Name() {
	case "number":
		r.result = password.GenNumber(r.length)
	case "symbol":
		r.result = password.GenSymbol(r.length)
	case "uppercase":
		r.result = password.GenUpper(r.length)
	case "lowercase":
		r.result = password.GenLower(r.length)
	case "random":
		r.result = password.GeneratePassword(r.length, r.lower, r.upper, r.symbol, r.number)
	}
	fmt.Println(r.result)
}
