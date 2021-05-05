// bash package implements simple bash commands that are sometimes useful when writing go scripts
package bash

import (
	"os"
	"strings"

	"github.com/edotau/goFish/simpleio"
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

func GetColumnCount(line string, delim byte) int {
	return len(strings.Split(line, string(delim)))
}

// Mkdir will create all parent directories of a provided path
func Mkdir(path string, perm os.FileMode) {
	//0755
	err := os.MkdirAll("/Users/temp", perm)
	simpleio.FatalErr(err)
}
