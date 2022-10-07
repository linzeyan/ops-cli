package test_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/linzeyan/ops-cli/cmd/validator"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	const subCommand = cmd.CommandConvert
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, cmd.CommandYaml2JSON, "-i", "assets/proxy.yaml", "-o", "testy.json"}, "assets/proxy.json"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandYaml2Toml, "-i", "assets/proxy.yaml", "-o", "testy.toml"}, "assets/proxy.toml"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandToml2JSON, "-i", "assets/proxy.toml", "-o", "testt.json"}, "assets/proxy.json"},
	}
	if validator.IsWindows() {
		if err := common.Dos2Unix("assets/proxy.json"); err != nil {
			t.Error(err)
		}
		if err := common.Dos2Unix("assets/proxy.toml"); err != nil {
			t.Error(err)
		}
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			err := exec.Command(mainCommand, testCases[i].input...).Run()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.FileExists(t, testCases[i].input[7])

			expected, err := os.ReadFile(testCases[i].expected)
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			got, err := os.ReadFile(testCases[i].input[7])
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.Equal(t, expected, got)
			_ = os.Remove(testCases[i].input[7])
		})
	}
}

func TestConvertBinary(t *testing.T) {
	const subCommand = cmd.CommandConvert
	testCommand := []string{
		cmd.CommandJSON2Csv,
		cmd.CommandJSON2Toml,
		cmd.CommandJSON2Yaml,
		cmd.CommandYaml2JSON,
		cmd.CommandYaml2Toml,
		cmd.CommandYaml2Csv,
	}

	for _, cmd := range testCommand {
		slice := strings.Split(cmd, "2")
		input := fmt.Sprintf("assets/proxy.%s", slice[0])
		output := fmt.Sprintf("out.%s", slice[1])
		t.Run(cmd, func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, cmd, "-i", input, "-o", output).Run(); err != nil {
				t.Error(cmd, err)
			}
			_ = os.Remove(output)
		})
	}
}
