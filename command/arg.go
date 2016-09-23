// 22 october 2015
package command

import (
	"os"
	"io"
	"reflect"
	"encoding/hex"

	"github.com/andlabs/reallymine/disk"
)

type argout struct {
	obj			reflect.Value
	deferfunc		func()
}

// This complicated structure allows us to define a fixed set of Arg objects and disallow nil at the same time, reducing the number of things that need validation.

type argiface interface {
	name() string
	desc() []string
	argtype() reflect.Type
	prepare(arg string) (out *argout, err error)
}

type arg struct {
	a	argiface
}

type Arg arg

// for usage information
var validArgs []Arg

func addarg(a argiface) Arg {
	aa := Arg{a}
	validArgs = append(validArgs, aa)
	return aa
}

var (
	typeDisk = reflect.TypeOf((*disk.Disk)(nil))
	typeWriter = reflect.TypeOf((*io.Writer)(nil)).Elem()
	typeFile = reflect.TypeOf((*os.File)(nil))
)

type argDiskType struct{}

func (argDiskType) name() string {
	return "disk"
}

func (argDiskType) desc() []string {
	return []string{
		"a filename of a disk device or disk image;",
		"must exist and be an even number of sectors long",
	}
}

func (argDiskType) argtype() reflect.Type {
	return typeDisk
}

func (argDiskType) prepare(arg string) (out *argout, err error) {
	d, err := disk.Open(arg)
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

func (argOutFileType) desc() []string {
	return []string{
		"either a file to perform a binary dump to",
		"or - to perform a hexdump on stdout",
	}
}

func (argOutFileType) argtype() reflect.Type {
	return typeWriter
}

func (argOutFileType) prepare(arg string) (out *argout, err error) {
	var of io.WriteCloser

	if arg == "-" {
		of = hex.Dumper(os.Stdout)
	} else {
		f, err := os.Open(arg)
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

func (argOutImageType) desc() []string {
	return []string{
		"a filename to write the output disk image to;",
		"must not exist (reallymine will not overwrite existing files or media)",
	}
}

func (argOutImageType) argtype() reflect.Type {
	return typeFile
}

func (argOutImageType) prepare(arg string) (out *argout, err error) {
	f, err := os.OpenFile(arg, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}
	out = new(argout)
	out.obj = reflect.ValueOf(f)
	out.deferfunc = func() {
		f.Close()
	}
	return out, nil
}

var argOutImage = addarg(&argOutImageType{})
var ArgOutImage Arg = argOutImage

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

func formatDescription(desc []string, args []Arg) string {
	ai := make([]interface{}, len(args))
	for i, a := range args {
		ai[i] = a.a.name()
	}
	return usageL2(true, desc, ai...)
}

// for reallymine to use directly

func ArgUsage() string {
	out := ""
	for _, a := range validArgs {
		out += usageL1(a.a.name())
		out += usageL2(false, a.a.desc())
	}
	return out
}
