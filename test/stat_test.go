package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/linzeyan/ops-cli/cmd/validator"
)

func TestStatBinary(t *testing.T) {
	if validator.IsWindows() {
		return
	}
	const subCommand = cmd.CommandStat
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, ".").Run(); err != nil {
			t.Error(err)
		}
	})
}
