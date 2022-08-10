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
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"log"

	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "encode and decode string",
	Run:   enc.Run,
}

var enc encodeFlag

func init() {
	rootCmd.AddCommand(encodeCmd)
}

type encodeFlag struct{}

func (e *encodeFlag) Run(cmd *cobra.Command, args []string) {

}

func (e *encodeFlag) Base32HexEncode(s string) string {
	return base32.HexEncoding.EncodeToString([]byte(s))
}

func (e *encodeFlag) Base32HexDecode(s string) string {
	out, err := base32.HexEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(out)
}

func (e *encodeFlag) Base32StdEncode(s string) string {
	return base32.StdEncoding.EncodeToString([]byte(s))
}

func (e *encodeFlag) Base32StdDecode(s string) string {
	out, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(out)
}

func (e *encodeFlag) Base64StdEncode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (e *encodeFlag) Base64StdDecode(s string) string {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(out)
}

func (e *encodeFlag) Base64URLEncode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (e *encodeFlag) Base64URLDecode(s string) string {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(out)
}

func (e *encodeFlag) HexEncode(s string) string {
	return hex.EncodeToString([]byte(s))
}

func (e *encodeFlag) HexDecode(s string) string {
	out, err := hex.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(out)
}
