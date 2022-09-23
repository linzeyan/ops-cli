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
	"strings"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/spf13/cobra"
)

func init() {
	var qrcodeFlag QrcodeFlag
	var qrcodeCmd = &cobra.Command{
		Use:   CommandQrcode,
		Short: "Read or generate QR Code",
		Run:   func(cmd *cobra.Command, _ []string) { _ = cmd.Help() },

		DisableFlagsInUseLine: true,
	}

	var qrcodeSubCmdRead = &cobra.Command{
		Use:   CommandRead,
		Args:  cobra.ExactArgs(1),
		Short: "Read QR code and print message",
		RunE: func(_ *cobra.Command, args []string) error {
			if !validator.ValidFile(args[0]) {
				return common.ErrInvalidArg
			}
			result, err := common.ReadQRCode(args[0])
			if err != nil {
				return err
			}
			PrintString(result)
			return nil
		},
		Example: common.Examples(`# Read QR code and print message
qrcode.png`, CommandQrcode, CommandRead),
		DisableFlagsInUseLine: true,
	}

	var qrcodeSubCmdText = &cobra.Command{
		Use:   CommandText,
		Args:  cobra.MinimumNArgs(1),
		Short: "Generate QR code with text",
		RunE:  qrcodeFlag.GenerateRunE,
		Example: common.Examples(`# Generate QR code with text
https://www.google.com -o out.png
https://www.google.com -o out.png -s 500`, CommandQrcode, CommandText),
	}

	var qrcodeSubCmdOtp = &cobra.Command{
		Use:   CommandOtp,
		Short: "Generate OTP QR code",
		RunE:  qrcodeFlag.GenerateRunE,
		Example: common.Examples(`# Generate OTP QR code
--otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`,
			CommandQrcode, CommandOtp),
	}

	var qrcodeSubCmdWifi = &cobra.Command{
		Use:   CommandWiFi,
		Short: "Generate WiFi QR code",
		RunE:  qrcodeFlag.GenerateRunE,
		Example: common.Examples(`# Generate WiFi QR code
--wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid -o wifi.png`,
			CommandQrcode, CommandWiFi),
	}
	rootCmd.AddCommand(qrcodeCmd)

	/* output arguments */
	qrcodeCmd.PersistentFlags().StringVarP(&qrcodeFlag.output, "output", "o", "qrcode.png", common.Usage("Output QRCode file path"))
	qrcodeCmd.PersistentFlags().IntVarP(&qrcodeFlag.size, "size", "s", 600, common.Usage("Specify QRCode generate size"))
	qrcodeCmd.AddCommand(qrcodeSubCmdRead)
	qrcodeCmd.AddCommand(qrcodeSubCmdText)
	qrcodeCmd.AddCommand(qrcodeSubCmdOtp)
	qrcodeCmd.AddCommand(qrcodeSubCmdWifi)

	/* Type: wifi */
	qrcodeSubCmdWifi.Flags().StringVarP(&qrcodeFlag.wifiPass, "wifi-pass", "", "", common.Usage("Specify password"))
	qrcodeSubCmdWifi.Flags().StringVarP(&qrcodeFlag.wifiSsid, "wifi-ssid", "", "", common.Usage("Specify SSID"))
	qrcodeSubCmdWifi.Flags().StringVarP(&qrcodeFlag.wifiType, "wifi-type", "", "WPA", common.Usage("WPA/WEP/nopass"))
	qrcodeSubCmdWifi.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid")
	/* Type: otp */
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeFlag.otpAccount, "otp-account", "", "", common.Usage("Specify account"))
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeFlag.otpIssuer, "otp-issuer", "", "", common.Usage("Specify issuer"))
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeFlag.otpSecret, "otp-secret", "", "", common.Usage("Specify secret"))
	qrcodeSubCmdOtp.Flags().StringVarP(&qrcodeFlag.otpAlgorithm, "otp-algorithm", "", "SHA1", common.Usage("Specify algorithm"))
	qrcodeSubCmdOtp.Flags().IntVarP(&qrcodeFlag.otpDigits, "otp-digits", "", 6, common.Usage("Specify digits"))
	qrcodeSubCmdOtp.Flags().IntVarP(&qrcodeFlag.otpPeriod, "otp-period", "", 30, common.Usage("Specify period"))
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
	otpAccount, otpSecret, otpIssuer, otpAlgorithm string
	otpPeriod, otpDigits                           int
}

func (qr *QrcodeFlag) GenerateRunE(cmd *cobra.Command, args []string) error {
	switch cmd.Name() {
	case CommandText:
		for i := range args {
			qr.text += args[i]
		}
		return common.GenerateQRCode(qr.text, qr.size, qr.output)
	case CommandOtp:
		if qr.otpSecret == "" {
			return common.ErrInvalidArg
		}
		switch qr.otpDigits {
		case 6, 7, 8:
		default:
			return common.ErrInvalidArg
		}
		alg := strings.ToUpper(qr.otpAlgorithm)
		if common.HashAlgorithm(alg) == nil {
			return common.ErrInvalidArg
		}
		switch qr.otpPeriod {
		case 15, 30, 60:
		default:
			return common.ErrInvalidArg
		}
		qr.text = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s&period=%d&algorithm=%s&digits=%d`,
			qr.otpIssuer, qr.otpAccount, qr.otpSecret, qr.otpIssuer, qr.otpPeriod, alg, qr.otpDigits)
		return common.GenerateQRCode(qr.text, qr.size, qr.output)
	case CommandWiFi:
		if qr.wifiSsid == "" {
			return common.ErrInvalidArg
		}
		qr.text = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`,
			qr.wifiSsid, qr.wifiType, qr.wifiPass)
		return common.GenerateQRCode(qr.text, qr.size, qr.output)
	}
	return common.ErrInvalidArg
}
