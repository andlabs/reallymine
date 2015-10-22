// 22 october 2015
package main

import (
	"io"
	"crypto/aes"
)

// These do the actual work of recovery.
// Pseudocode:
// Open encrypted medium
// Seek to the end of the medium, get its position (media size)
// FindKeySectorAndBridge(medium, media size), assume it succeeded
// If the bridge needs a KEK
// 	TryDefaultKEK()
// 	While that fails
// 		Ask the user for their password
// 		TryUserPassword()
// Else
// 	TryWithoutKEK()
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
	return nil, nil		// no key sector found :(
}

// I don't know when this is used, but have it here it anyway
var DefaultKEK128 = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x2B, 0x99, 0x2D, 0xDF, 0xA2, 0x32, 0x49, 0xD6,
}

var DefaultKEK = []byte{
	0x03, 0x14, 0x15, 0x92, 0x65, 0x35, 0x89, 0x79, 0x32, 0x38, 0x46, 0x26, 0x43, 0x38, 0x32, 0x79,
	0xFC, 0xEB, 0xEA, 0x6D, 0x9A, 0xCA, 0x76, 0x86, 0xCD, 0xC7, 0xB9, 0xD9, 0xBC, 0xC7, 0xCD, 0x86,
}

func tryKEK(bridge Bridge, keySector []byte, kek []byte) *aes.Cipher {
	return bridge.CreateDecrypter(keySector, kek)
}

func TryDefaultKEK(bridge Bridge, keySector []byte) *aes.Cipher {
	return tryKEK(bridge, keySector, DefaultKEK)
}

func TryUserPassword(bridge Bridge, keySector []byte, password []byte) *aes.Cipher {
	BUG("TODO UNIMPLEMENTED")
	panic("unreachable")
}

func TryWithoutKEK(bridge Bridge, keySector []byte) *aes.Cipher {
	return tryKEK(bridge, keySector, nil)
}
