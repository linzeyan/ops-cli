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
	"log"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh-keygen",
	Short: "Generate SSH keypair",
	Run:   sshCmdGlobalVar.Run,
}

var sshCmdGlobalVar sshFlag

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().IntVarP(&sshCmdGlobalVar.bit, "bits", "b", 4096, "Specifies the number of bits in the key to create.")
	sshCmd.Flags().StringVarP(&sshCmdGlobalVar.path, "file", "f", "./id_rsa", "Specify the file path to generate.")
}

type sshFlag struct {
	bit  int
	path string
}

func (r *sshFlag) Init() {
	if r.bit <= 2048 {
		r.bit = 2048
	}
	if validator.ValidFile(r.path) || validator.ValidFile(r.path+".pub") {
		PrintString(r.path + " exist.")
		r.path += "_new"
		PrintString("Use " + r.path)
	}
}

func (r *sshFlag) Run(_ *cobra.Command, _ []string) {
	r.Init()
	rsaFile := r.path
	pubFile := r.path + ".pub"
	/* Generate rsa keypair. */
	key, err := rsa.GenerateKey(rand.Reader, r.bit)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	/* Encode and write private key to file. */
	privateKey, err := Encoder.PemEncode(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = os.WriteFile(rsaFile, []byte(privateKey), common.FileModeROwner); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	/* Marshal public key. */
	pub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	publicKey := ssh.MarshalAuthorizedKey(pub)
	if err = os.WriteFile(pubFile, publicKey, common.FileModeROwner); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
