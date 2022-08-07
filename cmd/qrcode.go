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

	"github.com/linzeyan/qrcode"
	"github.com/spf13/cobra"
)

var qrcodeCmd = &cobra.Command{
	Use:     "qrcode",
	Aliases: []string{"qr"},
	Short:   "Read or generate QR Code",
	Run:     func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },
}

var qrcodeSubCmdRead = &cobra.Command{
	Use:   "read",
	Args:  cobra.ExactArgs(1),
	Short: "Read QR code and print message",
	Run: func(_ *cobra.Command, args []string) {
		qr.text = args[0]
		if !ValidFile(qr.text) {
			log.Println("file not found")
			return
		}
		result, err := qrcode.ReadQRCode(qr.text)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(result)
	},
	Example: Examples(`# Read QR code and print message
ops-cli qrcode qrcode.png`),
	DisableFlagsInUseLine: true,
}

var qrcodeSubCmdText = &cobra.Command{
	Use:   "text",
	Args:  cobra.MinimumNArgs(1),
	Short: "Generate QR code with text",
	Run: func(_ *cobra.Command, args []string) {
		for i := range args {
			qr.text += args[i]
		}
		err := qr.Generate()
		if err != nil {
			log.Println(err)
		}
	},
	Example: Examples(`# Generate QR code with text
ops-cli qrcode text https://www.google.com -o out.png
ops-cli qrcode text https://www.google.com -o out.png -s 500`),
}

var qrcodeSubCmdOtp = &cobra.Command{
	Use:   "otp",
	Short: "Generate OTP QR code",
	Run: func(cmd *cobra.Command, _ []string) {
		if qr.otpSecret == "" {
			_ = cmd.Help()
			return
		}
		qr.text = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s`, qr.otpIssuer, qr.otpAccount, qr.otpSecret, qr.otpIssuer)
		err := qr.Generate()
		if err != nil {
			log.Println(err)
		}
	},
	Example: Examples(`# Generate OTP QR code
ops-cli qrcode otp --otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`),
}

var qrcodeSubCmdWifi = &cobra.Command{
	Use:   "wifi",
	Short: "Generate WiFi QR code",
	Run: func(cmd *cobra.Command, _ []string) {
		if qr.wifiSsid == "" {
			_ = cmd.Help()
			return
		}
		qr.text = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`, qr.wifiSsid, qr.wifiType, qr.wifiPass)
		err := qr.Generate()
		if err != nil {
			log.Println(err)
		}
	},
	Example: Examples(`# Generate WiFi QR code
ops-cli qrcode wifi --wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid -o wifi.png`),
}

var qr qrcodeFlag

func init() {
	rootCmd.AddCommand(qrcodeCmd)

	/* output arguments */
	qrcodeCmd.PersistentFlags().StringVarP(&qr.output, "output", "o", "./qrcode.png", "Output QRCode file path")
	qrcodeCmd.PersistentFlags().IntVarP(&qr.size, "size", "s", 600, "Specify QRCode generate size")
	qrcodeCmd.AddCommand(qrcodeSubCmdRead)
	qrcodeCmd.AddCommand(qrcodeSubCmdText)
	qrcodeCmd.AddCommand(qrcodeSubCmdOtp)
	qrcodeCmd.AddCommand(qrcodeSubCmdWifi)

	/* Type: wifi */
	qrcodeSubCmdWifi.Flags().StringVarP(&qr.wifiPass, "wifi-pass", "", "", "Specify password")
	qrcodeSubCmdWifi.Flags().StringVarP(&qr.wifiSsid, "wifi-ssid", "", "", "Specify SSID")
	qrcodeSubCmdWifi.Flags().StringVarP(&qr.wifiType, "wifi-type", "", "WPA", "WPA/WEP/nopass")
	qrcodeSubCmdWifi.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid")
	/* Type: otp */
	qrcodeSubCmdOtp.Flags().StringVarP(&qr.otpAccount, "otp-account", "", "", "Specify account")
	qrcodeSubCmdOtp.Flags().StringVarP(&qr.otpIssuer, "otp-issuer", "", "", "Specify issuer")
	qrcodeSubCmdOtp.Flags().StringVarP(&qr.otpSecret, "otp-secret", "", "", "Specify secret")
	qrcodeSubCmdOtp.MarkFlagsRequiredTogether("otp-account", "otp-issuer", "otp-secret")
}

type qrcodeFlag struct {
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

func (q qrcodeFlag) Generate() error {
	if qr.text == "" {
		return errors.New("text is empty")
	}
	err := qrcode.GenerateQRCode(qr.text, qr.size, qr.output)
	if err != nil {
		return err
	}
	return nil
}
