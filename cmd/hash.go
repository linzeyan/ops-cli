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

	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "hash string",
	Run:   Hasher.Run,
}

var Hasher HashFlag

func init() {
	rootCmd.AddCommand(hashCmd)
}

type HashFlag struct{}

func (h *HashFlag) Run(cmd *cobra.Command, args []string) {

}

func (h *HashFlag) Md5Hash(i any) (string, error) {
	hasher := md5.New()
	return h.Write(hasher, i)
}

func (h *HashFlag) Sha1Hash(i any) (string, error) {
	hasher := sha1.New()
	return h.Write(hasher, i)
}

func (h *HashFlag) Sha256Hash(i any) (string, error) {
	hasher := sha256.New()
	return h.Write(hasher, i)
}

func (h *HashFlag) Sha512Hash(i any) (string, error) {
	hasher := sha512.New()
	return h.Write(hasher, i)
}

func (h *HashFlag) Write(hasher hash.Hash, i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
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
