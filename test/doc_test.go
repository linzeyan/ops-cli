package test_test

import (
	"os"
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
		{[]string{runCommand, mainGo, subCommand, "man", "-d", "doc"}, "doc/ops-cli.3"},
		{[]string{runCommand, mainGo, subCommand, "markdown", "-d", "doc"}, "doc/ops-cli.md"},
		{[]string{runCommand, mainGo, subCommand, "rest", "-d", "doc"}, "doc/ops-cli.rst"},
		{[]string{runCommand, mainGo, subCommand, "yaml", "-d", "doc"}, "doc/ops-cli.yaml"},
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
	_ = os.RemoveAll("doc")
}

func TestBinaryDoc(t *testing.T) {
	const subCommand = "doc"
	args := []string{"man", "markdown", "rest", "yaml"}

	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i], "-d", "doc").Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
	_ = os.RemoveAll("doc")
}
