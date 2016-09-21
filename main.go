// 22 october 2015
package main

import (
	"fmt"
	"os"
	"flag"
	"strings"

	"github.com/andlabs/reallymine/command"
//	"github.com/andlabs/reallymine/disk"
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
	// TODO
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

func usage() {
	errf("usage: %s [options] command [args...]\n", os.Args[0])
	errf("options:\n")
	flag.PrintDefaults()
	errf("commands:\n")
	errf("%s", command.FormatUsage(Commands))
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
			die("error running %s: %v\n", c.Name, err)
		}
		// all good; return successfully
		return
	}

	errf("error: unknown command %q\n", cmd)
	usage()
}
