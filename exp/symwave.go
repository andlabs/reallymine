// 23 october 2015
package main

import (
	"crypto/cipher"
	"bytes"
	"encoding/binary"

	"github.com/mendsley/gojwe"

"fmt"
"encoding/hex"
"os"
"io/ioutil"
)

type Symwave struct{}

func (Symwave) Name() string {
	return "Symwave"
}

func (Symwave) Is(keySector []byte) bool {
	return keySector[0] == 'S' &&
		keySector[1] == 'Y' &&
		keySector[2] == 'M' &&
		keySector[3] == 'W'
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
	// This is definitely correct; the 68000 is big endian.
	err := binary.Read(r, binary.BigEndian, &ks)
	if err != nil {
		BUG("error reading key sector into structure in Symwave.CreateDecrypter(): %v", err)
	}

	kek, err = gojwe.AesKeyUnwrap(symwaveKEKWrappingKey, ks.WrappedKEK[:])
	if err != nil {
		BUG("error unwrapping KEK in Symwave.CreateDecrypter(): %v", err)
	}

	dek1, err := gojwe.AesKeyUnwrap(kek, ks.WrappedDEK1[:])
	if err != nil {
		BUG("error unwrapping DEK part 1 in Symwave.CreateDecrypter(): %v", err)
	}

	dek2, err := gojwe.AesKeyUnwrap(kek, ks.WrappedDEK2[:])
	if err != nil {
		BUG("error unwrapping DEK part 2 in Symwave.CreateDecrypter(): %v", err)
	}

fmt.Printf("dek1:\n")
fmt.Print(hex.Dump(dek1))
fmt.Printf("dek2:\n")
fmt.Print(hex.Dump(dek2))
panic("not done yet")
}

func (Symwave) Decrypt(c cipher.Block, b []byte) {
	// TODO
}

/*
func init() {
	Bridges = append(Bridges, Symwave{})
}
*/

func main() {
	b, _ := ioutil.ReadAll(os.Stdin)
	Symwave{}.CreateDecrypter(b, nil)
}
