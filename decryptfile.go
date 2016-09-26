// 25 september 2016
package main

import (
	"fmt"
	"io"
	"crypto/aes"

	"github.com/andlabs/reallymine/command"
	"github.com/andlabs/reallymine/decryptloop"
)

// TODO merge the research commands into a single file

func cDecryptFile(in io.Reader, out io.Writer, dek []byte, steps decryptloop.StepList) error {
	cipher, err := aes.NewCipher(dek)
	if err != nil {
		return err
	}
	dl := decryptloop.New(steps, cipher, out)
	_, err = io.Copy(dl, in)
	if err != nil {
		return err
	}
	if dl.StillPendingData() {
		return fmt.Errorf("input file does not end with a complete block")
	}
	return nil
}

var decryptfile = &command.Command{
	Name:		"decryptfile",
	Args:		[]command.Arg{command.ArgInFile, command.ArgOutFile, command.ArgDEK, command.ArgDecryptionSteps},
	Description:	"Decrypts %s to %s using the provided %s and %s.",
	Do:			cDecryptFile,
}
