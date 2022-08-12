package test_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDig(t *testing.T) {
	const subCommand = "dig"
	testCases := []struct {
		input    []string
		expected interface{}
	}{
		{[]string{runCommand, mainGo, subCommand, "1.1.1.1", "PTR"}, "one.one.one.one."},
		{[]string{runCommand, mainGo, subCommand, "apple.com", "@8.8.8.8"}, "17.253.144.10"},
		{[]string{runCommand, mainGo, subCommand, "@8.8.8.8", "google.com", "CNAME"}, ""},
		{[]string{runCommand, mainGo, subCommand, "CNAME", "@8.8.8.8", "tw.yahoo.com"}, "fp-ycpi.g03.yahoodns.net."},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			out, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(err)
			}
			outString := string(out)
			split := strings.Fields(outString)
			var got string
			if l := len(split); l != 0 {
				got = split[l-1]
			} else if l == 0 && outString == "" {
				got = outString
			}
			assert.Equal(t, testCases[i].expected, got)
		})
	}
}
