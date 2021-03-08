// 8 March 2021
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/undeadbanegithub/reallymine/command"
	"github.com/undeadbanegithub/reallymine/decryptloop"
	"github.com/undeadbanegithub/reallymine/disk"
)

func errf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func die(format string, args ...interface{}) {
	errf(format, args...)
	errf("\n")
	os.Exit(1)
}

var Commands = []*command.Command{
	// TODO have a clear separation between user commands and research commands? if so, decrypt and getdek go above the rest
	dumplast,
	dumpkeysector,
	decryptkeysector,
	getdek,
	dumpfirst,
	decryptfile,
	decrypt,
}

func init() {
	problems := command.Validate(Commands)
	if len(problems) != 0 {
		errf("issues with reallymine commands:\n")
		errf("%s\n", strings.Join(problems, "\n"))
		errf("this means andlabs made a mistake; contact him\n")
		os.Exit(1)
	}
}

func init() {
	flag.Int64Var(
		&command.DiskSize,
		"disk-size",
		-1,
		command.ToFlagUsage(
			fmt.Sprintf(
				"Overrides the size of the disk to use, allowing you to pretend the disk is smaller than it "+
					"is. This value must be less than or equal to the disk's actual size, and must be a multiple of "+
					"the sector size (%d bytes). If the size is -1, the disk's actual size is used.",
				disk.SectorSize,
			),
		),
	)
	flag.IntVar(
		&command.WriteBufSize,
		"write-buf-size",
		10*1024*1024,
		command.ToFlagUsage(
			"Overrides the size of write buffer to be used during the encryption. The bigger the buffer, the "+
				"bigger chunks would be first decrypted into memory and then written to disk. You might want to set "+
				"the value higher, if the destination device is slow, to minimize the number of writes, else it "+
				"would not make much of a difference.",
		),
	)
}

func usage() {
	errf("usage: %s [options] command [args...]\n", os.Args[0])
	errf("options:\n")
	flag.PrintDefaults()
	errf("commands:\n")
	// TODO refine these names?
	errf("%s", command.FormatUsage(Commands))
	errf("command arguments:\n")
	errf("%s", command.ArgUsage())
	errf("decryption steps:\n")
	errf("%s", decryptloop.StepUsage())
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}
	cmd := flag.Arg(0)

	for _, c := range Commands {
		if cmd != c.Name {
			continue
		}
		args := flag.Args()[1:]
		err := c.Invoke(args)
		if err == command.ErrWrongArgCount {
			errf("error running %s: %v\n", c.Name, err)
			usage()
		}
		if err != nil {
			die("error running %s: %v", c.Name, err)
		}
		// all good; return successfully
		return
	}

	errf("error: unknown command %q\n", cmd)
	usage()
}
