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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash string or file",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var hashSubCmdMd5 = &cobra.Command{
	Use:   "md5",
	Args:  cobra.ExactArgs(1),
	Short: "Print MD5 Checksums",
	Run:   Hasher.Run,
}

var hashSubCmdSha1 = &cobra.Command{
	Use:   "sha1",
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-1 Checksums",
	Run:   Hasher.Run,
}

var hashSubCmdSha256 = &cobra.Command{
	Use:   "sha256",
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-256 Checksums",
	Run:   Hasher.Run,
}

var hashSubCmdSha512 = &cobra.Command{
	Use:   "sha512",
	Args:  cobra.ExactArgs(1),
	Short: "Print SHA-512 Checksums",
	Run:   Hasher.Run,
}

var Hasher HashFlag

func init() {
	rootCmd.AddCommand(hashCmd)

	hashCmd.AddCommand(hashSubCmdMd5, hashSubCmdSha1, hashSubCmdSha256, hashSubCmdSha512)
}

type HashFlag struct{}

func (h *HashFlag) Run(cmd *cobra.Command, args []string) {
	var out string
	var err error
	switch cmd.Name() {
	case "md5":
		out, err = h.Md5Hash(args[0])
	case "sha1":
		out, err = h.Sha1Hash(args[0])
	case "sha256":
		out, err = h.Sha256Hash(args[0])
	case "sha512":
		out, err = h.Sha512Hash(args[0])
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func (h *HashFlag) Md5Hash(i any) (string, error) {
	hasher := md5.New()
	return h.WriteString(hasher, i)
}

func (h *HashFlag) Sha1Hash(i any) (string, error) {
	hasher := sha1.New()
	return h.WriteString(hasher, i)
}

func (h *HashFlag) Sha256Hash(i any) (string, error) {
	hasher := sha256.New()
	return h.WriteString(hasher, i)
}

func (h *HashFlag) Sha512Hash(i any) (string, error) {
	hasher := sha512.New()
	return h.WriteString(hasher, i)
}

func (h *HashFlag) WriteString(hasher hash.Hash, i any) (string, error) {
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
