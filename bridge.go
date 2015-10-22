// 21 october 2015
package main

import (
	"crypto/aes"
)

type Bridge interface {
	Name() string
	Is(keySector []byte) bool
	NeedsKEK() bool
	CreateDecrypter(keySector []byte, kek []byte) (cipher *aes.Cipher)
	Decrypt(c *aes.Cipher, b []byte)
}

var Bridges []Bridge
