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
	"crypto/rand"
	"log"
	"math/big"
	mathRand "math/rand"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Args:  cobra.NoArgs,
	Short: "Generate random string",
	Run:   randomCmdGlobalVar.Run,
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
	Run:   randomCmdGlobalVar.Run,
	Example: Examples(` Generate a random string of lowercase letters
ops-cli random lowercase`),
}

var randomSubCmdNumber = &cobra.Command{
	Use:   "number",
	Short: "Generate a string consisting of numbers",
	Run:   randomCmdGlobalVar.Run,
	Example: Examples(`Generate a random string of numbers of length 100
ops-cli random number -l 100`),
}

var randomSubCmdSymbol = &cobra.Command{
	Use:   "symbol",
	Short: "Generate a string consisting of symbols",
	Run:   randomCmdGlobalVar.Run,
	Example: Examples(`# Generate a random string of symbols
ops-cli random symbol`),
}

var randomSubCmdUpper = &cobra.Command{
	Use:   "uppercase",
	Short: "Generate a string consisting of uppercase letters",
	Run:   randomCmdGlobalVar.Run,
	Example: Examples(`# Generate a random string of uppercase letters
ops-cli random uppercase`),
}

var randomCmdGlobalVar RandomFlag

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().IntVarP(&randomCmdGlobalVar.length, "length", "l", 24, "Specify the string length")
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.lower, "lower", "o", 4, "Number of lowercase letters to include in the string")
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.upper, "upper", "u", 4, "Number of uppercase letters to include in the string")
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.symbol, "symbol", "s", 4, "Number of symbols to include in the string")
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.number, "number", "n", 4, "Number of digits to include in the string")

	randomCmd.AddCommand(randomSubCmdLower)
	randomCmd.AddCommand(randomSubCmdNumber)
	randomCmd.AddCommand(randomSubCmdSymbol)
	randomCmd.AddCommand(randomSubCmdUpper)
}

type RandomFlag struct {
	/* Bind flags */
	length, lower, upper, symbol, number int

	/* Output string */
	result string
}

func (r *RandomFlag) Run(cmd *cobra.Command, _ []string) {
	var err error
	var p RandomString
	switch cmd.Name() {
	case "number":
		randomCmdGlobalVar.result, err = p.GenerateString(randomCmdGlobalVar.length, Numbers)
	case "symbol":
		randomCmdGlobalVar.result, err = p.GenerateString(randomCmdGlobalVar.length, Symbols)
	case "uppercase":
		randomCmdGlobalVar.result, err = p.GenerateString(randomCmdGlobalVar.length, UppercaseLetters)
	case "lowercase":
		randomCmdGlobalVar.result, err = p.GenerateString(randomCmdGlobalVar.length, LowercaseLetters)
	case "random":
		randomCmdGlobalVar.result, err = p.GenerateAll(randomCmdGlobalVar.length, randomCmdGlobalVar.lower, randomCmdGlobalVar.upper, randomCmdGlobalVar.symbol, randomCmdGlobalVar.number)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(randomCmdGlobalVar.result)
}

type RandomString struct{}

func (RandomString) GenerateString(length int, charSet RandomCharacter) (string, error) {
	var s strings.Builder
	var err error
	for i := int(0); i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		if err = s.WriteByte(charSet[n.Int64()]); err != nil {
			return "", err
		}
	}
	return s.String(), err
}

func (r RandomString) GenerateAll(length, minLower, minUpper, minSymbol, minNumber int) (string, error) {
	var err error
	var remain string
	leave := length - minLower - minUpper - minSymbol - minNumber
	if leave < 0 {
		return "", ErrInvalidLength
	}
	lower, err := r.GenerateString(minLower, LowercaseLetters)
	if err != nil {
		return lower, err
	}
	upper, err := r.GenerateString(minUpper, UppercaseLetters)
	if err != nil {
		return upper, err
	}
	symbol, err := r.GenerateString(minSymbol, Symbols)
	if err != nil {
		return symbol, err
	}
	num, err := r.GenerateString(minNumber, Numbers)
	if err != nil {
		return num, err
	}
	if leave != 0 {
		remain, err = r.GenerateString(leave, AllSet)
		if err != nil {
			return remain, err
		}
	}
	result := []byte(lower + upper + symbol + num + remain)
	mathRand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return string(result), err
}
