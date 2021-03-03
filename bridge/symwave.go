// 23 october 2015
package bridge

import (
	"bytes"
	"encoding/binary"

	"github.com/mendsley/gojwe"
	"github.com/undeadbanegithub/reallymine/byteops"
	"github.com/undeadbanegithub/reallymine/decryptloop"
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

type SymwaveKeySector struct {
	raw []byte

	// The DEK is stored as two separately-wrapped halves.
	// The KEK is only stored as one.
	d struct {
		Magic       [4]byte
		Unknown     [0xC]byte
		WrappedDEK1 [0x28]byte
		WrappedDEK2 [0x28]byte
		WrappedKEK  [0x28]byte
	}
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

func (Symwave) DecryptKeySector(keySector []byte, kek []byte) (KeySector, error) {
	return &SymwaveKeySector{
		raw: byteops.DupBytes(keySector),
	}, nil
}

func (ks *SymwaveKeySector) Raw() []byte {
	return ks.raw
}

func (ks *SymwaveKeySector) DEK() (dek []byte, err error) {
	r := bytes.NewReader(ks.raw)
	// Again, stored as little endian for some reason; this is a 68000 system so it should be big endian...
	err = binary.Read(r, binary.LittleEndian, &(ks.d))
	if err != nil {
		return nil, err
	}

	// And again with the endianness stuff...
	wrapped := ks.d.WrappedKEK[:]
	byteops.SwapLongs(wrapped)
	kek, err := gojwe.AesKeyUnwrap(symwaveKEKWrappingKey, wrapped)
	if err != nil {
		return nil, err
	}

	wrapped = ks.d.WrappedDEK1[:]
	byteops.SwapLongs(wrapped)
	dek1, err := gojwe.AesKeyUnwrap(kek, wrapped)
	if err != nil {
		return nil, err
	}

	wrapped = ks.d.WrappedDEK2[:]
	byteops.SwapLongs(wrapped)
	dek2, err := gojwe.AesKeyUnwrap(kek, wrapped)
	if err != nil {
		return nil, err
	}

	dek = byteops.DupBytes(dek1)
	_ = dek2 // doesn't seem to be used
	// And finally we just need one last endian correction...
	byteops.SwapLongs(dek)
	return dek, nil
}

func (Symwave) DecryptLoopSteps() decryptloop.StepList {
	return decryptloop.StepList{
		// ...and we can just decrypt the encrypted blocks as-is!
		decryptloop.StepDecrypt,
	}
}

func init() {
	Bridges = append(Bridges, Symwave{})
}
