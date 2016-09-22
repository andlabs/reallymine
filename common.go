// 22 october 2015
package main

import (
	"fmt"
	"encoding/hex"

	"github.com/andlabs/reallymine/disk"
	"github.com/andlabs/reallymine/bridge"
)

type foundKeySector struct {
	sector	[]byte
	pos		int64
	bridge	bridge.Bridge
}

// TODO allow a way to hook in every so often if the search takes too long
func findKeySector(d *disk.Disk, startAt int64) (fks *foundKeySector, err error) {
	fks = new(foundKeySector)
	iter, err := d.ReverseIter(startAt)
	if err != nil {
		return nil, err
	}
	for iter.Next() {
		fks.sector = iter.Sectors()
		fks.pos = iter.Pos()
		fks.bridge = bridge.IdentifyKeySector(fks.sector)
		if fks.bridge != nil {
			break
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if fks.bridge == nil {
		return nil, fmt.Errorf("key sector not found")
	}
	return fks, nil
}

func askForPassword() (password string, err error) {
	panic("TODO")
}

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
