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

var sshCmd = &cobra.Command{
	Use:   "ssh-keygen",
	Short: "Generate SSH keypair",
	RunE:  sshCmdGlobalVar.RunE,
}

var sshCmdGlobalVar sshFlag

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().IntVarP(&sshCmdGlobalVar.bit, "bits", "b", 4096, "Specifies the number of bits in the key to create.")
	sshCmd.Flags().StringVarP(&sshCmdGlobalVar.path, "file", "f", "id_rsa", "Specify the file path to generate.")
}

type sshFlag struct {
	bit  int
	path string
}

/* Init checks bit and file exist or not, and return files name. */
func (r *sshFlag) Init() (string, string) {
	if r.bit < 4096 {
		r.bit = 4096
	}
	privateKeyFile := r.path
	if validator.ValidFile(privateKeyFile) {
		PrintString(privateKeyFile + " exist.")
		privateKeyFile += "_new"
		PrintString("Use " + privateKeyFile)
	}
	publicKeyFile := privateKeyFile + ".pub"
	return privateKeyFile, publicKeyFile
}

func (r *sshFlag) RunE(_ *cobra.Command, _ []string) error {
	rsaFile, pubFile := r.Init()
	/* Generate rsa keypair. */
	key, err := rsa.GenerateKey(rand.Reader, r.bit)
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
	if err = os.WriteFile(rsaFile, []byte(privateKey), common.FileModeROwner); err != nil {
		return err
	}
	PrintString(rsaFile + " generated")
	/* Marshal public key. */
	pub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	publicKey := ssh.MarshalAuthorizedKey(pub)
	err = os.WriteFile(pubFile, publicKey, common.FileModeROwner)
	if err != nil {
		return err
	}
	PrintString(pubFile + " generated")
	return err
}