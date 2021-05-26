package main

import (
	"bufio"
	"bytes"
	"log"
	"strings"
)

// Dna bases are each encoded as a single byte
type Dna byte

// The following are valid Dna bases
const (
	A     Dna = 'A'
	C     Dna = 'C'
	G     Dna = 'G'
	T     Dna = 'T'
	N     Dna = 'N'
	MaskA Dna = 'a'
	MaskC Dna = 'c'
	MaskG Dna = 'g'
	MaskT Dna = 't'
	MaskN Dna = 'n'
	Gap   Dna = '-'
)

// ByteToDna is a helper function to convert a single byte into DNA constants pre-defined, which will detect invalid bases
func ByteToDna(b byte) Dna {
	switch b {
	case 'A':
		return A
	case 'C':
		return C
	case 'G':
		return G
	case 'T':
		return T
	case 'N':
		return N
	case 'a':
		return MaskA
	case 'c':
		return MaskC
	case 'g':
		return MaskG
	case 't':
		return MaskT
	case 'n':
		return MaskN
	case '-':
		return Gap
	case '*':
		return Gap
	default:
		log.Fatalf("Error: unexpected character in dna %c\n", b)
		// This return statement is for the sole purpose of compling.
		// log.Fatal will exit the function a and will never return this N if we encounter the error
		return N
	}
}

// DnaToByte will convert a single Dna base back to a byte
func DnaToByte(b Dna) byte {
	switch b {
	case A:
		return 'A'
	case C:
		return 'C'
	case G:
		return 'G'
	case T:
		return 'T'
	case N:
		return 'N'
	case MaskA:
		return 'a'
	case MaskC:
		return 'c'
	case MaskG:
		return 'g'
	case MaskT:
		return 't'
	case MaskN:
		return 'n'
	case Gap:
		return '-'
	default:
		log.Fatalf("Error: unexpected character in dna %c\n", b)
		return 'N'
	}
}

func parseDnaBases(reader *bufio.Reader) []Dna {
	var answer []Dna
	for nextBytes, err := reader.Peek(1); len(nextBytes) > 0 && !bytes.HasPrefix(nextBytes, []byte(">")) && err == nil; nextBytes, err = reader.Peek(1) {
		line, _ := reader.ReadBytes('\n')
		answer = append(answer, ToDna(line[:len(line)-1])...)
	}
	return answer
}

// ToDna will convert a slice of bytes into a slice of Bases with no Maskcase bases.
func ToDna(b []byte) []Dna {
	var answer []Dna = make([]Dna, len(b))
	for i, byteValue := range b {
		answer[i] = ByteToDna(byteValue)
	}
	return answer
}

// ToString prints fasta record to stdout
// Used for debuging and writing fasta records
func DnaToString(bases []Dna) string {
	buf := &strings.Builder{}
	buf.Grow(len(bases))
	for i := 0; i < len(bases); i++ {
		buf.WriteByte(DnaToByte(bases[i]))
	}
	return buf.String()
}
