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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func initSSHKeyGen() *cobra.Command {
	var flags struct {
		bit  int
		path string
		key  string
	}
	var sshkeygenCmd = &cobra.Command{
		Use:   CommandSSH,
		Short: "Generate SSH keypair",
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			var k SSHKeygen
			if flags.key != "" {
				return k.GetPub(flags.key)
			}
			if flags.bit < 4096 {
				flags.bit = 4096
			}
			return k.Generate(flags.bit, flags.path)
		},
	}
	sshkeygenCmd.Flags().IntVarP(&flags.bit, "bits", "b", 4096, common.Usage("Specifies the number of bits in the key to create"))
	sshkeygenCmd.Flags().StringVarP(&flags.path, "file", "f", "id_rsa", common.Usage("Specify the file path to generate"))
	sshkeygenCmd.Flags().StringVarP(&flags.key, "generate", "g", "", common.Usage("Get public key from private key"))
	return sshkeygenCmd
}

type SSHKeygen struct{}

/* Prepare checks file exist or not, and return files name. */
func (*SSHKeygen) Prepare(path string) (string, string) {
	privateKeyFile := path
	if validator.ValidFile(privateKeyFile) {
		PrintString(privateKeyFile + " exist.")
		privateKeyFile += "_new"
		PrintString("Use " + privateKeyFile)
	}
	publicKeyFile := privateKeyFile + ".pub"
	return privateKeyFile, publicKeyFile
}

func (r *SSHKeygen) Generate(bit int, path string) error {
	rsaFile, pubFile := r.Prepare(path)
	/* Generate rsa keypair. */
	key, err := rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		return err
	}
	/* Encode and write private key to file. */
	privateKey, err := Encoder.PemEncode(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return err
	}
	if err = os.WriteFile(rsaFile, []byte(privateKey), FileModeROwner); err != nil {
		return err
	}
	PrintString(rsaFile + " generated")
	/* Marshal public key. */
	pub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	publicKey := ssh.MarshalAuthorizedKey(pub)
	err = os.WriteFile(pubFile, publicKey, FileModeROwner)
	if err != nil {
		return err
	}
	PrintString(pubFile + " generated")
	return err
}

func (*SSHKeygen) GetPub(keyFile string) error {
	f, err := os.ReadFile(keyFile)
	if err != nil {
		return err
	}
	decode, err := Encoder.PemDecode(f)
	if err != nil {
		return err
	}
	key, err := x509.ParsePKCS1PrivateKey(decode)
	if err != nil {
		return err
	}
	pub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	PrintString(ssh.MarshalAuthorizedKey(pub))
	return err
}
