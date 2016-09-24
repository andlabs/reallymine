// 24 september 2016
package main

import (
	"fmt"
	"os"
	"io"
	"encoding/hex"

	"github.com/andlabs/reallymine/disk"
	"github.com/andlabs/reallymine/bridge"
	"github.com/andlabs/reallymine/kek"
	"github.com/hashicorp/vault/helper/password"
)

type Decrypter struct {
	d	*disk.Disk
	out	io.Writer
}

func askForPassword() (pw string, err error) {
	fmt.Print("Enter WD password: ")
	return password.Read(os.Stdin)
}

func (d *Decrypter) Decrypt() error {
	var sector []byte
	var pos int64
	var bridge bridge.Bridge

	// TODO allow a way to hook in every so often if the search takes too long
	iter, err := d.ReverseIter(startAt)
	if err != nil {
		return nil, err
	}
	for iter.Next() {
		sector = iter.Sectors()
		pos = iter.Pos()
		bridge = bridge.IdentifyKeySector(fks.sector)
		if bridge != nil {
			break
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if bridge == nil {
		return nil, fmt.Errorf("key sector not found")
	}

	var ks bridge.KeySector
	var curkek []byte
	var dek []byte

	try := func() {
		ks, err = b.DecryptKeySector(sector, curkek)
		if err == nil {
			dek, err = ks.DEK()
		}
	}

	if !b.NeedsKEK() {
		dek, err = ks.DEK()
		if err != nil {
			return nil, err
		}
	} else {
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
			if err != nil {		// includes cancelled
				return nil, err
			}
			curkek = kek.FromPassword(pw)
			try()
		}
		if err != nil {
			return nil, err
		}
	}

	// TODO refine or allow custom buffer sizes?
	iter, err = d.Iter(0, 1)
	if err != nil {
		return err
	}
	for iter.Next() {
		s := iter.Sectors()
		bridge.Decrypt(s)
		_, err = out.Write(s)
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
