// 22 october 2015
package main

import (
	"fmt"
	"io"
	"bytes"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/disk"
)

var zeroSector [disk.SectorSize]byte

func cDumpLast(d *disk.Disk, out io.Writer) error {
	var sector []byte

	pos := d.Size()
	iter, err := d.ReverseIter(pos)
	if err != nil {
		return err
	}

	found := false
	for iter.Next() {
		sector = iter.Sectors()
		pos = iter.Pos()
		if !bytes.Equal(sector, zeroSector[:]) {
			found = true
			break
		}
	}
	if err = iter.Err(); err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("non-empty sector not found")
	}

	fmt.Printf("sector at 0x%X\n", pos)
	_, err = out.Write(sector)
	return err
}

var dumplast = &command.Command{
	Name:		"dumplast",
	Args:		[]command.Arg{command.ArgDisk, command.ArgOutFile},
	Description:	"Dumps the last non-zero sector on %s to %s.",
	Do:			cDumpLast,
}

const firstSectorsCount = 64

func cDumpFirst(d *disk.Disk, out io.Writer) error {
	iter, err := d.Iter(0, firstSectorsCount)
	if err != nil {
		return err
	}
	if !iter.Next() {
		// if iter.Err() == nil, d.Size() == 0, so just treat it as success
		return iter.Err()
	}
	_, err = out.Write(iter.Sectors())
	return err
}

var dumpfirst = &command.Command{
	Name:		"dumpfirst",
	Args:		[]command.Arg{command.ArgDisk, command.ArgOutFile},
	Description:	"Dumps the first " +
		fmt.Sprintf("%d sectors (%d bytes)", firstSectorsCount, firstSectorsCount * disk.SectorSize) +
		" on %s to %s without decrypting.",
	Do:			cDumpFirst,
}
