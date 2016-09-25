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

	KEK			[]byte
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
	if err := iter.Err(); err != nil {
		return err
	}
	if d.Bridge == nil {
		return fmt.Errorf("key sector not found")
	}
	return nil
}

func (d *Decrypter) DecryptKeySector(a *kek.Asker) error {
	var ks bridge.KeySector
	var dek []byte

	if !b.NeedsKEK() {
		d.KEK, err = ks.DEK()
		if err != nil {
			return err
		}
	} else {
		wrong := false
		for a.Ask() {
			d.KEK = a.KEK()
			d.KeySector, err = b.DecryptKeySector(d.EncryptedKeySector, d.KEK)
			if err != nil {
				return err
			}
			d.DEK, err = ks.DEK()
			if err == bridge.ErrWrongKEK {
				wrong = true
				continue
			}
			if err != nil {
				return err
			}
			wrong = false
			break
		}
		if err := a.Err(); err != nil {
			return err
		}
		// preserve bridge.ErrWrongKEK if we asked to use a specific KEK or used -askonce
		if wrong {
			return bridge.ErrWrongKEK
		}
	}

	return nil
}

func (d *Decrypter) DecryptDisk() error {
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
