// 23 october 2015
package bridge

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"encoding/binary"
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

type InitioKeySector struct {
	raw	[]byte
	d	struct {		// d for "DEK block"
		Magic   [4]byte // 27 5D BA 35
		Unknown [8]byte
		Key     [32]byte // stored as little-endian longs
	}
}

func (Initio) DecryptKeySector(keySector []byte, kek []byte) (KeySector, error) {
	// copy to avoid clobbering
	keySector = DupBytes(keySector)
	kek = DupBytes(kek)

	SwapHalves(kek)
	Reverse(kek)
	kekcipher, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(keySector); i += 16 {
		block := keySector[i : i+16]
		SwapLongs(block)
		kekcipher.Decrypt(block, block)
		// Don't swap back; it'll be correct as-is.
	}

	return &InitioKeySector{
		raw:		keySector,
	}, nil
}

func (ks *InitioKeySector) Raw() []byte {
	return ks.raw
}

func (ks *InitioKeySector) valid() bool {
	return ks.d.Magic[0] == 0x27 &&
		ks.d.Magic[1] == 0x5D &&
		ks.d.Magic[2] == 0xBA &&
		ks.d.Magic[3] == 0x35
}

// Unlike the JMicron one, the Initio DEK block is at a fixed offset
// into the key sector.
const initioDEKOffset = 0x190

func (ks *InitioKeySector) extractDEKBlock() error {
	r := bytes.NewReader(ks.raw[initioDEKOffset:])
	// The endianness is most likely right, but unimportant since every field is [...]byte.
	return binary.Read(r, binary.LittleEndian, &(ks.d))
}

func (ks *InitioKeySector) ExtractDEK() (dek []byte, err error) {
	err = extractDEKBlock()
	if err != nil {
		return nil, err
	}
	if !ks.valid() {
		return nil, ErrWrongKEK
	}

	// make a copy to avoid altering ks.d
	dek = DupBytes(ks.d.Key[:])
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
	Bridges = append(Bridges, Initio{})
}
