package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/linzeyan/ops-cli/cmd/common"
)

func TestDfBinary(t *testing.T) {
	if common.IsDarwin() {
		return
	}
	const subCommand = cmd.CommandDf
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})
}
