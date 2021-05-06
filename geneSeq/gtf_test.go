package geneSeq

import (
	"fmt"
	"testing"
)

func TestGtfReading(t *testing.T) {
	gf := ReadGtf("testdata/denovo.reference.transcrtipt.assembly.annotated.gtf")
	for i := 0; i < 10; i++ {
		fmt.Printf("%v\n", (gf[i].ToString()))
	}
}

func BenchmarkGtfSimpleReading(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {

		ReadGtf("testdata/denovo.reference.transcrtipt.assembly.annotated.gtf")

		//for i := 0; i < 10; i++ {
		//		fmt.Printf("%v\n", (gf[i].ToString()))
		//	}
	}
}
