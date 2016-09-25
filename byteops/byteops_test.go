// 25 september 2016
package byteops

import (
	"bytes"
	"testing"
)

func runOne(t *testing.T, input []byte, expected []byte, name string, op func([]byte)) {
	if expected == nil {
		return
	}
	b := DupBytes(input)
	op(b)
	if !bytes.Equal(b, expected) {
		t.Errorf("%X: %s failed: expected %X, got %X\n", input, name, expected, b)
	}
}

type testset []struct {
	Input		[]byte
	Reverse		[]byte
	SwapWords	[]byte
	SwapLongs	[]byte
	SwapHalves	[]byte
}

var tests = testset{
	{
		Input:		[]byte{},
		Reverse:		[]byte{},
		SwapWords:	[]byte{},
		SwapLongs:	[]byte{},
		SwapHalves:	[]byte{},
	},
	{
		Input:		[]byte{ 0xEE, 0xFF },
		Reverse:		[]byte{ 0xFF, 0xEE },
		SwapWords:	[]byte{ 0xFF, 0xEE },
		SwapLongs:	nil,		// not an even number of longs
		SwapHalves:	[]byte{ 0xFF, 0xEE },
	},
	{
		Input:		[]byte{ 0xCC, 0xDD, 0xEE, 0xFF },
		Reverse:		[]byte{ 0xFF, 0xEE, 0xDD, 0xCC },
		SwapWords:	[]byte{ 0xDD, 0xCC, 0xFF, 0xEE },
		SwapLongs:	[]byte{ 0xFF, 0xEE, 0xDD, 0xCC },
		SwapHalves:	[]byte{ 0xEE, 0xFF, 0xCC, 0xDD },
	},
	{
		Input:		[]byte{ 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF },
		Reverse:		[]byte{ 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA },
		SwapWords:	[]byte{ 0xBB, 0xAA, 0xDD, 0xCC, 0xFF, 0xEE },
		SwapLongs:	nil,		// not an even number of longs
		SwapHalves:	[]byte{ 0xDD, 0xEE, 0xFF, 0xAA, 0xBB, 0xCC },
	},
	{
		Input:		[]byte{ 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF },
		Reverse:		[]byte{ 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88 },
		SwapWords:	[]byte{ 0x99, 0x88, 0xBB, 0xAA, 0xDD, 0xCC, 0xFF, 0xEE },
		SwapLongs:	[]byte{ 0xBB, 0xAA, 0x99, 0x88, 0xFF, 0xEE, 0xDD, 0xCC },
		SwapHalves:	[]byte{ 0xCC, 0xDD, 0xEE, 0xFF, 0x88, 0x99, 0xAA, 0xBB },
	},
	// TODO more tests
}

func TestPackage(t *testing.T) {
	for _, s := range tests {
		runOne(t, s.Input, s.Reverse, "Reverse()", Reverse)
		runOne(t, s.Input, s.SwapWords, "SwapWords()", SwapWords)
		runOne(t, s.Input, s.SwapLongs, "SwapLongs()", SwapLongs)
		runOne(t, s.Input, s.SwapHalves, "SwapHalves()", SwapHalves)
	}
}
