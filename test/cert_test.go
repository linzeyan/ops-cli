package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestCert(t *testing.T) {
	const subCommand = cmd.CommandCert
	testCases := []struct {
		input    []string
		expected any
	}{
		{[]string{runCommand, mainGo, subCommand, "www.google.com", "--dns"}, "[\n  \"www.google.com\"\n]\n"},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "--dns"}, "[\n  \"cloudflare-dns.com\",\n  \"*.cloudflare-dns.com\",\n  \"one.one.one.one\"\n]\n"},
		{[]string{runCommand, mainGo, subCommand, "www.google.com", "--issuer"}, "CN=WR2,O=Google Trust Services,C=US\n"},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "--issuer"}, "CN=DigiCert Global G2 TLS RSA SHA256 2020 CA1,O=DigiCert Inc,C=US\n"},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "--ip"}, "1.1.1.1:443\n"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}

func TestCertBinary(t *testing.T) {
	const subCommand = cmd.CommandCert
	t.Run(testHost, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, testHost).Run(); err != nil {
			t.Error(err)
		}
	})

	args := []string{"--days", "--dns", "--expiry", "--ip", "--issuer"}
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, testHost, args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
