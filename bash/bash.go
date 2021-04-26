// bash package implements simple bash commands that are sometimes useful when writing go scripts
package bash

import (
	"strings"
)

// Cut mimics the basic bash command to cut a column by any delim and returns the fields indices you specify
// indices start at 0
func Cut(line string, d byte, f ...int) []string {
	columns := strings.Split(line, string(d))
	var ans []string
	for _, i := range f {
		ans = append(ans, columns[i])
	}
	return ans
}
