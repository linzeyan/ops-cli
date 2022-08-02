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
			switch {
			case ValidFile(input):
				result, err := qrcode.ReadQRCode(input)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(result)
				return
			case strings.ToLower(input) == "wifi":
				qr.message = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`, qr.wifiSsid, qr.wifiType, qr.wifiPass)
			case strings.ToLower(input) == "otp":
				qr.message = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s`, qr.otpIssuer, qr.otpAccount, qr.otpSecret, qr.otpIssuer)
			}
			fallthrough
		default:
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

var qr qrcodeGeneratePng

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	/* output arguments */
	qrcodeCmd.Flags().StringVarP(&qr.output, "output", "o", "./qrcode.png", "Output QRCode file path")
	qrcodeCmd.Flags().IntVarP(&qr.size, "size", "s", 600, "Specify QRCode generate size")

	/* Type: normal */
	qrcodeCmd.Flags().StringVarP(&qr.message, "message", "m", "", "Input message")
	/* Type: wifi */
	qrcodeCmd.Flags().StringVar(&qr.wifiPass, "wifi-pass", "", "Specify password")
	qrcodeCmd.Flags().StringVar(&qr.wifiSsid, "wifi-ssid", "", "Specify SSID")
	qrcodeCmd.Flags().StringVar(&qr.wifiType, "wifi-type", "WPA", "WPA/WEP/nopass")
	/* Type: otp */
	qrcodeCmd.Flags().StringVar(&qr.otpAccount, "otp-account", "", "Specify account")
	qrcodeCmd.Flags().StringVar(&qr.otpIssuer, "otp-issuer", "", "Specify issuer")
	qrcodeCmd.Flags().StringVar(&qr.otpSecret, "otp-secret", "", "Specify secret")

	qrcodeCmd.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid")
	qrcodeCmd.MarkFlagsRequiredTogether("otp-account", "otp-issuer", "otp-secret")
}

type qrcodeGeneratePng struct {
	/* Bind flags */
	/* QR Code generate file path */
	output string
	/* QR Code generate file size */
	size int
	/* QR Code generate message */
	message string
	/* WiFi */
	wifiType, wifiPass, wifiSsid string
	/* OTP */
	otpAccount, otpSecret, otpIssuer string
}

func (q qrcodeGeneratePng) Generate() error {
	if qr.message == "" {
		return errors.New("message is empty")
	}
	err := qrcode.GenerateQRCode(qr.message, qr.size, qr.output)
	if err != nil {
		return err
	}
	return nil
}
