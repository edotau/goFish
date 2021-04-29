package bash

import (
	"testing"

	"github.com/vertgenlab/gonomics/fileio"
)

func TestGunzipReader(t *testing.T) {
	var line int = 0
	reader := NewGunzipReader("testdata/atacseq_simple_pool.vcf.gz")
	for i, err := ReadLine(reader); !err; i, err = ReadLine(reader) {
		i.String()
		line++
	}
	if line != 80000 {
		t.Errorf("Error: Gunzip reader is not reading the correct number of lines...\n")
	}
}

func BenchmarkGunzipReader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		reader := NewGunzipReader("testdata/atacseq_simple_pool.vcf.gz")
		for _, err := ReadLine(reader); !err; _, err = ReadLine(reader) {

		}

	}
}
func BenchmarkSimpleReader(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		reader := fileio.NewByteReader("testdata/atacseq_simple_pool.vcf.gz")
		for _, err := fileio.ReadLine(reader); !err; _, err = fileio.ReadLine(reader) {

		}
	}
}

/*
func BenchmarkPzip(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data := simpleio.ReadFromFile("testdata/pgzip.vcf.gz")
		NewPzip("testdata/output.pgzip.tmp.gz", data)
	}
}

func BenchmarkSimpleWriter(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data := simpleio.ReadFromFile("testdata/atacseq_simple_pool.vcf.gz")
		simpleio.WriteToFile("testdata/output.simplewriter.tmp.gz", data)
	}
}
*/
