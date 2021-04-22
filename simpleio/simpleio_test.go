package simpleio

import (
	"testing"
)

func TestSimpleReader(t *testing.T) {
	reader := NewReader("testdata/atacseq_simple_pool.vcf.gz")
	for _, done := ReadLine(reader); !done; _, done = ReadLine(reader) {
	}
}

func BenchmarkSimpleReading(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		reader := NewReader("testdata/atac_seq_table.txt")
		for _, done := ReadLine(reader); !done; _, done = ReadLine(reader) {
			//Nothing to assign, testing pure reading of the file
		}
	}
}

func BenchmarkSimpleReaderGz(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		reader := NewReader("testdata/simplePoolTest.vcf.gz")
		for _, done := ReadLine(reader); !done; _, done = ReadLine(reader) {
			//Nothing to assign, testing pure reading of the file
		}
	}
}
