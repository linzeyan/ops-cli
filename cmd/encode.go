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
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode and decode string or file",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var encodeSubCmdBase32Hex = &cobra.Command{
	Use:   "base32hex",
	Args:  cobra.ExactArgs(1),
	Short: "base32 hex encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdBase32Std = &cobra.Command{
	Use:   "base32std",
	Args:  cobra.ExactArgs(1),
	Short: "base32 standard encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdBase64Std = &cobra.Command{
	Use:   "base64std",
	Args:  cobra.ExactArgs(1),
	Short: "base64 standard encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdBase64URL = &cobra.Command{
	Use:   "base64url",
	Args:  cobra.ExactArgs(1),
	Short: "base64 url encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdHex = &cobra.Command{
	Use:   "hex",
	Args:  cobra.ExactArgs(1),
	Short: "hexadecimal encoding or decoding",
	Run:   Encoder.Run,
}

var Encoder EncodeFlag

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.PersistentFlags().BoolVarP(&Encoder.decode, "decode", "d", false, "Decodes input")
	encodeCmd.AddCommand(encodeSubCmdBase32Hex, encodeSubCmdBase32Std)
	encodeCmd.AddCommand(encodeSubCmdBase64Std, encodeSubCmdBase64URL)
	encodeCmd.AddCommand(encodeSubCmdHex)
}

type EncodeFlag struct {
	decode bool
}

func (e *EncodeFlag) Run(cmd *cobra.Command, args []string) {
	if e.decode {
		e.RunDecode(cmd, args)
		return
	}
	var err error
	var out string
	var data any
	switch ValidFile(args[0]) {
	case true:
		data, err = os.ReadFile(args[0])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	case false:
		data = args[0]
	}
	switch cmd.Name() {
	case "base32hex":
		out, err = e.Base32HexEncode(data)
	case "base32std":
		out, err = e.Base32StdEncode(data)
	case "base64std":
		out, err = e.Base64StdEncode(data)
	case "base64url":
		out, err = e.Base64URLEncode(data)
	case "hex":
		out, err = e.HexEncode(data)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func (e *EncodeFlag) RunDecode(cmd *cobra.Command, args []string) {
	var err error
	var out []byte
	switch cmd.Name() {
	case "base32hex":
		out, err = e.Base32HexDecode(args[0])
	case "base32std":
		out, err = e.Base32StdDecode(args[0])
	case "base64std":
		out, err = e.Base64StdDecode(args[0])
	case "base64url":
		out, err = e.Base64URLDecode(args[0])
	case "hex":
		out, err = e.HexDecode(args[0])
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func (e *EncodeFlag) Base32HexEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base32.HexEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base32.HexEncoding.EncodeToString(data), err
	default:
		return "", ErrInvalidVar
	}
}

func (e *EncodeFlag) Base32HexDecode(s string) ([]byte, error) {
	return base32.HexEncoding.DecodeString(s)
}

func (e *EncodeFlag) Base32StdEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base32.StdEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base32.StdEncoding.EncodeToString(data), err
	default:
		return "", ErrInvalidVar
	}
}

func (e *EncodeFlag) Base32StdDecode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (e *EncodeFlag) Base64StdEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base64.StdEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base64.StdEncoding.EncodeToString(data), err
	default:
		return "", ErrInvalidVar
	}
}

func (e *EncodeFlag) Base64StdDecode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (e *EncodeFlag) Base64URLEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base64.URLEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base64.URLEncoding.EncodeToString(data), err
	default:
		return "", ErrInvalidVar
	}
}

func (e *EncodeFlag) Base64URLDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

func (e *EncodeFlag) HexEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return hex.EncodeToString([]byte(data)), err
	case []byte:
		return hex.EncodeToString(data), err
	default:
		return "", ErrInvalidVar
	}
}

func (e *EncodeFlag) HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (e *EncodeFlag) JSONEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(i)
	return buf.String(), err
}

func (e *EncodeFlag) JSONDecode(r io.Reader, i any) (any, error) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(i)
	return i, err
}

func (e *EncodeFlag) PemEncode(b *pem.Block) (string, error) {
	var buf bytes.Buffer
	err := pem.Encode(&buf, b)
	return buf.String(), err
}

func (e *EncodeFlag) PemDecode(b []byte) (*pem.Block, error) {
	p, _ := pem.Decode(b)
	if p == nil {
		return p, ErrFileType
	}
	return p, nil
}

func (e *EncodeFlag) XMLEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	err := encoder.Encode(i)
	return buf.String(), err
}

func (e *EncodeFlag) XMLDecode(r io.Reader, i any) (any, error) {
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(i)
	return i, err
}

func (e *EncodeFlag) YamlEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(i)
	return buf.String(), err
}

func (e *EncodeFlag) YamlDecode(r io.Reader, i any) (any, error) {
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(i)
	return i, err
}
