package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestTcpingBinary(t *testing.T) {
	const subCommand = cmd.CommandTcping
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, "--help").Run(); err != nil {
			t.Error(err)
		}
	})
	t.Run(testHost, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, testHost, "443").Run(); err != nil {
			t.Error(err)
		}
	})
	t.Run(testHost, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand, testHost, "443", "-p", "udp").Run(); err != nil {
			t.Error(err)
		}
	})
}
