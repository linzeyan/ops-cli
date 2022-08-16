package test_test

import (
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestHashFile(t *testing.T) {
	if isWindows() {
		if err := dos2unix(license); err != nil {
			t.Error(err)
		}
	}
	testCases := struct {
		input    string
		expected map[string]string
	}{
		license,
		map[string]string{
			"md5":    "3b83ef96387f14655fc854ddc3c6bd57",
			"sha1":   "2b8b815229aa8a61e483fb4ba0588b8b6c491890",
			"sha256": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
			"sha512": "98f6b79b778f7b0a15415bd750c3a8a097d650511cb4ec8115188e115c47053fe700f578895c097051c9bc3dfb6197c2b13a15de203273e1a3218884f86e90e8",
		},
	}
	t.Run("md5", func(t *testing.T) {
		got, err := cmd.Hasher.Md5Hash(testCases.input)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, testCases.expected["md5"], got)
	})
	t.Run("sha1", func(t *testing.T) {
		got, err := cmd.Hasher.Sha1Hash(testCases.input)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, testCases.expected["sha1"], got)
	})
	t.Run("sha256", func(t *testing.T) {
		got, err := cmd.Hasher.Sha256Hash(testCases.input)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, testCases.expected["sha256"], got)
	})
	t.Run("sha512", func(t *testing.T) {
		got, err := cmd.Hasher.Sha512Hash(testCases.input)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, testCases.expected["sha512"], got)
	})
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
		t.Run("md5", func(t *testing.T) {
			got, err := cmd.Hasher.Md5Hash(testCase.input)
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
		t.Run("sha1", func(t *testing.T) {
			got, err := cmd.Hasher.Sha1Hash(testCase.input)
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
		t.Run("sha256", func(t *testing.T) {
			got, err := cmd.Hasher.Sha256Hash(testCase.input)
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
		t.Run("sha512", func(t *testing.T) {
			got, err := cmd.Hasher.Sha512Hash(testCase.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCase.expected, got)
		})
	}
}
