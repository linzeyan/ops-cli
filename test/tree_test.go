package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestTreeBinary(t *testing.T) {
	const subCommand = cmd.CommandTree
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})
	args := [][]string{
		{subCommand, "-l", "1"},
		{subCommand, "-l", "2"},
		{subCommand, "-m"},
		{subCommand, "-p"},
		{subCommand, "-u"},
		{subCommand, "-g"},
		{subCommand, "--inodes"},
		{subCommand, "--device"},
		{subCommand, "-a"},
		{subCommand, "-d"},
		{subCommand, "-c"},
		{subCommand, "-f"},
		{subCommand, "-s"},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
