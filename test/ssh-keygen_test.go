package test_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestBinarySSHKeygen(t *testing.T) {
	const subCommand = cmd.CommandSSH
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})
	files := []string{"test", "test2", "id_rsa", "id_rsa_new"}
	args := [][]string{
		{subCommand, "-b", "4096"},
		{subCommand, "-b", "4096", "-f", files[0]},
		{subCommand, "-f", files[1]},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
	for _, v := range files {
		_ = os.Remove(v)
		_ = os.Remove(v + ".pub")
	}
}
