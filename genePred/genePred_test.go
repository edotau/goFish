package genePred

import (
	"github.com/goFish/simpleio"
	"testing"
)

var tests []string = []string{
	"testdata/denovo_final.gp",
	"testdata/rna-seq_ensembl_psl_FINAL.gp",
	"testdata/rna-seq.genes.mapped.ensembl.gp",
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
