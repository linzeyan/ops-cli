package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestBinaryVersion(t *testing.T) {
	const subCommand = cmd.CommandVersion
	args := []string{"", "-y", "-j"}
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
