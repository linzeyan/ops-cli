package test_test

import (
	"os/exec"
	"testing"
)

func TestBinaryVersion(t *testing.T) {
	const subCommand = "version"
	args := []string{"-c", "-y", "-j"}
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i]).Run(); err != nil {
				t.Error(err)
			}
		})
	}
}
