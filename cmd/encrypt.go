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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

var Encryptor Encrypt

func init() {
	var flags struct {
		Key string `json:"key"`

		decrypt bool
		mode    string
	}
	var encryptCmd = &cobra.Command{
		Use:   CommandEncrypt,
		Short: "Encrypt or decrypt",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	var encryptSubCmdFile = &cobra.Command{
		Use:   CommandFile,
		Args:  cobra.ExactArgs(1),
		Short: "Encrypt or decrypt file",
		RunE: func(_ *cobra.Command, args []string) error {
			/* Read key in the config. */
			if rootConfig != "" {
				if err := ReadConfig(CommandEncrypt, &flags); err != nil {
					return common.ErrInvalidArg
				}
			}

			var err error
			filename := args[0]
			if flags.decrypt {
				if Encryptor.CheckSecret(flags.Key) == nil {
					return common.ErrInvalidArg
				}
				err = Encryptor.DecryptFile(flags.Key, filename, flags.mode)
			} else {
				err = Encryptor.EncryptFile(flags.Key, filename, flags.mode)
			}
			if err != nil {
				return err
			}
			return os.Rename(filename+tempFileExtension, filename)
		},
		Example: common.Examples(`# Encrypt file
~/README.md
~/README.md --config ~/config.toml
~/README.md -k '45984614e8f7d6c5'
~/README.md -k key.txt

# Decrypt file
~/README.md -d -k ~/README.md.key
~/README.md -d --config ~/config.toml
~/README.md -k '45984614e8f7d6c5' -d
~/README.md -k key.txt -d`, CommandEncrypt, CommandFile),
	}

	var encryptSubCmdString = &cobra.Command{
		Use:   CommandString,
		Args:  cobra.ExactArgs(1),
		Short: "Encrypt or decrypt string",
		RunE: func(_ *cobra.Command, args []string) error {
			text := args[0]
			/* Read key in the config. */
			if rootConfig != "" {
				if err := ReadConfig(CommandEncrypt, &flags); err != nil {
					return common.ErrInvalidArg
				}
			}
			if Encryptor.CheckSecret(flags.Key) == nil {
				return common.ErrInvalidArg
			}

			var err error
			var out string
			if flags.decrypt {
				out, err = Encryptor.DecryptString(flags.Key, text, flags.mode)
			} else {
				out, err = Encryptor.EncryptString(flags.Key, text, flags.mode)
			}
			if err != nil {
				return err
			}
			PrintString(out)
			return err
		},
		Example: common.Examples(`# Encrypt string
"Hello World!" --config ~/config.toml
"Hello World!" -k '45984614e8f7d6c5'
"Hello World!" -k key.txt

# Decrypt string
"Hello World!" -d --config ~/config.toml
"Hello World!" -k '45984614e8f7d6c5' -d
"Hello World!" -k key.txt -d`, CommandEncrypt, CommandString),
	}
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.PersistentFlags().BoolVarP(&flags.decrypt, "decrypt", "d", false, common.Usage("Decrypt"))
	encryptCmd.PersistentFlags().StringVarP(&flags.mode, "mode", "m", "CTR", common.Usage("Encrypt mode(CFB/OFB/CTR/GCM)"))
	encryptCmd.PersistentFlags().StringVarP(&flags.Key, "key", "k", "", common.Usage("Specify the encrypt key text or key file"))

	encryptCmd.AddCommand(encryptSubCmdFile)
	encryptCmd.AddCommand(encryptSubCmdString)
}

type Encrypt struct{}

func (*Encrypt) CheckSecret(secret string) []byte {
	/* Read key file. */
	if validator.ValidFile(secret) {
		key, err := os.ReadFile(secret)
		if err != nil {
			return nil
		}
		switch len(key) {
		case 16, 24, 32:
			return key
		}
	}
	/* If secret is a key. */
	byteKey := []byte(secret)
	switch len(byteKey) {
	case 16, 24, 32:
		return byteKey
	}
	return nil
}

func (e *Encrypt) GetKey(secret, filename string, perm os.FileMode) []byte {
	if key := e.CheckSecret(secret); key != nil {
		return key
	}
	f := filepath.Clean(filename) + keyFileExtension
	if validator.ValidFile(f) {
		return nil
	}
	/* If secret is not a file and not a valid key, generate a new key. */
	key := Randoms.GenerateString(32, AllSet)
	err := os.WriteFile(f, key, perm)
	if err != nil {
		return nil
	}
	return key
}

func (*Encrypt) stream(b cipher.Block, iv []byte, mode string, decrypt bool) cipher.Stream {
	switch mode {
	case EncryptModeCFB:
		if decrypt {
			return cipher.NewCFBDecrypter(b, iv)
		}
		return cipher.NewCFBEncrypter(b, iv)
	case EncryptModeCTR:
		return cipher.NewCTR(b, iv)
	case EncryptModeOFB:
		return cipher.NewOFB(b, iv)
	default:
		return cipher.NewCTR(b, iv)
	}
}

func (e *Encrypt) EncryptFile(secret, filename, mode string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(e.GetKey(secret, filename, fInfo.Mode()))
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	out, err := os.OpenFile(filename+tempFileExtension, os.O_RDWR|os.O_CREATE, fInfo.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 1024)
	stream := e.stream(block, iv, mode, false)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			_, wrErr := out.Write(buf[:n])
			if wrErr != nil {
				return wrErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	_, err = out.Write(iv)
	return err
}

func (e *Encrypt) DecryptFile(secret, filename, mode string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(e.GetKey(secret, filename, fInfo.Mode()))
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	fLen := fInfo.Size() - int64(len(iv))
	_, err = f.ReadAt(iv, fLen)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(filename+tempFileExtension, os.O_RDWR|os.O_CREATE, fInfo.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 1024)
	stream := e.stream(block, iv, mode, true)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			if n > int(fLen) {
				n = int(fLen)
			}
			fLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			_, wrErr := out.Write(buf[:n])
			if wrErr != nil {
				return wrErr
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (e *Encrypt) EncryptString(secret, text, mode string) (string, error) {
	var out []byte
	var err error
	key := []byte(secret)
	plainText := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	switch mode {
	case EncryptModeCFB, EncryptModeCTR, EncryptModeOFB:
		out = make([]byte, aes.BlockSize+len(plainText))
		iv := out[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}
		stream := e.stream(block, iv, mode, false)
		stream.XORKeyStream(out[aes.BlockSize:], plainText)
	case EncryptModeGCM:
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
		nonce := make([]byte, gcm.NonceSize())
		out = gcm.Seal(nonce, nonce, []byte(text), nil)
	default:
		return "", err
	}
	return Encoder.Base64URLEncode(out)
}

func (e *Encrypt) DecryptString(secret, text, mode string) (string, error) {
	var out []byte
	var err error
	key := []byte(secret)
	cipherText, err := Encoder.Base64URLDecode(text)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	switch mode {
	case EncryptModeCFB, EncryptModeCTR, EncryptModeOFB:
		iv := cipherText[:aes.BlockSize]
		stream := e.stream(block, iv, mode, true)
		out = cipherText[aes.BlockSize:]
		stream.XORKeyStream(out, out)
	case EncryptModeGCM:
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
		nonceSize := gcm.NonceSize()
		nonce, enc := cipherText[:nonceSize], cipherText[nonceSize:]
		out, err = gcm.Open(nil, nonce, enc, nil)
		if err != nil {
			return "", err
		}
	default:
		return "", err
	}
	return string(out), err
}
