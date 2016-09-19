// 21 october 2015
package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"encoding/binary"
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

func jmicronDecryptKeySector(keySector []byte, kek []byte) error {
	Reverse(kek)
	kekcipher, err := aes.NewCipher(kek)
	if err != nil {
		return err
	}
	for i := 0; i < len(keySector); i += 16 {
		block := keySector[i : i+16]
		Reverse(block)
		kekcipher.Decrypt(block, block)
		Reverse(block)
	}
	return nil
}

// the DEK can be anywhere in the decrypted key sector
func jmicronFindDEK(keySector []byte) (offset int) {
	for i := 0; i < len(keySector)-4; i++ {
		if keySector[i+0] == 'D' &&
			keySector[i+1] == 'E' &&
			keySector[i+2] == 'K' &&
			keySector[i+3] == '1' {
			return i
		}
	}
	return -1 // not found; this isn't the right KEK
}

// The names Key3EE2, Key3EF2, and Key3F02 are from the paper.
// But I recognize the hex numbers as addresses in the JMicron chip's
// RAM. These RAM addresses followed me around throughout
// disassembly, and I *knew* they were suspicious, damnit!
type jmicronDEKBlock struct {
	Magic     [4]byte // 'DEK1'
	Checksum  uint16
	Unknown   uint16
	Random1   uint32
	Key3EE2   [16]byte // This is the first half of the AES-256 key.
	Random2   uint32
	Key3EF2   [16]byte // This is the second half of the AES-256 key.
	Random3   uint32
	Key3F02   [32]byte // I don't know what this is but I highly doubt it's a key.
	Random4   uint32
	KeySize   byte
	Remaining [1 + 4 + 2]byte
}

// TODO make this more like the initio one
func jmicronExtractDEK(keySector []byte, offset int) ([]byte, error) {
	var dekblock jmicronDEKBlock

	r := bytes.NewReader(keySector[offset:])
	// The endianness is likely wrong. We don't use any of
	// the endian-dependent fields, though. I can figure the
	// correct endianness from the disassembly if they're ever
	// actually needed.
	err := binary.Read(r, binary.BigEndian, &dekblock)
	if err != nil {
		return nil, err
	}

	if dekblock.KeySize != 0x20 {
		return nil, IncompleteImplementation("The size of the encryption key in your JMicron sector (%d) is not known.", dekblock.KeySize)
	}

	dek := make([]byte, 32)
	copy(dek[:16], dekblock.Key3EE2[:])
	copy(dek[16:], dekblock.Key3EF2[:])
	Reverse(dek)
	return dek, nil
}

func (JMicron) ExtractDEK(keySector []byte, kek []byte) (dek []byte, err error) {
	// make a copy of these so the originals aren't touched
	keySector = DupBytes(keySector)
	kek = DupBytes(kek)

	err = jmicronDecryptKeySector(keySector, kek)
	if err != nil {
		return nil, err
	}
	offset := jmicronFindDEK(keySector)
	if offset == -1 { // wrong KEK
		return nil, ErrWrongKEK
	}
	return jmicronExtractDEK(keySector, offset)
}

func (JMicron) Decrypt(c cipher.Block, b []byte) {
	for i := 0; i < len(b); i += 16 {
		block := b[i : i+16]
		Reverse(block)
		c.Decrypt(block, block)
		Reverse(block)
	}
}

func init() {
	Bridges = append(Bridges, JMicron{})
}
