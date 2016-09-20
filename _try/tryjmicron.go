// 20 october 2015
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"crypto/aes"
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

var pi = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x32, 0x38, 0x46, 0x26, 0x43, 0x38, 0x32, 0x79,
	0xFC, 0xEB, 0xEA, 0x6D, 0x9A, 0xCA, 0x76, 0x86, 0xCD, 0xC7, 0xB9, 0xD9, 0xBC, 0xC7, 0xCD, 0x86,
}

func reverse(b []byte) {
	c := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		n := len(b) - i - 1
		c[i] = b[n]
	}
	copy(b, c)
}

func decryptDEK(dek []byte, key []byte) error {
	reverse(key)
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	for i := 0; i < len(dek); i += 16 {
		reverse(dek[i:i + 16])
		c.Decrypt(dek[i:], dek[i:])
		reverse(dek[i:i + 16])
	}
	return nil
}

type JMicron struct {
	Magic		[4]byte		// 'DEK1'
	Checksum	uint16
	Unknown		uint16
	Random1		uint32
	Key3EE2		[16]byte
	Random2		uint32
	Key3EF2		[16]byte
	Random3		uint32
	Key3F02		[32]byte
	Random4		uint32
	KeySize		byte
	Remaining	[1 + 4 + 2]byte
}

var errWrongChip = fmt.Errorf("wrong chip")

func ReadJMicron(dek []byte) (j *JMicron, err error) {
	i := 0
	for ; i < len(dek) - 4; i++ {
		if dek[i] == 'D' &&
			dek[i + 1] == 'E' &&
			dek[i + 2] == 'K' &&
			dek[i + 3] == '1' {
			break
		}
	}
	if i >= len(dek) - 4 {
		return nil, errWrongChip
	}

	j = new(JMicron)
	r := bytes.NewReader(dek[i:])
	// the true endianness isn't known; we deal with endian-agnostic data only
	err = binary.Read(r, binary.BigEndian, j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func swap(b []byte) {
	for i := 0; i < len(b); i += 4 {
		b[i + 0], b[i + 3] = b[i + 3], b[i + 0]
		b[i + 1], b[i + 2] = b[i + 2], b[i + 1]
	}
}

func main() {
	fdek, _ := os.Open(os.Args[1])
	fblock, _ := os.Open(os.Args[2])
	dek, _ := ioutil.ReadAll(fdek)
	block, _ := ioutil.ReadAll(fblock)
	fdek.Close()
	fblock.Close()

	decryptDEK(dek, pi)
	j, err := ReadJMicron(dek)
	if err != nil { panic(err) }
	fmt.Printf("%x\n", j.KeySize)
	key := make([]byte, 32)
	copy(key, j.Key3EE2[:])
	copy(key[16:], j.Key3EF2[:])
	reverse(key)
	c, err := aes.NewCipher(key)
	if err != nil { panic(err) }
	for i := 0; i < len(block); i += 16 {
		reverse(block[i:i + 16])
		c.Decrypt(block[i:], block[i:])
		reverse(block[i:i + 16])
	}
	fmt.Println(hex.Dump(block))
}
