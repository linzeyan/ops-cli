package test_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	const subCommand = "encrypt"
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, "aes", "-f", "../.gitignore", "-k", "32449939618748684094059431382108"}, "../.gitignore"},
		{[]string{runCommand, mainGo, subCommand, "aes", "-f", "../.gitignore", "-k", "32449939618748684094059431382108", "-d"}, "../.gitignore"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			err := exec.Command(mainCommand, testCases[i].input...).Run()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.FileExists(t, testCases[i].expected)
		})
	}
}

func TestBinaryEncrypt(t *testing.T) {
	const subCommand = "encrypt"
	args := [][]string{
		{subCommand, "aes", "-f", "../.gitignore", "-k", "32449939618748684094059431382108"},
		{subCommand, "aes", "-f", "../.gitignore", "-k", "32449939618748684094059431382108", "-d"},
	}

	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
