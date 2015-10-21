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
	for i := 0; i < len(b) / 2; i++ {
		n := len(b) - i - 1
		b[i], b[n] = b[n], b[i]
	}
}

func swap(b []byte) {
	for i := 0; i < len(b); i += 4 {
		b[i + 0], b[i + 3] = b[i + 3], b[i + 0]
		b[i + 1], b[i + 2] = b[i + 2], b[i + 1]
	}
}

func fliphalves(b []byte) {
	c := make([]byte, len(b))
	copy(c, b[16:])
	copy(c[16:], b)
	copy(b, c)
}

func decryptDEK(dek []byte, key []byte) error {
	fliphalves(key)
	reverse(key)
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	for i := 0; i < len(dek); i += 16 {
		swap(dek[i:i + 16])
		c.Decrypt(dek[i:], dek[i:])
		// don't swap back; it'll be correct this way
	}
	return nil
}

type Initio struct {
	Magic		[4]byte		// 27 5D BA 35
	Unknown		[8]byte
	Key			[32]byte		// stored as little-endian longs
}

var errWrongChip = fmt.Errorf("wrong chip")

func ReadInitio(dek []byte) (j *Initio, err error) {
	i := 0x190
	if dek[i] == 0x27 &&
		dek[i + 1] == 0x5D &&
		dek[i + 2] == 0xBA &&
		dek[i + 3] == 0x35 {
		// do nothing
	} else {
		return nil, errWrongChip
	}

	j = new(Initio)
	r := bytes.NewReader(dek[i:])
	err = binary.Read(r, binary.LittleEndian, j)
	if err != nil {
		return nil, err
	}

	swap(j.Key[:])

	return j, nil
}

func main() {
	fdek, _ := os.Open(os.Args[1])
	fblock, _ := os.Open(os.Args[2])
	dek, _ := ioutil.ReadAll(fdek)
	block, _ := ioutil.ReadAll(fblock)
	fdek.Close()
	fblock.Close()

	decryptDEK(dek, pi)
	j, err := ReadInitio(dek)
	if err != nil { panic(err) }
	key := make([]byte, len(j.Key))
	copy(key, j.Key[:])
	fliphalves(key)
	reverse(key)
	c, _ := aes.NewCipher(key)
	for i := 0; i < len(block); i += 16 {
//		reverse(block[i:i + 16])
		swap(block[i:i + 16])
		c.Decrypt(block[i:], block[i:])
		// we DO need to swap after this though
		swap(block[i:i + 16])
	}
	fmt.Println(hex.Dump(block))
}
