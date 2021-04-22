package bam

import (
	"log"
)

func ConsumesReference(b byte) bool {
	switch b {
	case 'M':
		return true
	case 'I':
		return false
	case 'D':
		return true
	case 'N':
		return true
	case 'S':
		return false
	case 'H':
		return false
	case 'P':
		return false
	case '=':
		return true
	case 'X':
		return true
	}
	log.Fatalf("Invalid Byte: %c", b)
	return false
}

func GetIndel(c []ByteCigar) []ByteCigar {
	var ans []ByteCigar
	for _, i := range c {
		if i.Op == Insertion || i.Op == Deletion {
			ans = append(ans, i)
		}
	}
	return ans
}
