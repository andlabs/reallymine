// 25 september 2016
package main

import (
	"fmt"
	"io"

	"github.com/undeadbanegithub/reallymine/command"
	"github.com/undeadbanegithub/reallymine/disk"
	"github.com/undeadbanegithub/reallymine/kek"
)

func runUntilDEK(d *disk.Disk, out io.Writer) (dec *Decrypter, err error) {
	dec = &Decrypter{
		Disk: d,
		Out:  out,
	}
	err = dec.FindKeySector()
	if err != nil {
		return nil, err
	}
	asker, err := kek.NewAsker(kek.AskReal)
	if err != nil {
		return nil, err
	}
	err = dec.ExtractDEK(asker)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func cGetDEK(d *disk.Disk) error {
	dec, err := runUntilDEK(d, nil)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", formatBridge(dec.Bridge))
	fmt.Printf("DEK: %s\n", formatKey(dec.DEK))
	// TODO make this a format function?
	fmt.Printf("decryption steps: %v\n", dec.Bridge.DecryptLoopSteps())
	return nil
}

var getdek = &command.Command{
	Name:        "getdek",
	Args:        []command.Arg{command.ArgDisk},
	Description: "Gets the DEK and decryption steps to use on %s and prints it on stdout.",
	Do:          cGetDEK,
}

func cDecrypt(d *disk.Disk, out io.Writer) error {
	dec, err := runUntilDEK(d, out)
	if err != nil {
		return err
	}
	return dec.DecryptDisk()
}

var decrypt = &command.Command{
	Name:        "decrypt",
	Args:        []command.Arg{command.ArgDisk, command.ArgOutImage},
	Description: "Decrypts the entire disk %s to %s.",
	Do:          cDecrypt,
}
