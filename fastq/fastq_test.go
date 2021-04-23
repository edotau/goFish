package fastq

import (
	"fmt"
	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/simpleio"
	"github.com/vertgenlab/gonomics/fastq"
	"testing"
)

var toyFq *Fastq = &Fastq{
	Name: "A00257:502:HWWM3DRXX:1:2101:1000:1031 1:N:0:NCTAGGAG+NGTTGCAA",
	Seq:  code.ToDna([]byte("NCTAGGAG")),
	Qual: []byte("#FFFFFFF"),
}

/*
func TestQual(t *testing.T) {
	fqs := FastqReader("testdata/Q16-5_S15_L001_R1_001.fastq.gz")
	fmt.Printf("ReadName\tLength\tQuality\n")
	for i := range fqs {
		MetricsTable(&i)
	}
}*/

func TestToyFastq(t *testing.T) {
	fmt.Printf("%v\n", toyFq)
	reader := simpleio.NewReader("testdata/toy.fastq")
	defer reader.Close()
	i, done := GunzipFastq(reader)
	if done {
		t.Errorf("Error: was not able to process any fastq records...\n")
	}
	if !Equal(i, toyFq) {
		t.Errorf("Error: io processing did not match original fastq toy example...\n\n%s\n\n%s\n", ToString(i), ToString(toyFq))
	}
}

func BenchmarkIoSpeed(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Read("testdata/toy.fastq")
	}
}

func BenchmarkGonomics(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		fastq.Read("testdata/toy.fastq")
	}
}
