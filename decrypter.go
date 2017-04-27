// 24 september 2016
package main

import (
	"fmt"
	"io"
	"crypto/cipher"
	"crypto/aes"
	"sync"
	"runtime"

	"github.com/andlabs/reallymine/disk"
	"github.com/andlabs/reallymine/bridge"
	"github.com/andlabs/reallymine/kek"
	"github.com/andlabs/reallymine/decryptloop"
)

// TODO rename this type
type Decrypter struct {
	Disk		*disk.Disk
	Out		io.WriteSeeker

	EncryptedKeySector		[]byte
	KeySectorPos			int64
	Bridge				bridge.Bridge

	KEK			[]byte
	KeySector		bridge.KeySector
	DEK			[]byte

	blockPool		*sync.Pool
	blocksIn		chan blockOp
	blocksOut		chan blockOp
	blockCipher	cipher.Block
	blockSteps	decryptloop.StepList
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

// TODO refine to allow custom buffer sizes
const NumSectorsAtATime = 102400

// To allow the disk iterator to continue reading while we're
// decrypting, we use a sync.Pool of sector buffers that are
// decrypted. This function creates each new buffer. Note that
// the pool stores pointers to slices. This means we don't
// alter len(*b); instead, refer to the size field of blockOp below.
func newPoolBuffer() interface{} {
	b := make([]byte, disk.SectorSize * NumSectorsAtATime)
	return &b
}

// This structure communicates blocks to decrypt and to write.
type blockOp struct {
	block	*[]byte
	pos		int64
	size		int
}

// This goroutine reads blocks from the disk and sends them across
// d.blocksIn to be read by the decrypters. When finished, it closes
// d.blocksIn and sends the error, if any, to errChan.
func (d *Decrypter) decryptInLoop(errChan chan<- error) {
	iter, err := d.Disk.Iter(0, NumSectorsAtATime)
	if err != nil {
		close(d.blocksIn)
		errChan <- err
		return
	}
	for iter.Next() {
		s := iter.Sectors()
		b := d.blockPool.Get().(*[]byte)
		copy(*b, s)
		d.blocksIn <- blockOp{
			block:	b,
			pos:		iter.Pos(),
			size:		len(s),
		}
	}
	close(d.blocksIn)
	errChan <- iter.Err()
}

// This goroutine reads decrypted blocks and writes them out.
// It stops when d.blocksOut is closed, which is done by
// d.DecryptDisk() itself below.
// TODO report progress in this goroutine
func (d *Decrypter) decryptOutLoop(errChan chan<- error) {
	for bo := range d.blocksOut {
		_, err := d.Out.Seek(bo.pos, io.SeekStart)
		if err != nil {
			errChan <- err
			return
		}
		b := *(bo.block)
		b = b[:bo.size]
		_, err = d.Out.Write(b)
		if err != nil {
			errChan <- err
			return
		}
		d.blockPool.Put(bo.block)
	}
	errChan <- nil
}

// This goroutine does the actual decryption. It runs until d.blocksIn
// is closed, at which point it calls wg.Done(); once the wait group is
// completely finished, d.DecryptDisk() below will close d.blocksOut.
// This allows multiple copies of this goroutine to run simultaneously.
func (d *Decrypter) decryptDecryptLoop(wg *sync.WaitGroup) {
	for bo := range d.blocksIn {
		b := *(bo.block)
		b = b[:bo.size]
		d.blockSteps.DecryptBlock(d.blockCipher, b)
		d.blocksOut <- bo
	}
	wg.Done()
}

func (d *Decrypter) DecryptDisk() error {
	cipher, err := aes.NewCipher(d.DEK)
	if err != nil {
		return err
	}
	d.blockCipher = cipher
	d.blockSteps = d.Bridge.DecryptLoopSteps()
	// TODO remove when PLX is finished
	if len(d.blockSteps) == 0 {
		return fmt.Errorf("** The %s bridge's decryption scheme is not yet known. Please contact andlabs to help contribute it to reallymine.", d.Bridge.Name())
	}

	d.blockPool = &sync.Pool{
		New:		newPoolBuffer,
	}
	d.blocksIn = make(chan blockOp)
	d.blocksOut = make(chan blockOp)

	blocksInErr := make(chan error)
	go d.decryptInLoop(blocksInErr)

	blocksOutErr := make(chan error)
	go d.decryptOutLoop(blocksOutErr)

	n := runtime.NumCPU()
	// TODO n - 2 to account for the above goroutines?
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go d.decryptDecryptLoop(wg)
	}
	go func() {
		wg.Wait()
		close(d.blocksOut)
	}()

	done := 0
	for done != 2 {
		select {
		case err := <-blocksInErr:
			if err != nil {
				return err
			}
			done++
		case err := <-blocksOutErr:
			if err != nil {
				return err
			}
			done++
		}
	}
	return nil
}
