package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestLINEBinary(t *testing.T) {
	const subCommand = cmd.CommandLINE
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})
}
