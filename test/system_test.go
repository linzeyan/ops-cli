package test_test

import (
	"os/exec"
	"testing"
)

func TestBinarySystem(t *testing.T) {
	const subCommand = "system"
	args := [][]string{
		{subCommand, "cpu"},
		{subCommand, "disk"},
		{subCommand, "host"},
		{subCommand, "load"},
		{subCommand, "memory"},
		{subCommand, "network"},
		{subCommand, "network", "-a"},
		{subCommand, "network", "-i"},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
