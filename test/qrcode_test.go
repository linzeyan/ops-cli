package test_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQrcode(t *testing.T) {
	const subCommand = "qrcode"
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, "otp", "--otp-account", "my@gmail.com", "--otp-secret", "fqowefilkjfoqwie", "--otp-issuer", "aws", "-o", "otp.png"}, "otp.png"},
		{[]string{runCommand, mainGo, subCommand, "wifi", "--wifi-type", "WPA", "--wifi-pass", "your_password", "--wifi-ssid", "your_wifi_ssid", "-o", "wifi.png", "-s", "500"}, "wifi.png"},
		{[]string{runCommand, mainGo, subCommand, "text", "https://www.google.com", "-o", "text.png"}, "text.png"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			err := exec.Command(mainCommand, testCases[i].input...).Run()
			if err != nil {
				t.Error(err)
			}
			assert.FileExists(t, testCases[i].expected)
		})
	}
}

func TestQrcodeRead(t *testing.T) {
	const subCommand = "qrcode"
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, "read", "assets/example.png"}, "WIFI:S:your_wifi_ssid;T:WPA;P:your_password;;\n"},
		{[]string{runCommand, mainGo, subCommand, "read", "otp.png"}, "otpauth://totp/aws:my@gmail.com?secret=fqowefilkjfoqwie&issuer=aws\n"},
		{[]string{runCommand, mainGo, subCommand, "read", "text.png"}, "https://www.google.com\n"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(err)
			}
			assert.FileExists(t, testCases[i].input[4])
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}

func TestBinaryQrcode(t *testing.T) {
	const subCommand = "qrcode"
	args := [][]string{
		{subCommand, "otp", "--otp-account", "my@gmail.com", "--otp-secret", "fqowefilkjfoqwie", "--otp-issuer", "aws", "-o", "otp.png"},
		{subCommand, "wifi", "--wifi-type", "WPA", "--wifi-pass", "your_password", "--wifi-ssid", "your_wifi_ssid", "-o", "wifi.png", "-s", "500"},
		{subCommand, "text", "https://www.google.com", "-o", "text.png"},
	}
	t.Run("read", func(t *testing.T) {
		if err := exec.Command(binaryCommand, []string{subCommand, "read", "assets/example.png"}...).Run(); err != nil {
			t.Error(err)
		}
	})
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(err)
			}
		})
	}
	_ = exec.Command("rm", "-f", "text.png", "wifi.png", "otp.png").Run()
}