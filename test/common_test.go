package test_test

import (
	"fmt"
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
		{1024, "1.00KB"},
		{65530, "63.99KB"},
		{65530000, "62.49MB"},
		{89505395658828411324, "77.63EB"},
		{89505395658828.411324, "81.40TB"},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%.2f", testCase.input), func(t *testing.T) {
			got := common.ByteSize(testCase.input).String()
			assert.Equal(t, testCase.expected, got)
		})
	}
}
