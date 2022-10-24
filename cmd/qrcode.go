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

func initQrcode() *cobra.Command {
	var flags struct {
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
			if !validator.IsFile(args[0]) {
				return common.ErrInvalidArg
			}
			result, err := common.ReadQRCode(args[0])
			if err != nil {
				return err
			}
			printer.Printf(result)
			return nil
		},
		Example: common.Examples(`# Read QR code and print message
qrcode.png`, CommandQrcode, CommandRead),
		DisableFlagsInUseLine: true,
	}

	runE := func(cmd *cobra.Command, args []string) error {
		switch cmd.Name() {
		case CommandText:
			flags.text = common.SliceStringToString(args, " ")
			flags.text = strings.TrimRight(flags.text, " ")
		case CommandOTP:
			if flags.otpSecret == "" {
				return common.ErrInvalidArg
			}
			switch flags.otpDigits {
			case 6, 7, 8:
			default:
				return common.ErrInvalidArg
			}
			switch flags.otpPeriod {
			case 15, 30, 60:
			default:
				return common.ErrInvalidArg
			}
			alg := strings.ToUpper(flags.otpAlgorithm)
			switch alg {
			case "SHA1", "SHA256", "SHA512":
			default:
				return common.ErrInvalidArg
			}

			flags.text = fmt.Sprintf(`otpauth://totp/%s:%s?secret=%s&issuer=%s&period=%d&algorithm=%s&digits=%d`,
				flags.otpIssuer, flags.otpAccount, flags.otpSecret, flags.otpIssuer, flags.otpPeriod, alg, flags.otpDigits)
		case CommandWiFi:
			if flags.wifiSsid == "" {
				return common.ErrInvalidArg
			}

			flags.text = fmt.Sprintf(`WIFI:S:%s;T:%s;P:%s;;`,
				flags.wifiSsid, flags.wifiType, flags.wifiPass)
		}
		return common.GenerateQRCode(flags.text, flags.size, flags.output)
	}

	var qrcodeSubCmdText = &cobra.Command{
		Use:   CommandText,
		Args:  cobra.MinimumNArgs(1),
		Short: "Generate QR code with text",
		RunE:  runE,
		Example: common.Examples(`# Generate QR code with text
https://www.google.com -o out.png
https://www.google.com -o out.png -s 500`, CommandQrcode, CommandText),
	}

	var qrcodeSubCmdOtp = &cobra.Command{
		Use:   CommandOTP,
		Short: "Generate OTP QR code",
		RunE:  runE,
		Example: common.Examples(`# Generate OTP QR code
--otp-account my@gmail.com --otp-secret fqowefilkjfoqwie --otp-issuer aws`,
			CommandQrcode, CommandOTP),
	}

	var qrcodeSubCmdWifi = &cobra.Command{
		Use:   CommandWiFi,
		Short: "Generate WiFi QR code",
		RunE:  runE,
		Example: common.Examples(`# Generate WiFi QR code
--wifi-type WPA --wifi-pass your_password --wifi-ssid your_wifi_ssid -o wifi.png`,
			CommandQrcode, CommandWiFi),
	}

	/* output arguments */
	qrcodeCmd.PersistentFlags().StringVarP(&flags.output, "output", "o", "qrcode.png", common.Usage("Output QRCode file path"))
	qrcodeCmd.PersistentFlags().IntVarP(&flags.size, "size", "s", 600, common.Usage("Specify QRCode generate size"))
	qrcodeCmd.AddCommand(qrcodeSubCmdRead)
	qrcodeCmd.AddCommand(qrcodeSubCmdText)
	qrcodeCmd.AddCommand(qrcodeSubCmdOtp)
	qrcodeCmd.AddCommand(qrcodeSubCmdWifi)

	/* Type: wifi */
	qrcodeSubCmdWifi.Flags().StringVarP(&flags.wifiPass, "wifi-pass", "", "", common.Usage("Specify password"))
	qrcodeSubCmdWifi.Flags().StringVarP(&flags.wifiSsid, "wifi-ssid", "", "", common.Usage("Specify SSID"))
	qrcodeSubCmdWifi.Flags().StringVarP(&flags.wifiType, "wifi-type", "", "WPA", common.Usage("WPA/WEP/nopass"))
	qrcodeSubCmdWifi.MarkFlagsRequiredTogether("wifi-pass", "wifi-ssid")
	/* Type: otp */
	qrcodeSubCmdOtp.Flags().StringVarP(&flags.otpAccount, "otp-account", "", "", common.Usage("Specify account"))
	qrcodeSubCmdOtp.Flags().StringVarP(&flags.otpIssuer, "otp-issuer", "", "", common.Usage("Specify issuer"))
	qrcodeSubCmdOtp.Flags().StringVarP(&flags.otpSecret, "otp-secret", "", "", common.Usage("Specify secret"))
	qrcodeSubCmdOtp.Flags().StringVarP(&flags.otpAlgorithm, "otp-algorithm", "", "SHA1", common.Usage("Specify algorithm"))
	qrcodeSubCmdOtp.Flags().IntVarP(&flags.otpDigits, "otp-digits", "", 6, common.Usage("Specify digits"))
	qrcodeSubCmdOtp.Flags().IntVarP(&flags.otpPeriod, "otp-period", "", 30, common.Usage("Specify period"))
	qrcodeSubCmdOtp.MarkFlagsRequiredTogether("otp-account", "otp-issuer", "otp-secret")
	return qrcodeCmd
}
