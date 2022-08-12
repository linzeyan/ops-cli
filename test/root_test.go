package test_test

import (
	"fmt"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
)

const (
	mainCommand = "go"
	mainGo      = "../main.go"
	runCommand  = "run"
)

func TestDomain(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected bool
	}{
		{"1.1.1.1", false},
		{"example.com", true},
		{"Hello world", false},
		{11111, false},
		{"dns-admin.google.com", true},
		{"dns-admin.google.com.", true},
	}
	for i := range testCases {
		t.Run(fmt.Sprintf("%s", testCases[i].input), func(t *testing.T) {
			if got := cmd.ValidDomain(testCases[i].input); got != testCases[i].expected {
				t.Errorf("Expected %t, got %t", testCases[i].expected, got)
			}
		})
	}
}

func TestFile(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", false},
		{"../vendor", true},
		{"command.json", false},
		{"root_test.go", true},
		{"/dev/null", true},
	}
	for i := range testCases {
		t.Run(testCases[i].input, func(t *testing.T) {
			if got := cmd.ValidFile(testCases[i].input); got != testCases[i].expected {
				t.Errorf("Expected %t, got %t", testCases[i].expected, got)
			}
		})
	}
}

func TestIP(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", true},
		{"999.1.1.1", false},
		{"260.2.3.4", false},
		{"example.com", false},
		{"2404:6800:4008:c01::65", true},
	}
	for i := range testCases {
		t.Run(testCases[i].input, func(t *testing.T) {
			if got := cmd.ValidIP(testCases[i].input); got != testCases[i].expected {
				t.Errorf("Expected %t, got %t", testCases[i].expected, got)
			}
		})
	}
}

func TestIPv4(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", true},
		{"999.1.1.1", false},
		{"260.2.3.4", false},
		{"example.com", false},
		{"2404:6800:4008:c01::65", false},
	}
	for i := range testCases {
		t.Run(testCases[i].input, func(t *testing.T) {
			if got := cmd.ValidIPv4(testCases[i].input); got != testCases[i].expected {
				t.Errorf("Expected %t, got %t", testCases[i].expected, got)
			}
		})
	}
}

func TestIPv6(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", false},
		{"999.1.1.1", false},
		{"260.2.3.4", false},
		{"example.com", false},
		{"2404:6800:4008:c01::65", true},
	}
	for i := range testCases {
		t.Run(testCases[i].input, func(t *testing.T) {
			if got := cmd.ValidIPv6(testCases[i].input); got != testCases[i].expected {
				t.Errorf("Expected %t, got %t", testCases[i].expected, got)
			}
		})
	}
}

func TestUrl(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.1.1.1", false},
		{"999.1.1.1", false},
		{"2404:6800:4008:c01::65", false},
		{"https://1.1.1.1", true},
		{"example.com", false},
		{"https://example.com", true},
		{"https://example.com/?", true},
		{"https://example.com/api/v1/add?user=1", true},
	}
	for i := range testCases {
		t.Run(testCases[i].input, func(t *testing.T) {
			if got := cmd.ValidURL(testCases[i].input); got != testCases[i].expected {
				t.Errorf("Expected %t, got %t", testCases[i].expected, got)
			}
		})
	}
}