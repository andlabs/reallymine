// 23 october 2015
package initio

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"encoding/binary"

	"github.com/andlabs/reallymine/bridge"
)

type Initio struct{}

func (Initio) Name() string {
	return "Initio"
}

func (Initio) Is(keySector []byte) bool {
	return keySector[0] == 'W' &&
		keySector[1] == 'D' &&
		keySector[2] == 0x01 &&
		keySector[3] == 0x14
}

func (Initio) NeedsKEK() bool {
	return true
}

func decryptKeySector(keySector []byte, kek []byte) error {
	SwapHalves(kek)
	Reverse(kek)
	kekcipher, err := aes.NewCipher(kek)
	if err != nil {
		return err
	}
	for i := 0; i < len(keySector); i += 16 {
		block := keySector[i : i+16]
		SwapLongs(block)
		kekcipher.Decrypt(block, block)
		// Don't swap back; it'll be correct as-is.
	}
	return nil
}

type dekBlock struct {
	Magic   [4]byte // 27 5D BA 35
	Unknown [8]byte
	Key     [32]byte // stored as little-endian longs
}

func (d *dekBlock) valid() bool {
	return d.Magic[0] == 0x27 &&
		d.Magic[1] == 0x5D &&
		d.Magic[2] == 0xBA &&
		d.Magic[3] == 0x35
}

// Unlike the JMicron one, the Initio DEK block is at a fixed offset
// into the key sector.
const dekOffset = 0x190

func extractDEKBlock(keySector []byte) (*dekBlock, error) {
	dekblock := new(dekBlock)
	r := bytes.NewReader(keySector[initioDEKOffset:])
	// The endianness is most likely right, but unimportant since every field is [...]byte.
	err := binary.Read(r, binary.LittleEndian, dekblock)
	if err != nil {
		return nil, err
	}
	return dekblock, nil
}

func (Initio) ExtractDEK(keySector []byte, kek []byte) (dek []byte, err error) {
	// make a copy of these so the originals aren't touched
	keySector = DupBytes(keySector)
	kek = DupBytes(kek)

	err = decryptKeySector(keySector, kek)
	if err != nil {
		return nil, err
	}
	dekblock, err := extractDEKBlock(keySector)
	if err != nil {
		return nil, err
	}
	if !dekblock.valid() {
		return nil, ErrWrongKEK
	}

	dek = dekblock.Key[:]
	SwapLongs(dek) // undo the little-endian-ness
	SwapHalves(dek)
	Reverse(dek)
	return dek, nil
}

func (Initio) Decrypt(c cipher.Block, b []byte) {
	for i := 0; i < len(b); i += 16 {
		block := b[i : i+16]
		SwapLongs(block)
		c.Decrypt(block, block)
		// We DO need to swap again after this, though!
		SwapLongs(block)
	}
}

func init() {
	bridge.Bridges = append(bridge.Bridges, Initio{})
}
