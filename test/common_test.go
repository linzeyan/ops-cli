package test_test

import (
	"testing"

	"github.com/linzeyan/ops-cli/cmd/common"
	"github.com/stretchr/testify/assert"
)

func TestCommonBytes(t *testing.T) {
	testCases := []struct {
		input    float64
		expected string
	}{
		{1, "1.00B"},
		{1024, "1.00KiB"},
		{65530, "63.99KiB"},
		{65530000, "62.49MiB"},
		{89505395658828411324, "77.63EiB"},
		{89505395658828.411324, "81.40TiB"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.expected, func(t *testing.T) {
			got := common.ByteSize(testCase.input)
			assert.Equal(t, testCase.expected, got)
		})
	}
}

func FuzzCommonBytes(t *testing.F) {
	t.Fuzz(func(t *testing.T, input float64) {
		a := assert.New(t)
		v := input / (1 << 10)
		switch {
		case input >= 1<<80:
			a.Contains(common.ByteSize(input), "YiB")
			if v >= 1<<80 {
				a.Contains(common.ByteSize(v), "YiB")
			} else {
				a.Contains(common.ByteSize(v), "ZiB")
			}
		case input >= 1<<70:
			a.Contains(common.ByteSize(input), "ZiB")
			a.Contains(common.ByteSize(v), "EiB")
		case input >= 1<<60:
			a.Contains(common.ByteSize(input), "EiB")
			a.Contains(common.ByteSize(v), "PiB")
		case input >= 1<<50:
			a.Contains(common.ByteSize(input), "PiB")
			a.Contains(common.ByteSize(v), "TiB")
		case input >= 1<<40:
			a.Contains(common.ByteSize(input), "TiB")
			a.Contains(common.ByteSize(v), "GiB")
		case input >= 1<<30:
			a.Contains(common.ByteSize(input), "GiB")
			a.Contains(common.ByteSize(v), "MiB")
		case input >= 1<<20:
			a.Contains(common.ByteSize(input), "MiB")
			a.Contains(common.ByteSize(v), "KiB")
		case input >= 1<<10:
			a.Contains(common.ByteSize(input), "KiB")
			a.Contains(common.ByteSize(v), "B")
		case input < 0:
			a.Equal(common.ByteSize(input), "")
			a.Equal(common.ByteSize(v), "")
		default:
			a.Contains(common.ByteSize(input), "B")
			a.Contains(common.ByteSize(v), "B")
		}
	})
}
