// 12 january 2015
package main

import (
	"fmt"
	"crypto/sha256"
)

func brokensum(b []byte) []byte {
	h := sha256.New()
	if len(b) > 64 {
		first := b[:64]
		last := b[len(b) - (len(b) % 64):]
		// oops, this is really what Unlock.exe does
		// the memcpy() that copies from b into the blocking buffer is only run once, outside the loop
		// and since the blocking function doesn't overwrite the blocking buffer...
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

func main() {
	p := make([]byte, 64+64+32, 64+64+32)
	p[64+12] = 4
	s := brokensum(p)
	i := 999
	for {
		s = brokensum(s)
		i--
		if i == 0 {
			break
		}
	}
	fmt.Printf("%x\n", s)
}
