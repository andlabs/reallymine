// 22 october 2015
package main

import (
	"fmt"
	"io"
	"encoding/hex"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/disk"
)

func cDumpKSDec(d *disk.Disk, out io.Writer) error {
	// TODO add -fakesize option of sorts
	last := d.Size()
	fks, err := findKeySector(d, last)
	if err != nil {
		return err
	}

	raw, dek, err := tryGetKeySectorAndDEK(fks.bridge, fks.sector)
	if err != nil {
		return err
	}
	if raw == nil {		// cancelled
		// TODO exit with error?
		return nil
	}

	fks.sector = raw
	fmt.Print(fks.dump())
	fmt.Printf("DEK: %s\n", hex.EncodeToString(dek))
	return nil
}

var dumpksdec = &command.Command{
	Name:		"dumpksdec",
	Args:		[]command.Arg{command.ArgDisk, command.ArgOutFile},
	Description:	"Identifies, decrypts, and dumps the key sector on %s to %s.",
	Do:			cDumpKSDec,
}
