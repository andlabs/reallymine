// 22 october 2015
package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"encoding/hex"
	"bytes"
)

func errf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func die(format string, args ...interface{}) {
	errf(format, args...)
	errf("\n")
	os.Exit(1)
}

type Command struct {
	Name		string
	Args			[]string
	Description	string
	Do			func([]string) error
}

var zeroSector [disk.SectorSize]byte

func cDumpLast(args []string) error {
	var sector []byte

	d, err := disk.Open(args[0])
	if err != nil {
		return err
	}
	defer d.Close()

	// TODO add -fakesize option of sorts
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
	if !found {		// not found
		return fmt.Errorf("non-empty sector not found")
	}

	fmt.Printf("sector starting at %d\n", pos)
	fmt.Printf("%s", hex.Dump(sector))
	return nil
}

var dumplast = &Command{
	Name:		"dumplast",
	Args:		[]string{"file"},
	Description:	"Hexdumps the last non-zero sector in file.",
	Do:			cDumpLast,
}

type foundKeySector struct {
	sector	[]byte
	pos		int64
	bridge	bridge.Bridge
}

func findKeySector(d *Disk, startAt int64) (fks *foundKeySector, err error) {
	fks = new(foundKeySector)

	iter, err := disk.ReverseIter(startAt)
	if err != nil {
		return nil, err
	}

	for iter.Next() {
		fks.sector = iter.Sectors()
		fks.pos = iter.Pos()
		fks.bridge = bridge.IdentifyKeySector(sector)
		if fks.bridge != nil {
			break
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if fks.bridge == nil {		// not found
		return nil, nil
	}

	return fks, nil
}

func cDumpKSRaw(args []string) error {
	d, err := OpenDisk(args[0])
	if err != nil {
		return err
	}
	defer d.Close()

	// TODO add -fakesize option of sorts
	last := d.Size()

	fks, err := findKeySector(d, last)
	if err != nil {
		return err
	}
	if fks == nil {
		return fmt.Errorf("key sector not found")
	}

	fmt.Printf("sector starting at %d\n", fks.pos)
	fmt.Printf("%s", hex.Dump(fks.sector))
	fmt.Printf("bridge type %s\n", fks.bridge.Name())
	return nil
}

var dumpksraw = &Command{
	Name:		"dumpksraw",
	Args:		[]string{"file"},
	Description:	"Identifies and hexdumps the key sector in file",
	Do:			cDumpKSRaw,
}

var Commands = []*Command{
	dumplast,
	dumpksraw,
}

func usage() {
/*
	errf("usage: %s encrypted decrypted\n", os.Args[0])
	errf("	encrypted must exist; should not be a device\n")
	errf("	decrypted must NOT exist\n")
	os.Exit(1)
*/
	errf("usage: %s [options] command [args...]\n", os.Args[0])
	errf("options:\n")
	flag.PrintDefaults()
	errf("commands:\n")
	for _, c := range Commands {
		// See package flag's source for details on this formatting.
		errf("  %s %s\n", c.Name, strings.Join(c.Args, " "))
		errf("    	%s\n", c.Description)
	}
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}
	cmd := flag.Arg(0)

	for _, c := range Commands {
		if cmd != c.Name {
			continue
		}
		args := flag.Args()[1:]
		if len(args) != len(c.Args) {
			errf("error: incorrect number of arguments for command %s\n", c.Name)
			usage()
		}
		err := c.Do(args)
		if err != nil {
			die("error running %s: %v\n", c.Name, err)
		}
		// all good; return successfully
		return
	}

	errf("error: unknown command %q\n", cmd)
	usage()
}

/*

	infname := os.Args[1]
	outfname := os.Args[2]

	in, err := os.Open(infname)
	if err != nil {
		die("error opening encrypted file %s: %v", infname, err)
	}
	defer in.Close()

	// TODO make sure infile is not a device
	// we must outright forbid it because we aren't running sector-to-sector anymore

	insize, err := in.Seek(0, 2)
	if err != nil {
		errf("error finding size of encrypted file %s: %v", infname, err)
	}

	fmt.Printf("Finding key sector...\n")
	keySector, bridge := FindKeySectorAndBridge(in, insize)
	if bridge == nil {
		errf("Sorry, we couldn't find the key sector.\n")
		errf("Either the drive isn't a complete image,\n")
		errf("or the encryption isn't supported yet.\n")
		os.Exit(1)
	}
	fmt.Printf("Found %s.\n", bridge.Name())
	if !bridge.NeedsKEK() {
		fmt.Printf("You will not need to enter your password\n")
		fmt.Printf("for this bridge chip.\n")
	} else {
		fmt.Printf("Trying without a password...\n")
	}

	c := TryGetDecrypter(keySector, bridge, func(firstTime bool) (password string, cancelled bool) {
		if firstTime {
			fmt.Printf("The drive's password is needed to decrypt your drive.\n")
			fmt.Printf("Please enter it now.\n")
		} else {
			fmt.Printf("Password incorrect.\n")
		}
		// TODO
		os.Exit(2)
		panic("unreachable")
	})
	if c == nil {
		fmt.Printf("User aborted operation.\n")
		os.Exit(1)
	}

	// TODO decrypt a few sectors to verify the partition table

	_, err = in.Seek(0, 0)
	if err != nil {
		die("error seeking back to start of decrypted file %s: %v", infname, err)
	}

	out, err := os.OpenFile(outfname, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			errf("Error creating decrypted file %s: %v\n", outfname, err)
			errf("%s will not overwrite a file that already exists.\n", os.Args[0])
			errf("In particular, %s does not allow in-place decryption.\n", os.Args[0])
			os.Exit(1)
		}
		die("error creating decrypted file %s: %v", outfname, err)
	}

	fmt.Printf("Beginning decryption!\n")
	sectors := make([]byte, NumSectorsAtATime * SectorSize)
	n := int64(0)
	inmb := insize / 1024 / 1024
	for DecryptNext(in, out, bridge, c, sectors) {
		n += NumMBAtATime
		fmt.Printf("%d MB / %d MB complete.\n", n, inmb)
	}

	fmt.Printf("Completed successfully!\n")
}
*/
