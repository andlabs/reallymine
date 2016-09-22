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

func askForPassword() (pw string, err error) {
	fmt.Print("Enter WD password: ")
	return password.Read(os.Stdin)
}

func tryGetKeySectorAndDEK(b bridge.Bridge, sector []byte) (raw []byte, dek []byte, err error) {
	var ks bridge.KeySector
	var curkek []byte

	try := func() {
		ks, err = b.DecryptKeySector(sector, curkek)
		if err == nil {
			dek, err = ks.DEK()
		}
	}

	if !b.NeedsKEK() {
		dek, err = ks.DEK()
		if err != nil {
			return nil, nil, err
		}
		return ks.Raw(), dek, nil
	}

	curkek = kek.Default
	first := true
	try()
	for err == bridge.ErrWrongKEK {
		if first {
			fmt.Printf("You need the WD password to decrypt this drive.\n")
			first = false
		} else {
			fmt.Printf("Wrong WD password.\n")
		}
		pw, err := askForPassword()
		if err == password.ErrInterrupted {		// cancelled
			return nil, nil, nil
		}
		if err != nil {
			return nil, nil, err
		}
		curkek = kek.FromPassword(pw)
		try()
	}
	if err != nil {
		return nil, nil, err
	}
	return ks.Raw(), dek, nil
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
