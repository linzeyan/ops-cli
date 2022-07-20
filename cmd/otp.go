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
	"time"

	"github.com/spf13/cobra"
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Calculate passcode",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if !otpGenerateSecret && otpSecret != "" {
			var result, err = otpTOTP(otpSecret)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(result)
			return
		}
		if otpGenerateSecret {
			var secret, err = otpGenSecret()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(secret)
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Calculate the passcode for the specified secret
ops-cli otp -s 6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7

# Generate OTP and specify a period of 15 seconds
ops-cli otp -g -p 15

# Generate OTP and specify SHA256 algorithm
ops-cli otp -g -a sha256

# Generate OTP and specify SHA512 algorithm, the period is 15 seconds
ops-cli otp -g -a sha512 -p 15

# Calculate the passcode of the specified secret, the period is 15 seconds, and the number of digits is 7
ops-cli otp -s T7L756M2FEL6CHISIXVSGT4VUDA4ZLIM -p 15 -d 7`),
}

var otpSecret string
var otpGenerateSecret bool
var otpAlgorithm string
var otpPeriod, otpDigits int8

func init() {
	rootCmd.AddCommand(otpCmd)

	otpCmd.Flags().StringVarP(&otpSecret, "secret", "s", "", "Specify TOTP secret key")
	otpCmd.Flags().BoolVarP(&otpGenerateSecret, "generate", "g", false, "Generate secret key")
	otpCmd.Flags().StringVarP(&otpAlgorithm, "algorithm", "a", "SHA1", "The hash algorithm used by the credential(SHA1/SHA256/SHA512)")
	otpCmd.Flags().Int8VarP(&otpPeriod, "period", "p", 30, "The period parameter defines a validity period in seconds for the TOTP code(15/30/60)")
	otpCmd.Flags().Int8VarP(&otpDigits, "digits", "d", 6, "The number of digits in a one-time password(6/7/8)")
}

func otpSetTimeInterval() int64 {
	t := time.Now().Local().Unix()
	switch otpPeriod {
	case 15:
		t = t / 15
	case 60:
		t = t / 60
	default:
		t = t / 30
	}
	return t
}

func otpSetDigits() [2]int {
	var digits [2]int
	switch otpDigits {
	case 7:
		digits = [2]int{7, 10000000}
	case 8:
		digits = [2]int{8, 100000000}
	default:
		digits = [2]int{6, 1000000}
	}
	return digits
}

func otpSetAlgorithm() func() hash.Hash {
	var alg func() hash.Hash
	switch strings.ToLower(otpAlgorithm) {
	case "sha256":
		alg = sha256.New
	case "sha512":
		alg = sha512.New
	default:
		alg = sha1.New
	}
	return alg
}

func otpHOTP(secret string, timeInterval int64) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeInterval))
	alg := otpSetAlgorithm()
	hasher := hmac.New(alg, key)
	hasher.Write(buf)
	h := hasher.Sum(nil)
	offset := h[len(h)-1] & 0xf
	r := bytes.NewReader(h[offset : offset+4])

	var data uint32
	err = binary.Read(r, binary.BigEndian, &data)
	if err != nil {
		return "", err
	}
	var digits [2]int = otpSetDigits()
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

func otpTOTP(secret string) (string, error) {
	return otpHOTP(secret, otpSetTimeInterval())
}

func otpGenSecret() (string, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, otpSetTimeInterval())
	if err != nil {
		return "", err
	}
	alg := otpSetAlgorithm()
	hasher := hmac.New(alg, buf.Bytes())
	secret := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return secret, nil
}

// func otpVerify(secret string, input string) (bool, error) {
// 	passcode, err := otpTOTP(secret)
// 	if err != nil {
// 		return passcode == input, err
// 	}
// 	return passcode == input, nil
// }
