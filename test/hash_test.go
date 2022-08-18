package test_test

import (
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestHashFile(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]string
	}{
		{"../LICENSE",
			map[string]string{
				cmd.HashMd5:    "3b83ef96387f14655fc854ddc3c6bd57",
				cmd.HashSha1:   "2b8b815229aa8a61e483fb4ba0588b8b6c491890",
				cmd.HashSha256: "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
				cmd.HashSha512: "98f6b79b778f7b0a15415bd750c3a8a097d650511cb4ec8115188e115c47053fe700f578895c097051c9bc3dfb6197c2b13a15de203273e1a3218884f86e90e8",
			}},
		{"assets/proxy.yaml",
			map[string]string{
				cmd.HashMd5:    "e18b4913dc4b2d61e39c6bd53c9d5f4f",
				cmd.HashSha1:   "1207ddb9d9d65835d43587d246a1a6a0bed13d02",
				cmd.HashSha256: "4cd53bc452432e850e83f9e3db5215d8fa6d9fb10d5de01157abe011ac69d90f",
				cmd.HashSha512: "1a2a649cf401e6c6762698364af83e02324f2c8ee7d3c42605c8ac67dc8230622c38ef10833309033bab4b13077555e16e41fdc777b381deeddced3f26b9d0ec",
			}},
	}
	for _, testCase := range testCases {
		if isWindows() {
			if err := cmd.Dos2Unix(testCase.input); err != nil {
				t.Error(err)
			}
		}
		for k, v := range testCase.expected {
			t.Run(k, func(t *testing.T) {
				got, err := cmd.Hasher.Hash(cmd.HashAlgorithm(k), testCase.input)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, v, got)
			})
		}
	}
}

func TestMd5(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{testHost, "1d5920f4b44b27a802bd77c4f0536f5a"},
		{"ops-cli", "532aeb6eb18c6bfb8750f46a60d28a13"},
		{"https://github.com", "3097fca9b1ec8942c4305e550ef1b50a"},
	}
	for _, testCase := range testCases {
		t.Run(cmd.HashMd5, func(t *testing.T) {
			got, err := cmd.Hasher.Hash(cmd.HashAlgorithm(cmd.HashMd5), testCase.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestSha1(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{testHost, "baea954b95731c68ae6e45bd1e252eb4560cdc45"},
		{"ops-cli", "cc32cdd66c03e92415c44cc701191afcb82bd08d"},
		{"https://github.com", "84b7e44aa54d002eac8d00f5bfa9cc93410f2a48"},
	}
	for _, testCase := range testCases {
		t.Run(cmd.HashSha1, func(t *testing.T) {
			got, err := cmd.Hasher.Hash(cmd.HashAlgorithm(cmd.HashSha1), testCase.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestSha256(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{testHost, "d4c9d9027326271a89ce51fcaf328ed673f17be33469ff979e8ab8dd501e664f"},
		{"ops-cli", "e95c20437deefde8b32499f10bccc4d8f813664de63d43ae24db9e4036584d8c"},
		{"https://github.com", "996e1f714b08e971ec79e3bea686287e66441f043177999a13dbc546d8fe402a"},
	}
	for _, testCase := range testCases {
		t.Run(cmd.HashSha256, func(t *testing.T) {
			got, err := cmd.Hasher.Hash(cmd.HashAlgorithm(cmd.HashSha256), testCase.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestSha512(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{
			testHost,
			"a5b5955a4db31736f9dfd45c89c12331e0370074fc7fec0ac4d189a62391bf7060287f957ce67cf3adcac7a4353a7a8241e33084a9b543cbb3f39770970a41b2",
		},
		{
			"ops-cli",
			"be514862160b204fb36df048c567a5ffcb96c46ea8e64c8375c5598da0dec6a78eb8a17fb1b0252415292e99d0a3e4d017cc8a9fe6ef467e2d92f66636651ab5",
		},
		{
			"https://github.com",
			"44679c4abfe5ecb67d21e6069ade2745f7607b1ae2d8bf8ad245994917331e8539689f43d2aa0f445e4b2f86875742c751570f6550006673a0a2edbbd8877fb9",
		},
	}
	for _, testCase := range testCases {
		t.Run(cmd.HashSha512, func(t *testing.T) {
			got, err := cmd.Hasher.Hash(cmd.HashAlgorithm(cmd.HashSha512), testCase.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCase.expected, got)
		})
	}
}
