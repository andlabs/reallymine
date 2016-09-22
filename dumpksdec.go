// 22 october 2015
package main

import (
	"fmt"
	"encoding/hex"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/disk"
)

func cDumpKSDec(d *disk.Disk) error {
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
	Args:		[]string{"disk"},
	Description:	"Identifies, decrypts, and hexdumps the key sector in disk.",
	Do:			cDumpKSDec,
}
