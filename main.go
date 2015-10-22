// 22 october 2015
package main

import (
	"fmt"
	"os"
	"crypto/aes"
)

func main() {
	var c *aes.Cipher

	f, _ := os.Open(os.Args[1])
	fout, _ := os.Create(os.Args[2])

	size, _ := f.Seek(0, 2)
	keySector, bridge := FindKeySectorAndBridge(f, size)
	if keySector == nil {
		fmt.Println("no key sector found")
		return
	}
	fmt.Println("found " + bridge.Name())

	if bridge.NeedsKEK() {
		c = TryDefaultKEK(bridge, keySector)
		for c == nil {
			panic("TODO")
		}
	} else {
		c = GetWithoutKEK(bridge, keySector)
	}

	_, err := f.Seek(0, 0)
	if err != nil {
		// TODO
		panic(err)
	}
	sector := make([]byte, SectorSize)
	for {
		_, err := f.Read(sector)
		if err != nil { break }
		bridge.Decrypt(c, sector)
		fout.Write(sector)
	}
}
