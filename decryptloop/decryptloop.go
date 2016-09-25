// 25 september 2016
package decryptloop

import (
	"fmt"
	"io"
	"strings"
	"crypto/cipher"
)

type DecryptLoop struct {
	steps	[]Step
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

type UnknownStepNameError string

func (e UnknownStepNameError) Error() string {
	return fmt.Sprintf("unknown decrypt loop step name %q", string(e))
}

func FromString(s string, c cipher.Block, out io.Writer) (*DecryptLoop, error) {
	names := strings.Split(s, " ")
	steps := make([]Step, len(names))
	for i, name := range names {
		step, ok := stepsByName[name]
		if !ok {
			return nil, UnknownStepNameError(name)
		}
		steps[i] = step
	}
	return New(steps, c, out), nil
}

func (dl *DecryptLoop) String() string {
	names := make([]string, len(dl.steps))
	for i, step := range dl.steps {
		names[i] = step.name()
	}
	return strings.Join(names, " ")
}

func (dl *DecryptLoop) writeBlock() (n int, err error) {
	for _, step := range dl.steps {
		step.do(dl.c, dl.buf)
	}
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
