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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/linzeyan/qrcode"
	"github.com/spf13/cobra"
)

var qrcodeCmd = &cobra.Command{
	Use:   "qrcode",
	Short: "Read or generate QR Code",
	Run: func(cmd *cobra.Command, args []string) {
		switch l := len(args); {
		case l != 1:
			_ = cmd.Help()
			return
		case l == 1:
			input := args[0]
			_, err := os.Stat(input)
			switch {
			case err == nil:
				result, err := qrcode.ReadQRCode(input)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(result)
				return
			case strings.ToLower(input) == "wifi":
				qrcodeMessage = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`, qrcodeWIFISsid, qrcodeWIFIType, qrcodeWIFIPass)
			case strings.ToLower(input) == "otp":
				qrcodeMessage = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s`, qrcodeOtpIssuer, qrcodeOtpAccount, qrcodeOtpSecret, qrcodeOtpIssuer)
			}
			fallthrough
		default:
			var qr qrcodeGeneratePng
			err := qr.Generate()
			if err != nil {
				log.Println(err)
			}
		}
	},
	Example: Examples(`# Read QR code and print message
ops-cli qrcode qrcode.png

# Generate QR code with message
ops-cli qrcode msg -m https://www.google.com -o out.png
ops-cli qrcode msg -m https://www.google.com -o out.png -s 500

# Generate WiFi QR code
ops-cli qrcode wifi --wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid -o wifi.png

# Generate OTP QR code
ops-cli qrcode otp --otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`),
}

var qrcodeMessage, qrcodeOutput string
var qrcodeSize int

var qrcodeWIFIType, qrcodeWIFIPass, qrcodeWIFISsid string
var qrcodeOtpAccount, qrcodeOtpSecret, qrcodeOtpIssuer string

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	/* output arguments */
	qrcodeCmd.Flags().StringVarP(&qrcodeOutput, "output", "o", "./qrcode.png", "Output QRCode file path")
	qrcodeCmd.Flags().IntVarP(&qrcodeSize, "size", "s", 600, "Specify QRCode generate size")

	/* Type: normal */
	qrcodeCmd.Flags().StringVarP(&qrcodeMessage, "message", "m", "", "Input message")
	/* Type: wifi */
	qrcodeCmd.Flags().StringVar(&qrcodeWIFIPass, "wifi-pass", "", "Specify password")
	qrcodeCmd.Flags().StringVar(&qrcodeWIFISsid, "wifi-ssid", "", "Specify SSID")
	qrcodeCmd.Flags().StringVar(&qrcodeWIFIType, "wifi-type", "WPA", "WPA/WEP/nopass")
	/* Type: otp */
	qrcodeCmd.Flags().StringVar(&qrcodeOtpAccount, "otp-account", "", "Specify account")
	qrcodeCmd.Flags().StringVar(&qrcodeOtpIssuer, "otp-issuer", "", "Specify issuer")
	qrcodeCmd.Flags().StringVar(&qrcodeOtpSecret, "otp-secret", "", "Specify secret")

	qrcodeCmd.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid")
	qrcodeCmd.MarkFlagsRequiredTogether("otp-account", "otp-issuer", "otp-secret")
}

type qrcodeGeneratePng struct{}

func (q qrcodeGeneratePng) Generate() error {
	if qrcodeMessage == "" {
		return errors.New("message is empty")
	}
	err := qrcode.GenerateQRCode(qrcodeMessage, qrcodeSize, qrcodeOutput)
	if err != nil {
		return err
	}
	return nil
}
