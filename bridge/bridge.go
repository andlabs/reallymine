// 21 october 2015
package bridge

import (
	"fmt"
	"crypto/cipher"
)

var ErrWrongKEK = fmt.Errorf("wrong KEK")

type Bridge interface {
	Name() string
	Is(keySector []byte) bool
	NeedsKEK() bool
	ExtractDEK(keySector []byte, kek []byte) (dek []byte, err error)
	Decrypt(c cipher.Block, b []byte)
}

var Bridges []Bridge

func IdentifyKeySector(possibleKeySector []byte) Bridge {
	for _, b := range Bridges {
		if b.Is(possibleKeySector) {
			return b
		}
	}
	return nil // not a (known) key sector
}

type IncompleteImplementationError string

func IncompleteImplementation(format string, args ...interface{}) IncompleteImplementationError {
	return IncompleteImplementationError(fmt.Sprintf(format, args...))
}

func (i IncompleteImplementationError) Error() string {
	return string(i)
}
