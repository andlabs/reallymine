// 22 october 2015
package command

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"io/ioutil"
	"reflect"
	"encoding/hex"

	"github.com/andlabs/reallymine/disk"
	"github.com/andlabs/reallymine/kek"
	"github.com/andlabs/reallymine/decryptloop"
)

// DiskSize is passed as the size parameter to disk.Open() when an
// argument of type ArgDisk is processed.
var DiskSize int64 = -1

type argout struct {
	obj			reflect.Value
	deferfunc		func()
}

// This complicated structure allows us to define a fixed set of Arg objects and disallow nil at the same time, reducing the number of things that need validation.

type argiface interface {
	name() string
	desc() string
	argtype() reflect.Type
	prepare(arg string) (out *argout, err error)
}

type arg struct {
	a	argiface
}

type Arg arg

// for usage information
var validArgs []Arg

// TODO complain if any arguments have duplicated names
func addarg(a argiface) Arg {
	aa := Arg{a}
	validArgs = append(validArgs, aa)
	return aa
}

var (
	typeDisk = reflect.TypeOf((*disk.Disk)(nil))
	typeWriter = reflect.TypeOf((*io.Writer)(nil)).Elem()
	typeFile = reflect.TypeOf((*os.File)(nil))
	typeAsker = reflect.TypeOf((*kek.Asker)(nil))
	typeByteSlice = reflect.TypeOf([]byte(nil))
	typeReader = reflect.TypeOf((*io.Reader)(nil)).Elem()
	typeStepList = reflect.TypeOf(decryptloop.StepList(nil))
)

type argDiskType struct{}

func (argDiskType) name() string {
	return "disk"
}

func (argDiskType) desc() string {
	return fmt.Sprintf("The filename of a disk device or disk image. The file must exist and must have a size which is a multiple of the sector size (%d bytes).", disk.SectorSize)
}

func (argDiskType) argtype() reflect.Type {
	return typeDisk
}

func (argDiskType) prepare(arg string) (out *argout, err error) {
	d, err := disk.Open(arg, DiskSize)
	if err != nil {
		return nil, err
	}
	out = new(argout)
	out.obj = reflect.ValueOf(d)
	out.deferfunc = func() {
		d.Close()
	}
	return out, nil
}

var argDisk = addarg(&argDiskType{})
var ArgDisk Arg = argDisk

type argOutFileType struct{}

func (argOutFileType) name() string {
	return "outfile"
}

func (argOutFileType) desc() string {
	return "Either the name of a file to dump the raw data to or - to perform a hexdump on stdout."
}

func (argOutFileType) argtype() reflect.Type {
	return typeWriter
}

func (argOutFileType) prepare(arg string) (out *argout, err error) {
	var of io.WriteCloser

	if arg == "-" {
		of = hex.Dumper(os.Stdout)
	} else {
		f, err := os.Create(arg)
		if err != nil {
			return nil, err
		}
		of = f
	}
	out = new(argout)
	out.obj = reflect.ValueOf(of)
	out.deferfunc = func() {
		// TODO catch the error in the case of stdout?
		// TODO we need to worry about multiplexing then
		of.Close()
	}
	return out, nil
}

var argOutFile = addarg(&argOutFileType{})
var ArgOutFile Arg = argOutFile

type argOutImageType struct{}

func (argOutImageType) name() string {
	return "outimage"
}

func (argOutImageType) desc() string {
	return "The name of a file to write the output disk image to. This file must not exist already; reallymine will not overwrite an existing file or drive."
}

// TODO use typeFile?
func (argOutImageType) argtype() reflect.Type {
	return typeWriter
}

func (argOutImageType) prepare(arg string) (out *argout, err error) {
	f, err := os.OpenFile(arg, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}
	// Use 10MB write buffer (should be configurable)
	bufw := bufio.NewWriterSize(f,10485760)
	out = new(argout)
	out.obj = reflect.ValueOf(bufw)
	out.deferfunc = func() {
		bufw.Flush()
		f.Close()
	}
	return out, nil
}

var argOutImage = addarg(&argOutImageType{})
var ArgOutImage Arg = argOutImage

type argKEKType struct{}

func (argKEKType) name() string {
	return "kek"
}

func (argKEKType) desc() string {
	return kek.AskerDescription
}

func (argKEKType) argtype() reflect.Type {
	return typeAsker
}

func (argKEKType) prepare(arg string) (out *argout, err error) {
	asker, err := kek.NewAsker(arg)
	if err != nil {
		return nil, err
	}
	out = new(argout)
	out.obj = reflect.ValueOf(asker)
	out.deferfunc = func() {}
	return out, nil
}

var argKEK = addarg(&argKEKType{})
var ArgKEK Arg = argKEK

type argDEKType struct{}

func (argDEKType) name() string {
	return "dek"
}

func (argDEKType) desc() string {
	return "A hexadecimal string to use as the DEK."
}

func (argDEKType) argtype() reflect.Type {
	return typeByteSlice
}

func (argDEKType) prepare(arg string) (out *argout, err error) {
	b, err := hex.DecodeString(arg)
	if err != nil {
		return nil, err
	}
	out = new(argout)
	out.obj = reflect.ValueOf(b)
	out.deferfunc = func() {}
	return out, nil
}

var argDEK = addarg(&argDEKType{})
var ArgDEK Arg = argDEK

type argInFileType struct{}

func (argInFileType) name() string {
	return "infile"
}

func (argInFileType) desc() string {
	return "Either the name of a file to read from or - to read from stdin."
}

func (argInFileType) argtype() reflect.Type {
	return typeReader
}

func (argInFileType) prepare(arg string) (out *argout, err error) {
	var inf io.ReadCloser

	if arg == "-" {
		// don't /actually/ close os.Stdin
		inf = ioutil.NopCloser(os.Stdin)
	} else {
		f, err := os.Open(arg)
		if err != nil {
			return nil, err
		}
		inf = f
	}
	out = new(argout)
	out.obj = reflect.ValueOf(inf)
	out.deferfunc = func() {
		inf.Close()
	}
	return out, nil
}

var argInFile = addarg(&argInFileType{})
var ArgInFile Arg = argInFile

type argDecryptionStepsType struct{}

func (argDecryptionStepsType) name() string {
	return "decryption-steps"
}

func (argDecryptionStepsType) desc() string {
	return "A space-delimited list of decryption steps. Must not be empty. Because this is space-delimited, wrap this argument in quotes to have your shell treat the list as one argument."
}

func (argDecryptionStepsType) argtype() reflect.Type {
	return typeStepList
}

func (argDecryptionStepsType) prepare(arg string) (out *argout, err error) {
	steps, err := decryptloop.StepListFromString(arg)
	if err != nil {
		return nil, err
	}
	out = new(argout)
	out.obj = reflect.ValueOf(steps)
	out.deferfunc = func() {}
	return out, nil
}

var argDecryptionSteps = addarg(&argDecryptionStepsType{})
var ArgDecryptionSteps Arg = argDecryptionSteps

// for command.go

func (a Arg) argtype() reflect.Type {
	return a.a.argtype()
}

// TODO rename argout and fields to something more sane for command.go
func (a Arg) prepare(arg string) (*argout, error) {
	return a.a.prepare(arg)
}

func arglist(args []Arg) string {
	list := ""
	for _, a := range args {
		list += " " + a.a.name()
	}
	return list
}

func formatDescription(desc string, args []Arg) string {
	ai := make([]interface{}, len(args))
	for i, a := range args {
		ai[i] = a.a.name()
	}
	return usageL2(desc, ai...)
}

// for reallymine to use directly

func ArgUsage() string {
	out := ""
	for _, a := range validArgs {
		out += usageL1(a.a.name())
		out += usageL2("%s", a.a.desc())
	}
	return out
}
