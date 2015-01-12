// 12 january 2015
package main

import (
	"fmt"
	"crypto/sha256"
)

func brokensum(b []byte) []byte {
	h := sha256.New()
	// TODO actually implement the brokenness
/*	// TODO see if this is right
	if len(b) > 32 {
		first := b[:32]
		last := b[len(b) - 32:]
		for i := 0; i < len(b) >> 9; i++ {
			h,Write(first)
		}
		h.Write(last)
	} else {
*/
		h.Write(b)
//	}
	return h.Sum(nil)
}

func main() {
	p := make([]byte, 0, 0)
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
