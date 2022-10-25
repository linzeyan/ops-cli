package test_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestWhois(t *testing.T) {
	const subCommand = cmd.CommandWhois
	testCases := []struct {
		input    []string
		expected string
	}{
		{
			[]string{runCommand, mainGo, subCommand, testHost, "--ns"},
			`[
	"NS1.GOOGLE.COM",
	"NS2.GOOGLE.COM",
	"NS3.GOOGLE.COM",
	"NS4.GOOGLE.COM"
]`,
		},
		{
			[]string{runCommand, mainGo, subCommand, "apple.com", "--ns"},
			`[
	"A.NS.APPLE.COM",
	"B.NS.APPLE.COM",
	"C.NS.APPLE.COM",
	"D.NS.APPLE.COM"
]`,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.JSONEq(t, testCases[i].expected, string(got))
		})
		time.Sleep(time.Second * 2)
	}
}

func TestWhoisRegistrar(t *testing.T) {
	const subCommand = cmd.CommandWhois
	testCases := []struct {
		input    []string
		expected string
	}{
		{
			[]string{runCommand, mainGo, subCommand, "godaddy.com", "--registrar"},
			"GoDaddy.com\n",
		},
		{
			[]string{runCommand, mainGo, subCommand, "cloudflare.com", "--registrar"},
			"CloudFlare\n",
		},
	}

	for i := range testCases {
		t.Run(testCases[i].input[3], func(t *testing.T) {
			got, err := exec.Command(mainCommand, testCases[i].input...).Output()
			if err != nil {
				t.Error(testCases[i].input, err)
			}
			assert.Equal(t, testCases[i].expected, string(got))
		})
		time.Sleep(time.Second * 2)
	}
}

func TestWhoisBinary(t *testing.T) {
	const subCommand = cmd.CommandWhois
	args := []string{"-d", "-e", "-n", "-r"}

	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, testHost, args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
		time.Sleep(time.Second * 2)
	}
}
