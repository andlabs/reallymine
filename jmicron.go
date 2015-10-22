// 21 october 2015
package main

import (
	"crypto/aes"
)

type JMicron struct{}

func (JMicron) Name() string {
	return "JMicron"
}

func (JMicron) Is(keySector []byte) bool {
	return keySector[0] == 'W' &&
		keySector[1] == 'D' &&
		keySector[2] == 'v' &&
		keySector[3] == '1'
}

func (JMicron) NeedsKEK() bool {
	return true
}

func (JMicron) CreateDecrypter(keySector []byte, kek []byte) (cipher *aes.Cipher) {
	var i int

	// decrypt the key sector
	Reverse(kek)
	kekcipher, err := aes.NewCipher(kek)
	if err != nil {
		panic("[BUG] error creating KEK decrypter in JMicron.CreateDecrypter()")
	}
	for i := 0; i < len(keySector); i += 16 {
		Reverse(keySector[i:i + 16])
		kekcipher.Decrypt(keySector[i:i + 16], keySector[i:i + 16])
		Reverse(keySector[i:i + 16])
	}

	// find the DEK itself
	// this can be anywhere in the block
	found := false
	for i = 0; i < len(keySector) - 4; i++ {
		found = keySector[i + 0] == 'D' &&
			keySector[i + 1] == 'E' &&
			keySector[i + 2] == 'K' &&
			keySector[i + 3] == '1'
		if found {
			break
		}
	}
	if !found {			// wrong KEK, sorry
		return nil
	}

	// the first half of the key is the 16 bytes at i + 0xC
	// the second half of the key is the 16 bytes at i + 0x20
	// the key size is the byte at i + 0x58
	if keySector[i + 0x58] != 0x20 {
		panic("[HELP] TODO write this")
	}
	dek := make([]byte, 32)
	copy(dek[:16], keySector[i + 0xC:])
	copy(dek[16:], keySector[i + 0x20:])
	Reverse(dek)
	dekcipher, err := aes.NewCipher(dek)
	if err != nil {
		panic("[BUG] TODO")
	}
	return dekcipher
}

func (JMicron) Decrypt(c *aes.Cipher, b []byte) {
}

func init() {
	Bridges = append(Bridges, JMicron{})
}
