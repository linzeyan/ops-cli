package test_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/linzeyan/ops-cli/cmd"
	"github.com/stretchr/testify/require"
)

func TestIPBinary(t *testing.T) {
	const subCommand = cmd.CommandIP
	args := []string{"all"}
	for i := range args {
		t.Run(args[i], func(t *testing.T) {
			if err := exec.Command(binaryCommand, subCommand, args[i]).Run(); err != nil {
				t.Error(args[i], err)
			}
		})
	}
}

func Test_IP_Convert(t *testing.T) {
	inputs := []string{
		"16843009",
		"1.65793",
		"1.1.257",
		"1.1.1.1",
	}

	for _, s := range inputs {
		ip, err := cmd.ParseAnyIPv4Netip(s)
		require.NoError(t, err)
		val, f2, f3, f4 := cmd.IPv4NetipToAllForms(ip)
		fmt.Printf("[%s] => 整數: %s, a.b: %s, a.b.c: %s, a.b.c.d: %s\n", s, val, f2, f3, f4)
		require.Equal(t, inputs[0], val)
		require.Equal(t, inputs[1], f2)
		require.Equal(t, inputs[2], f3)
		require.Equal(t, inputs[3], f4)
	}
}
