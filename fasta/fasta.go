// Package fasta contains data structures for processes and functions that operate on fasta files
// containing base sequences
package fasta

import (
	"bytes"
	"strings"

	"github.com/goFish/code"
	"github.com/goFish/simpleio"
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
