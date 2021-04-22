package simpleio

import (
	"strconv"
	"strings"
)

type ByteCigar struct {
	RunLen uint32
	Op     byte
}

// CigarByte is a light weight cigar stuct
// that stores the runlength as an int (not int64) and Op as a byte.
type CigarByte struct {
	Len int
	Op  byte
}

type CigarOp byte

const (
	Match     byte = 'M'
	Insertion byte = 'I'
	Deletion  byte = 'D'
	N         byte = 'N'
	SoftClip  byte = 'S'
	HardClip  byte = 'H'
	Padded    byte = 'P'
	Equal     byte = '='
	Mismatch  byte = 'X'
	Unknown   byte = '*'
)

// ByteCigarToString will process the cigar byte struct and parse and/or convert the data into a string.
func BytesToCigar(cigar []byte) []ByteCigar {
	if cigar[0] == '*' {
		return nil
	}
	var ans []ByteCigar = make([]ByteCigar, 0, 1)
	var lastNum int = 0
	//var check CigarOp
	for i := 0; i < len(cigar); i++ {
		//check = LookUpByte[cigar[i]]
		if IsValidCigar(cigar[i]) {
			ans = append(ans, ByteCigar{RunLen: StringToUInt32(string(cigar[lastNum:i])), Op: cigar[i]})
			lastNum = i + 1
		}
	}
	return ans
}

func IsValidCigar(op byte) bool {
	switch op {
	case 'M':
		return true
	case 'I':
		return true
	case 'D':
		return true
	case 'N':
		return true
	case 'S':
		return true
	case 'H':
		return true
	case 'P':
		return true
	case '=':
		return true
	case 'X':
		return true
	case '*':
		return true
	default:
		return false
	}
	return false
}

// LookUpCigByte will return the number that maps to each cigar byte
// In the sam/bam specs: CIGAR: op len<<4|op. ‘MIDNSHP=X’→‘012345678’

//var LookUpCigar []CigarOp = []CigarOp{Match, Insertion, Deletion, N, SoftClip, HardClip, Padded, Equal, Mismatch}

func ByteCigarString(cigar []ByteCigar) string {
	if cigar == nil {
		return "*"
	}
	var str strings.Builder
	var err error
	for _, c := range cigar {
		_, err = str.WriteString(strconv.Itoa(int(c.RunLen)))
		ErrorHandle(err)
		err = str.WriteByte(c.Op)
		ErrorHandle(err)
	}
	return str.String()
}

func ByteMatrixTrace(a int64, b int64, c int64) (int64, byte) {
	if a >= b && a >= c {
		return a, 'M'
	} else if b >= c {
		return b, 'I'
	} else {
		return c, 'D'
	}
}

func ReverseBytesCigar(alpha []ByteCigar) {
	var i, off int
	for i, off = len(alpha)/2-1, len(alpha)-1-i; i >= 0; i, off = i-1, len(alpha)-1-i {
		alpha[i], alpha[off] = alpha[off], alpha[i]
	}
}

func ReversePath(alpha []uint32) {
	var i, off int
	for i, off = len(alpha)/2-1, len(alpha)-1-i; i >= 0; i-- {
		alpha[i], alpha[off] = alpha[off], alpha[i]
	}
}

//QueryLength calculates the length of the query read from a slice of Cigar structs.
func QueryRunLen(c []ByteCigar) int {
	if c[0].Op == '*' {
		return 0
	}
	var ans uint32
	for _, v := range c {
		if v.Op == Match || v.Op == Insertion || v.Op == SoftClip || v.Op == Equal || v.Op == Mismatch {
			ans = ans + v.RunLen
		}
	}
	return int(ans)
}

func CatByteCigar(cigs []ByteCigar, newCigs []ByteCigar) []ByteCigar {
	if len(newCigs) == 0 || newCigs == nil {
		return cigs
	} else if len(cigs) == 0 {
		return newCigs
	} else {
		cigs = AddCigar(cigs, newCigs[0])
		cigs = append(cigs, newCigs[1:]...)
		return cigs
	}
}

func AddCigar(cigs []ByteCigar, newCig ByteCigar) []ByteCigar {
	if len(cigs) == 0 {
		cigs = append(cigs, newCig)
	} else if cigs[len(cigs)-1].Op == newCig.Op {
		cigs[len(cigs)-1].RunLen += newCig.RunLen
	} else {
		cigs = append(cigs, newCig)
	}
	return cigs
}

func MatrixSetup(size int) ([][]int64, [][]byte) {
	m := make([][]int64, size)
	trace := make([][]byte, size)
	for idx := range m {
		m[idx] = make([]int64, size)
		trace[idx] = make([]byte, size)
	}
	return m, trace
}

var LookUpByte = map[byte]byte{
	'M': Match,
	'I': Insertion,
	'D': Deletion,
	'N': N,
	'S': SoftClip,
	'H': HardClip,
	'P': Padded,
	'=': Equal,
	'X': Mismatch,
	'*': Unknown,
}
