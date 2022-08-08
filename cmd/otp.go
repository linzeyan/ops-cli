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
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Calculate passcode or generate secret",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var otpSubCmdCalculate = &cobra.Command{
	Use:   "calculate",
	Args:  cobra.MinimumNArgs(1),
	Short: "Calculate passcode",
	Run:   of.Run,
	Example: Examples(`# Calculate the passcode for the specified secret
ops-cli otp calculate 6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7
ops-cli otp calculate 6BDR T7AT RRCZ V5IS FLOH AHQL YF4Z ORG7

# Calculate the passcode of the specified secret, the period is 15 seconds, and the number of digits is 7
ops-cli otp calculate T7L756M2FEL6CHISIXVSGT4VUDA4ZLIM -p 15 -d 7`),
}

var otpSubCmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generate otp secret",
	Run:   of.Run,
	Example: Examples(`# Generate OTP and specify a period of 15 seconds
ops-cli otp generate -p 15

# Generate OTP and specify SHA256 algorithm
ops-cli otp generate -a sha256

# Generate OTP and specify SHA512 algorithm, the period is 15 seconds
ops-cli otp generate -a sha512 -p 15`),
}

var of otpFlag

func init() {
	rootCmd.AddCommand(otpCmd)

	otpCmd.PersistentFlags().StringVarP(&of.alg, "algorithm", "a", "SHA1", "The hash algorithm used by the credential(SHA1/SHA256/SHA512)")
	otpCmd.PersistentFlags().Int8VarP(&of.period, "period", "p", 30, "The period parameter defines a validity period in seconds for the TOTP code(15/30/60)")
	otpCmd.PersistentFlags().Int8VarP(&of.digit, "digits", "d", 6, "The number of digits in a one-time password(6/7/8)")

	otpCmd.AddCommand(otpSubCmdCalculate, otpSubCmdGenerate)
}

type otpFlag struct {
	/* Bind flags */
	/* The period parameter defines a validity period in seconds */
	period int8
	/* The number of digits in a one-time password */
	digit int8
	/* The algorithm of OTP */
	alg string
}

func (o *otpFlag) SetTimeInterval() int64 {
	var timeInterval int64
	switch o.period {
	case 15:
		timeInterval = rootNow.Unix() / 15
	case 60:
		timeInterval = rootNow.Unix() / 60
	default:
		timeInterval = rootNow.Unix() / 30
	}
	return timeInterval
}

func (o *otpFlag) SetDigits() [2]int {
	var digits [2]int
	switch o.digit {
	case 7:
		digits = [2]int{7, 10000000}
	case 8:
		digits = [2]int{8, 100000000}
	default:
		digits = [2]int{6, 1000000}
	}
	return digits
}

func (o *otpFlag) SetAlgorithm() func() hash.Hash {
	var algorithm func() hash.Hash
	switch strings.ToLower(o.alg) {
	case "sha256":
		algorithm = sha256.New
	case "sha512":
		algorithm = sha512.New
	default:
		algorithm = sha1.New
	}
	return algorithm
}

func (o *otpFlag) HOTP(secret string, timeInterval int64) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeInterval))
	hasher := hmac.New(o.SetAlgorithm(), key)
	hasher.Write(buf)
	h := hasher.Sum(nil)
	offset := h[len(h)-1] & 0xf
	r := bytes.NewReader(h[offset : offset+4])

	var data uint32
	err = binary.Read(r, binary.BigEndian, &data)
	if err != nil {
		return "", err
	}
	var digits = o.SetDigits()
	h12 := (int(data) & 0x7fffffff) % digits[1]
	passcode := strconv.Itoa(h12)

	length := len(passcode)
	if length == digits[0] {
		return passcode, nil
	}
	for i := (digits[0] - length); i > 0; i-- {
		passcode = "0" + passcode
	}
	return passcode, nil
}

func (o *otpFlag) TOTP(secret string) (string, error) {
	return o.HOTP(secret, o.SetTimeInterval())
}

func (o *otpFlag) GenSecret() (string, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, o.SetTimeInterval())
	if err != nil {
		return "", err
	}
	hasher := hmac.New(o.SetAlgorithm(), buf.Bytes())
	secret := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return secret, nil
}

func (o *otpFlag) Verify(secret string, input string) (bool, error) {
	passcode, err := o.TOTP(secret)
	if err != nil {
		return passcode == input, err
	}
	return passcode == input, nil
}

func (o *otpFlag) Run(cmd *cobra.Command, args []string) {
	switch cmd.Name() {
	case "calculate":
		var secret string
		switch l := len(args); {
		case l == 1:
			secret = args[0]
		/* Merge all strings. */
		case l > 1:
			for i := range args {
				secret += args[i]
			}
		}
		var result, err = o.TOTP(secret)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(result)
	case "generate":
		var secret, err = o.GenSecret()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(secret)
	}
}
