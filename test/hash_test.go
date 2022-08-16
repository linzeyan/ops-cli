package test_test

import (
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestHashFile(t *testing.T) {
	if isWindows() {
		return
	}
	testCases := struct {
		input    string
		expected map[string]string
	}{
		mainGo,
		map[string]string{
			"md5":    "1b48671ec88f2b498820450f802097a6",
			"sha1":   "b2828f7ce1eb4872e5543617f9bf2b7ef28e7c61",
			"sha256": "74d680e9a561929551611bcecf9fc1704c75a8491e9aec00065adbfaefa36905",
			"sha512": "9aa41fc10c66de39e30dcf7e35be7c15878322df8e318922dc8623f140464c51b964953d950d8dc1ef4b0d4b614f01c065ef50d0c0e686cdd21ff4580fec3ce7",
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
