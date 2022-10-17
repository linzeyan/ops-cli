package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestPingBinary(t *testing.T) {
	const subCommand = cmd.CommandPing
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, "--help").Run(); err != nil {
			t.Error(err)
		}
	})
}
