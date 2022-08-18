package test_test

import (
	"os/exec"
	"testing"
)

func TestBinaryOtp(t *testing.T) {
	const subCommand = "otp"
	calculateArgs := [][]string{
		{subCommand, "calculate", "6BDR", "T7AT", "RRCZ", "V5IS", " FLOH", "AHQL", "YF4Z", " ORG7", "-p", "15", "-d", "7"},
		{subCommand, "calculate", "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7", "-p", "30", "-d", "8"},
		{subCommand, "calculate", "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7", "-d", "8"},
		{subCommand, "calculate", "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7", "-p", "60"},
		{subCommand, "calculate", "6BDRT7ATRRCZV5ISFLOHAHQLYF4ZORG7"},
	}
	generateArgs := [][]string{
		{subCommand, "generate", "-p", "15"},
		{subCommand, "generate", "-a", "sha256"},
		{subCommand, "generate", "-a", "sha512", "-p", "60"},
	}

	for i := range calculateArgs {
		t.Run(calculateArgs[i][2], func(t *testing.T) {
			if err := exec.Command(binaryCommand, calculateArgs[i]...).Run(); err != nil {
				t.Error(calculateArgs[i], err)
			}
		})
	}
	for i := range generateArgs {
		t.Run(generateArgs[i][2], func(t *testing.T) {
			if err := exec.Command(binaryCommand, generateArgs[i]...).Run(); err != nil {
				t.Error(generateArgs[i], err)
			}
		})
	}
}
