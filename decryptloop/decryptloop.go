// 25 september 2016
package decryptloop

import (
	"crypto/cipher"
	"io"
)

type DecryptLoop struct {
	steps StepList
	c     cipher.Block
	buf   []byte
	pos   int
	out   io.Writer
}

func New(steps StepList, c cipher.Block, out io.Writer) *DecryptLoop {
	return &DecryptLoop{
		steps: steps,
		c:     c,
		buf:   make([]byte, c.BlockSize()),
		out:   out,
	}
}

// report only how much of Write()'s b we consumed, not how much was written to the underlying writer
func (dl *DecryptLoop) writeBlock() (err error) {
	dl.pos = 0
	dl.steps.runBlock(dl.c, dl.buf)
	_, err = dl.out.Write(dl.buf)
	return err
}

func (dl *DecryptLoop) writeIter(b []byte) (n int, err error) {
	fullAvailable := len(b) >= len(dl.buf[dl.pos:])
	n = copy(dl.buf[dl.pos:], b)
	if fullAvailable {
		dl.pos += n
	}
	if dl.pos != len(dl.buf) {
		return n, nil
	}
	return n, dl.writeBlock()
}

// TODO write a test suite to ensure this is working properly
func (dl *DecryptLoop) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		n2, err := dl.writeIter(b)
		n += n2
		if err != nil {
			return n, err
		}
		b = b[n2:]
	}
	return n, nil
}

func (dl *DecryptLoop) StillPendingData() bool {
	return dl.pos != 0
}
