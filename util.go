// 21 october 2015
package main

import (
	"crypto/aes"
)

// Byte shuffling.

func Reverse(b []byte) {
	if len(b)%2 == 1 {
		BUG("Reverse() called with odd-sized slice")
	}
	for i := 0; i < len(b)/2; i++ {
		n := len(b) - i - 1
		b[i], b[n] = b[n], b[i]
	}
}

func SwapLongs(b []byte) {
	if len(b)%4 != 0 {
		BUG("SwapLongs() called with len(b) not a multiple of 4")
	}
	for i := 0; i < len(b); i += 4 {
		b[i+0], b[i+3] = b[i+3], b[i+0]
		b[i+1], b[i+2] = b[i+2], b[i+1]
	}
}

func SwapHalves(b []byte) {
	n := len(b)
	if n%2 == 1 {
		BUG("SwapHalves() called with odd-sized slice")
	}
	n /= 2
	c := make([]byte, len(b))
	copy(c, b[n:])
	copy(c[n:], b)
	copy(b, c)
}

// Other functions.

func DupBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func NewAES(key []byte) *aes.Cipher {
	c, err := aes.NewCipher(key)
	if err != nil {
		BUG("error creating AES cipher in NewAES(): %v", err)
	}
	return c
}
