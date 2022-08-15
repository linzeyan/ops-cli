package test_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	const subCommand = "random"
	testCases := []struct {
		input          []string
		unexpectedText []cmd.RandomCharacter
		expectedLength int
	}{
		{
			[]string{runCommand, mainGo, subCommand, "number"},
			[]cmd.RandomCharacter{cmd.LowercaseLetters, cmd.UppercaseLetters, cmd.Symbols},
			24,
		},
		{
			[]string{runCommand, mainGo, subCommand, "number", "-l", "64"},
			[]cmd.RandomCharacter{cmd.LowercaseLetters, cmd.UppercaseLetters, cmd.Symbols},
			64,
		},
		{
			[]string{runCommand, mainGo, subCommand, "symbol"},
			[]cmd.RandomCharacter{cmd.LowercaseLetters, cmd.UppercaseLetters, cmd.Numbers},
			24,
		},
		{
			[]string{runCommand, mainGo, subCommand, "lowercase"},
			[]cmd.RandomCharacter{cmd.Numbers, cmd.UppercaseLetters, cmd.Symbols},
			24,
		},
		{
			[]string{runCommand, mainGo, subCommand, "uppercase", "-l", "200"},
			[]cmd.RandomCharacter{cmd.Numbers, cmd.LowercaseLetters, cmd.Symbols},
			200,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			out, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(err)
			}
			got := string(out)
			assert.Len(t, strings.TrimRight(got, "\n"), testCases[i].expectedLength)
			for _, textType := range testCases[i].unexpectedText {
				if strings.ContainsAny(got, string(textType)) {
					t.Errorf("Expect %s, got %v", testCases[i].input[3], textType)
				}
			}
		})
	}
}

func TestBinaryRandom(t *testing.T) {
	const subCommand = "random"
	args := [][]string{
		{subCommand, "lowercase", "-l", "30"},
		{subCommand, "uppercase", "-l", "40"},
		{subCommand, "number", "-l", "50"},
		{subCommand, "symbol", "-l", "60"},
		{subCommand, "-l", "70"},
		{subCommand, "-s", "10"},
		{subCommand, "-o", "10", "-s", "10", "-l", "32"},
		{subCommand, "-u", "10", "-s", "10", "-o", "10", "-l", "64"},
		{subCommand, "-n", "15", "-s", "10", "-o", "10", "-l", "64", "-u", "10"},
	}
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand).Run(); err != nil {
			t.Error(err)
		}
	})
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(err)
			}
		})
	}
}
