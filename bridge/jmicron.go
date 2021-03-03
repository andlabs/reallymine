// 21 october 2015
package bridge

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"

	"github.com/undeadbanegithub/reallymine/byteops"
	"github.com/undeadbanegithub/reallymine/decryptloop"
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

type JMicronKeySector struct {
	raw []byte

	// The names Key3EE2, Key3EF2, and Key3F02 are from the
	// paper. But I recognize the hex numbers as addresses in the
	// JMicron chip's RAM. These RAM addresses followed me
	// around throughout disassembly, and I *knew* they were
	// suspicious, damnit!
	d struct { // d for "DEK block"
		Magic     [4]byte // 'DEK1'
		Checksum  uint16  // TODO check this too?
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
}

func (JMicron) DecryptKeySector(keySector []byte, kek []byte) (KeySector, error) {
	// copy these to avoid overwriting them
	keySector = byteops.DupBytes(keySector)
	kek = byteops.DupBytes(kek)

	byteops.Reverse(kek)
	kekcipher, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(keySector); i += 16 {
		block := keySector[i : i+16]
		byteops.Reverse(block)
		kekcipher.Decrypt(block, block)
		byteops.Reverse(block)
	}

	return &JMicronKeySector{
		raw: keySector,
	}, nil
}

func (ks *JMicronKeySector) Raw() []byte {
	return ks.raw
}

// the DEK can be anywhere in the decrypted key sector
func (ks *JMicronKeySector) findDEK() (offset int) {
	for i := 0; i < len(ks.raw)-4; i++ {
		if ks.raw[i+0] == 'D' &&
			ks.raw[i+1] == 'E' &&
			ks.raw[i+2] == 'K' &&
			ks.raw[i+3] == '1' {
			return i
		}
	}
	return -1 // not found; this isn't the right KEK
}

func (ks *JMicronKeySector) DEK() (dek []byte, err error) {
	offset := ks.findDEK()
	if offset == -1 {
		return nil, ErrWrongKEK
	}

	r := bytes.NewReader(ks.raw[offset:])
	// The endianness is likely wrong. We don't use any of
	// the endian-dependent fields, though. I can figure the
	// correct endianness from the disassembly if they're ever
	// actually needed.
	err = binary.Read(r, binary.BigEndian, &(ks.d))
	if err != nil {
		return nil, err
	}

	if ks.d.KeySize != 0x20 {
		return nil, incompleteImpl("The size of the encryption key in your JMicron sector (%d) is not known.", ks.d.KeySize)
	}

	dek = make([]byte, 32)
	copy(dek[:16], ks.d.Key3EE2[:])
	copy(dek[16:], ks.d.Key3EF2[:])
	byteops.Reverse(dek)
	return dek, nil
}

func (JMicron) DecryptLoopSteps() decryptloop.StepList {
	return decryptloop.StepList{
		decryptloop.StepReverse,
		decryptloop.StepDecrypt,
		decryptloop.StepReverse,
	}
}

func init() {
	Bridges = append(Bridges, JMicron{})
}
