// 22 october 2015
package main

import (
	"crypto/cipher"
	"io"
)

// These do the actual work of recovery.
// Pseudocode:
// Open encrypted medium
// Seek to the end of the medium, get its position (media size)
// FindKeySectorAndBridge(medium, media size), assume it succeeded
// Write a function to ask for the user password
// 	It takes a bool; if true, this is the first time; if false, the password was wrong
// 	It should return nil, true if the user cancelled the operation or non-nil, false otherwise
// TryGetDecrypter(that function)
// If that returns nil, the user aborted the operation; stop
// Seek back to start
// While there are sectors to read
// 	Read a sector
// 	Decrypt sector using the bridge's Decrypt() method
// 	Write it back

// TODO make this stop early, giving the user the option to continue
func FindKeySectorAndBridge(media io.ReaderAt, startAt int64) (keySector []byte, bridge Bridge) {
	sector := make([]byte, SectorSize)
	pos := startAt - SectorSize
	for pos >= 0 {
		_, err := media.ReadAt(sector, pos)
		// io.ReaderAt specifies that EOF may be returned when reading right at the end of the file
		if err != nil && err != io.EOF {
			BUG("error reading sector in FindKeySectorAndBridge(): %v", err)
		}
		bridge = IdentifyKeySector(sector)
		if bridge != nil {
			return sector, bridge
		}
		// not the key sector; keep going
		pos -= SectorSize
	}
	return nil, nil // no key sector found :(
}

// I don't know when this is used, but have it here it anyway
var DefaultKEK128 = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x2B, 0x99, 0x2D, 0xDF, 0xA2, 0x32, 0x49, 0xD6,
}

var DefaultKEK = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x32, 0x38, 0x46, 0x26, 0x43, 0x38, 0x32, 0x79,
	0xFC, 0xEB, 0xEA, 0x6D, 0x9A, 0xCA, 0x76, 0x86, 0xCD, 0xC7, 0xB9, 0xD9, 0xBC, 0xC7, 0xCD, 0x86,
}

func TryGetDecrypter(keySector []byte, bridge Bridge, askPassword func(firstTime bool) (password []byte, cancelled bool)) (c cipher.Block) {
	try := func(keySector []byte, bridge Bridge, kek []byte) cipher.Block {
		return bridge.CreateDecrypter(keySector, kek)
	}

	if !bridge.NeedsKEK() {
		return try(keySector, bridge, nil) // should not return nil
	}

	c = try(keySector, bridge, DefaultKEK)
	firstTime := true
	for c == nil { // whlie the default KEK didn't work or the user password is wrong
		password, cancelled := askPassword(firstTime)
		if cancelled { // user aborted
			return nil
		}
		// TODO
		_ = password
		firstTime = false // in case the password was wrong
	}
	return c
}
