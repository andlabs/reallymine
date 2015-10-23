// 21 october 2015
package main

import (
	"crypto/cipher"
)

type Bridge interface {
	Name() string
	Is(keySector []byte) bool
	NeedsKEK() bool
	CreateDecrypter(keySector []byte, kek []byte) (c cipher.Block)
	Decrypt(c cipher.Block, b []byte)
}

var Bridges []Bridge

func IdentifyKeySector(possibleKeySector []byte) Bridge {
	for _, b := range Bridges {
		if b.Is(possibleKeySector) {
			return b
		}
	}
	return nil // not a (known) key sector
}
