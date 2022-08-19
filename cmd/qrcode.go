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
	"fmt"
	"log"
	"os"

	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/linzeyan/qrcode"
	"github.com/spf13/cobra"
)

var qrcodeCmd = &cobra.Command{
	Use:   "qrcode",
	Short: "Read or generate QR Code",
	Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

	DisableFlagsInUseLine: true,
}

var qrcodeSubCmdRead = &cobra.Command{
	Use:   "read",
	Args:  cobra.ExactArgs(1),
	Short: "Read QR code and print message",
	Run: func(_ *cobra.Command, args []string) {
		qrcodeCmdGlobalVar.text = args[0]
		if !validator.ValidFile(qrcodeCmdGlobalVar.text) {
			log.Println(ErrFileNotFound)
			os.Exit(1)
		}
		result, err := qrcode.ReadQRCode(qrcodeCmdGlobalVar.text)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		PrintString(result)
	},
	Example: Examples(`# Read QR code and print message
ops-cli qrcode read qrcode.png`),
	DisableFlagsInUseLine: true,
}

var qrcodeSubCmdText = &cobra.Command{
	Use:   "text",
	Args:  cobra.MinimumNArgs(1),
	Short: "Generate QR code with text",
	Run: func(_ *cobra.Command, args []string) {
		for i := range args {
			qrcodeCmdGlobalVar.text += args[i]
		}
		if err := qrcodeCmdGlobalVar.Generate(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
	Example: Examples(`# Generate QR code with text
ops-cli qrcode text https://www.google.com -o out.png
ops-cli qrcode text https://www.google.com -o out.png -s 500`),
}

var qrcodeSubCmdOtp = &cobra.Command{
	Use:   "otp",
	Short: "Generate OTP QR code",
	Run: func(_ *cobra.Command, _ []string) {
		if qrcodeCmdGlobalVar.otpSecret == "" {
			log.Println(ErrArgNotFound)
			os.Exit(1)
		}
		qrcodeCmdGlobalVar.text = fmt.Sprintf(
			`otpauth://totp/%s:%s?secret=%s&issuer=%s`,
			qrcodeCmdGlobalVar.otpIssuer,
			qrcodeCmdGlobalVar.otpAccount,
			qrcodeCmdGlobalVar.otpSecret,
			qrcodeCmdGlobalVar.otpIssuer,
		)
		if err := qrcodeCmdGlobalVar.Generate(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
	Example: Examples(`# Generate OTP QR code
ops-cli qrcode otp --otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`),
}

var qrcodeSubCmdWifi = &cobra.Command{
	Use:   "wifi",
	Short: "Generate WiFi QR code",
	Run: func(_ *cobra.Command, _ []string) {
		if qrcodeCmdGlobalVar.wifiSsid == "" {
			log.Println(ErrArgNotFound)
			os.Exit(1)
		}
		qrcodeCmdGlobalVar.text = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`, qrcodeCmdGlobalVar.wifiSsid, qrcodeCmdGlobalVar.wifiType, qrcodeCmdGlobalVar.wifiPass)
		if err := qrcodeCmdGlobalVar.Generate(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
	Example: Examples(`# Generate WiFi QR code
ops-cli qrcode wifi --wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid -o wifi.png`),
}

var qrcodeCmdGlobalVar QrcodeFlag

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	/* output arguments */
	qrcodeCmd.PersistentFlags().StringVarP(&qrcodeCmdGlobalVar.output, "output", "o", "./qrcode.png", "Output QRCode file path")
	qrcodeCmd.PersistentFlags().IntVarP(&qrcodeCmdGlobalVar.size, "size", "s", 600, "Specify QRCode generate size")
	qrcodeCmd.AddCommand(qrcodeSubCmdRead)
	qrcodeCmd.AddCommand(qrcodeSubCmdText)
	qrcodeCmd.AddCommand(qrcodeSubCmdOtp)
	qrcodeCmd.AddCommand(qrcodeSubCmdWifi)

	/* Type: wifi */
	qrcodeSubCmdWifi.Flags().StringVarP(&qrcodeCmdGlobalVar.wifiPass, "wifi-pass", "", "", "Specify password")
	qrcodeSubCmdWifi.Flags().StringVarP(&qrcodeCmdGlobalVar.wifiSsid, "wifi-ssid", "", "", "Specify SSID")
	qrcodeSubCmdWifi.Flags().StringVarP(&qrcodeCmdGlobalVar.wifiType, "wifi-type", "", "WPA", "WPA/WEP/nopass")
	qrcodeSubCmdWifi.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid")
	/* Type: otp */
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeCmdGlobalVar.otpAccount, "otp-account", "", "", "Specify account")
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeCmdGlobalVar.otpIssuer, "otp-issuer", "", "", "Specify issuer")
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeCmdGlobalVar.otpSecret, "otp-secret", "", "", "Specify secret")
	qrcodeSubCmdOtp.MarkFlagsRequiredTogether("otp-account", "otp-issuer", "otp-secret")
}

type QrcodeFlag struct {
	/* Bind flags */
	/* QR Code generate file path */
	output string
	/* QR Code generate file size */
	size int
	/* QR Code generate text */
	text string
	/* WiFi */
	wifiType, wifiPass, wifiSsid string
	/* OTP */
	otpAccount, otpSecret, otpIssuer string
}

func (q *QrcodeFlag) Generate() error {
	if q.text == "" {
		return ErrArgNotFound
	}
	return qrcode.GenerateQRCode(q.text, q.size, q.output)
}
