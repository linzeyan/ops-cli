package test_test

import (
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNetmask(t *testing.T) {
	const subCommand = cmd.CommandNetmask

	testCases := []struct {
		input    []string
		expected string
	}{
		{
			[]string{subCommand, "-b", "3.4.5.6/7"},
			"00000010 00000000 00000000 00000000 / 11111110 00000000 00000000 00000000\n",
		},
		{[]string{subCommand, "-o", "3.4.5.6/7"}, "2 0 0 0 / 376 0 0 0\n"},
		{[]string{subCommand, "-d", "3.4.5.6/7"}, "2.0.0.0 / 254.0.0.0\n"},
		{[]string{subCommand, "-x", "3.4.5.6/7"}, "2 0 0 0 / fe 0 0 0\n"},
		{[]string{subCommand, "-i", "3.4.5.6/7"}, "2.0.0.0 / 1.255.255.255\n"},
		{[]string{subCommand, "-r", "3.4.5.6/7"}, "2.0.0.0 -> 3.255.255.255 (33554432)\n"},
		{
			[]string{subCommand, "-r", "fe80::aede:48ff:fe00:1122/64"},
			"fe80:: -> fe80::ffff:ffff:ffff:ffff (18446744073709551616)\n",
		},
		{[]string{subCommand, "-c", "3.4.5.0-3.4.5.1"}, "3.4.5.0/31\n"},
		{
			[]string{subCommand, "-c", "fe80::aede:48ff:fe00:1122-fe80::aede:48ff:fe00:112a"},
			"fe80::aede:48ff:fe00:1122/127\nfe80::aede:48ff:fe00:1124/126\nfe80::aede:48ff:fe00:1128/127\nfe80::aede:48ff:fe00:112a/128\n",
		},
		{[]string{subCommand, "-i", "3.4.5.0-3.4.5.1"}, "3.4.5.0 / 0.0.0.1\n"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.input[1], func(t *testing.T) {
			got, err := exec.Command(binaryCommand, testCase.input...).Output()
			if err != nil {
				t.Error(testCase.input[1], err)
			}
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestNetmaskBinary(t *testing.T) {
	const subCommand = cmd.CommandNetmask
	t.Run(subCommand, func(t *testing.T) {
		if err := exec.Command(binaryCommand, subCommand).Run(); err != nil {
			t.Error(err)
		}
	})

	args := [][]string{
		{subCommand, "-b", "3.4.5.6/7"},
		{subCommand, "-o", "3.4.5.6/7"},
		{subCommand, "-d", "3.4.5.6/7"},
		{subCommand, "-x", "3.4.5.6/7"},
		{subCommand, "-i", "3.4.5.6/7"},
		{subCommand, "-r", "3.4.5.6/7"},
		{subCommand, "-c", "1.1.1.1-2.2.2.2"},
		{subCommand, "-i", "1.1.1.1-2.2.2.2"},
	}
	for i := range args {
		t.Run(args[i][1], func(t *testing.T) {
			if err := exec.Command(binaryCommand, args[i]...).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}
