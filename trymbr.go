// 11 january 2015
package main

import (
	"fmt"
	"os"
	"encoding/hex"
	"crypto/aes"
)

const keyblockoff = 1000202059776
const mbroff = 0
const blocksize = 512
//const keysize = 16		// AES-128
//const keysize = 24		// AES-192
const keysize = 32		// AES-256

const firstPotentialOff = 0x50
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

	for i := firstPotentialOff; i + keysize <= lastPotentialOff; i++ {
		key := keyblock[i:i + keysize]
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
		hexout := hex.Dump(mbrout)
		fmt.Printf("0x%X\n%s\n", i, hexout)
	}
}
