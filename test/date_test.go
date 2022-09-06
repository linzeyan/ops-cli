package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestDateBinary(t *testing.T) {
	const subCommand = cmd.CommandDate
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})

	args := [][]string{
		{subCommand, "-u"},
		{subCommand, "-D"},
		{subCommand, "-T"},
		{subCommand, "-D", "-u"},
		{subCommand, "-T", "-u"},
		{subCommand, "-s"},
		{subCommand, "milli", "-s"},
		{subCommand, "micro", "-s"},
		{subCommand, "nano", "-s"},
		{subCommand, "2006", "-f"},
		{subCommand, "01", "-f"},
		{subCommand, "02", "-f"},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
