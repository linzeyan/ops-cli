package test_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCert(t *testing.T) {
	const subCommand = "cert"
	testCases := []struct {
		input    []string
		expected interface{}
	}{
		{[]string{runCommand, mainGo, subCommand, "www.google.com", "--dns"}, "[\n  \"www.google.com\"\n]\n\n"},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "--dns"}, "[\n  \"cloudflare-dns.com\",\n  \"*.cloudflare-dns.com\",\n  \"one.one.one.one\"\n]\n\n"},
		{[]string{runCommand, mainGo, subCommand, "www.google.com", "--issuer"}, "CN=GTS CA 1C3,O=Google Trust Services LLC,C=US\n"},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "--issuer"}, "CN=DigiCert TLS Hybrid ECC SHA384 2020 CA1,O=DigiCert Inc,C=US\n"},
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "--ip"}, "1.1.1.1:443\n"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}
