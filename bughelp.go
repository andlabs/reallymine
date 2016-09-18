// 22 october 2015
package main

import (
	"fmt"
	"os"
)

func BUG(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[BUG] ")
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\nPlease report to andlabs on github.com/andlabs/reallymine.\n")
	os.Exit(1)
}
