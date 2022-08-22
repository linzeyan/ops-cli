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

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type EncodeType string

const (
	Base32    EncodeType = "base32"
	Base64    EncodeType = "base64"
	Hex       EncodeType = "hex"
	Std       EncodeType = "std"
	URL       EncodeType = "url"
	Base32Hex EncodeType = Base32 + Hex
	Base32Std EncodeType = Base32 + Std
	Base64Std EncodeType = Base64 + Std
	Base64URL EncodeType = Base64 + URL
)

func (e EncodeType) String() string {
	return string(e)
}

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode and decode string or file",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var encodeSubCmdBase32Hex = &cobra.Command{
	Use:   Base32Hex.String(),
	Args:  cobra.ExactArgs(1),
	Short: "Base32 hex encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdBase32Std = &cobra.Command{
	Use:   Base32Std.String(),
	Args:  cobra.ExactArgs(1),
	Short: "Base32 standard encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdBase64Std = &cobra.Command{
	Use:   Base64Std.String(),
	Args:  cobra.ExactArgs(1),
	Short: "Base64 standard encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdBase64URL = &cobra.Command{
	Use:   Base64URL.String(),
	Args:  cobra.ExactArgs(1),
	Short: "Base64 url encoding or decoding",
	Run:   Encoder.Run,
}

var encodeSubCmdHex = &cobra.Command{
	Use:   Hex.String(),
	Args:  cobra.ExactArgs(1),
	Short: "Hexadecimal encoding or decoding",
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
	switch validator.ValidFile(args[0]) {
	case true:
		data, err = os.ReadFile(args[0])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	case false:
		data = args[0]
	}
	switch EncodeType(cmd.Name()) {
	case Base32Hex:
		out, err = e.Base32HexEncode(data)
	case Base32, Base32Std:
		out, err = e.Base32StdEncode(data)
	case Base64, Base64Std, Std:
		out, err = e.Base64StdEncode(data)
	case Base64URL, URL:
		out, err = e.Base64URLEncode(data)
	case Hex:
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
	switch EncodeType(cmd.Name()) {
	case Base32Hex:
		out, err = e.Base32HexDecode(args[0])
	case Base32, Base32Std:
		out, err = e.Base32StdDecode(args[0])
	case Base64, Base64Std, Std:
		out, err = e.Base64StdDecode(args[0])
	case URL, Base64URL:
		out, err = e.Base64URLDecode(args[0])
	case Hex:
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
		return "", common.ErrInvalidArg
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
		return "", common.ErrInvalidArg
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
		return "", common.ErrInvalidArg
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
		return "", common.ErrInvalidArg
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
		return "", common.ErrInvalidArg
	}
}

func (e *EncodeFlag) HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (e *EncodeFlag) JSONEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", common.IndentTwoSpaces)
	err := encoder.Encode(i)
	return buf.String(), err
}

func (e *EncodeFlag) JSONDecode(r io.Reader, i any) (any, error) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(i)
	return i, err
}

func (e *EncodeFlag) JSONMarshaler(src, dst any) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return err
	}
	return nil
}

func (e *EncodeFlag) PemEncode(i any, t ...string) (string, error) {
	var err error
	var block = &pem.Block{Type: ""}
	if len(t) != 0 {
		for _, arg := range t {
			block.Type += arg
		}
	}
	switch data := i.(type) {
	case string:
		block.Bytes = []byte(data)
	case []byte:
		block.Bytes = data
	case *pem.Block:
		block = data
	default:
		return "", common.ErrInvalidArg
	}
	var buf bytes.Buffer
	err = pem.Encode(&buf, block)
	return buf.String(), err
}

func (e *EncodeFlag) PemDecode(b []byte) ([]byte, error) {
	p, _ := pem.Decode(b)
	if p == nil {
		return nil, ErrFileType
	}
	return p.Bytes, nil
}

func (e *EncodeFlag) XMLEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", common.IndentTwoSpaces)
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
