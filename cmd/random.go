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
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"regexp"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Args:  cobra.NoArgs,
	Short: "Generate random string",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`# Generate a random string
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
	Example: common.Examples(` Generate a random string of lowercase letters
ops-cli random lowercase`),
}

var randomSubCmdNumber = &cobra.Command{
	Use:   "number",
	Short: "Generate a string consisting of numbers",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`Generate a random string of numbers of length 100
ops-cli random number -l 100`),
}

var randomSubCmdSymbol = &cobra.Command{
	Use:   "symbol",
	Short: "Generate a string consisting of symbols",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`# Generate a random string of symbols
ops-cli random symbol`),
}

var randomSubCmdUpper = &cobra.Command{
	Use:   "uppercase",
	Short: "Generate a string consisting of uppercase letters",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`# Generate a random string of uppercase letters
ops-cli random uppercase`),
}

var randomSubCmdBootstrap = &cobra.Command{
	Use:   "bootstrap-token",
	Short: "Generate a bootstrap token",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`# Generate a bootstrap token
ops-cli random bootstrap-token`),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
}

var randomSubCmdBase64 = &cobra.Command{
	Use:   common.Base64,
	Short: "Generate a base64 string",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`# Generate a base64 string
ops-cli random base64 -l 100`),
}
var randomSubCmdHex = &cobra.Command{
	Use:   common.Hex,
	Short: "Generate a hexadecimal string",
	Run:   randomCmdGlobalVar.Run,
	Example: common.Examples(`# Generate a hexadecimal string
ops-cli random hex -l 30`),
}

var randomCmdGlobalVar RandomFlag

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().IntVarP(&randomCmdGlobalVar.length, "length", "l", 24, common.Usage("Specify the string length"))
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.lower, "lower", "o", 4, common.Usage("Number of lowercase letters to include in the string"))
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.upper, "upper", "u", 4, common.Usage("Number of uppercase letters to include in the string"))
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.symbol, "symbol", "s", 4, common.Usage("Number of symbols to include in the string"))
	randomCmd.Flags().IntVarP(&randomCmdGlobalVar.number, "number", "n", 4, common.Usage("Number of digits to include in the string"))

	randomCmd.AddCommand(randomSubCmdLower)
	randomCmd.AddCommand(randomSubCmdNumber)
	randomCmd.AddCommand(randomSubCmdSymbol)
	randomCmd.AddCommand(randomSubCmdUpper)
	randomCmd.AddCommand(randomSubCmdBase64)
	randomCmd.AddCommand(randomSubCmdHex)
	randomCmd.AddCommand(randomSubCmdBootstrap)
}

type RandomFlag struct {
	/* Bind flags */
	length, lower, upper, symbol, number int
}

func (r *RandomFlag) Run(cmd *cobra.Command, _ []string) {
	var p RandomString
	switch cmd.Name() {
	case "number":
		p = p.GenerateString(randomCmdGlobalVar.length, Numbers)
	case "symbol":
		p = p.GenerateString(randomCmdGlobalVar.length, Symbols)
	case "uppercase":
		p = p.GenerateString(randomCmdGlobalVar.length, UppercaseLetters)
	case "lowercase":
		p = p.GenerateString(randomCmdGlobalVar.length, LowercaseLetters)
	case "random":
		p = p.GenerateAll(randomCmdGlobalVar.length, randomCmdGlobalVar.lower, randomCmdGlobalVar.upper, randomCmdGlobalVar.symbol, randomCmdGlobalVar.number)
	case "bootstrap-token":
		r1 := p.Rand(3)
		r2 := p.Rand(8)
		id, _ := Encoder.HexEncode(r1)
		token, _ := Encoder.HexEncode(r2)
		PrintString(id + "." + token)
		return
	case common.Base64:
		b := p.Rand(randomCmdGlobalVar.length)
		encode, _ := Encoder.PemEncode(b)
		re := regexp.MustCompile("-.*-\n")
		out := re.ReplaceAllString(encode, "")
		PrintString(out)
		return
	case common.Hex:
		b := p.Rand(randomCmdGlobalVar.length)
		out, _ := Encoder.HexEncode(b)
		PrintString(out)
		return
	}
	PrintString(p.String())
}

type RandomCharacter string

const (
	LowercaseLetters RandomCharacter = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters RandomCharacter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Symbols          RandomCharacter = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	Numbers          RandomCharacter = "0123456789"
	AllSet           RandomCharacter = LowercaseLetters + UppercaseLetters + Symbols + Numbers
)

type RandomString []byte

func (RandomString) GenerateString(length int, charSet RandomCharacter) RandomString {
	if length <= 0 {
		return nil
	}
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return nil
		}
		if err = buf.WriteByte(charSet[n.Int64()]); err != nil {
			return nil
		}
	}
	return buf.Bytes()
}

func (r RandomString) GenerateAll(length, minLower, minUpper, minSymbol, minNumber int) RandomString {
	if length <= 0 {
		return nil
	}
	var result RandomString
	if minLower < 0 {
		minLower = 0
	}
	if minUpper < 0 {
		minUpper = 0
	}
	if minSymbol < 0 {
		minSymbol = 0
	}
	if minNumber < 0 {
		minNumber = 0
	}
	leave := length - minLower - minUpper - minSymbol - minNumber
	if leave < 0 {
		return nil
	}
	lower := r.GenerateString(minLower, LowercaseLetters)
	result = append(result, lower...)

	upper := r.GenerateString(minUpper, UppercaseLetters)
	result = append(result, upper...)

	symbol := r.GenerateString(minSymbol, Symbols)
	result = append(result, symbol...)

	num := r.GenerateString(minNumber, Numbers)
	result = append(result, num...)

	if leave != 0 {
		remain := r.GenerateString(leave, AllSet)
		result = append(result, remain...)
	}
	mathRand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

func (r RandomString) Rand(length int) []byte {
	if length <= 0 {
		return nil
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}

func (r RandomString) String() string {
	if r == nil {
		return "<nil>"
	}
	return string(r)
}

func (r RandomString) Bytes() []byte {
	return []byte(r)
}
