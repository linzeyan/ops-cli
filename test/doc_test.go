package test_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	const subCommand = cmd.CommandDoc
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, cmd.CommandMan, "-d", "doc"}, "doc/ops-cli.3"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandMarkdown, "-d", "doc"}, "doc/ops-cli.md"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandReST, "-d", "doc"}, "doc/ops-cli.rst"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandYaml, "-d", "doc"}, "doc/ops-cli.yaml"},
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

func TestDocBinary(t *testing.T) {
	const subCommand = cmd.CommandDoc
	args := []string{cmd.CommandMan, cmd.CommandMarkdown, cmd.CommandReST, cmd.CommandYaml}

	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i], "-d", "doc").Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
	_ = os.RemoveAll("doc")
}
