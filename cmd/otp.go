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
	"encoding/binary"
	"hash"
	"strconv"
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/spf13/cobra"
)

func initOtp() *cobra.Command {
	var flags struct {
		/* Bind flags */
		/* The period parameter defines a validity period in seconds */
		period int8
		/* The number of digits in a one-time password */
		digit int8
		/* The algorithm of OTP */
		alg string
	}
	var otpCmd = &cobra.Command{
		Use:   CommandOtp,
		Short: "Calculate passcode or generate secret",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}
	runE := func(cmd *cobra.Command, args []string) error {
		o := OTP{
			Period:    flags.period,
			Digit:     flags.digit,
			Algorithm: flags.alg,
		}
		var result string
		var err error
		switch cmd.Name() {
		case CommandCalculate:
			var secret string
			switch l := len(args); {
			case l == 1:
				secret = common.StringRemoveSpaces(args[0])
			/* Merge all strings. */
			case l > 1:
				secret = common.SliceStringToStringNoSpaces(args)
			}
			result, err = o.TOTP(secret)
		case CommandGenerate:
			result, err = o.GenSecret()
		}
		if err != nil {
			return err
		}
		PrintString(result)
		return err
	}

	var otpSubCmdCalculate = &cobra.Command{
		Use:   CommandCalculate,
		Args:  cobra.MinimumNArgs(1),
		Short: "Calculate passcode",
		RunE:  runE,
		Example: common.Examples(`# Calculate the passcode for the specified secret
6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7
6BDR T7AT RRCZ V5IS FLOH AHQL YF4Z ORG7

# Calculate the passcode of the specified secret, the period is 15 seconds, and the number of digits is 7
T7L756M2FEL6CHISIXVSGT4VUDA4ZLIM -p 15 -d 7`, CommandOtp, CommandCalculate),
	}

	var otpSubCmdGenerate = &cobra.Command{
		Use:   CommandGenerate,
		Short: "Generate otp secret",
		RunE:  runE,
		Example: common.Examples(`# Generate OTP and specify a period of 15 seconds
-p 15

# Generate OTP and specify SHA256 algorithm
-a sha256

# Generate OTP and specify SHA512 algorithm, the period is 15 seconds
-a sha512 -p 15`, CommandOtp, CommandGenerate),
	}

	otpCmd.PersistentFlags().StringVarP(&flags.alg, "algorithm", "a", "SHA1", common.Usage("The hash algorithm used by the credential(SHA1/SHA256/SHA512)"))
	otpCmd.PersistentFlags().Int8VarP(&flags.period, "period", "p", 30, common.Usage("The period parameter defines a validity period in seconds for the TOTP code(15/30/60)"))
	otpCmd.PersistentFlags().Int8VarP(&flags.digit, "digits", "d", 6, common.Usage("The number of digits in a one-time password(6/7/8)"))

	otpCmd.AddCommand(otpSubCmdCalculate, otpSubCmdGenerate)
	return otpCmd
}

type OTP struct {
	Period    int8
	Digit     int8
	Algorithm string
}

func (o *OTP) SetTimeInterval() int64 {
	switch o.Period {
	case 15:
		return common.TimeNow.Unix() / 15
	case 60:
		return common.TimeNow.Unix() / 60
	default:
		return common.TimeNow.Unix() / 30
	}
}

func (o *OTP) SetDigits() [2]int {
	switch o.Digit {
	case 7:
		return [2]int{7, 10000000}
	case 8:
		return [2]int{8, 100000000}
	default:
		return [2]int{6, 1000000}
	}
}

func (o *OTP) SetAlgorithm() func() hash.Hash {
	switch strings.ToLower(o.Algorithm) {
	case HashSha256:
		return sha256.New
	case HashSha512:
		return sha512.New
	default:
		return sha1.New
	}
}

func (o *OTP) HOTP(secret string, timeInterval int64) (string, error) {
	key, err := Encoder.Base32StdDecode(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeInterval))
	hasher := hmac.New(o.SetAlgorithm(), key)
	_, err = hasher.Write(buf)
	if err != nil {
		return "", err
	}
	h := hasher.Sum(nil)
	offset := h[len(h)-1] & 0xf
	r := bytes.NewReader(h[offset : offset+4])

	var data uint32
	if err = binary.Read(r, binary.BigEndian, &data); err != nil {
		return "", err
	}
	var digits = o.SetDigits()
	h12 := (int(data) & 0x7fffffff) % digits[1]
	passcode := strconv.Itoa(h12)

	length := len(passcode)
	if length == digits[0] {
		return passcode, err
	}
	for i := (digits[0] - length); i > 0; i-- {
		passcode = "0" + passcode
	}
	return passcode, err
}

func (o *OTP) TOTP(secret string) (string, error) {
	return o.HOTP(secret, o.SetTimeInterval())
}

func (o *OTP) GenSecret() (string, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, o.SetTimeInterval())
	if err != nil {
		return "", err
	}
	hasher := hmac.New(o.SetAlgorithm(), buf.Bytes())
	return Encoder.Base32StdEncode(hasher.Sum(nil))
}

func (o *OTP) Verify(input, secret string) (bool, error) {
	passcode, err := o.TOTP(secret)
	return passcode == input, err
}
