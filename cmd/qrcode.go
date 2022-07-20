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
	Short: "Read or output QR Code",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if !qrcodeGenerate && qrcodeFileInput != "" {
			result, err := qrcode.ReadQRCode(qrcodeFileInput)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(result)
			return
		}

		if qrcodeGenerate {
			switch strings.ToLower(qrcodeType) {
			case "wifi":
				qrcodeMessage = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`, qrcodeWIFISsid, qrcodeWIFIType, qrcodeWIFIPass)
				err := qrcodeGenPng()
				if err != nil {
					log.Println(err)
				}
			case "otp":
				qrcodeMessage = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s`, qrcodeOtpIssuer, qrcodeOtpAccount, qrcodeOtpSecret, qrcodeOtpIssuer)
				err := qrcodeGenPng()
				if err != nil {
					log.Println(err)
				}
			default:
				err := qrcodeGenPng()
				if err != nil {
					log.Println(err)
				}
			}
			return
		}
		cmd.Help()
	},
	Example: Examples(`# Read QR code and print message
ops-cli qrcode -f qrcode.png

# Generate QR code with message
ops-cli qrcode -g -m https://www.google.com -o out.png
ops-cli qrcode -g -m https://www.google.com -o out.png -s 500

# Generate WiFi QR code
ops-cli qrcode -g -t wifi --wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid

# Generate OTP QR code
ops-cli qrcode -g -t otp --otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`),
}

var qrcodeGenerate bool

var qrcodeFileInput, qrcodeMessage, qrcodeOutput string
var qrcodeSize int

var qrcodeType string
var qrcodeWIFIType, qrcodeWIFIPass, qrcodeWIFISsid string
var qrcodeOtpAccount, qrcodeOtpSecret, qrcodeOtpIssuer string

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	qrcodeCmd.Flags().BoolVarP(&qrcodeGenerate, "generate", "g", false, "Generate QRcode")
	qrcodeCmd.Flags().StringVarP(&qrcodeFileInput, "file", "f", "", "Specify file path to read")

	qrcodeCmd.Flags().StringVarP(&qrcodeType, "type", "t", "normal", "Generate type(normal, otp, WiFi)")
	/* Type: normal */
	qrcodeCmd.Flags().StringVarP(&qrcodeMessage, "message", "m", "", "Input message")
	/* Type: wifi */
	qrcodeCmd.Flags().StringVarP(&qrcodeWIFIPass, "wifi-pass", "", "", "Specify password")
	qrcodeCmd.Flags().StringVarP(&qrcodeWIFISsid, "wifi-ssid", "", "", "Specify SSID")
	qrcodeCmd.Flags().StringVarP(&qrcodeWIFIType, "wifi-type", "", "WPA", "WPA/WEP/nopass")
	/* Type: otp */
	qrcodeCmd.Flags().StringVarP(&qrcodeOtpAccount, "otp-account", "", "", "Specify account")
	qrcodeCmd.Flags().StringVarP(&qrcodeOtpIssuer, "otp-issuer", "", "", "Specify issuer")
	qrcodeCmd.Flags().StringVarP(&qrcodeOtpSecret, "otp-secret", "", "", "Specify secret")
	/* output arguments */
	qrcodeCmd.Flags().StringVarP(&qrcodeOutput, "output", "o", "./qrcode.png", "Output QRCode file path")
	qrcodeCmd.Flags().IntVarP(&qrcodeSize, "size", "s", 600, "Specify QRCode generate size")
	qrcodeCmd.MarkFlagsRequiredTogether("output", "size")
	qrcodeCmd.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid", "wifi-type")
	qrcodeCmd.MarkFlagsRequiredTogether("otp-account", "otp-issuer", "otp-secret")
}

func qrcodeGenPng() error {
	if qrcodeMessage == "" {
		return errors.New("message is empty")
	}
	err := qrcode.GenerateQRCode(qrcodeMessage, qrcodeSize, qrcodeOutput)
	if err != nil {
		return err
	}
	return nil
}
