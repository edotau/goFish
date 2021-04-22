package fasta

import (
	"fmt"
	"testing"
)

func TestFasta(t *testing.T) {
	fa := Read("testdata/small.fa")
	for _, i := range fa {
		fmt.Printf("name=%s len=%d\n", i.Name, len(i.Seq))
	}

}

func BenchmarkFastaReading(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		Read("testdata/small.fa")
	}
}
