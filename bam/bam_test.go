package bam

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/goFish/code"
)

//readBamTests are the files pairs used to test our binary reader functions
var readBamTests = []struct {
	bam string
	sam string
}{
	{"testdata/tenXbarcodeTest.bam", "testdata/tenXbarcodeTest.sam"},
}

func BenchmarkNewSamReader(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		ReadSamRecord("testdata/tenXbarcodeTest.sam")
	}
}

//TestBamToSamReader will convert a bam file into a sam record and
//perform a comparison with the same sam file that was created using samtools view.
func TestBamToSamReader(t *testing.T) {
	for _, test := range readBamTests {
		_, bamFile := BasicRead(test.bam)
		samFile := ReadSamRecord(test.sam)

		if len(bamFile) != len(samFile) {
			t.Errorf("Error: File lines are not equal...\n")
		}
		for i := 0; i < len(bamFile); i++ {
			if !IsEqualDebug(bamFile[i], samFile[i]) {
				t.Fatalf("Error: Did not create the same sam file as samtools view...\n")
			}
		}
	}
}

//BenchmarkBamReader will benchmark the speed of decoding a bam file and convert the data into sam records.
//(This is used to compare with reading a sam file)
func BenchmarkBamReader(b *testing.B) {
	var bamFile []*Sam
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, test := range readBamTests {
			_, bamFile = BasicRead(test.bam)
		}

	}
	if bamFile == nil {
		b.Errorf("Error: There was a problem reading in the sam file after %d lines...\n", len(bamFile))
	}
}

func IsEqualDebug(a *Sam, b *Sam) bool {
	if strings.Compare(a.QName, b.QName) != 0 {
		log.Printf("Names are not equal: %s != %s\n", a.QName, b.QName)
		return false
	}
	if a.Flag != b.Flag {
		log.Printf("Flags are not equal: %d != %d\n", a.Flag, b.Flag)
		return false
	}
	if strings.Compare(a.RName, b.RName) != 0 {
		log.Printf("Names are not equal: %s != %s\n", a.RName, b.RName)
		return false
	}
	if a.Pos != b.Pos {
		log.Printf("Positions are not equal: %d != %d\n", a.Pos, b.Pos)
		return false
	}
	if a.MapQ != b.MapQ {
		log.Printf("Mapping Quals are not equal: %d != %d\n", a.MapQ, b.MapQ)
		return false
	}
	if strings.Compare(ByteCigarToString(a.Cigar), ByteCigarToString(b.Cigar)) != 0 {
		log.Printf("Cigars are not equal: %s != %s\n", ByteCigarToString(a.Cigar), ByteCigarToString(b.Cigar))
		return false
	}
	if strings.Compare(a.MateRef, b.MateRef) != 0 {
		log.Printf("RNext are not equal: %s != %s\n", a.MateRef, b.MateRef)
		return false
	}
	if a.MatePos != b.MatePos {
		log.Printf("PNext are not equal: %d != %d\n", a.MatePos, b.MatePos)
		return false
	}
	if a.TmpLen != b.TmpLen {
		log.Printf("TLen are not equal: %d != %d\n", a.TmpLen, b.TmpLen)
		return false
	}
	if string(a.Seq) != string(b.Seq) {
		log.Printf("Sequences are different:\n%s\n%s\n", code.ToString(a.Seq), code.ToString(b.Seq))
		return false
	}
	if !bytes.Equal(a.Qual, b.Qual) {
		log.Printf("Qual Scores are not the same...\n%s != %s", a.Qual, b.Qual)
		return false
	}
	if strings.Compare(a.Aux, b.Aux) != 0 {
		log.Printf("Need to reformat notes: \n%s\n%s\n", a.Aux, b.Aux)
		return false
	}

	return true
}
