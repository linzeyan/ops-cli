package test_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	const subCommand = "doc"
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, "man", "-d", "/tmp/doc"}, "/tmp/doc/ops-cli.3"},
		{[]string{runCommand, mainGo, subCommand, "markdown", "-d", "/tmp/doc"}, "/tmp/doc/ops-cli.md"},
		{[]string{runCommand, mainGo, subCommand, "rest", "-d", "/tmp/doc"}, "/tmp/doc/ops-cli.rst"},
		{[]string{runCommand, mainGo, subCommand, "yaml", "-d", "/tmp/doc"}, "/tmp/doc/ops-cli.yaml"},
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
