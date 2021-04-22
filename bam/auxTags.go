package bam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

// ASCII is a printable ASCII character included in an Aux tag.
type ASCII byte

// Hex is a byte slice represented as a hex string in an Aux tag.
type Hex []byte

// Text is a byte slice represented as a string in an Aux tag.
type Text []byte

// An Aux represents an auxiliary data field from a SAM alignment record.
type Aux []byte

// A Tag represents an auxiliary or header tag label.
type Tag [2]byte

var (
	headerTag     = Tag{'H', 'D'}
	versionTag    = Tag{'V', 'N'}
	sortOrderTag  = Tag{'S', 'O'}
	groupOrderTag = Tag{'G', 'O'}

	refDictTag       = Tag{'S', 'Q'}
	refNameTag       = Tag{'S', 'N'}
	refLengthTag     = Tag{'L', 'N'}
	alternativeLocus = Tag{'A', 'H'}
	assemblyIDTag    = Tag{'A', 'S'}
	md5Tag           = Tag{'M', '5'}
	speciesTag       = Tag{'S', 'P'}
	uriTag           = Tag{'U', 'R'}

	readGroupTag    = Tag{'R', 'G'}
	centerTag       = Tag{'C', 'N'}
	descriptionTag  = Tag{'D', 'S'}
	dateTag         = Tag{'D', 'T'}
	flowOrderTag    = Tag{'F', 'O'}
	keySequenceTag  = Tag{'K', 'S'}
	libraryTag      = Tag{'L', 'B'}
	insertSizeTag   = Tag{'P', 'I'}
	platformTag     = Tag{'P', 'L'}
	platformUnitTag = Tag{'P', 'U'}
	sampleTag       = Tag{'S', 'M'}

	programTag      = Tag{'P', 'G'}
	idTag           = Tag{'I', 'D'}
	programNameTag  = Tag{'P', 'N'}
	commandLineTag  = Tag{'C', 'L'}
	previousProgTag = Tag{'P', 'P'}
	progDesc        = Tag{'D', 'S'}

	commentTag = Tag{'C', 'O'}
)

// Value returns v containing the value of the auxiliary tag.
func (a Aux) Value() interface{} {
	switch t := a.Type(); t {
	case 'A':
		return a[3]
	case 'c':
		return int8(a[3])
	case 'C':
		return uint8(a[3])
	case 's':
		return int16(binary.LittleEndian.Uint16(a[3:5]))
	case 'S':
		return binary.LittleEndian.Uint16(a[3:5])
	case 'i':
		return int32(binary.LittleEndian.Uint32(a[3:7]))
	case 'I':
		return binary.LittleEndian.Uint32(a[3:7])
	case 'f':
		return math.Float32frombits(binary.LittleEndian.Uint32(a[3:7]))
	case 'Z': // Z and H Require that parsing stops before the terminating zero.
		return string(a[3:])
	case 'H':
		return []byte(a[3:])
	case 'B':
		length := int32(binary.LittleEndian.Uint32(a[4:8]))
		switch t := a[3]; t {
		case 'c':
			c := a[8:]
			return *(*[]int8)(unsafe.Pointer(&c))
		case 'C':
			return []uint8(a[8:])
		case 's':
			Bs := make([]int16, length)
			err := binary.Read(bytes.NewBuffer(a[8:]), binary.LittleEndian, &Bs)
			if err != nil {
				panic(fmt.Sprintf("sam: binary.Read of s field failed: %v", err))
			}
			return Bs
		case 'S':
			BS := make([]uint16, length)
			err := binary.Read(bytes.NewBuffer(a[8:]), binary.LittleEndian, &BS)
			if err != nil {
				panic(fmt.Sprintf("sam: binary.Read of S field failed: %v", err))
			}
			return BS
		case 'i':
			Bi := make([]int32, length)
			err := binary.Read(bytes.NewBuffer(a[8:]), binary.LittleEndian, &Bi)
			if err != nil {
				panic(fmt.Sprintf("sam: binary.Read of i field failed: %v", err))
			}
			return Bi
		case 'I':
			BI := make([]uint32, length)
			err := binary.Read(bytes.NewBuffer(a[8:]), binary.LittleEndian, &BI)
			if err != nil {
				panic(fmt.Sprintf("sam: binary.Read of I field failed: %v", err))
			}
			return BI
		case 'f':
			Bf := make([]float32, length)
			err := binary.Read(bytes.NewBuffer(a[8:]), binary.LittleEndian, &Bf)
			if err != nil {
				panic(fmt.Sprintf("sam: binary.Read of f field failed: %v", err))
			}
			return Bf
		default:
			return fmt.Errorf("%%B!(UNKNOWN ARRAY type=%c)", t)
		}
	default:
		return fmt.Errorf("%%?!(UNKNOWN type=%c)", t)
	}
}

func (a Aux) Type() byte { return a[2] }
