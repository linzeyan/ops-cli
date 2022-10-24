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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

var Hasher Hash

func initHash() *cobra.Command {
	var flags struct {
		check bool
		list  bool
	}
	var hashCmd = &cobra.Command{
		Use:   CommandHash + " [-c check|-l list] [checksums.txt|string or file]",
		Args:  cobra.ExactArgs(1),
		Short: "Hash string or file",
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case flags.check:
				Hasher.CheckFile(args[0])
				return
			case flags.list:
				Hasher.ListAll(args[0])
				return
			default:
				_ = cmd.Help()
			}
		},
		DisableFlagsInUseLine: true,
	}

	runE := func(cmd *cobra.Command, args []string) error {
		var err error
		hasher := HashAlgorithm(cmd.Name())
		out, err := Hasher.Hash(hasher, args[0])
		if err != nil {
			return err
		}
		printer.Printf(out)
		return err
	}

	var hashSubCmdMd5 = &cobra.Command{
		Use:   HashMd5 + " [string|file]",
		Args:  cobra.ExactArgs(1),
		Short: "Print MD5 Checksums",
		RunE:  runE,

		DisableFlagsInUseLine: true,
	}

	var hashSubCmdSha1 = &cobra.Command{
		Use:   HashSha1 + " [string|file]",
		Args:  cobra.ExactArgs(1),
		Short: "Print SHA-1 Checksums",
		RunE:  runE,

		DisableFlagsInUseLine: true,
	}

	var hashSubCmdSha256 = &cobra.Command{
		Use:   HashSha256 + " [string|file]",
		Args:  cobra.ExactArgs(1),
		Short: "Print SHA-256 Checksums",
		RunE:  runE,

		DisableFlagsInUseLine: true,
	}

	var hashSubCmdSha512 = &cobra.Command{
		Use:   HashSha512 + " [string|file]",
		Args:  cobra.ExactArgs(1),
		Short: "Print SHA-512 Checksums",
		RunE:  runE,

		DisableFlagsInUseLine: true,
	}

	hashCmd.Flags().BoolVarP(&flags.check, "check", "c", false, common.Usage("Read SHA sums from the file and check them"))
	hashCmd.Flags().BoolVarP(&flags.list, "list", "l", false, common.Usage("List multiple SHA sums for the specify input"))
	hashCmd.AddCommand(hashSubCmdMd5, hashSubCmdSha1, hashSubCmdSha256, hashSubCmdSha512)
	return hashCmd
}

type Hash struct{}

func (h *Hash) Hash(hasher hash.Hash, i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		if common.IsFile(data) {
			return h.WriteFile(hasher, data)
		}
		_, err = hasher.Write([]byte(data))
	case []byte:
		_, err = hasher.Write(data)
	default:
		return "", common.ErrInvalidArg
	}
	if err != nil {
		return "", err
	}
	return Encoder.HexEncode(hasher.Sum(nil))
}

func (h *Hash) WriteFile(hasher hash.Hash, filename string) (string, error) {
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

func (h *Hash) CheckFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		printer.Error(err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewScanner(f)
	for reader.Scan() {
		var hasher hash.Hash
		slice := strings.Fields(reader.Text())
		switch len([]byte(slice[0])) {
		case 32:
			hasher = HashAlgorithm(HashMd5)
		case 40:
			hasher = HashAlgorithm(HashSha1)
		case 64:
			hasher = HashAlgorithm(HashSha256)
		case 128:
			hasher = HashAlgorithm(HashSha512)
		}
		got, err := h.WriteFile(hasher, slice[1])
		if err != nil {
			printer.Error(err)
			continue
		}
		if got == slice[0] {
			printer.Printf("%s: OK", slice[1])
		}
		hasher.Reset()
	}
}

func (h *Hash) ListAll(s string) {
	algs := []string{HashMd5, HashSha1, HashSha256, HashSha512}
	m := make(map[string]string)
	for _, alg := range algs {
		hasher := HashAlgorithm(alg)
		out, err := h.Hash(hasher, s)
		if err != nil {
			printer.Error(err)
			continue
		}
		m[strings.ToUpper(alg)] = out
	}
	printer.Printf(printer.SetYamlAsDefaultFormat(rootOutputFormat), m)
}

func HashAlgorithm(alg string) hash.Hash {
	m := map[string]hash.Hash{
		HashMd5:        md5.New(),
		HashSha1:       sha1.New(),
		HashSha224:     sha256.New224(),
		HashSha256:     sha256.New(),
		HashSha384:     sha512.New384(),
		HashSha512:     sha512.New(),
		HashSha512_224: sha512.New512_224(),
		HashSha512_256: sha512.New512_256(),
	}
	if h, ok := m[alg]; ok {
		return h
	}
	return nil
}
