package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestDfBinary(t *testing.T) {
	const subCommand = cmd.CommandDf
	t.Run(subCommand, func(t *testing.T) {
		out, err := exec.Command(binaryCommand, subCommand).Output()
		if err != nil {
			t.Error(string(out))
			t.Error(err)
		}
	})
}
