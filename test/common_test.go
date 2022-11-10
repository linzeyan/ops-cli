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
			a.Contains(common.ByteSize(input), common.YiB.String())
			if v >= 1<<80 {
				a.Contains(common.ByteSize(v), common.YiB.String())
			} else {
				a.Contains(common.ByteSize(v), common.ZiB.String())
			}
		case input >= 1<<70:
			a.Contains(common.ByteSize(input), common.ZiB.String())
			a.Contains(common.ByteSize(v), common.EiB.String())
		case input >= 1<<60:
			a.Contains(common.ByteSize(input), common.EiB.String())
			a.Contains(common.ByteSize(v), common.PiB.String())
		case input >= 1<<50:
			a.Contains(common.ByteSize(input), common.PiB.String())
			a.Contains(common.ByteSize(v), common.TiB.String())
		case input >= 1<<40:
			a.Contains(common.ByteSize(input), common.TiB.String())
			a.Contains(common.ByteSize(v), common.GiB.String())
		case input >= 1<<30:
			a.Contains(common.ByteSize(input), common.GiB.String())
			a.Contains(common.ByteSize(v), common.MiB.String())
		case input >= 1<<20:
			a.Contains(common.ByteSize(input), common.MiB.String())
			a.Contains(common.ByteSize(v), common.KiB.String())
		case input >= 1<<10:
			a.Contains(common.ByteSize(input), common.KiB.String())
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
