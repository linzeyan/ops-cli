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

	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash string or file",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var hashSubCmdMd5 = &cobra.Command{
	Use:   HashMd5,
	Args:  cobra.ExactArgs(1),
	Short: "Print MD5 Checksums",
	Run:   Hasher.Run,
}

var hashSubCmdSha1 = &cobra.Command{
	Use:   HashSha1,
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-1 Checksums",
	Run:   Hasher.Run,
}

var hashSubCmdSha256 = &cobra.Command{
	Use:   HashSha256,
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-256 Checksums",
	Run:   Hasher.Run,
}

var hashSubCmdSha512 = &cobra.Command{
	Use:   HashSha512,
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-512 Checksums",
	Run:   Hasher.Run,
}

var Hasher HashFlag

func init() {
	rootCmd.AddCommand(hashCmd)

	hashCmd.PersistentFlags().BoolVarP(&Hasher.check, "check", "c", false, "Read SHA sums from the file and check them")
	hashCmd.AddCommand(hashSubCmdMd5, hashSubCmdSha1, hashSubCmdSha256, hashSubCmdSha512)
}

type HashFlag struct {
	check bool
}

func (h *HashFlag) Run(cmd *cobra.Command, args []string) {
	var hasher hash.Hash
	var out string
	var err error
	hasher = HashAlgorithm(cmd.Name())
	if h.check {
		h.CheckFile(hasher, args[0])
		return
	}
	out, err = h.Hash(hasher, args[0])
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
		if ValidFile(data) {
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

func (h *HashFlag) CheckFile(hasher hash.Hash, filename string) {
	var err error
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewScanner(f)
	for reader.Scan() {
		slice := strings.Fields(reader.Text())
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
