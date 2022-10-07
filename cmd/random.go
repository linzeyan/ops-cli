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

var Randoms Random

func initRandom() *cobra.Command {
	var flags struct {
		/* Bind flags */
		length, lower, upper, symbol, number int
	}

	run := func(cmd *cobra.Command, _ []string) {
		var out []byte
		switch cmd.Name() {
		case CommandNumber:
			out = Randoms.GenerateString(flags.length, Numbers)
		case CommandSymbol:
			out = Randoms.GenerateString(flags.length, Symbols)
		case CommandUppercase:
			out = Randoms.GenerateString(flags.length, UppercaseLetters)
		case CommandLowercase:
			out = Randoms.GenerateString(flags.length, LowercaseLetters)
		case CommandRandom:
			out = Randoms.GenerateAll(flags.length, flags.lower, flags.upper, flags.symbol, flags.number)
		case CommandBootstrap:
			r1 := Randoms.Rand(3)
			r2 := Randoms.Rand(8)
			id, _ := Encoder.HexEncode(r1)
			token, _ := Encoder.HexEncode(r2)
			PrintString(id + "." + token)
			return
		case CommandBase64:
			b := Randoms.Rand(flags.length)
			encode, _ := Encoder.PemEncode(b)
			re := regexp.MustCompile("-.*-\n")
			out := re.ReplaceAllString(encode, "")
			PrintString(out)
			return
		case CommandHex:
			b := Randoms.Rand(flags.length)
			out, _ := Encoder.HexEncode(b)
			PrintString(out)
			return
		}
		PrintString(out)
	}

	var randomCmd = &cobra.Command{
		Use:   CommandRandom,
		Args:  cobra.NoArgs,
		Short: "Generate random string",
		Run:   run,
		Example: common.Examples(`# Generate a random string of length 32
-l 32

# Generate a random string of length 32 consisting of 10 symbols, 10 lowercase letters, 10 uppercase letters, 2 numbers
-l 32 -s 10 -o 10 -u 10 -n 2`, CommandRandom),
	}

	var randomSubCmdLower = &cobra.Command{
		Use:   CommandLowercase,
		Short: "Generate a string consisting of lowercase letters",
		Run:   run,
		Example: common.Examples(` Generate a random string of lowercase letters
-l 10`, CommandRandom, CommandLowercase),
	}

	var randomSubCmdNumber = &cobra.Command{
		Use:   CommandNumber,
		Short: "Generate a string consisting of numbers",
		Run:   run,
		Example: common.Examples(`Generate a random string of numbers of length 100
-l 100`, CommandRandom, CommandNumber),
	}

	var randomSubCmdSymbol = &cobra.Command{
		Use:   CommandSymbol,
		Short: "Generate a string consisting of symbols",
		Run:   run,
	}

	var randomSubCmdUpper = &cobra.Command{
		Use:   CommandUppercase,
		Short: "Generate a string consisting of uppercase letters",
		Run:   run,
	}

	var randomSubCmdBootstrap = &cobra.Command{
		Use:   CommandBootstrap,
		Short: "Generate a bootstrap token",
		Run:   run,

		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
	}

	var randomSubCmdBase64 = &cobra.Command{
		Use:   CommandBase64,
		Short: "Generate a base64 string",
		Run:   run,
		Example: common.Examples(`# Generate a base64 string
-l 100`, CommandRandom, CommandBase64),
	}

	var randomSubCmdHex = &cobra.Command{
		Use:   CommandHex,
		Short: "Generate a hexadecimal string",
		Run:   run,
		Example: common.Examples(`# Generate a hexadecimal string
-l 30`, CommandRandom, CommandHex),
	}
	randomCmd.PersistentFlags().IntVarP(&flags.length, "length", "l", 24, common.Usage("Specify the string length"))
	randomCmd.Flags().IntVarP(&flags.lower, "lower", "o", 4, common.Usage("Number of lowercase letters to include in the string"))
	randomCmd.Flags().IntVarP(&flags.upper, "upper", "u", 4, common.Usage("Number of uppercase letters to include in the string"))
	randomCmd.Flags().IntVarP(&flags.symbol, "symbol", "s", 4, common.Usage("Number of symbols to include in the string"))
	randomCmd.Flags().IntVarP(&flags.number, "number", "n", 4, common.Usage("Number of digits to include in the string"))

	randomCmd.AddCommand(randomSubCmdLower)
	randomCmd.AddCommand(randomSubCmdNumber)
	randomCmd.AddCommand(randomSubCmdSymbol)
	randomCmd.AddCommand(randomSubCmdUpper)
	randomCmd.AddCommand(randomSubCmdBase64)
	randomCmd.AddCommand(randomSubCmdHex)
	randomCmd.AddCommand(randomSubCmdBootstrap)
	return randomCmd
}

type RandomCharacter string

const (
	LowercaseLetters RandomCharacter = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters RandomCharacter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Symbols          RandomCharacter = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	Numbers          RandomCharacter = "0123456789"
	AllSet           RandomCharacter = LowercaseLetters + UppercaseLetters + Symbols + Numbers
)

type Random []byte

func (Random) GenerateString(length int, charSet RandomCharacter) Random {
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

func (r Random) GenerateAll(length, minLower, minUpper, minSymbol, minNumber int) Random {
	if length <= 0 {
		return nil
	}
	var result Random
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

func (Random) Rand(length int) []byte {
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

func (r Random) String() string {
	if r == nil {
		return "<nil>"
	}
	return string(r)
}
