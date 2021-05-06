package geneSeq

import (
	"log"
	"strings"
	"testing"

	"github.com/edotau/goFish/simpleio"
)

var gtf_tests []string = []string{
	"testdata/gasAcu_small.BROADS1.104.gtf",
	"testdata/denovo.reference.transcrtipt.assembly.annotated.gtf",
	"testdata/geneFeat.gzipTest.gtf.gz",
}

func TestGtfReading(t *testing.T) {
	for _, test := range gtf_tests {
		ans := simpleio.ReadFromFile(test)
		gf := ReadGtf(test)
		if len(ans) != len(gf) {
			t.Errorf("Error: line numbers between read in functions are different...\n")
		} else {
			for i := 0; i < len(gf); i++ {
				if !strings.HasPrefix(ans[i], "#") && gf[i].ToString() != ans[i] {
					log.Fatalf("Error: to string method is not producing the same text as the input...\n")
				}
			}
		}

	}
}

func BenchmarkGtfSimpleReading(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ReadGtf("testdata/geneFeat.gzipTest.gtf.gz")
	}
}
