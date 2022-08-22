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
	"log"
	"os"
	"path"

	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt or decrypt file",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var encryptSubCmdAes = &cobra.Command{
	Use:   "aes",
	Short: "Encrypt or decrypt",
	Run:   Encryptor.AesRun,
}

var Encryptor EncrytpFlag

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.PersistentFlags().BoolVarP(&Encryptor.decrypt, "decrypt", "d", false, "Decrypt")
	encryptCmd.PersistentFlags().StringVarP(&Encryptor.key, "key", "k", "", "Specify the encrypt key text or key file")
	encryptCmd.PersistentFlags().StringVarP(&Encryptor.file, "file", "f", "", "Specify the file")
	encryptCmd.PersistentFlags().StringVarP(&Encryptor.mode, "mode", "m", "CTR", "Specify the encrypt mode")

	encryptCmd.AddCommand(encryptSubCmdAes)
}

type EncrytpFlag struct {
	decrypt bool
	key     string
	file    string
	mode    string
}

func (e *EncrytpFlag) AesRun(cmd *cobra.Command, _ []string) {
	var err error

	if e.decrypt {
		if err = e.AesDecrypt(e.key, e.file); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}
	if err = e.AesEncrypt(e.key, e.file); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (e *EncrytpFlag) getKey(secret, filename string, perm os.FileMode) []byte {
	const keyFileExtension = ".key"
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

func (e *EncrytpFlag) AesEncrypt(secret, filename string) error {
	err := e.aesEncrypt(secret, filename)
	if err != nil {
		return err
	}
	return os.Rename(filename+".bin", filename)
}

func (e *EncrytpFlag) stream(b cipher.Block, iv []byte) cipher.Stream {
	switch e.mode {
	case "CFB":
		if e.decrypt {
			return cipher.NewCFBDecrypter(b, iv)
		}
		return cipher.NewCFBEncrypter(b, iv)
	case "OFB":
		return cipher.NewOFB(b, iv)
	default:
		return cipher.NewCTR(b, iv)
	}
}

func (e *EncrytpFlag) aesEncrypt(secret, filename string) error {
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
	out, err := os.OpenFile(filename+".bin", os.O_RDWR|os.O_CREATE, fInfo.Mode())
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

func (e *EncrytpFlag) AesDecrypt(secret, filename string) error {
	err := e.aesDecrypt(secret, filename)
	if err != nil {
		return err
	}
	return os.Rename(filename+".raw", filename)
}

func (e *EncrytpFlag) aesDecrypt(secret, filename string) error {
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
	out, err := os.OpenFile(filename+".raw", os.O_RDWR|os.O_CREATE, fInfo.Mode())
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
