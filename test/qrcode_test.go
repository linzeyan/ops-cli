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
		{[]string{runCommand, mainGo, subCommand, "otp", "--otp-account", "my@gmail.com", "--otp-secret", "fqowefilkjfoqwie", "--otp-issuer", "aws", "-o", "/tmp/otp.png"}, "/tmp/otp.png"},
		{[]string{runCommand, mainGo, subCommand, "wifi", "--wifi-type", "WPA", "--wifi-pass", "your_password", "--wifi-ssid", "your_wifi_ssid", "-o", "/tmp/wifi.png", "-s", "500"}, "/tmp/wifi.png"},
		{[]string{runCommand, mainGo, subCommand, "text", "https://www.google.com", "-o", "/tmp/text.png"}, "/tmp/text.png"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			_, err := exec.Command(mainCommand, testCases[i].input...).Output()
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
		{[]string{runCommand, mainGo, subCommand, "read", "/tmp/otp.png"}, "otpauth://totp/aws:my@gmail.com?secret=fqowefilkjfoqwie&issuer=aws\n"},
		{[]string{runCommand, mainGo, subCommand, "read", "/tmp/text.png"}, "https://www.google.com\n"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
			// assert.FileExists(t, testCases[i].expected)
		})
	}
}
