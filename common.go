// 22 october 2015
package main

import (
	"fmt"
	"os"
	"encoding/hex"

	"github.com/andlabs/reallymine/disk"
	"github.com/andlabs/reallymine/bridge"
	"github.com/andlabs/reallymine/kek"
	"github.com/hashicorp/vault/helper/password"
)

func tryGetKeySectorAndDEK(b bridge.Bridge, sector []byte) (raw []byte, dek []byte, err error) {}

func dumpSector(sector []byte, pos int64) string {
	s := ""
	if pos >= 0 {
		s = fmt.Sprintf("sector at 0x%X\n", pos)
	}
	s += hex.Dump(sector)
	return s
}

func (fks *foundKeySector) dump() string {
	s := dumpSector(fks.sector, fks.pos)
	s += fmt.Sprintf("bridge type %s\n", fks.bridge.Name())
	return s
}
