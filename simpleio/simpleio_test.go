package simpleio

import (
	"testing"
)

func TestSimpleReader(t *testing.T) {
	reader := NewReader("testdata/atacseq_simple_pool.vcf.gz")
	for _, done := ReadLine(reader); !done; _, done = ReadLine(reader) {
	}
}

func TestSimpleWriter(t *testing.T) {
	ans := ReadFromFile("testdata/atacseq_simple_pool.vcf.gz")
	writer := "testdata/simplewriter_test.gz"
	WriteToFile(writer, ans)
	test := ReadFromFile(writer)
	if len(ans) == len(test) {
		for i := 0; i < len(ans); i++ {
			if ans[i] != test[i] {
				t.Errorf("Error: simple writer did not produce an identical file...\n")
			}
		}
	} else {
		t.Errorf("Error: number of lines read do not match...\n")
	}
	Remove("testdata/simplewriter_test.gz")
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
		reader := NewReader("testdata/atacseq_simple_pool.vcf.gz")
		for _, done := ReadLine(reader); !done; _, done = ReadLine(reader) {
			//Nothing to assign, testing pure reading of the file
		}
	}
}
