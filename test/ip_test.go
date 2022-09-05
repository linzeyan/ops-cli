package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestIPBinary(t *testing.T) {
	const subCommand = cmd.CommandIP
	args := []string{"all"}
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
