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

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
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
	RunE: func(_ *cobra.Command, args []string) error {
		if !validator.ValidFile(args[0]) {
			return ErrFileNotFound
		}
		result, err := common.ReadQRCode(args[0])
		if err != nil {
			return err
		}
		PrintString(result)
		return nil
	},
	Example: common.Examples(`# Read QR code and print message
ops-cli qrcode read qrcode.png`),
	DisableFlagsInUseLine: true,
}

var qrcodeSubCmdText = &cobra.Command{
	Use:   "text",
	Args:  cobra.MinimumNArgs(1),
	Short: "Generate QR code with text",
	RunE:  qrcodeCmdGlobalVar.GenerateRunE,
	Example: common.Examples(`# Generate QR code with text
ops-cli qrcode text https://www.google.com -o out.png
ops-cli qrcode text https://www.google.com -o out.png -s 500`),
}

var qrcodeSubCmdOtp = &cobra.Command{
	Use:   "otp",
	Short: "Generate OTP QR code",
	RunE:  qrcodeCmdGlobalVar.GenerateRunE,
	Example: common.Examples(`# Generate OTP QR code
ops-cli qrcode otp --otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`),
}

var qrcodeSubCmdWifi = &cobra.Command{
	Use:   "wifi",
	Short: "Generate WiFi QR code",
	RunE:  qrcodeCmdGlobalVar.GenerateRunE,
	Example: common.Examples(`# Generate WiFi QR code
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

func (qr *QrcodeFlag) GenerateRunE(cmd *cobra.Command, args []string) error {
	var err error
	switch cmd.Name() {
	case "text":
		for i := range args {
			qr.text += args[i]
		}
		err = common.GenerateQRCode(qr.text, qr.size, qr.output)
	case "otp":
		if qr.otpSecret == "" {
			log.Println(ErrArgNotFound)
			os.Exit(1)
		}
		qr.text = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s`,
			qr.otpIssuer, qr.otpAccount, qr.otpSecret, qr.otpIssuer)
		err = common.GenerateQRCode(qr.text, qr.size, qr.output)
	case "wifi":
		if qr.wifiSsid == "" {
			log.Println(ErrArgNotFound)
			os.Exit(1)
		}
		qr.text = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`,
			qr.wifiSsid, qr.wifiType, qr.wifiPass)
		err = common.GenerateQRCode(qr.text, qr.size, qr.output)
	}
	return err
}
