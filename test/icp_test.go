package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestIcpBinary(t *testing.T) {
	const subCommand = cmd.CommandIcp
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, "-h").Run(); err != nil {
			t.Error(err)
		}
	})
}
