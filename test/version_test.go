package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestVersionBinary(t *testing.T) {
	const subCommand = cmd.CommandVersion
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})
	args := []string{"yaml", "json"}
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, "--output", args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
