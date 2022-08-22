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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"
	"path"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

const (
	EncryptModeCBC = "CBC"
	EncryptModeCFB = "CFB"
)
const (
	keyFileExtension  = ".key"
	tempFileExtension = ".temp"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt or decrypt",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var encryptSubCmdFile = &cobra.Command{
	Use:   "file",
	Args:  cobra.ExactArgs(1),
	Short: "Encrypt or decrypt file",
	Run:   Encryptor.FileRun,
	Example: common.Examples(`# Encrypt file
ops-cli encrypt file ~/README.md
ops-cli encrypt file ~/README.md -k '45984614e8f7d6c5'
ops-cli encrypt file ~/README.md -k key.txt

# Decrypt file
ops-cli encrypt file ~/README.md -d -k ~/README.md.key
ops-cli encrypt file ~/README.md -k '45984614e8f7d6c5' -d
ops-cli encrypt file ~/README.md -k key.txt -d`),
}

var encryptSubCmdString = &cobra.Command{
	Use:   "string",
	Args:  cobra.ExactArgs(1),
	Short: "Encrypt or decrypt string",
	Run:   Encryptor.StringRun,
}

var Encryptor EncrytpFlag

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.PersistentFlags().BoolVarP(&Encryptor.decrypt, "decrypt", "d", false, "Decrypt")
	encryptCmd.PersistentFlags().StringVarP(&Encryptor.mode, "mode", "m", "CTR", "Specify the encrypt mode(CFB/OFB/CTR)")
	encryptCmd.PersistentFlags().StringVarP(&Encryptor.key, "key", "k", "", "Specify the encrypt key text or key file")

	encryptCmd.AddCommand(encryptSubCmdFile)
	encryptCmd.AddCommand(encryptSubCmdString)
}

type EncrytpFlag struct {
	decrypt bool
	key     string
	mode    string
}

func (e *EncrytpFlag) FileRun(cmd *cobra.Command, args []string) {
	var err error
	filename := args[0]
	if e.decrypt {
		err = e.DecryptFile(e.key, filename)
	} else {
		err = e.EncryptFile(e.key, filename)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = os.Rename(filename+tempFileExtension, filename); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (e *EncrytpFlag) getKey(secret, filename string, perm os.FileMode) []byte {
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
		return nil
	}
	/* secret is a key */
	byteKey := []byte(secret)
	switch len(byteKey) {
	case 16, 24, 32:
		return byteKey
	}
	if e.decrypt {
		return nil
	}
	/* If secret is not a file and not a valid key, generate a new key. */
	var p RandomString
	key := p.GenerateString(32, AllSet)
	err := os.WriteFile(path.Base(filename)+keyFileExtension, key, perm)
	if err != nil {
		return nil
	}
	return key
}

func (e *EncrytpFlag) stream(b cipher.Block, iv []byte) cipher.Stream {
	switch e.mode {
	case EncryptModeCFB:
		if e.decrypt {
			return cipher.NewCFBDecrypter(b, iv)
		}
		return cipher.NewCFBEncrypter(b, iv)
	case "CTR":
		return cipher.NewCTR(b, iv)
	case "OFB":
		return cipher.NewOFB(b, iv)
	default:
		return cipher.NewCTR(b, iv)
	}
}

func (e *EncrytpFlag) EncryptFile(secret, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(e.getKey(secret, filename, fInfo.Mode()))
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
	stream := e.stream(block, iv)
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

func (e *EncrytpFlag) DecryptFile(secret, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(e.getKey(secret, filename, fInfo.Mode()))
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
	stream := e.stream(block, iv)
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
			break
		}
		if err != nil {
			return err
		}
	}
	return err
}

func (e *EncrytpFlag) StringRun(cmd *cobra.Command, args []string) {
	var err error
	var out string
	text := args[0]
	if e.decrypt {
		out, err = e.DecryptString(e.key, text)
	} else {
		out, err = e.EncryptString(e.key, text)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	PrintString(out)
}

func (e *EncrytpFlag) EncryptString(secret, text string) (string, error) {
	var out []byte
	var err error
	key := []byte(secret)
	plainText := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	switch e.mode {
	case "CBC":
		out = make([]byte, aes.BlockSize+len(plainText))
		iv := out[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}
		padding := aes.BlockSize - len(plainText)%aes.BlockSize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		plainText = append(plainText, padtext...)
		cbc := cipher.NewCBCEncrypter(block, iv)
		cbc.CryptBlocks(out, plainText)
	case EncryptModeCFB:
		out = make([]byte, aes.BlockSize+len(plainText))
		iv := out[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}
		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(out[aes.BlockSize:], plainText)
	case "GCM":
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
		nonce := make([]byte, gcm.NonceSize())
		out = gcm.Seal(nonce, nonce, []byte(text), nil)
	}
	return Encoder.Base64URLEncode(out)
}

func (e *EncrytpFlag) DecryptString(secret, text string) (string, error) {
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
	switch e.mode {
	case "CBC":
		// cbc := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
	case EncryptModeCFB:
		iv := cipherText[:aes.BlockSize]
		stream := cipher.NewCFBDecrypter(block, iv)
		out = cipherText[aes.BlockSize:]
		stream.XORKeyStream(out, out)
	case "GCM":
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
	}
	return string(out), err
}
