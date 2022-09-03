package test_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestEncryptAes(t *testing.T) {
	input := mainGo
	expected, err := os.ReadFile(input)
	if err != nil {
		t.Error(err)
	}
	if err := cmd.Encryptor.EncryptFile("84815131446564008011748691915873", input); err != nil {
		t.Error(err)
	}
	if err := cmd.Encryptor.DecryptFile("84815131446564008011748691915873", input); err != nil {
		t.Error(err)
	}
	got, err := os.ReadFile(input)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, got)
	_ = os.Remove(mainGo + ".temp")
}

func TestEncrypt(t *testing.T) {
	const subCommand = cmd.CommandEncrypt
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, cmd.CommandFile, "../.gitignore", "-k", "32449939618748684094059431382108"}, "../.gitignore"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandFile, "../.gitignore", "-k", "32449939618748684094059431382108", "-d"}, "../.gitignore"},
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
	const subCommand = cmd.CommandEncrypt
	args := [][]string{
		{subCommand, cmd.CommandFile, "../.gitignore", "-k", "32449939618748684094059431382108"},
		{subCommand, cmd.CommandFile, "../.gitignore", "-k", "32449939618748684094059431382108", "-d"},
	}

	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
