// 24 september 2016

// +build TODO

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
	Disk		*disk.Disk
	Out		io.Writer

	EncryptedKeySector		[]byte
	KeySectorPos			int64
	Bridge				bridge.Bridge

	KeySector		bridge.KeySector
	DEK			[]byte
}

func (d *Decrypter) FindKeySector() error {
	// TODO allow a way to hook in every so often if the search takes too long
	iter, err := d.Disk.ReverseIter(startAt)
	if err != nil {
		return err
	}
	for iter.Next() {
		d.EncryptedKeySector = iter.Sectors()
		d.Pos = iter.Pos()
		d.Bridge = bridge.IdentifyKeySector(fks.sector)
		if d.Bridge != nil {
			break
		}
	}
	if err = iter.Err(); err != nil {
		return err
	}
	if d.Bridge == nil {
		return fmt.Errorf("key sector not found")
	}
	return nil
}

{
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
