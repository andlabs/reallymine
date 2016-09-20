// 24 october 2015
package main

import (
	"crypto/cipher"
	"bytes"
	"encoding/binary"

"fmt"
"encoding/hex"
"os"
"io/ioutil"
)

type PLX struct{}

func (PLX) Name() string {
	return "PLX (Oxford Semiconductor)"
}

func (PLX) Is(keySector []byte) bool {
	return keySector[0] == 'S' &&
		keySector[1] == 'I' &&
		keySector[2] == 'n' &&
		keySector[3] == 'E'
}

func (PLX) NeedsKEK() bool {
	return true
}

type plxKeySector struct {
	Magic		[4]byte
	Unknown		[0x10]byte
	EncryptedDEK	[32]byte
}

// MAJOR TODO
// Locate the second copy of the SInE block and test that one too.
// I need to find out if it will /always/ be at the same place or nt

// TODO how does the chip know the KEK was valid?
func (PLX) CreateDecrypter(keySector []byte, kek []byte) (c cipher.Block) {
	var ks plxKeySector

	// make a copy of these so the originals aren't touched
	keySector = DupBytes(keySector)
	kek = DupBytes(kek)

	SwapLongs(kek)
	Reverse(kek)
	kekcipher := NewAES(kek)

	r := bytes.NewReader(keySector)
	// This is definitely correct; the 68000 is big endian.
	err := binary.Read(r, binary.BigEndian, &ks)
	if err != nil {
		BUG("error reading key sector into structure in PLX.CreateDecrypter(): %v", err)
	}

	dek := ks.EncryptedDEK[:]
fmt.Printf("encrypted:\n")
fmt.Print(hex.Dump(dek))
	kekcipher.Decrypt(dek[:16], dek[:16])
	kekcipher.Decrypt(dek[16:], dek[16:])

fmt.Printf("decrypted:\n")
fmt.Print(hex.Dump(dek))
panic("not done yet")
}

func (PLX) Decrypt(c cipher.Block, b []byte) {
	// TODO
}

/*
func init() {
	Bridges = append(Bridges, PLX{})
}
*/

func main() {
	b, _ := ioutil.ReadAll(os.Stdin)
	b = b[0xAC:]
	kek := []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x32, 0x38, 0x46, 0x26, 0x43, 0x38, 0x32, 0x79,
	0xFC, 0xEB, 0xEA, 0x6D, 0x9A, 0xCA, 0x76, 0x86, 0xCD, 0xC7, 0xB9, 0xD9, 0xBC, 0xC7, 0xCD, 0x86,
}
	PLX{}.CreateDecrypter(b, kek)
}
