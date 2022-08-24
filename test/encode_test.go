package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/linzeyan/ops-cli/cmd/common"
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
	const subCommand = "encode"
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{runCommand, mainGo, subCommand, common.Base32Hex, testHost}, "CTNMUPRCCKN66RRD\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base32Std, testHost}, "M5XW6Z3MMUXGG33N\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base64Std, testHost}, "Z29vZ2xlLmNvbQ==\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base64URL, testHost}, "Z29vZ2xlLmNvbQ==\n"},
		{[]string{runCommand, mainGo, subCommand, common.Hex, testHost}, "676f6f676c652e636f6d\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base32Hex, "-d", "ALQLICHPEH636KHJ"}, "UuY29tL3R3\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base32Std, "-d", "JBGTMTDZHEZWI==="}, "HM6Ly93d\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base64Std, "-d", "aWxsZWdhbA=="}, "illegal\n"},
		{[]string{runCommand, mainGo, subCommand, common.Base64URL, "-d", "aHR0cHM6Ly9naXRodWIuY29t"}, "https://github.com\n"},
		{[]string{runCommand, mainGo, subCommand, common.Hex, "-d", "64617461"}, "data\n"},
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

func TestBinaryEncode(t *testing.T) {
	const subCommand = "encode"
	args := [][]string{
		{subCommand, common.Base32Hex, testHost},
		{subCommand, common.Base32Std, testHost},
		{subCommand, common.Base64Std, testHost},
		{subCommand, common.Base64URL, testHost},
		{subCommand, common.Hex, testHost},
		{subCommand, common.Base32Hex, "C5O70R35", "-d"},
		{subCommand, common.Base32Std, "MFYHA3DF", "-d"},
		{subCommand, common.Base64Std, "YXBwbGU=", "-d"},
		{subCommand, common.Base64URL, "aHR0cHM6Ly93d3cuYXBwbGUuY29tL3R3", "-d"},
		{subCommand, common.Hex, "6170706c65", "-d"},
	}

	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
