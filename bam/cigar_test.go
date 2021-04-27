package bam

import (
	"testing"
)

var b1 ByteCigar = ByteCigar{RunLen: 35, Op: 'M'}
var b2 ByteCigar = ByteCigar{RunLen: 2, Op: 'I'}
var b3 ByteCigar = ByteCigar{RunLen: 16, Op: 'D'}

func TestBytesToCigar(t *testing.T) {
	var cigarbytes = []byte("35M2I16D")
	bc := ReadToBytesCigar(cigarbytes)
	t.Logf("%s\n", ByteCigarToString(bc))
}

func BenchmarkBytesToCigar(b *testing.B) {
	b.ReportAllocs()
	var cigarbytes = []byte("35M2I16D")
	for n := 0; n < b.N; n++ {
		ReadToBytesCigar(cigarbytes)
	}
}
