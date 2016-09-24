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
	r		*io.SectionReader
	close	func() error
}

func Open(filename string, size int64) (d *Disk, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	errfail := func(err error) (*Disk, error) {
		f.Close()
		return nil, err
	}
	realsize, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return errfail(err)
	}
	if size == -1 {
		size = realsize
	} else if size > realsize {
		return errfail(fmt.Errorf("requested disk size larger than actual disk size"))
	}
	if size % SectorSize != 0 {
		return errfail(fmt.Errorf("disk size is not a multiple of the sector size; this is likely not a disk"))
	}
	return &Disk{
		r:		io.NewSectionReader(f, 0, size),
		close:	f.Close,
	}, nil
}

func (d *Disk) Close() error {
	return d.close()
}

func (d *Disk) Size() int64 {
	return d.r.Size()
}

func (d *Disk) ReadSectorsAt(sectors []byte, pos int64) (int, error) {
	if len(sectors) % SectorSize != 0 {
		return 0, io.ErrShortBuffer		// TODO better error?
	}
	n, err := d.r.ReadAt(sectors, pos)
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
	s.sectors = make([]byte, countPer * SectorSize)
	if reverse {
		// The first call to Next() will push s.pos to the last block.
		s.pos = startAt
		s.incr = -1
	} else {
		// The first call to Next() will increment s.pos to startAt.
		s.pos = startAt - int64(len(s.sectors))
		s.incr = 1
	}
	return s, nil
}

func (d *Disk) Iter(startAt int64, countPer int) (*SectorIter, error) {
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
	s.pos += int64(s.incr * len(s.sectors))
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
