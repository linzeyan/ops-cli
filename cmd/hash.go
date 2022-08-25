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
	"bufio"
	"hash"
	"io"
	"log"
	"os"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   common.CommandHash + " [string|file]",
	Args:  cobra.ExactArgs(1),
	Short: "Hash string or file",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case Hasher.check:
			Hasher.CheckFile(args[0])
			return
		case Hasher.list:
			Hasher.ListAll(args[0])
			return
		default:
			_ = cmd.Help()
		}
	},
}

var hashSubCmdMd5 = &cobra.Command{
	Use:   common.HashMd5 + " [string|file]",
	Args:  cobra.ExactArgs(1),
	Short: "Print MD5 Checksums",
	Run:   Hasher.Run,

	DisableFlagsInUseLine: true,
}

var hashSubCmdSha1 = &cobra.Command{
	Use:   common.HashSha1 + " [string|file]",
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-1 Checksums",
	Run:   Hasher.Run,

	DisableFlagsInUseLine: true,
}

var hashSubCmdSha256 = &cobra.Command{
	Use:   common.HashSha256 + " [string|file]",
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-256 Checksums",
	Run:   Hasher.Run,

	DisableFlagsInUseLine: true,
}

var hashSubCmdSha512 = &cobra.Command{
	Use:   common.HashSha512 + " [string|file]",
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-512 Checksums",
	Run:   Hasher.Run,

	DisableFlagsInUseLine: true,
}

var Hasher HashFlag

func init() {
	rootCmd.AddCommand(hashCmd)

	hashCmd.Flags().BoolVarP(&Hasher.check, "check", "c", false, common.Usage("Read SHA sums from the file and check them"))
	hashCmd.Flags().BoolVarP(&Hasher.list, "list", "l", false, common.Usage("List multiple SHA sums for the specify input"))
	hashCmd.AddCommand(hashSubCmdMd5, hashSubCmdSha1, hashSubCmdSha256, hashSubCmdSha512)
}

type HashFlag struct {
	check bool
	list  bool
}

func (h *HashFlag) Run(cmd *cobra.Command, args []string) {
	hasher := common.HashAlgorithm(cmd.Name())
	out, err := h.Hash(hasher, args[0])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func (h *HashFlag) Hash(hasher hash.Hash, i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		if validator.ValidFile(data) {
			return h.WriteFile(hasher, data)
		}
		_, err = hasher.Write([]byte(data))
	case []byte:
		_, err = hasher.Write(data)
	default:
		return "", ErrInvalidVar
	}
	if err != nil {
		return "", err
	}
	return Encoder.HexEncode(hasher.Sum(nil))
}

func (h *HashFlag) WriteFile(hasher hash.Hash, filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return "", err
	}
	return Encoder.HexEncode(hasher.Sum(nil))
}

func (h *HashFlag) CheckFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewScanner(f)
	for reader.Scan() {
		var hasher hash.Hash
		slice := strings.Fields(reader.Text())
		switch len([]byte(slice[0])) {
		case 32:
			hasher = common.HashAlgorithm(common.HashMd5)
		case 40:
			hasher = common.HashAlgorithm(common.HashSha1)
		case 64:
			hasher = common.HashAlgorithm(common.HashSha256)
		case 128:
			hasher = common.HashAlgorithm(common.HashSha512)
		}
		got, err := h.WriteFile(hasher, slice[1])
		if err != nil {
			PrintString(err)
			continue
		}
		if got == slice[0] {
			PrintString(slice[1] + ": OK")
		}
		hasher.Reset()
	}
}

func (h *HashFlag) ListAll(s string) {
	algs := []string{common.HashMd5, common.HashSha1, common.HashSha256, common.HashSha512}
	m := make(map[string]string)
	for _, alg := range algs {
		hasher := common.HashAlgorithm(alg)
		out, err := h.Hash(hasher, s)
		if err != nil {
			PrintString(err)
			continue
		}
		m[strings.ToUpper(alg)] = out
	}
	OutputDefaultYAML(m)
}
