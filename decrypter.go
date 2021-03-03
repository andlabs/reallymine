// 24 september 2016
package main

import (
	"crypto/aes"
	"fmt"
	"io"

	"github.com/undeadbanegithub/reallymine/bridge"
	"github.com/undeadbanegithub/reallymine/decryptloop"
	"github.com/undeadbanegithub/reallymine/disk"
	"github.com/undeadbanegithub/reallymine/kek"
)

// TODO rename this type
type Decrypter struct {
	Disk *disk.Disk
	Out  io.Writer

	EncryptedKeySector []byte
	KeySectorPos       int64
	Bridge             bridge.Bridge

	KEK       []byte
	KeySector bridge.KeySector
	DEK       []byte
}

func (d *Decrypter) FindKeySector() error {
	iter, err := d.Disk.ReverseIter(d.Disk.Size())
	if err != nil {
		return err
	}
	for iter.Next() {
		d.EncryptedKeySector = iter.Sectors()
		d.KeySectorPos = iter.Pos()
		d.Bridge = bridge.IdentifyKeySector(d.EncryptedKeySector)
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

func (d *Decrypter) decryptKeySector() (err error) {
	d.KeySector, err = d.Bridge.DecryptKeySector(d.EncryptedKeySector, d.KEK)
	if err != nil {
		return err
	}
	d.DEK, err = d.KeySector.DEK()
	if err != nil {
		return err
	}
	return nil
}

func (d *Decrypter) ExtractDEK(a *kek.Asker) (err error) {
	if !d.Bridge.NeedsKEK() {
		return d.decryptKeySector()
	}

	for a.Ask() {
		d.KEK = a.KEK()
		err = d.decryptKeySector()
		if err == bridge.ErrWrongKEK {
			continue
		}
		if err != nil {
			return err
		}
		break
	}
	// preserve bridge.ErrWrongKEK if we asked to use a specific KEK or used -askonce or -default
	wrong := err == bridge.ErrWrongKEK
	// but return this error first
	if err := a.Err(); err != nil {
		return err
	}
	if wrong {
		return bridge.ErrWrongKEK
	}
	return nil
}

const NumSectorsAtATime = 102400

func (d *Decrypter) DecryptDisk() error {
	cipher, err := aes.NewCipher(d.DEK)
	if err != nil {
		return err
	}
	steps := d.Bridge.DecryptLoopSteps()
	// TODO remove when PLX is finished
	if len(steps) == 0 {
		return fmt.Errorf("** The %s bridge's decryption scheme is not yet known. Please contact andlabs to help contribute it to reallymine.", d.Bridge.Name())
	}
	dl := decryptloop.New(steps, cipher, d.Out)
	// TODO refine or allow custom buffer sizes?
	iter, err := d.Disk.Iter(0, NumSectorsAtATime)
	if err != nil {
		return err
	}
	for iter.Next() {
		// TODO report progress in MB
		s := iter.Sectors()
		_, err = dl.Write(s)
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
