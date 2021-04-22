package vcf

import (
	"testing"
)

func TestVcfReader(t *testing.T) {
	ReadVcfs("testdata/small.vcf.gz")
}

var testVcf string = "testdata/small.vcf.gz"

func BenchmarkUnmarshalVcf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {

		unmarshalVcf()
	}
}

func unmarshalVcf() {
	reader := NewReader(testVcf)
	ReadHeader(reader)
	for _, done := UnmarshalVcf(reader); !done; _, done = UnmarshalVcf(reader) {

	}
}
