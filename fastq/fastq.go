// Package fastq contains data structures and functions to process sequencing data coming off the sequencer
package fastq

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/goFish/simpleio"
)

type Fastq struct {
	Name string
	Seq  []byte
	Qual []byte
}

func Read(filename string) []Fastq {
	var ans []Fastq
	reader := simpleio.NewReader(filename)
	defer reader.Close()
	for i, done := GunzipFastq(reader); !done; i, done = GunzipFastq(reader) {
		ans = append(ans, *i)
	}
	return ans
}

func GunzipFastq(reader *simpleio.SimpleReader) (*Fastq, bool) {
	ans := Fastq{}
	line, done := simpleio.ReadLine(reader)
	if done {
		return nil, true
	}

	ans.Name = line.String()[1:]
	line, done = simpleio.ReadLine(reader)
	if done {
		return nil, true
	}
	ans.Seq = make([]byte, len(line.Bytes()))
	copy(ans.Seq, line.Bytes())

	line, done = simpleio.ReadLine(reader)
	if line.String() != "+" {
		log.Fatalf("Error: This line should be a + (plus) sign \n")
	}

	line, done = simpleio.ReadLine(reader)
	if done {
		return nil, true
	}
	ans.Qual = make([]byte, len(line.Bytes()))
	copy(ans.Qual, line.Bytes())

	return &ans, false
}

func ToString(fq *Fastq) string {
	var buffer strings.Builder

	_, err := buffer.WriteString(fq.Name)
	simpleio.ErrorHandle(err)
	err = buffer.WriteByte('\n')
	simpleio.ErrorHandle(err)

	_, err = buffer.Write(fq.Seq)
	simpleio.ErrorHandle(err)
	err = buffer.WriteByte('\n')
	simpleio.ErrorHandle(err)

	err = buffer.WriteByte('+')
	simpleio.ErrorHandle(err)
	err = buffer.WriteByte('\n')
	simpleio.ErrorHandle(err)

	_, err = buffer.Write(fq.Qual)
	simpleio.ErrorHandle(err)
	err = buffer.WriteByte('\n')
	simpleio.ErrorHandle(err)

	return buffer.String()
}

func Equal(x, y *Fastq) bool {
	if x.Name != y.Name {
		return false
	}
	if !SameByteSlice(x.Seq, y.Seq) {
		return false
	}
	if !SameByteSlice(x.Qual, y.Qual) {
		return false
	}
	return true
}

func SameByteSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func FastqReader(filename string) <-chan Fastq {
	fqs := make(chan Fastq)
	reader := simpleio.NewReader(filename)
	var wg sync.WaitGroup
	wg.Add(1)
	go ToChannel(reader, fqs, &wg)
	go func() {
		wg.Wait()
		close(fqs)
	}()
	return fqs
}

func ToChannel(reader *simpleio.SimpleReader, fq chan<- Fastq, wg *sync.WaitGroup) {
	for i, done := GunzipFastq(reader); !done; i, done = GunzipFastq(reader) {
		fq <- *i
	}
	wg.Done()
}

func MetricsTable(fq *Fastq) string {
	//fmt.Printf("ReadName\tLength\tQuality\n")
	buffer := &strings.Builder{}
	//for each := range fqs {
	//	buffer.Reset()
	buffer.WriteString(PrintName(GetInfo(fq.Name)))
	buffer.WriteByte('\t')
	buffer.WriteString(simpleio.IntToString(len(fq.Seq)))
	buffer.WriteByte('\t')
	buffer.WriteString(fmt.Sprintf("%f\n", FindAveQuality(fq.Qual)))
	//fmt.Print(buffer.WriteString(buffer.String()))
	//fmt.Print(buffer.String())
	//}
	return buffer.String()
}

func FindAveReadLength(fqs <-chan Fastq) float64 {
	var total int = 0
	var length int = 0
	for each := range fqs {
		length += len(each.Seq)
		total++
	}
	return float64(length) / float64(total)
}

func FindAveQuality(fqs []byte) float64 {
	var total int = 0
	var length int = 0
	for _, each := range fqs {
		length += int(each)
		total++
	}
	return float64(length) / float64(total)
}

type Info struct {
	Name  string
	Read  byte
	Index string
}

func PrintName(fq Info) string {
	return fq.Name
}
func GetInfo(s string) Info {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		log.Fatalf("Error: did not parse space character correctly...\n")
	}
	ans := Info{
		Name: parts[0],
	}
	parts = strings.Split(parts[1], ":")
	if parts[0] == "1" {
		ans.Read = '+'
	} else {
		ans.Read = '-'
	}
	ans.Index = parts[len(parts)-1]
	return ans

}
