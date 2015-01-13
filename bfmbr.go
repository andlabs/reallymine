// 12 january 2015
// from trymbr.go 11 january 2015
package main

import (
	"fmt"
	"os"
	"crypto/aes"
	"bytes"
)

const blocksize = 512
const keysize = 16		// AES-128
//const keysize = 24		// AES-192
//const keysize = 32		// AES-256
//const keysize = 8			// DES

const firstPotentialOff = 0//0x50
const lastPotentialOff = blocksize

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Seek(keyblockoff, 0)
	if err != nil {
		panic(err)
	}
	keyblock := make([]byte, blocksize)
	n, err := f.Read(keyblock)
	if err != nil {
		panic(err)
	} else if n != blocksize {
		panic(n)
	}

	_, err = f.Seek(mbroff, 0)
	if err != nil {
		panic(err)
	}
	mbr := make([]byte, blocksize)
	n, err = f.Read(mbr)
	if err != nil {
		panic(err)
	} else if n != blocksize {
		panic(n)
	}

	mbrout := make([]byte, blocksize)

	try := func(key []byte, xfmt string, xargs ...interface{}) {
		cipher, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}
		cbs := cipher.BlockSize()
		if (lastPotentialOff - firstPotentialOff) % cbs != 0 {
			panic("uneven blocks")
		}
		for j := 0; j < blocksize; j += cbs {
			cipher.Decrypt(mbrout[j:], mbr[j:])
		}
		if bytes.Contains(mbrout, []byte("pera")) {
			fmt.Printf("%x ", key)
			fmt.Printf(xfmt, xargs...)
			fmt.Printf("\n")
		}
	}

	tryrev := func(key []byte, xfmt string, xargs ...interface{}) {
		for rev := 0; rev < keysize; rev += 4 {
			key[rev + 0], key[rev + 3] = key[rev + 3], key[rev + 0]
			key[rev + 1], key[rev + 2] = key[rev + 2], key[rev + 1]
		}
		try(key, "u32 byteswap " + xfmt, xargs...)
	}

	for i := firstPotentialOff; i + keysize <= lastPotentialOff; i++ {
		key := make([]byte, keysize)
		copy(key, keyblock[i:])
		try(key, "normal")
		tryrev(key, "normal")
	}

	for i := firstPotentialOff; i + keysize <= lastPotentialOff; i++ {
		for j := firstPotentialOff; j + keysize <= lastPotentialOff; j++ {
			key := make([]byte, keysize)
			for k := 0; k < keysize; k++ {
				key[k] = keyblock[i + k] ^ keyblock[j + k]
			}
			try(key, "xor %d %d", i, j)
			tryrev(key, "xor %d %d", i, j)
		}
	}
}
