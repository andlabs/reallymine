// 22 october 2015
package main

import (
	"fmt"
	"io"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/disk"
)

func cDumpKSEnc(d *disk.Disk, out io.Writer) error {
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
	Args:		[]command.Arg{command.ArgDisk, command.ArgOutFile},
	Description:	[]string{"Identifies and dumps the key sector in %s to %s without decrypting it."},
	Do:			cDumpKSEnc,
}
