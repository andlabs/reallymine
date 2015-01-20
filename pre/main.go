// 20 january 2015
package main

import (
	"fmt"
	"os"
)

// TODO:
// - colorize error output?
// - do an extra check to make sure in isn't a symlink to out? or isn't the same as out?
// - make sure all cases of errf() have a trailing newline

var inname string
var in *os.File
var instat os.FileInfo
var outname string
var out *os.File
var outstat os.FileInfo

func die() {
	if out != nil {
		action := "closing"
		err := out.Close()
		if err == nil {
			action = "removing"
			err = os.Remove(outname)
		}
		if err != nil {
			errf("\nError %s output file %s: %v\n", action, outname, err)
			errf("You will need to remove this file manually before running reallymine again.\n")
		}
	}
	os.Exit(1)
}

func errf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func dd_rescuenote() {
	errf(`
Note to Linux users: Some distributions, such as
Debian and Ubuntu, provide both a non-GNU
ddrescue/"dd_rescue" and a GNU ddrescue/"gddrescue"
in their package repositories. The non-GNU one is
less robust and and less feature-ful than the
GNU one, and thus should not be used anymore.
Of particular importance is the log file feature,
which can be used later to aid recovering from
damaged drives.
For more information, see 
http://askubuntu.com/a/211579/257298
`)
}

func main() {
	var err error

	if len(os.Args) != 3 {
		errf(`usage: %s infile outfile
    infile must be a regular file containing a
        raw dump of the encrypted hard drive,
        such as one produced by GNU ddrescue
    outfile must not exist
`, os.Args[0])
		dd_rescuenote()
		die()
	}

	inname = os.Args[1]
	in, err = os.Open(inname)
	if err != nil {
		errf("Error opening input file %s: %v\n", inname, err)
		die()
	}
	instat, err = in.Stat()
	if err != nil {
		errf("Error getting information about input file %s: %v\n", inname, err)
		die()
	}
	inkind := instat.Mode() & os.ModeType
	if (inkind & os.ModeDevice) != 0 {
		errf(`Error: the input file %s appears to be a device file.
This likely means you are trying to run reallymine
directly on the hard drive you want to decrypt.
Please don't. This tool will likely strain that
drive. It is not designed to read out entire hard
drives. Instead, make an image of the drive using
a tool such as GNU ddrescue. These tools are more
suited to reading entire hard drives all at once.
Once you have an image of the drive, pass that to
reallymine instead to get started.
`, inname)
		dd_rescuenote()
		die()
	}
	// TODO check if it's a regular file?

	outname = os.Args[2]
	out, err = os.OpenFile(outname, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if os.IsExist(err) {
		errf(`Error: output file %s already exists.
reallymine will not blindly overwrite files that
already exist. It will also not overwrite the input
file with the unencrypted output. Please specify a
file that does not exist yet to get started.
`, outname)
		die()
	} else if err != nil {
		errf("Error creating output file %s: %v\n", outname, err)
		die()
	}
	outstat, err = out.Stat()
	if err != nil {
		errf("Error getting information about output file %s: %v\n", err)
		die()
	}
}
