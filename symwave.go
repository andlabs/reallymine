// 23 october 2015
package main

import (
	"crypto/cipher"
	"bytes"
	"encoding/binary"

	"github.com/mendsley/gojwe"
)

type Symwave struct{}

func (Symwave) Name() string {
	return "Symwave"
}

func (Symwave) Is(keySector []byte) bool {
	// note: stored little endian despite being a big endian system
	return keySector[3] == 'S' &&
		keySector[2] == 'Y' &&
		keySector[1] == 'M' &&
		keySector[0] == 'W'
}

func (Symwave) NeedsKEK() bool {
	return false
}

// The DEK is stored as two separately-wrapped halves.
// The KEK is only stored as one.
type symwaveKeySector struct {
	Magic		[4]byte
	Unknown		[0xC]byte
	WrappedDEK1	[0x28]byte
	WrappedDEK2	[0x28]byte
	WrappedKEK	[0x28]byte
}

// This is hardcoded into the Symwave firmware.
var symwaveKEKWrappingKey = []byte{
	0x29, 0xA2, 0x60, 0x7A,
	0xEA, 0x0B, 0x64, 0xAB,
	0x7B, 0xB3, 0xB9, 0xAB,
	0xA5, 0x69, 0x8B, 0x40,
	0x2E, 0x47, 0x93, 0xA6,
	0x81, 0x45, 0xC9, 0xCC,
	0x79, 0x94, 0x6A, 0x01,
	0x84, 0x0B, 0x34, 0xFE,
}

func (Symwave) CreateDecrypter(keySector []byte, kek []byte) (c cipher.Block) {
	var ks symwaveKeySector

	r := bytes.NewReader(keySector)
	// Again, stored as little endian for some reason; this is a 68000 system so it should be big endian...
	err := binary.Read(r, binary.LittleEndian, &ks)
	if err != nil {
		BUG("error reading key sector into structure in Symwave.CreateDecrypter(): %v", err)
	}

	// And again with the endianness stuff...
	wrapped := ks.WrappedKEK[:]
	SwapLongs(wrapped)
	kek, err = gojwe.AesKeyUnwrap(symwaveKEKWrappingKey, wrapped)
	if err != nil {
		BUG("error unwrapping KEK in Symwave.CreateDecrypter(): %v", err)
	}

	wrapped = ks.WrappedDEK1[:]
	SwapLongs(wrapped)
	dek1, err := gojwe.AesKeyUnwrap(kek, wrapped)
	if err != nil {
		BUG("error unwrapping DEK part 1 in Symwave.CreateDecrypter(): %v", err)
	}

	wrapped = ks.WrappedDEK2[:]
	SwapLongs(wrapped)
	dek2, err := gojwe.AesKeyUnwrap(kek, wrapped)
	if err != nil {
		BUG("error unwrapping DEK part 2 in Symwave.CreateDecrypter(): %v", err)
	}

	_ = dek2
	// And finally we just need one last endian correction...
	SwapLongs(dek1)
	return NewAES(dek1)
}

func (Symwave) Decrypt(c cipher.Block, b []byte) {
	for i := 0; i < len(b); i += 16 {
		block := b[i : i+16]
		// ...and we can just use block as-is!
		c.Decrypt(block, block)
	}
}

func init() {
	Bridges = append(Bridges, Symwave{})
}
