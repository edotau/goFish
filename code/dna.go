package code

import (
	"log"
	"strings"
)

type Dna byte

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
	Gap   Dna = '-'
)

// ByteToDna converts a byte into a dna.Base if it matches one of the acceptable DNA characters.
// Notes: It will also mask the lower case values and return dna.Base as uppercase bases.
// Note: '*', used by VCF to denote deleted alleles, becomes a Gap in DNA.
func ByteToDnaNoMask(b byte) Dna {
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
		return A
	case 'c':
		return C
	case 'g':
		return G
	case 't':
		return T
	case '-':
		return Gap
	case '*':
		return Gap
	default:
		log.Fatalf("Error: unexpected character in dna %c\n", b)
		return N
	}
}

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
	case '-':
		return Gap
	case '*':
		return Gap
	default:
		log.Fatalf("Error: unexpected character in dna %c\n", b)
		return N
	}
}

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
	case Gap:
		return '-'
	default:
		log.Fatalf("Error: unexpected character in dna %c\n", b)
		return 'N'
	}
}

func DnaToByteNoMask(b Dna) byte {
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
		return 'A'
	case MaskC:
		return 'C'
	case MaskG:
		return 'G'
	case MaskT:
		return 'T'
	case Gap:
		return '-'
	default:
		log.Fatalf("Error: unexpected character in dna %c\n", b)
		return 'N'
	}
}

// ToDna will convert a slice of bytes into a slice of Bases with no lowercase bases.
func ToDna(b []byte) []Dna {
	var answer []Dna = make([]Dna, len(b))
	for i, byteValue := range b {
		answer[i] = ByteToDna(byteValue)
	}
	return answer
}

func ToString(bases []Dna) string {
	buffer := strings.Builder{}
	buffer.Grow(len(bases))
	for _, i := range bases {
		buffer.WriteByte(DnaToByteNoMask(i))
	}
	return buffer.String()
}
