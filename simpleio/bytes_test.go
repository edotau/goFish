package simpleio

/*
import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/vertgenlab/gonomics/cigar"
	"github.com/vertgenlab/gonomics/dna"
	"io"
	"log"
	"strings"
	"testing"
	"unsafe"
)

var c1 cigar.Cigar = cigar.Cigar{RunLength: 35, Op: 'M'}
var c2 cigar.Cigar = cigar.Cigar{RunLength: 2, Op: 'I'}
var c3 cigar.Cigar = cigar.Cigar{RunLength: 16, Op: 'D'}

var strs = []string{"a", "b", "c", "d", "aaa", "bbb", "ccc", "ddd", "aaaa", "bbbb", "cccc", "dddd"}
var b1 ByteCigar = ByteCigar{RunLen: 35, Op: 'M'}
var b2 ByteCigar = ByteCigar{RunLen: 2, Op: 'I'}
var b3 ByteCigar = ByteCigar{RunLen: 16, Op: 'D'}

var dnaByte []byte = []byte{'a', 'a', 'a', 'a', 'a', 'C', 'C', 'C', 'C', 'T', 'T', 'T', 'a', 'a', 'a', 'a', 'a'}
var dnaString string = "aaaaaCCCCTTTaaaaa"

func TestBytesDna(t *testing.T) {
	dnaB := ByteSliceToDnaBases(dnaByte)
	fmt.Printf("%s\n", "aaaaaCCCCTTTaaaaa")
	fmt.Printf("%s\n", dna.BasesToString(dnaB))
}

func BenchmarkCopyDna(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dnaOne := ByteSliceToDnaBases([]byte("ATGGAGGAGTTAGAGGAGA"))
		dnaClone := make([]dna.Base, len(dnaOne))
		copy(dnaClone, dnaOne)
		if len(dnaClone) != len(dnaOne) {
			log.Printf("Error: did not copy over the same number of bases...\n")
		}
	}
}

func BenchmarkAppendDna(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		dnaOne := ByteSliceToDnaBases([]byte("ATGGAGGAGTTAGAGGAGA"))
		dnaClone := append(dnaOne[:0:0], dnaOne...)
		if len(dnaClone) != len(dnaOne) {
			log.Printf("Error: did not copy over the same number of bases...\n")
		}
	}
}

func BenchmarkReadindDnaBytes(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		ByteSliceToDnaBases(dnaByte)
	}
}

func BenchmarkStringToDnaBase(b *testing.B) {
	s := string(dnaByte)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		dna.StringToBases(s)
	}
}

func TestBytesToCigar(t *testing.T) {
	var cigarbytes = []byte("35M2I16D")
	bc := BytesToCigar(cigarbytes)
	fmt.Printf("%s\n", ByteCigarString(bc))
}

func TestCig(t *testing.T) {
	var cigarsString = "35M2I16D"
	testCig := []ByteCigar{b1, b2, b3}
	if s := ByteCigarString(testCig); s == cigarsString {
		fmt.Printf("%s\n", s)
	}
	fmt.Println(unsafe.Alignof(testCig))
}

func BenchmarkBytesToCigar(b *testing.B) {
	b.ReportAllocs()
	var cigarbytes = []byte("35M2I16D")
	for n := 0; n < b.N; n++ {
		BytesToCigar(cigarbytes)
	}
}

func BenchmarkStringToCigar(b *testing.B) {
	b.ReportAllocs()
	var cigarstring string = "35M2I16D"
	for n := 0; n < b.N; n++ {

		cigar.FromString(cigarstring)
	}
}

func TestOldStringCig(t *testing.T) {
	var cigarsString = "35M2I16D"
	testCig := []ByteCigar{b1, b2, b3}
	if s := cigar.ToString([]*cigar.Cigar{&c1, &c2, &c3}); s == cigarsString {
		fmt.Printf("%s\n", s)
	}
	fmt.Println(unsafe.Alignof(testCig))
}

func BenchmarkByteCigar(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		ByteCigarString([]ByteCigar{b1, b2, b3})
	}
}

func BenchmarkOldCig(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		cigar.ToString([]*cigar.Cigar{&c1, &c2, &c3})
	}
}

func concatB(strs []string, w io.Writer) {

	bw := bufio.NewWriter(w)
	bw.WriteString("START")
	bw.WriteString(strs[0])
	for _, v := range strs {
		bw.WriteString("-")
		bw.WriteString(v)
	}
	bw.WriteString("END")
	// Flush
	bw.Flush()
}

func BenchmarkConcatBytesBuffer(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var b bytes.Buffer
		concatA(strs, &b)
	}
}

func BenchmarkConcatStringsJoin(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var b bytes.Buffer
		concatC(strs, &b)
	}
}
func BenchmarkConcatWriter(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var b bytes.Buffer
		concatD(strs, &b)
	}
}

func BenchmarkConcatCopy(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var b bytes.Buffer
		concatE(strs, &b)
	}
}

func concatA(strs []string, w io.Writer) {

	var buf bytes.Buffer
	buf.WriteString("START")
	buf.WriteString(strs[0])
	for _, v := range strs {
		buf.WriteString("-")
		buf.WriteString(v)
	}
	buf.WriteString("END")
	w.Write(buf.Bytes())
}

func concatC(strs []string, w io.Writer) {

	var buf bytes.Buffer
	buf.WriteString("START")
	buf.WriteString(strings.Join(strs, "-"))
	buf.WriteString("END")
	w.Write(buf.Bytes())
}

func concatD(strs []string, w io.Writer) {

	w.Write([]byte("START"))

	w.Write([]byte(strs[0]))
	for _, v := range strs {
		w.Write([]byte("-"))
		w.Write([]byte(v))
	}
	w.Write([]byte("END"))
}

func concatE(strs []string, w io.Writer) {

	n := len("START") + len("END")
	for _, s := range strs {
		n += len(s) + 1
	}
	buf := make([]byte, n)

	n = copy(buf, "START")
	for _, s := range strs {
		n += copy(buf[n:], s)
		buf[n] = '-'
		n++
	}
	copy(buf[n:], "END")
	w.Write(buf)
}*/
