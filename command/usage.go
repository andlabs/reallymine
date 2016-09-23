// 23 september 2016
package command

import (
	"fmt"
	"strings"
)

// See package flag's source for details on this formatting.

func usageL1(s string) string {
	return mt.Sprintf("  %s\n", s)
}

func usageL2(format bool, s []string, args ...interface{}) string {
	// don't modify s
	t := make([]string, len(s))
	copy(t, s)
	for i, _ := range t {
		if format {
			t[i] = fmt.Sprintf(t[i], args...)
		}
		t[i] = "    	" + t[i] + "\n"
	}
	return strings.Join(t, "")
}
