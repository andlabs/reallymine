// 23 october 2015
package kek

import (
	"crypto/sha256"
	"unicode/utf16"
)

// I don't know when this is used, but have it here it anyway
var Default128 = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x2B, 0x99, 0x2D, 0xDF, 0xA2, 0x32, 0x49, 0xD6,
}

var Default = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x32, 0x38, 0x46, 0x26, 0x43, 0x38, 0x32, 0x79,
	0xFC, 0xEB, 0xEA, 0x6D, 0x9A, 0xCA, 0x76, 0x86, 0xCD, 0xC7, 0xB9, 0xD9, 0xBC, 0xC7, 0xCD, 0x86,
}

func brokensum(b []byte) []byte {
	h := sha256.New()
	if len(b) > 64 {
		first := b[:64]
		last := b[len(b)-(len(b)%64):]
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

// Yes, that's right folks, Unlock.exe gives us this:
// 1) Constant salt.
// 2) The salt is a string.
// 3) The combined string must be UTF-16 encoded.
// Since this is Windows software we're talking about, that's
// UTF-16 little-endian.
// TODO this function should be improved somehow...
func saltAndUTF16(password string) []byte {
	sp := "WDC." + password
	u16 := utf16.Encode([]rune(sp))
	p := make([]byte, 2*len(u16))
	for i := 0; i < len(u16); i++ {
		u := u16[i]
		p[2*i] = byte(u & 0xFF)
		p[2*i+1] = byte((u >> 8) & 0xFF)
	}
	return p
}

// Don't mind the weird loop here; we're imitating Unlock.exe.
func FromPassword(password string) []byte {
	kek := brokensum(saltAndUTF16(password))
	i := 999
	for {
		kek = brokensum(kek)
		i--
		if i == 0 {
			break
		}
	}
	return kek
}
