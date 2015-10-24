// 23 october 2015
package main

import (
	"crypto/sha256"
"fmt"
"encoding/hex"
"unicode/utf16"
)

// I don't know when this is used, but have it here it anyway
var DefaultKEK128 = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x2B, 0x99, 0x2D, 0xDF, 0xA2, 0x32, 0x49, 0xD6,
}

var DefaultKEK = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x32, 0x38, 0x46, 0x26, 0x43, 0x38, 0x32, 0x79,
	0xFC, 0xEB, 0xEA, 0x6D, 0x9A, 0xCA, 0x76, 0x86, 0xCD, 0xC7, 0xB9, 0xD9, 0xBC, 0xC7, 0xCD, 0x86,
}

func brokensum(b []byte) []byte {
	h := sha256.New()
	if len(b) > 64 {
		first := b[:64]
		last := b[len(b) - (len(b) % 64):]
		// Oops, this is really what Unlock.exe does.
		// The memcpy() that copies from b into the blocking
		// buffer is only run once, outside the loop.
		// And since the blocking function doesn't overwrite
		// the blocking buffer...
		l := len(b)
		for l > 64 {
			h.Write(first)
			l -= 64
		}
		h.Write(last)
	} else {
		h.Write(b)
	}
	return h.Sum(nil)
}

func KEKFromPassword(password []byte) []byte {
ppw:=utf16.Encode([]rune(string(password)))
password=make([]byte,2*len(ppw))
for j:=0;j<len(ppw);j++{
u16:=ppw[j]
password[2*j]=byte(u16)
password[2*j+1]=byte(u16>>8)
}
	salted := make([]byte, 8 + len(password))
	salted[0] = 'W'
	salted[2] = 'D'
	salted[4] = 'C'
	salted[6] = '.'
	copy(salted[8:], password)
fmt.Println(hex.Dump(salted))

	kek := brokensum(salted)
	i := 999
	for {
		kek = brokensum(kek)
		i--
		if i == 0 {
			break
		}
	}

fmt.Println(hex.Dump(kek))
	return kek
}
