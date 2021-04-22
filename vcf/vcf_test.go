package vcf

import (
	"testing"
	//"strings"
	"github.com/vertgenlab/gonomics/vcf"
	//"log"
)

func TestVcfReader(t *testing.T) {
	ReadVcfs("testdata/small.vcf")
	//for _, i := range ans {
	//	fmt.Printf("%s\n", ToString(&i))
	//format := len(strings.Split(i.Format, ":"))
	//samples := strings.Split(i.Samples, "\t")
	//fmt.Printf("ref=%s, alt=%s, ", i.Ref, i.Alt)

	//for j := 0; j < len(samples); j ++ {

	//	gt := strings.Split(samples[j], ":")
	//	if len(samples) != len(gt) {
	//log.Fatalf("format=%s, samples=%s\n",i.Format, gt)
	//log.Fatalf("genotype=%s, AD=%s, ", gt[0], gt[1])
	//	}

	//}
	//fmt.Printf("\n")
	//}
}

var testVcf string = "testdata/rabsTHREEbepa_chainNets.vcf.gz"

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
		//fmt.Printf("%v\n", v)
	}
	//if done {
	//	if v != nil {
	//		log.Fatalf("Error: last returned value should be nil...\n")
	//	}
	//}
}

func BenchmarkOldVcfReader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = vcf.Read(testVcf)
	}
}
