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
}

func OpenDisk(filename string) (d *Disk, err error) {
	d = new(Disk)
	d.f, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Disk) Close() error {
	return d.f.Close()
}

func (d *Disk) Size() (int64, error) {
	return d.f.Seek(0, io.SeekEnd)
}

// TODO write a function to make this stop early, giving the user the option to continue

type SearchFunc func(sector []byte) (found bool)

func (d *Disk) ReverseSearch(startAt int64, f SearchFunc) (sector []byte, pos int64, err error) {
	sector = make([]byte, SectorSize)
	pos = startAt - SectorSize
	for pos >= 0 {
		// TODO is this correct usage of ReadAt()?
		_, err := d.f.ReadAt(sector, pos)
		// io.ReaderAt specifies that EOF may be returned when reading right at the end of the file
		if err != nil && err != io.EOF {
			return nil, 0, err
		}
		found := f(sector)
		if found {
			return sector, pos, nil
		}
		pos -= SectorSize
	}
	return nil, 0, nil
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

var ErrCancelled = fmt.Errorf("cancelled")

type ForEachFunc func(sectors []byte) (cancel bool)

func (d *Disk) ForEach(nSectorsPer int, f ForEachFunc) error {
	_, err := d.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	sectors := make([]byte, nSectorsPer * SectorSize)
	for {
		n, err := d.f.Read(sectors)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// handle the final block properly if it's shorter
		sectors = sectors[:n]
		cancel := f(sectors)
		if cancel {
			return ErrCancelled
		}
	}
	return nil
}

func (d *Disk) ReadSectorAt(pos int64) ([]byte, error) {
	sector := make([]byte, SectorSize)
	// TODO see if we can just use d.f.ReadAt()
	_, err := d.f.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}
	_, err = io.ReadFull(d.f, sector)
	if err != nil {
		return nil, err
	}
	return sector, err
}

type SectorIter struct {
	d		*Disk
	sectors	[]byte
	reverse	bool
	pos		int64
	err		error
}

func (d *Disk) Iter(startAt int64, countPer int) *SectorIter {
	return &SectorIter{
		d:		d,
		sectors:	make([]byte, countPer * SectorSize),
		pos:		startAt,
	}
}

func (d *Disk) ReverseIter(startAt int64, countPer int) *SectorIter {
	s := d.Iter(startAt, countPer)
	s.reverse = true
	return s
}

func (s *SectorIter) Next() bool {
	// TODO
}

func (s *SectorIter) Sectors() []byte {
	return s.sectors
}

func (s *SectorIter) Err() error {
	return s.err
}
