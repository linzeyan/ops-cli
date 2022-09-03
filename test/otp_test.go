package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

func TestOtpBinary(t *testing.T) {
	const subCommand = cmd.CommandOtp
	calculateArgs := [][]string{
		{subCommand, cmd.CommandCalculate, "6BDR", "T7AT", "RRCZ", "V5IS", " FLOH", "AHQL", "YF4Z", " ORG7", "-p", "15", "-d", "7"},
		{subCommand, cmd.CommandCalculate, "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7", "-p", "30", "-d", "8"},
		{subCommand, cmd.CommandCalculate, "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7", "-d", "8"},
		{subCommand, cmd.CommandCalculate, "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7", "-p", "60"},
		{subCommand, cmd.CommandCalculate, "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7"},
	}
	for i := range calculateArgs {
		t.Run(calculateArgs[i][2], func(t *testing.T) {
			if err := exec.Command(binaryCommand, calculateArgs[i]...).Run(); err != nil {
				t.Error(calculateArgs[i], err)
			}
		})
	}

	generateArgs := [][]string{
		{subCommand, cmd.CommandGenerate, "-p", "15"},
		{subCommand, cmd.CommandGenerate, "-a", "sha256"},
		{subCommand, cmd.CommandGenerate, "-a", "sha512", "-p", "60"},
	}
	for i := range generateArgs {
		t.Run(generateArgs[i][2], func(t *testing.T) {
			if err := exec.Command(binaryCommand, generateArgs[i]...).Run(); err != nil {
				t.Error(generateArgs[i], err)
			}
		})
	}
}
