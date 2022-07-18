/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
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

// otpCmd represents the otp command
var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Calculate passcode",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {
		otpOptions()
	},
	Example: `ops-cli otp -s 6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7
ops-cli otp -g -p 15
ops-cli otp -g -a sha256
ops-cli otp -g -a sha512 -p 15
ops-cli otp -s T7L756M2FEL6CHISIXVSGT4VUDA4ZLIM -p 15 -d 7`,
}

var otpSecret string
var otpGenerateSecret bool
var otpAlgorithm string
var otpPeriod, otpDigits int8

func init() {
	rootCmd.AddCommand(otpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// otpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// otpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	otpCmd.Flags().StringVarP(&otpSecret, "secret", "s", "", "Specify TOTP secret key")
	otpCmd.Flags().BoolVarP(&otpGenerateSecret, "generate", "g", false, "Generate secret key")
	otpCmd.Flags().StringVarP(&otpAlgorithm, "algorithm", "a", "SHA1", "The hash algorithm used by the credential(SHA1/SHA256/SHA512)")
	otpCmd.Flags().Int8VarP(&otpPeriod, "period", "p", 30, "The period parameter defines a validity period in seconds for the TOTP code(15/30/60)")
	otpCmd.Flags().Int8VarP(&otpDigits, "digits", "d", 6, "The number of digits in a one-time password(6/7/8)")
}

func otpOptions() {
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
	}
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
