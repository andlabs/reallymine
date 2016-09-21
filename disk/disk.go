// 22 october 2015
package disk

import (
	"fmt"
	"os"
	"io"
)

const SectorSize = 512

// Disk is currently not safe for concurrent use.
type Disk struct {
	f	*os.File
	size	int64
}

func Open(filename string) (d *Disk, err error) {
	d = new(Disk)
	d.f, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	d.size, err = d.f.Seek(0, io.SeekEnd)
	if err != nil {
		d.f.Close()
		return nil, err
	}
	if d.size % SectorSize != 0 {
		d.f.Close()
		return nil, fmt.Errorf("disk size is not a multiple of the sector size; this is likely not a disk")
	}
	return d, nil
}

func (d *Disk) Close() error {
	return d.f.Close()
}

func (d *Disk) Size() int64 {
	return d.size
}

/* TODO
func TryGetDecrypter(keySector []byte, bridge Bridge, askPassword func(firstTime bool) (password string, cancelled bool)) (c cipher.Block) {
	try := func(keySector []byte, bridge Bridge, kek []byte) cipher.Block {
		return bridge.CreateDecrypter(keySector, kek)
	}

	if !bridge.NeedsKEK() {
		return try(keySector, bridge, nil) // should not return nil
	}

	c = try(keySector, bridge, DefaultKEK)
	firstTime := true
	for c == nil { // whlie the default KEK didn't work or the user password is wrong
		password, cancelled := askPassword(firstTime)
		if cancelled { // user aborted
			return nil
		}
		kek := KEKFromPassword(password)
		c = try(keySector, bridge, kek)
		firstTime = false // in case the password was wrong
	}
	return c
}
*/

func (d *Disk) ReadSectorsAt(sectors []byte, pos int64) (int64, error) {
	if len(sectors) % SectorSize != 0 {
		return 0, io.ErrShortBuffer		// TODO better error?
	}
	n, err := d.f.ReadAt(sectors, pos)
	if err == io.EOF {
		if n == 0 {		// this is truly the end of the disk
			return 0, io.EOF
		}
		if n % SectorSize != 0 {
			return n, io.ErrUnexpectedEOF
		}
		// Allow a short read at the end of the disk; the next call
		// to ReadSectorsAt() will return EOF with nothing read.
		// (Of course, if n == len(sectors), it isn't a *short* read,
		// but the point still applies.)
		return n, nil
	}
	// This handles the remaining cases.
	// If n != len(sectors) then err will not be nil by the
	// requirements of io.ReaderAt (and also not be io.EOF
	// due to the above code) so we're good on that part.
	// If n == len(sectors) then just pass err unchnaged; it
	// might be nil, so the success case is handled too.
	return n, err
}

type SectorIter struct {
	d		*Disk
	sectors	[]byte
	pos		int64
	incr		int		// in units of len(sectors)
	eof		bool
	err		error
}

func (d *Disk) mkiter(startAt int64, countPer int, reverse bool) (*SectorIter, error) {
	if startAt % SectorSize != 0 {
		return nil, fmt.Errorf("startAt must be sector-aligned")
	}
	s := new(SectorIter)
	s.d = d
	span := s.d.size - startAt
	if reverse {
		span = startAt
	}
	s.sectors = make([]byte, countPer * SectorSize)
	if reverse {
		// The first call to Next() will push s.pos to the last block.
		s.pos = startAt
		s.incr = -1
	} else {
		// The first call to Next() will increment s.pos to startAt.
		s.pos = startAt - len(sectors)
		s.incr = 1
	}
	return s, nil
}

func (d *Disk) Iter(startAt int64, countPer int) *SectorIter {
	return d.mkiter(startAt, countPer, false)
}

// TODO allow different sized iter blocks if we ever split this into its own package; this requires handling the last short read in Next()
func (d *Disk) ReverseIter(startAt int64) (*SectorIter, error) {
	return d.mkiter(startAt, 1, true)
}

func (s *SectorIter) Next() bool {
	if s.eof {
		return false
	}
	s.pos += s.incr * len(s.sectors)
	n, err := s.d.ReadSectorsAt(s.sectors, s.pos)
	if err == io.EOF {
		s.eof = true
		return false
	}
	if err != nil {
		s.err = err
		return false
	}
	s.sectors = s.sectors[:n]		// trim short last read
	return true
}

func (s *SectorIter) Sectors() []byte {
	return s.sectors
}

func (s *SectorIter) Pos() int64 {
	return s.pos
}

func (s *SectorIter) Err() error {
	return s.err
}
