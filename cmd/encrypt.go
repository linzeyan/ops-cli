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
	Run:   encryptCmdGlobalVar.Run,
}

var encryptCmdGlobalVar EncrytpFlag

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.PersistentFlags().BoolVarP(&encryptCmdGlobalVar.decrypt, "decrypt", "d", false, "Decrypt")
	encryptCmd.PersistentFlags().StringVarP(&encryptCmdGlobalVar.key, "key", "k", "", "Specify the encrypt key text or key file")
	encryptCmd.PersistentFlags().StringVarP(&encryptCmdGlobalVar.inFile, "input-file", "i", "", "Specify the input file")
	encryptCmd.PersistentFlags().StringVarP(&encryptCmdGlobalVar.outFile, "output-file", "o", "", "Specify the output file")

	encryptCmd.AddCommand(encryptSubCmdAes)
}

type EncrytpFlag struct {
	decrypt bool
	key     string
	inFile  string
	outFile string
}

func (e *EncrytpFlag) Run(cmd *cobra.Command, _ []string) {
	var err error
	switch cmd.Name() {
	case "aes":
		switch e.decrypt {
		case true:
			err = e.AesDecrypt()
		default:
			err = e.AesEncrypt()
		}
	default:
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (e *EncrytpFlag) getKey() []byte {
	const keyFileExtension = ".key"
	/* Read key file. */
	if ValidFile(e.key) {
		key, err := os.ReadFile(e.key)
		if err != nil {
			return nil
		}
		switch len(key) {
		case 16, 24, 32:
			return key
		}
		return nil
	}
	/* e.key is a key */
	byteKey := []byte(e.key)
	switch len(byteKey) {
	case 16, 24, 32:
		return byteKey
	}
	if e.decrypt {
		return nil
	}
	/* If e.key is not a file and not a valid key, generate a new key. */
	var p RandomString
	key, err := p.genString(32, AllSet)
	if err != nil {
		return nil
	}
	_, filename := path.Split(e.inFile)
	err = os.WriteFile(filename+keyFileExtension, []byte(key), os.ModePerm)
	if err != nil {
		return nil
	}
	return []byte(key)
}

func (e *EncrytpFlag) AesEncrypt() error {
	f, err := os.Open(e.inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	block, err := aes.NewCipher(e.getKey())
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	if e.outFile == "" {
		e.outFile = e.inFile + ".bin"
	}
	out, err := os.OpenFile(e.outFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
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

func (e *EncrytpFlag) AesDecrypt() error {
	f, err := os.Open(e.inFile)
	if err != nil {
		return err
	}
	defer f.Close()
	block, err := aes.NewCipher(e.getKey())
	if err != nil {
		return err
	}
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	fLen := fInfo.Size() - int64(len(iv))
	_, err = f.ReadAt(iv, fLen)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(e.outFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
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
