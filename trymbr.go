// 11 january 2015
package main

import (
	"fmt"
	"os"
	"io"
	"encoding/hex"
	"crypto/aes"
)

const blocksize = 512
const keysize = 16		// AES-128
//const keysize = 24		// AES-192
//const keysize = 32		// AES-256
//const keysize = 8			// DES

const firstPotentialOff = 0x50
const lastPotentialOff = blocksize

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	keyblock := make([]byte, blocksize)
	_, err = io.ReadFull(f, keyblock)
	if err != nil {
		panic(err)
	}
	f.Close()

	f, err = os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	mbr := make([]byte, blocksize)
	_, err = io.ReadFull(f, mbr)
	if err != nil {
		panic(err)
	}
	f.Close()

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
