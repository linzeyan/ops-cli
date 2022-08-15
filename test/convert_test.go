package test_test

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	const subCommand = "convert"
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, "yaml2json", "-i", "assets/proxy.yaml", "-o", "/tmp/testy.json"}, "assets/proxy.json"},
		{[]string{runCommand, mainGo, subCommand, "yaml2toml", "-i", "assets/proxy.yaml", "-o", "/tmp/testy.toml"}, "assets/proxy.toml"},
		{[]string{runCommand, mainGo, subCommand, "toml2json", "-i", "assets/proxy.toml", "-o", "/tmp/testt.json"}, "assets/proxy.json"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			err := exec.Command(mainCommand, testCases[i].input...).Run()
			if err != nil {
				t.Error(err)
			}
			assert.FileExists(t, testCases[i].input[7])
			expected, err := getHash(testCases[i].expected)
			if err != nil {
				t.Error(err)
			}
			got, err := getHash(testCases[i].input[7])
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, expected, got)
			_ = exec.Command("rm", "-f", testCases[i].input[7]).Run()
		})
	}
}

func TestBinaryConvert(t *testing.T) {
	const subCommand = "convert"
	testCommand := []string{"json2csv", "json2toml", "json2yaml", "yaml2json", "yaml2toml", "yaml2csv"}

	for _, cmd := range testCommand {
		slice := strings.Split(cmd, "2")
		input := fmt.Sprintf("assets/proxy.%s", slice[0])
		output := fmt.Sprintf("/tmp/out.%s", slice[1])
		t.Run(cmd, func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, cmd, "-i", input, "-o", output).Run(); err != nil {
				t.Error(err)
			}
			_ = exec.Command("rm", "-f", output).Run()
		})
	}
}

func getHash(filename string) (uint32, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	h := crc32.NewIEEE()
	_, err = io.Copy(h, f)
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
