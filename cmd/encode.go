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
	"os"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Encoder Encode

func initEncode() *cobra.Command {
	var flags struct {
		decode bool
	}
	var encodeCmd = &cobra.Command{
		Use:   CommandEncode,
		Short: "Encode and decode string or file",
		RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	runDecode := func(cmd *cobra.Command, args []string) {
		var err error
		var out []byte
		switch cmd.Name() {
		case CommandBase32Hex:
			out, err = Encoder.Base32HexDecode(args[0])
		case CommandBase32Std:
			out, err = Encoder.Base32StdDecode(args[0])
		case CommandBase64Std:
			out, err = Encoder.Base64StdDecode(args[0])
		case CommandBase64URL:
			out, err = Encoder.Base64URLDecode(args[0])
		case CommandHex:
			out, err = Encoder.HexDecode(args[0])
		}
		if err != nil {
			logger.Info(err.Error())
			printer.Error(err)
			return
		}
		printer.Printf("%s", out)
	}

	run := func(cmd *cobra.Command, args []string) {
		if flags.decode {
			runDecode(cmd, args)
			return
		}
		var err error
		var out string
		var data any
		switch common.IsFile(args[0]) {
		case true:
			data, err = os.ReadFile(args[0])
			if err != nil {
				logger.Info(err.Error())
				printer.Error(err)
				return
			}
		case false:
			data = args[0]
		}
		switch cmd.Name() {
		case CommandBase32Hex:
			out, err = Encoder.Base32HexEncode(data)
		case CommandBase32Std:
			out, err = Encoder.Base32StdEncode(data)
		case CommandBase64Std:
			out, err = Encoder.Base64StdEncode(data)
		case CommandBase64URL:
			out, err = Encoder.Base64URLEncode(data)
		case CommandHex:
			out, err = Encoder.HexEncode(data)
		}
		if err != nil {
			logger.Info(err.Error())
			printer.Error(err)
			return
		}
		printer.Printf(out)
	}

	var encodeSubCmdBase32Hex = &cobra.Command{
		Use:   CommandBase32Hex,
		Args:  cobra.ExactArgs(1),
		Short: "Base32 hex encoding or decoding",
		Run:   run,
	}

	var encodeSubCmdBase32Std = &cobra.Command{
		Use:   CommandBase32Std,
		Args:  cobra.ExactArgs(1),
		Short: "Base32 standard encoding or decoding",
		Run:   run,
	}

	var encodeSubCmdBase64Std = &cobra.Command{
		Use:   CommandBase64Std,
		Args:  cobra.ExactArgs(1),
		Short: "Base64 standard encoding or decoding",
		Run:   run,
	}

	var encodeSubCmdBase64URL = &cobra.Command{
		Use:   CommandBase64URL,
		Args:  cobra.ExactArgs(1),
		Short: "Base64 url encoding or decoding",
		Run:   run,
	}

	var encodeSubCmdHex = &cobra.Command{
		Use:   CommandHex,
		Args:  cobra.ExactArgs(1),
		Short: "Hexadecimal encoding or decoding",
		Run:   run,
	}

	encodeCmd.PersistentFlags().BoolVarP(&flags.decode, "decode", "d", false, common.Usage("Decodes input"))
	encodeCmd.AddCommand(encodeSubCmdBase32Hex, encodeSubCmdBase32Std)
	encodeCmd.AddCommand(encodeSubCmdBase64Std, encodeSubCmdBase64URL)
	encodeCmd.AddCommand(encodeSubCmdHex)
	return encodeCmd
}

type Encode struct{}

func (*Encode) Base32HexEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base32.HexEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base32.HexEncoding.EncodeToString(data), err
	default:
		logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(i))
		return "", common.ErrInvalidArg
	}
}

func (*Encode) Base32HexDecode(s string) ([]byte, error) {
	b, err := base32.HexEncoding.DecodeString(s)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(s))
		return nil, err
	}
	return b, err
}

func (*Encode) Base32StdEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base32.StdEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base32.StdEncoding.EncodeToString(data), err
	default:
		logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(i))
		return "", common.ErrInvalidArg
	}
}

func (*Encode) Base32StdDecode(s string) ([]byte, error) {
	b, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(s))
		return nil, err
	}
	return b, err
}

func (*Encode) Base64StdEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base64.StdEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base64.StdEncoding.EncodeToString(data), err
	default:
		logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(i))
		return "", common.ErrInvalidArg
	}
}

func (*Encode) Base64StdDecode(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(s))
		return nil, err
	}
	return b, err
}

func (*Encode) Base64URLEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return base64.URLEncoding.EncodeToString([]byte(data)), err
	case []byte:
		return base64.URLEncoding.EncodeToString(data), err
	default:
		logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(i))
		return "", common.ErrInvalidArg
	}
}

func (*Encode) Base64URLDecode(s string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(s))
		return nil, err
	}
	return b, err
}

func (*Encode) HexEncode(i any) (string, error) {
	var err error
	switch data := i.(type) {
	case string:
		return hex.EncodeToString([]byte(data)), err
	case []byte:
		return hex.EncodeToString(data), err
	default:
		logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(i))
		return "", common.ErrInvalidArg
	}
}

func (*Encode) HexDecode(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(s))
		return nil, err
	}
	return b, err
}

func (*Encode) JSONEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", IndentTwoSpaces)
	err := encoder.Encode(i)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(i))
	}
	return buf.String(), err
}

func (*Encode) JSONDecode(r io.Reader, i any) (any, error) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(i)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(i))
	}
	return i, err
}

func (*Encode) JSONMarshaler(src, dst any) error {
	var err error
	switch data := src.(type) {
	case string:
		if err = json.Unmarshal([]byte(data), dst); err != nil {
			logger.Debug(err.Error(),
				common.NewField("src", src),
				common.NewField("dst", dst),
			)
			return err
		}
	case []byte:
		if err = json.Unmarshal(data, dst); err != nil {
			logger.Debug(err.Error(),
				common.NewField("src", src),
				common.NewField("dst", dst),
			)
			return err
		}
	default:
		bytes, err := json.Marshal(data)
		if err != nil {
			logger.Debug(err.Error(), common.DefaultField(data))
			return err
		}
		if err = json.Unmarshal(bytes, dst); err != nil {
			logger.Debug(err.Error(),
				common.NewField("bytes", bytes),
				common.NewField("dst", dst),
			)
			return err
		}
	}
	return err
}

func (*Encode) PemEncode(i any, t ...string) (string, error) {
	var err error
	var block = &pem.Block{Type: ""}
	if len(t) != 0 {
		block.Type += common.SliceStringToString(t)
	}
	switch data := i.(type) {
	case string:
		block.Bytes = []byte(data)
	case []byte:
		block.Bytes = data
	case *pem.Block:
		block = data
	default:
		logger.Debug(common.ErrInvalidArg.Error(), common.DefaultField(i))
		return "", common.ErrInvalidArg
	}
	var buf bytes.Buffer
	err = pem.Encode(&buf, block)
	if err != nil {
		logger.Debug(err.Error(),
			common.NewField("buf", buf),
			common.NewField("block", block),
		)
	}
	return buf.String(), err
}

func (*Encode) PemDecode(b []byte) ([]byte, error) {
	p, _ := pem.Decode(b)
	if p == nil {
		logger.Debug(common.ErrInvalidFile.Error(), common.DefaultField(b))
		return nil, common.ErrInvalidFile
	}
	return p.Bytes, nil
}

func (*Encode) XMLEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", IndentTwoSpaces)
	err := encoder.Encode(i)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(i))
	}
	return buf.String(), err
}

func (*Encode) XMLDecode(r io.Reader, i any) (any, error) {
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(i)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(i))
	}
	return i, err
}

func (*Encode) YamlEncode(i any) (string, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(i)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(i))
	}
	return buf.String(), err
}

func (*Encode) YamlDecode(r io.Reader, i any) (any, error) {
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(i)
	if err != nil {
		logger.Debug(err.Error(), common.DefaultField(i))
	}
	return i, err
}
