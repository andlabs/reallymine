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

func (JMicron) decryptKeySector(keySector []byte, kek []byte) {
	Reverse(kek)
	kekcipher := NewAES(kek)
	for i := 0; i < len(keySector); i += 16 {
		block := keySector[i:i + 16]
		Reverse(block)
		kekcipher.Decrypt(block, block)
		Reverse(block)
	}
}

// the DEK can be anywhere in the decrypted key sector
func (JMicron) findDEK(keySector []byte) (offset int) {
	for i = 0; i < len(keySector) - 4; i++ {
		if keySector[i + 0] == 'D' &&
			keySector[i + 1] == 'E' &&
			keySector[i + 2] == 'K' &&
			keySector[i + 3] == '1' {
			return i
		}
	}
	return -1		// not found; this isn't the right KEK
}

func (JMicron) extractDEK(keySector []byte, offset int) []byte {
	keySector = keySector[offset:]
	// the first half of the key is the 16 bytes at offset 0xC
	// the second half of the key is the 16 bytes at offset 0x20
	// the key size is the byte at offset 0x58
	if keySector[0x58] != 0x20 {
		HELP("The size of the encryption key in your JMicron sector is not known.\nThis means your drive is new to reallymine, and support must be added.\nPlease help us!")
	}
	dek := make([]byte, 32)
	copy(dek[:16], keySector[0xC:])
	copy(dek[16:], keySector[0x20:])
	Reverse(dek)
	return dek
}

func (j JMicron) CreateDecrypter(keySector []byte, kek []byte) (cipher *aes.Cipher) {
	// make a copy of these so the originals aren't touched
	keySector = DupBytes(keySector)
	kek = DupBytes(kek)

	j.decryptKeySector(keySector, kek)
	offset := j.findDEK(keySector)
	if offset == -1 {		// wrong KEK
		return nil
	}
	return NewAES(j.extractDEK(keySector, offset))
}

func (JMicron) Decrypt(c *aes.Cipher, b []byte) {
	for i := 0; i < len(b); i += 16 {
		block := b[i:i + 16]
		Reverse(block)
		c.Decrypt(block, block)
		Reverse(block)
	}
}

func init() {
	Bridges = append(Bridges, JMicron{})
}
