package simpleio

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkBgzipReader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ReadBgzipFile("testdata/atacseq_simple_pool.vcf.gz")
	}

}

func BenchmarkSimpleReader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		SimpleReaderBench("testdata/atacseq_simple_pool.vcf.gz")

	}
}

func SimpleReaderBench(filename string) *bytes.Buffer {
	reader := NewReader(filename)
	var ans bytes.Buffer
	for curr, done := ReadLine(reader); !done; curr, done = ReadLine(reader) {
		curr.WriteByte('\n')
		io.Copy(&ans, curr)
		//Nothing to assign, testing pure reading of the file
	}
	return &ans
}
