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

func usageL2(format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)
	lines := wrapL2(s)
	s = ""
	for _, t := range lines {
		s += "    	" + t + "\n"
	}
	return s
}

// a safe bet, given 8 spaces according to the L2 prefix and 80 characters as an average terminal width
const maxWidth = 70

func wrapL2(s string) (lines []string) {
	words := strings.Split(s, " ")
	line := words[0]
	for _, word := range words[1:] {
		candidate := line + " " + word
		if len(candidate) > maxWidth {
			lines = append(lines, line)
			line = word	// start next line with next word
			continue
		}
		line = candidate
	}
	if line != "" {		// add the last line
		lines = append(lines, line)
	}
	return lines
}
