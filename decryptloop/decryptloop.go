// 25 september 2016
package decryptloop

import (
	"io"
	"crypto/cipher"
)

type DecryptLoop struct {
	steps	StepList
	c		cipher.Block
	buf		[]byte
	out		io.Writer
}

func New(steps []Step, c cipher.Block, out io.Writer) *DecryptLoop {
	return &DecryptLoop{
		steps:	steps,
		c:		c,
		buf:		make([]byte, 0, c.BlockSize()),
		out:		out,
	}
}

func FromString(s string, c cipher.Block, out io.Writer) (*DecryptLoop, error) {
	steps, err := stepListFromString(s)
	if err != nil {
		return nil, err
	}
	return New(steps, c, out), nil
}

func (dl *DecryptLoop) writeBlock() (n int, err error) {
	dl.steps.runBlock(dl.c, dl.buf)
	n, err = dl.out.Write(dl.buf)
	dl.buf = dl.buf[0:0]		// reuse dl.buf
	return n, err
}

func (dl *DecryptLoop) neededForNext() int {
	return dl.c.BlockSize() - len(dl.buf)
}

func (dl *DecryptLoop) Write(b []byte) (n int, err error) {
	for {
		needed := dl.neededForNext()
		if len(b) < needed {
			n += len(b)
			dl.buf = append(dl.buf, b...)
			return n, nil
		}
		n += needed
		dl.buf = append(dl.buf, b[:needed]...)
		n2, err := dl.writeBlock()
		if err != nil {
			// only count what of b was actually written
			n -= n2 - needed
			if n < 0 {
				// none of b was written
				n = 0
			}
			return n, err
		}
		b = b[needed:]
	}
	return n, nil
}

func (dl *DecryptLoop) StillPendingData() bool {
	return len(dl.buf) != 0
}
