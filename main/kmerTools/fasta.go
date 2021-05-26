package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

// Fasta is a struct/obj capturing a fasta record and contains a name and sequence coded as bytes
type Fasta struct {
	Name string
	Seq  []Dna
}

// ReadFasta takes a filename string as an input and output an array/list containing all fasta records inside the file
func ReadFasta(filename string) []Fasta {
	var ans []Fasta
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	for i, done := readFastaHelper(reader); done == nil; i, done = readFastaHelper(reader) {
		ans = append(ans, *i)
	}
	return ans
}

// readFastaHelper is a helper function that takes an buffered io reader and produces a Fasta struct
func readFastaHelper(reader *bufio.Reader) (*Fasta, error) {
	for line, err := reader.ReadBytes('\n'); err == nil; line, err = reader.ReadBytes('\n') {
		// Looks for lines that begins with ">" header.
		// Greater sign is trimmed and trailing end of line character is also trimmed
		if strings.HasPrefix(string(line), ">") {
			return &Fasta{
				Name: string(line[1 : len(line)-1]),
				Seq:  parseDnaBases(reader),
			}, nil
		}
	}
	return nil, errors.New("Error: there was an error processing Fasta...\n")
}

// ToString prints fasta record to stdout
// Used for debuging and writing fasta records
func (fa *Fasta) ToString() string {
	buf := &strings.Builder{}
	buf.WriteByte('>')
	buf.WriteString(fa.Name)
	buf.WriteByte('\n')
	for i := 0; i < len(fa.Seq); i += 50 {
		if i+50 > len(fa.Seq) {
			buf.Grow(len(fa.Seq[i:]))
			for _, i := range fa.Seq[i:] {
				buf.WriteByte(DnaToByte(i))
			}
			buf.WriteByte('\n')
		} else {
			buf.Grow(50)
			for _, i := range fa.Seq[i : i+50] {
				buf.WriteByte(DnaToByte(i))
			}
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
