package test_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

var files = []string{"text.png", "wifi.png", "otp.png"}

func TestQrcode(t *testing.T) {
	const subCommand = cmd.CommandQrcode
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, cmd.CommandOtp, "--otp-account", "my@gmail.com", "--otp-secret", "fqowefilkjfoqwie", "--otp-issuer", "aws", "-o", files[2]}, files[2]},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandWiFi, "--wifi-type", "WPA", "--wifi-pass", "your_password", "--wifi-ssid", "your_wifi_ssid", "-o", files[1], "-s", "500"}, files[1]},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandText, "https://www.google.com", "-o", files[0]}, files[0]},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			err := exec.Command(mainCommand, testCases[i].input...).Run()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.FileExists(t, testCases[i].expected)
		})
	}
}

func TestQrcodeRead(t *testing.T) {
	const subCommand = cmd.CommandQrcode
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, "read", "assets/example.png"}, "WIFI:S:your_wifi_ssid;T:WPA;P:your_password;;\n"},
		{[]string{runCommand, mainGo, subCommand, "read", files[2]}, "otpauth://totp/aws:my@gmail.com?secret=fqowefilkjfoqwie&issuer=aws\n"},
		{[]string{runCommand, mainGo, subCommand, "read", files[0]}, "https://www.google.com\n"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.FileExists(t, testCases[i].input[4])
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}

func TestBinaryQrcode(t *testing.T) {
	const subCommand = cmd.CommandQrcode
	t.Run("read", func(t *testing.T) {
		if err := exec.Command(binaryCommand, []string{subCommand, "read", "assets/example.png"}...).Run(); err != nil {
			t.Error(err)
		}
	})
	args := [][]string{
		{subCommand, cmd.CommandOtp, "--otp-account", "my@gmail.com", "--otp-secret", "fqowefilkjfoqwie", "--otp-issuer", "aws", "-o", files[2]},
		{subCommand, cmd.CommandWiFi, "--wifi-type", "WPA", "--wifi-pass", "your_password", "--wifi-ssid", "your_wifi_ssid", "-o", files[1], "-s", "500"},
		{subCommand, cmd.CommandText, "https://www.google.com", "-o", files[0]},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
	for _, v := range files {
		_ = os.Remove(v)
	}
}
