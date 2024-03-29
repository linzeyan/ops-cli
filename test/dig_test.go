package test_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestDig(t *testing.T) {
	const subCommand = cmd.CommandDig
	testCases := []struct {
		input    []string
		expected any
	}{
		{[]string{runCommand, mainGo, subCommand, "@8.8.8.8", "1.1.1.1", "PTR"}, "one.one.one.one."},
		{[]string{runCommand, mainGo, subCommand, "apple.com", "@8.8.8.8"}, "17.253.144.10"},
		{[]string{runCommand, mainGo, subCommand, "@8.8.8.8", testHost, "CNAME"}, ""},
		{[]string{runCommand, mainGo, subCommand, "CNAME", "@8.8.8.8", "tw.yahoo.com"}, "fp-ycpi.g03.yahoodns.net."},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			out, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
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

func TestDigBinary(t *testing.T) {
	const subCommand = cmd.CommandDig
	servers := []string{"@1.1.1.1", "@8.8.8.8"}
	for _, server := range servers {
		t.Run(server, func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, testHost, server).Run(); err != nil {
				t.Error(server, err)
			}
		})
	}

	args := []string{"A", "AAAA", "CNAME", "NS", "ANY"}
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, testHost, "@8.8.8.8", args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
