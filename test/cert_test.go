package test_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCert(t *testing.T) {
	testCases := []struct {
		input    []string
		expected interface{}
	}{
		{[]string{"run", mainGo, "cert", "www.google.com", "--dns"}, "[www.google.com]\n"},
		{[]string{"run", mainGo, "cert", "1.1.1.1", "--dns"}, "[cloudflare-dns.com *.cloudflare-dns.com one.one.one.one]\n"},
		{[]string{"run", mainGo, "cert", "www.google.com", "--issuer"}, "CN=GTS CA 1C3,O=Google Trust Services LLC,C=US\n"},
		{[]string{"run", mainGo, "cert", "1.1.1.1", "--issuer"}, "CN=DigiCert TLS Hybrid ECC SHA384 2020 CA1,O=DigiCert Inc,C=US\n"},
		{[]string{"run", mainGo, "cert", "1.1.1.1", "--ip"}, "1.1.1.1:443\n"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command("go", testCases[i].input...).Output()
			if err != nil {
				t.Errorf("%s exec error:%s", testCases[i].input[3], err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}
