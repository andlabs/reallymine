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
// 	It should return "", true if the user cancelled the operation or string, false otherwise
// 	And yes, the password is a string; see kek.go for details.
// TryGetDecrypter(that function)
// If that returns nil, the user aborted the operation; stop
// Seek back to start
// Allocate a buffer of the appropriate size of each block to decrypt at a time
// for DecryptNext(...)
// 	Update a progress bar or something

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

func TryGetDecrypter(keySector []byte, bridge Bridge, askPassword func(firstTime bool) (password string, cancelled bool)) (c cipher.Block) {
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
		kek := KEKFromPassword(password)
		c = try(keySector, bridge, kek)
		firstTime = false // in case the password was wrong
	}
	return c
}

func DecryptNext(from io.Reader, to io.Writer, bridge Bridge, c cipher.Block, sectors []byte) (more bool) {
	n, err := from.Read(sectors)
	if err == io.EOF {
		if n == 0 {
			return false		// no more
		}
		// otherwise pretend no error
		// the final read might be smaller than the buffer
		// the next one should return err == io.EOF and n == 0
		err = nil
	}
	if err != nil {
		BUG("error reading sectors in DecryptNextSector(): %v", err)
	}
	// handle the final block properly if it's shorter
	sectors = sectors[:n]
	bridge.Decrypt(c, sectors)
	_, err = to.Write(sectors)
	if err != nil {
		BUG("error writing decrypted sector in DecryptNext(): %v", err)
	}
	return true
}
