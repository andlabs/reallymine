// 24 september 2016
package kek

import (
	"fmt"
	"os"
	"encoding/hex"

	"github.com/hashicorp/vault/helper/password"
)

// TODO rearrange the parts of this package

type askfn func() (kek []byte, repeat bool, err error)

func askSpecific(b []byte) askfn {
	return func() ([]byte, bool, error) {
		return b, false, nil
	}
}

func askUser(note string, repeat bool) askfn {
	return func() ([]byte, bool, error) {
		kek, err := realAskUser(note)
		return kek, repeat, err
	}
}

type Asker struct {
	fns		[]askfn
	kek		[]byte
	err		error
}

func mkasker(fns ...askfn) (*Asker, error) {
	return &Asker{
		fns:		fns,
	}, nil
}

// TODO see if we can have AskReal, etc. be instances to avoid a useless error check
func NewAsker(cmdstr string) (a *Asker, err error) {
	switch cmdstr {
	case AskReal:
		return mkasker(
			askSpecific(Default),
			askUser(noteNeedsPassword, false),
			askUser(notePasswordWrong, true))
	case AskOnce:
		return mkasker(askUser("", false))
	case AskOnly:
		return mkasker(
			askUser("", false),
			askUser(notePasswordWrong, true))
	case AskDefault:
		return mkasker(askSpecific(Default))
	}
	kek, err := hex.DecodeString(cmdstr)
	if err != nil {
		return nil, err
	}
	return mkasker(askSpecific(kek))
}

const (
	AskReal = "-real"
	AskOnce = "-askonce"
	AskOnly = "-onlyask"
	AskDefault = "-default"
)

const (
	noteNeedsPassword = "You need the WD password to decrypt this drive."
	notePasswordWrong = "Wrong WD password."
)

// TODO how to get the secure insert icon in OS X Terminal?
func realAskUser(note string) (kek []byte, err error) {
	if note != "" {
		fmt.Printf("%s\n", note)
	}
	fmt.Print("Enter WD password: ")
	pw, err := password.Read(os.Stdin)
	fmt.Println()		// because password.Read() doesn't
	if err != nil {		// including cancelled
		return nil, err
	}
	return FromPassword(pw), nil
}

func (a *Asker) Ask() bool {
	if a.err != nil {
		return false
	}
	if len(a.fns) == 0 {
		return false
	}
	kek, repeat, err := a.fns[0]()
	a.kek = kek
	a.err = err
	if a.err != nil {
		return false
	}
	if !repeat {		// no more from this one, advance to the next one
		a.fns = a.fns[1:]
	}
	return true
}

func (a *Asker) KEK() []byte {
	return a.kek
}

func (a *Asker) Err() error {
	return a.err
}

const AskerDescription = "" +
	"This specifies a KEK to decrypt the key sector with. " +
	"This argument can be one of the following strings:\n" +
	"- " + AskReal + " to use the default KEK and then ask for a password until the correct one is used, just like the main decrypt command\n" +
	"- " + AskOnce + " to ask for a password once and only use the resultant KEK\n" +
	"- " + AskOnly + " to only ask for passwords until the correct one is used\n" +
	"- " + AskDefault + " to only use the default KEK\n" +
	"Any other string is taken as a hexadecimal string to use as the KEK."
