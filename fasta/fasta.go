// Package fasta contains data structures for processes and functions that operate on fasta files
// containing base sequences
package fasta

import (
	"bytes"
	"log"
	"strings"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/simpleio"
)

type Fasta struct {
	Name string
	Seq  []code.Dna
}

func FastaReader(reader *simpleio.SimpleReader) (*Fasta, bool) {
	var done bool
	reader.Buffer, done = simpleio.ReadLine(reader)
	if !done {
		if strings.HasPrefix(reader.Buffer.String(), ">") {
			return &Fasta{
				Name: reader.Buffer.String()[1:],
				Seq:  GetSequence(reader),
			}, false
		}

	}
	return nil, true
}

func GetSequence(reader *simpleio.SimpleReader) []code.Dna {
	var line *bytes.Buffer
	var answer []code.Dna
	for nextBytes, err := reader.Peek(1); len(nextBytes) > 0 && !bytes.HasPrefix(nextBytes, []byte(">")) && err == nil; nextBytes, err = reader.Peek(1) {
		line, _ = simpleio.ReadLine(reader)
		answer = append(answer, code.ToDna(line.Bytes())...)
	}
	return answer
}

func Read(filename string) []Fasta {
	reader := simpleio.NewReader(filename)
	var ans []Fasta
	for i, err := FastaReader(reader); !err; i, err = FastaReader(reader) {
		ans = append(ans, *i)
	}
	return ans
}

func FetchHttpToMap(link string) map[string][]code.Dna {
	hash := make(map[string][]code.Dna)
	stream := simpleio.NewReader(link)
	for i, err := FastaReader(stream); !err; i, err = FastaReader(stream) {
		if _, key := hash[i.Name]; !key {
			hash[i.Name] = i.Seq
		} else {
			log.Fatalf("Error: fasta files does not contain unique header names...\n")
		}
	}
	stream.Close()
	return hash
}

func (fa *Fasta) ToString() string {
	buf := &strings.Builder{}
	buf.WriteByte('>')
	buf.WriteString(fa.Name)
	buf.WriteByte('\n')
	for i := 0; i < len(fa.Seq); i += 50 {
		if i+50 > len(fa.Seq) {
			buf.Grow(len(fa.Seq[i:]))
			for _, i := range fa.Seq[i:] {
				buf.WriteByte(code.DnaToByte(i))
			}
			buf.WriteByte('\n')
		} else {
			buf.Grow(50)
			for _, i := range fa.Seq[i : i+50] {
				buf.WriteByte(code.DnaToByte(i))
			}
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
