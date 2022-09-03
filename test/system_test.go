package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestBinarySystem(t *testing.T) {
	const subCommand = cmd.CommandSystem
	args := [][]string{
		{subCommand, cmd.CommandCPU},
		{subCommand, cmd.CommandDisk},
		{subCommand, cmd.CommandHost},
		{subCommand, cmd.CommandLoad},
		{subCommand, cmd.CommandMemory},
		{subCommand, cmd.CommandNetwork},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
