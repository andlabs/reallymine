// 24 october 2015
package bridge

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"

	"github.com/undeadbanegithub/reallymine/byteops"
	"github.com/undeadbanegithub/reallymine/decryptloop"
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

type PLXKeySector struct {
	raw []byte
	d   struct {
		Magic        [4]byte
		Unknown      [0x10]byte
		EncryptedDEK [32]byte
	}
	dek []byte
}

// MAJOR TODO
// Locate the second copy of the SInE block and test that one too.
// I need to find out if it will /always/ be at the same place or nt

// TODO how does the chip know the KEK was valid?
func (PLX) DecryptKeySector(keySector []byte, kek []byte) (KeySector, error) {
	// make a copy of these so the originals aren't touched
	keySector = byteops.DupBytes(keySector)
	kek = byteops.DupBytes(kek)

	ks := new(PLXKeySector)
	ks.raw = keySector

	r := bytes.NewReader(ks.raw)
	// TODO copy comment from jmicron.go here; find out what endianness the ARM in the chip is
	err := binary.Read(r, binary.BigEndian, &(ks.d))
	if err != nil {
		return nil, err
	}

	byteops.SwapLongs(kek)
	byteops.Reverse(kek)
	kekcipher, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	ks.dek = ks.d.EncryptedDEK[:]
	kekcipher.Decrypt(ks.dek[:16], ks.dek[:16])
	kekcipher.Decrypt(ks.dek[16:], ks.dek[16:])
	return ks, nil
}

func (ks *PLXKeySector) Raw() []byte {
	return ks.raw
}

func (ks *PLXKeySector) DEK() ([]byte, error) {
	return ks.dek, nil
}

func (PLX) DecryptLoopSteps() decryptloop.StepList {
	return decryptloop.StepList{
		// TODO
	}
}

func init() {
	Bridges = append(Bridges, PLX{})
}
