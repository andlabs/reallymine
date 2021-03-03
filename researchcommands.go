// 22 october 2015
package main

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"io"

	"github.com/undeadbanegithub/reallymine/bridge"
	"github.com/undeadbanegithub/reallymine/command"
	"github.com/undeadbanegithub/reallymine/decryptloop"
	"github.com/undeadbanegithub/reallymine/disk"
	"github.com/undeadbanegithub/reallymine/kek"
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
	Name:        "dumplast",
	Args:        []command.Arg{command.ArgDisk, command.ArgOutFile},
	Description: "Dumps the last non-zero sector on %s to %s.",
	Do:          cDumpLast,
}

func cDecryptKeySector(d *disk.Disk, out io.Writer, a *kek.Asker) error {
	dec := &Decrypter{
		Disk: d,
		Out:  out,
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
	Name:        "dumpkeysector",
	Args:        []command.Arg{command.ArgDisk, command.ArgOutFile},
	Description: "Identifies and dumps the key sector on %s to %s.",
	Do:          cDumpKeySector,
}

var decryptkeysector = &command.Command{
	Name:        "decryptkeysector",
	Args:        []command.Arg{command.ArgDisk, command.ArgOutFile, command.ArgKEK},
	Description: "Identifies, decrypts, and dumps the key sector on %s to %s using %s.",
	Do:          cDecryptKeySector,
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
	Name: "dumpfirst",
	Args: []command.Arg{command.ArgDisk, command.ArgOutFile},
	Description: "Dumps the first " +
		fmt.Sprintf("%d sectors (%d bytes)", firstSectorsCount, firstSectorsCount*disk.SectorSize) +
		" on %s to %s without decrypting.",
	Do: cDumpFirst,
}

func cDecryptFile(in io.Reader, out io.Writer, dek []byte, steps decryptloop.StepList) error {
	cipher, err := aes.NewCipher(dek)
	if err != nil {
		return err
	}
	dl := decryptloop.New(steps, cipher, out)
	_, err = io.Copy(dl, in)
	if err != nil {
		return err
	}
	if dl.StillPendingData() {
		return fmt.Errorf("input file does not end with a complete block")
	}
	return nil
}

var decryptfile = &command.Command{
	Name:        "decryptfile",
	Args:        []command.Arg{command.ArgInFile, command.ArgOutFile, command.ArgDEK, command.ArgDecryptionSteps},
	Description: "Decrypts %s to %s using the provided %s and %s.",
	Do:          cDecryptFile,
}
