// 22 october 2015
package main

import (
	"fmt"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/disk"
)

func cDumpKSEnc(d *disk.Disk) error {
	// TODO add -fakesize option of sorts
	last := d.Size()
	fks, err := findKeySector(d, last)
	if err != nil {
		return err
	}

	fmt.Print(fks.dump())
	return nil
}

var dumpksenc = &command.Command{
	Name:		"dumpksenc",
	Args:		[]string{"disk"},
	Description:	"Identifies and hexdumps the key sector in disk without decrypting it.",
	Do:			cDumpKSEnc,
}
