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
			err := exec.Command(mainCommand, testCases[i].input...).Run()
			if err != nil {
				t.Error(err)
			}
			assert.FileExists(t, testCases[i].expected)
		})
	}
	_ = exec.Command("rm", "rf", "/tmp/doc").Run()
}

func TestBinaryDoc(t *testing.T) {
	const subCommand = "doc"
	dir := "/tmp/doc"
	args := []string{"man", "markdown", "rest", "yaml"}

	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i], "-d", dir).Run(); err != nil {
				t.Error(err)
			}
		})
	}
	_ = exec.Command("rm", "-rf", dir).Run()
}
