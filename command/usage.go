// 23 september 2016
package command

import (
	"fmt"
	"strings"
)

// See package flag's source for details on this formatting.
// We do one thing package flag doesn't: wrap descriptions.

func usageL1(s string) string {
	return fmt.Sprintf("  %s\n", s)
}

const prefixL2 = "    	"

func usageL2(format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)
	lines := wrapL2(s)
	s = ""
	for _, t := range lines {
		s += prefixL2 + t + "\n"
	}
	return s
}

// a safe bet, given 8 spaces according to the L2 prefix and 80 characters as an average terminal width
const maxWidth = 70

func wrapL2(s string) (wrapped []string) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		wrappedLine := words[0]
		// ensure bulleted lists are indented
		linePrefix := ""
		if words[0] == "-" {
			linePrefix = "  "
			words = words[1:]
			wrappedLine += " " + words[0]
		}
		for _, word := range words[1:] {
			candidate := wrappedLine + " " + word
			if len(candidate) > maxWidth {
				wrapped = append(wrapped, wrappedLine)
				// and start the next line with the next word
				candidate = linePrefix + word
			}
			wrappedLine = candidate
		}
		if wrappedLine != "" { // add the last line
			wrapped = append(wrapped, wrappedLine)
		}
	}
	return wrapped
}

func ToFlagUsage(s string) string {
	wrapped := wrapL2(s)
	return strings.Join(wrapped, "\n"+prefixL2)
}
