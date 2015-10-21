// 13 january 2015
// trymbr.go 11 january 2015
package main

import (
	"os"
	"io"
	"strconv"
)

const keyblockoff = 1000202059776
const blocksize = 512

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pos, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(int64(pos), 0)
	if err != nil {
		panic(err)
	}

	block := make([]byte, blocksize)
	_, err = io.ReadFull(f, block)
	if err != nil {
		panic(err)
	}
	_, err = os.Stdout.Write(block)
	if err != nil {
		panic(err)
	}
}
