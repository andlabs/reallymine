// 21 october 2015
package main

// Byte shuffling.

func Reverse(b []byte) {
	for i := 0; i < len(b) / 2; i++ {
		n := len(b) - i - 1
		b[i], b[n] = b[n], b[i]
	}
}

func SwapLongs(b []byte) {
	for i := 0; i < len(b); i += 4 {
		b[i + 0], b[i + 3] = b[i + 3], b[i + 0]
		b[i + 1], b[i + 2] = b[i + 2], b[i + 1]
	}
}

func SwapHalves(b []byte) {
	n := len(b)
	if n % 2 == 1 {
		panic("[BUG] odd length slice passed to SwapHalves()")
	}
	n /= 2
	c := make([]byte, len(b))
	copy(c, b[n:])
	copy(c[n:], b)
	copy(b, c)
}
