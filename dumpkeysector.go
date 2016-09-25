// 22 october 2015
package main

import (
	"fmt"
	"io"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/disk"
	"github.com/andlabs/reallymine/bridge"
	"github.com/andlabs/reallymine/kek"
)

func cDecryptKeySector(d *disk.Disk, out io.Writer, a *kek.Asker) error {
	dec := &Decrypter{
		Disk:		d,
		Out:		out,
	}
	err := dec.FindKeySector()
	if err != nil {
		return err
	}

	sector := dec.EncryptedKeySector
	if a != nil {
		err = dec.ExtractDEK(a)
		// ignore bridge.ErrWrongKEK; this command should still produce a dump even in the face of the wrong KEK (with -askonce, -default, or a specific KEK)
		if err != nil && err != bridge.ErrWrongKEK {
			return err
		}
		sector = dec.KeySector.Raw()
	}

	fmt.Printf("%s\n%s\n",
		formatSectorPos(dec.KeySectorPos),
		formatBridge(dec.Bridge))
	_, err = out.Write(sector)
	return err
}

func cDumpKeySector(d *disk.Disk, out io.Writer) error {
	return cDecryptKeySector(d, out, nil)
}

var dumpkeysector = &command.Command{
	Name:		"dumpkeysector",
	Args:		[]command.Arg{command.ArgDisk, command.ArgOutFile},
	Description:	"Identifies and dumps the key sector on %s to %s.",
	Do:			cDumpKeySector,
}

var decryptkeysector = &command.Command{
	Name:		"decryptkeysector",
	Args:		[]command.Arg{command.ArgDisk, command.ArgOutFile, command.ArgKEK},
	Description:	"Identifies, decrypts, and dumps the key sector on %s to %s using %s.",
	Do:			cDecryptKeySector,
}
