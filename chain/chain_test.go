package chain

import (
	"os"
	"testing"
)

var readWriteTests = []struct {
	filename string // input
}{
	{"testdata/axtTest.chain.gz"},
	{"testdata/big.chain.gz"},
	{"testdata/twoChainZ.chain.gz"},
}

func TestReadAndWrite(t *testing.T) {

	for _, test := range readWriteTests {
		input := Read(test.filename)
		name := test.filename + ".tmp"
		Write(name, input)
		output := Read(name)
		if len(input) == len(output) {
			for i := 0; i < len(output); i++ {
				if ToString(&input[i]) != ToString(&output[i]) {

					t.Errorf("Error: chain package read and write testing failed...\n")
				}
			}
		}
		os.Remove(name)
	}
}
