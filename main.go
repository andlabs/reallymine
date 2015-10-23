// 22 october 2015
package main

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open(os.Args[1])
	fout, _ := os.Create(os.Args[2])

	size, _ := f.Seek(0, 2)
	keySector, bridge := FindKeySectorAndBridge(f, size)
	if keySector == nil {
		fmt.Println("no key sector found")
		return
	}
	fmt.Println("found " + bridge.Name())

	c := TryGetDecrypter(keySector, bridge, func(firstTime bool) (password []byte, cancelled bool) {
		if firstTime {
			fmt.Println("We need the drive's password to decrypt your drive.")
		} else {
			fmt.Println("Password incorrect.")
		}
		// TODO
		return nil, true
	})
	if c == nil {
		fmt.Println("User aborted.")
		return
	}

	_, err := f.Seek(0, 0)
	if err != nil {
		// TODO
		panic(err)
	}
	for DecryptNextSector(f, fout, bridge, c) {
		// TODO
	}
}
