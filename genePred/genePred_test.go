package genePred

import (
	"testing"

	"github.com/goFish/simpleio"
)

var tests []string = []string{
	"testdata/rna-seq.genes.mapped.ensembl.gp.gz",
}

func TestGenePredRead(t *testing.T) {
	for _, i := range tests {
		ans := Read(i)
		textReader := simpleio.NewReader(i)
		var index int = 0
		for each, err := simpleio.ReadLine(textReader); !err; each, err = simpleio.ReadLine(textReader) {
			if ToString(&ans[index]) == each.String() {
				index++
			} else {
				t.Errorf("Error: genePred parsing did not match raw text...\n")
			}
		}
		if index != len(ans) {
			t.Errorf("Error: the number genePred structs and lines in the file do not match...\n")
		}
	}
}
