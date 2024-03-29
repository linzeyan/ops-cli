package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestBase32Hex(t *testing.T) {
	expected := testHost
	encode, err := cmd.Encoder.Base32HexEncode(expected)
	if err != nil {
		t.Error(err)
	}
	got, err := cmd.Encoder.Base32HexDecode(encode)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, string(got))
}

func TestBase32Std(t *testing.T) {
	expected := runCommand
	encode, err := cmd.Encoder.Base32StdEncode(expected)
	if err != nil {
		t.Error(err)
	}
	got, err := cmd.Encoder.Base32StdDecode(encode)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, string(got))
}

func TestBase64Std(t *testing.T) {
	expected := mainGo
	encode, err := cmd.Encoder.Base64StdEncode(expected)
	if err != nil {
		t.Error(err)
	}
	got, err := cmd.Encoder.Base64StdDecode(encode)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, string(got))
}

func TestBase64URL(t *testing.T) {
	expected := binaryCommand
	encode, err := cmd.Encoder.Base64URLEncode(expected)
	if err != nil {
		t.Error(err)
	}
	got, err := cmd.Encoder.Base64URLDecode(encode)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, string(got))
}

func TestHex(t *testing.T) {
	expected := mainCommand
	encode, err := cmd.Encoder.HexEncode(expected)
	if err != nil {
		t.Error(err)
	}
	got, err := cmd.Encoder.HexDecode(encode)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, string(got))
}

func TestPem(t *testing.T) {
	expected := mainCommand
	encode, err := cmd.Encoder.PemEncode(expected, "OPS KEY")
	if err != nil {
		t.Error(err)
	}
	got, err := cmd.Encoder.PemDecode([]byte(encode))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, string(got))
}

func TestEncode(t *testing.T) {
	const subCommand = cmd.CommandEncode
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase32Hex, testHost}, "CTNMUPRCCKN66RRD"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase32Std, testHost}, "M5XW6Z3MMUXGG33N"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase64Std, testHost}, "Z29vZ2xlLmNvbQ=="},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase64URL, testHost}, "Z29vZ2xlLmNvbQ=="},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandHex, testHost}, "676f6f676c652e636f6d"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase32Hex, "-d", "ALQLICHPEH636KHJ"}, "UuY29tL3R3"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase32Std, "-d", "JBGTMTDZHEZWI==="}, "HM6Ly93d"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase64Std, "-d", "aWxsZWdhbA=="}, "illegal"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandBase64URL, "-d", "aHR0cHM6Ly9naXRodWIuY29t"}, "https://github.com"},
		{[]string{runCommand, mainGo, subCommand, cmd.CommandHex, "-d", "64617461"}, "data"},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
	}
}

func TestEncodeBinary(t *testing.T) {
	const subCommand = cmd.CommandEncode
	args := [][]string{
		{subCommand, cmd.CommandBase32Hex, testHost},
		{subCommand, cmd.CommandBase32Std, testHost},
		{subCommand, cmd.CommandBase64Std, testHost},
		{subCommand, cmd.CommandBase64URL, testHost},
		{subCommand, cmd.CommandHex, testHost},
		{subCommand, cmd.CommandBase32Hex, "C5O70R35", "-d"},
		{subCommand, cmd.CommandBase32Std, "MFYHA3DF", "-d"},
		{subCommand, cmd.CommandBase64Std, "YXBwbGU=", "-d"},
		{subCommand, cmd.CommandBase64URL, "aHR0cHM6Ly93d3cuYXBwbGUuY29tL3R3", "-d"},
		{subCommand, cmd.CommandHex, "6170706c65", "-d"},
	}

	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
