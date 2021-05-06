package geneSeq

import (
	"log"
	"testing"

	"github.com/edotau/goFish/simpleio"
)

var gff_tests []string = []string{
	"testdata/fish.gff3",
}

func TestGff3Reader(t *testing.T) {
	for _, test := range gff_tests {
		ans := simpleio.ReadFromFile(test)
		featFmtThree := ReadGff3(test)
		if len(ans) != len(featFmtThree) {
			t.Errorf("Error: line numbers between read in functions are different...\n")
		} else {
			for i := 0; i < len(featFmtThree); i++ {
				if featFmtThree[i].ToString() != ans[i] {
					log.Fatalf("Error: to string method is not producing the same text as the input...\n")
				}
			}
		}
	}
}
