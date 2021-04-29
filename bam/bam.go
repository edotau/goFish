// Package bam is used to process binary alignment files and and decode data into human readable sam text.
package bam

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/simpleio"
)

// BamReader contains data fields used to read and process binary file.
type BamReader struct {
	File      *os.File
	Gunzip    *Bgzip
	data      []byte
	bytesRead int
	error     error
}

// BinaryDecoder: Uses the BAM format 4.2, defined in hts specs. Note: This is not a bam struct but place holders for pointers to decode binary file
type BinaryDecoder struct {
	RefID     int32
	Pos       int32
	Bai       uint16
	MapQ      uint8
	RNLength  uint8
	Flag      uint16
	NCigarOp  uint16
	LSeq      int32
	NextRefID int32
	NextPos   int32
	TLength   int32
	QName     string
	Cigar     []uint32
	Seq       []byte
	Qual      []byte
	Aux       []*BamAux
	BlockSize int32
}

// chrInfo is used to decode fields contained in the header
type chrInfo struct {
	text   string
	length int32
	numRef int32
}

// BamAux is a struct that organizes the extra tags at the end of sam/bam records
type BamAux struct {
	Tag   [2]byte
	Type  byte
	Value interface{}
}

// CigarByte is a light weight cigar stuct
// that stores the runlength as an int (not int64) and Op as a byte.
type CigarByte struct {
	Len int
	Op  byte
}

// Read will process a bam file and return a slice of sam records that were decoded from binary.
func Read(filename string) (*Header, <-chan Sam) {
	header := &Header{}
	sams := make(chan Sam)
	wg := sync.WaitGroup{}
	wg.Add(1)
	if strings.HasSuffix(filename, ".sam") {
		reader := NewSamReader(filename)
		header = ReadSamHeader(reader)
		go func() {
			for i, done := UnmarshalSam(reader); !done; i, done = UnmarshalSam(reader) {
				sams <- *i
			}
			wg.Done()
		}()
	}
	if strings.HasSuffix(filename, ".bam") {
		reader := NewBamReader(filename)
		header = ReadHeader(reader)
		binaryData := make(chan *BinaryDecoder)
		go BamToChannel(reader, binaryData)
		go func() {
			for each := range binaryData {
				sams <- *BamBlockToSam(header, each)
			}
			wg.Done()
		}()

	}
	go func() {
		wg.Wait()
		close(sams)
	}()
	return header, sams
}

func BasicRead(filename string) (*Header, []*Sam) {
	bamFile := NewBamReader(filename)
	defer bamFile.File.Close()
	h := ReadHeader(bamFile)
	binaryData := make(chan *BinaryDecoder)
	var ans []*Sam
	go BamToChannel(bamFile, binaryData)
	for each := range binaryData {
		ans = append(ans, BamBlockToSam(h, each))
	}
	return h, ans
}

// Read will process a bam file and return a slice of sam records that were decoded from binary.
func BamToSam(filename string) (*Header, <-chan Sam) {
	bamFile := NewBamReader(filename)
	h := ReadHeader(bamFile)
	binaryData := make(chan *BinaryDecoder)
	ans := make(chan Sam)
	var wg sync.WaitGroup
	wg.Add(1)
	go BamToChannel(bamFile, binaryData)

	go func() {
		for each := range binaryData {
			ans <- *BamBlockToSam(h, each)
		}
		wg.Done()
	}()
	go func() {
		wg.Wait()
		close(ans)
	}()
	return h, ans
}

func HttpReaderDev(url string) *Bgzip {
	resp, err := http.Get(url)
	simpleio.ErrorHandle(err)
	return NewBgzipReader(resp.Body)
}

// NewBamReader is similar to fileio.EasyVim/fileio.EasyReader
// which will allocate memory for the struct fields
// and is ready to start processing bam lines after calling this function.
func NewBamReader(filename string) *BamReader {
	var bamR *BamReader = &BamReader{}
	bamR.File = simpleio.Vim(filename)
	bamR.Gunzip = NewBgzipReader(bamR.File)
	return bamR
}

// ReadHeader will take a BamReader structure as an input
// performs a quick check to make sure the binary file is a valid bam
// then process header lines and returns a BamHeader (similar to samHeader).
func ReadHeader(reader *BamReader) *Header {
	bamHeader := MakeHeader()
	magic := make([]byte, 4)
	magic = BgzipBuffer(reader.Gunzip, magic)
	if string(magic) != "BAM\001" {
		log.Fatalf("Not a BAM file: %s\n", string(reader.data))
	}
	reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &bamHeader.bamDecoder.length)
	simpleio.ErrorHandle(reader.error)
	reader.data = make([]byte, bamHeader.bamDecoder.length)
	var i, j, k int = 0, 0, 0
	for i = 0; i < int(bamHeader.bamDecoder.length); {
		j, reader.error = reader.Gunzip.Read(reader.data[i:])
		simpleio.ErrorHandle(reader.error)
		i += j
	}
	bamHeader.Text.Write(reader.data)
	reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &bamHeader.bamDecoder.numRef)
	simpleio.ErrorHandle(reader.error)
	reader.data = make([]byte, bamHeader.bamDecoder.numRef)
	var lengthName, lengthSeq int32
	for i = 0; i < int(bamHeader.bamDecoder.numRef); i++ {
		lengthName, lengthSeq = 0, 0
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &lengthName)
		simpleio.ErrorHandle(reader.error)
		reader.data = make([]byte, lengthName)
		for j = 0; j < int(lengthName); {
			k, reader.error = reader.Gunzip.Read(reader.data[j:])
			simpleio.ErrorHandle(reader.error)
			j += k
		}
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &lengthSeq)
		simpleio.ErrorHandle(reader.error)
		bamHeader.Chroms = append(bamHeader.Chroms, ChromSize{Name: strings.Trim(string(reader.data), "\n\000"), Size: int(lengthSeq), Order: len(bamHeader.Chroms)})
	}
	return bamHeader
}

// MakeHeader allocates memory for new bam header.
func MakeHeader() *Header {
	return &Header{ChromSize: make(map[string]int), bamDecoder: &chrInfo{length: 0, numRef: 0}}
}

// BamToChannel is a goroutine that will use the binary reader to decode bam records
// and send them off to a channel that could be processed into a sam record further downstream.
func BamToChannel(reader *BamReader, binaryData chan<- *BinaryDecoder) {
	var blockSize int32
	var bitFlag uint32
	var stats uint32
	var i, j int
	var b byte
	var block *BinaryDecoder
	buf := bytes.NewBuffer([]byte{})
	// This loop will attempt to decode a bam block, which is equivalent to one sam line in this case
	for {
		block = &BinaryDecoder{}
		buf.Reset()
		// read block size
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &blockSize)
		if reader.error == io.EOF {
			close(binaryData)
			break
		}
		simpleio.ErrorHandle(reader.error)
		//read block data
		block.BlockSize = blockSize
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.RefID)
		simpleio.ErrorHandle(reader.error)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.Pos)
		simpleio.ErrorHandle(reader.error)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &stats)
		simpleio.ErrorHandle(reader.error)

		block.Bai = uint16((stats >> 16) & 0xffff)
		block.MapQ = uint8((stats >> 8) & 0xff)
		block.RNLength = uint8((stats >> 0) & 0xff)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &bitFlag)
		simpleio.ErrorHandle(reader.error)

		// get Flag and NCigarOp from bitFlag
		block.Flag = uint16(bitFlag >> 16)
		block.NCigarOp = uint16(bitFlag & 0xffff)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.LSeq)
		simpleio.ErrorHandle(reader.error)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.NextRefID)
		simpleio.ErrorHandle(reader.error)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.NextPos)
		simpleio.ErrorHandle(reader.error)

		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.TLength)
		simpleio.ErrorHandle(reader.error)

		// parse the read name
		block.QName = binaryByteToString(reader, b, buf)

		// parse cigar block
		block.Cigar = make([]uint32, block.NCigarOp)
		for i = 0; i < int(block.NCigarOp) && reader.error == nil; i++ {
			reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.Cigar[i])
			simpleio.ErrorHandle(reader.error)

		}
		// parse seq
		block.Seq = make([]byte, (block.LSeq+1)/2)
		for i = 0; i < int((block.LSeq+1)/2); i++ {
			reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.Seq[i])
			simpleio.ErrorHandle(reader.error)

		}
		block.Qual = make([]byte, block.LSeq)
		for i = 0; i < int(block.LSeq); i++ {
			reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &block.Qual[i])
			simpleio.ErrorHandle(reader.error)
		}
		// read auxiliary data
		j = 8*4 + int(block.RNLength) + 4*int(block.NCigarOp) + int((block.LSeq+1)/2) + int(block.LSeq)
		for i = 0; j+i < int(blockSize); {
			block.Aux = append(block.Aux, decodeAuxiliary(reader))
			i += reader.bytesRead
		}
		binaryData <- block
	}
}

// BamBlockToSam is a function that will convert a decoded
// (already processed binary structure) to a human readable sam data.
func BamBlockToSam(header *Header, bam *BinaryDecoder) *Sam {
	return &Sam{
		QName:   bam.QName,
		Flag:    bam.Flag,
		RName:   header.Chroms[bam.RefID].Name,
		Pos:     int(bam.Pos + 1),
		MapQ:    bam.MapQ,
		Cigar:   Uint32ToByteCigar(bam.Cigar),
		MateRef: setRNext(header, bam),
		MatePos: int(bam.NextPos + 1),
		TmpLen:  int(bam.TLength),
		Seq:     code.ToDna([]byte(BamSeq(bam.Seq))),
		Qual:    formatQual(bam.Qual),
		Aux:     auxToString(bam.Aux),
	}
}

// setRNext will process the reference name of the mate pair alignment,
//if the alignment is on the same fragment, then will set to "=".
func setRNext(header *Header, bam *BinaryDecoder) string {
	if bam.NextRefID == bam.RefID {
		return "="
	} else if bam.NextRefID > 0 {
		return header.Chroms[bam.NextRefID].Name //fmt.Sprintf("Chrom: %d", bam.NextRefID)
	} else {
		return "*"
	}
}

// BamSeq will convert raw bytes to a string which can be converted to dna.Base.
//TODO: Look into converting bytes straight to dna.Base
func BamSeq(seq []byte) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	t := []byte{'=', 'A', 'C', 'M', 'G', 'R', 'S', 'V', 'T', 'W', 'Y', 'H', 'K', 'D', 'B', 'N'}
	for i := 0; i < len(seq); i++ {
		b1 := seq[i] >> 4
		b2 := seq[i] & 0xf
		fmt.Fprintf(writer, "%c", t[b1])
		if b2 != 0 {
			fmt.Fprintf(writer, "%c", t[b2])
		}
	}
	writer.Flush()
	return buffer.String()
}

// qualToString will convert the aul bytes into a readable string.
func qualToString(qual []byte) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	fmt.Fprintf(writer, "%s", string(formatQual(qual)))
	writer.Flush()
	return buffer.String()
}

// formatQual is a helper function that will add the 33 offset to the ASCII values or set '*' if the qual scores do not exist in the bam.
func formatQual(q []byte) []byte {
	for _, v := range q {
		if v != 0xff {
			a := make([]byte, len(q))
			for i, p := range q {
				a[i] = p + 33
			}
			return a
		}
	}
	return []byte{'*'}
}

// axtToString will convert the sam/bam auxiliary struct into a human readable string.
func auxToString(aux []*BamAux) string {
	var ans []string
	for i := 0; i < len(aux); i++ {
		ans = append(ans, fmt.Sprintf("%c%c:%c:%v", aux[i].Tag[0], aux[i].Tag[1], aux[i].Type, aux[i].Value))
	}
	return strings.Join(ans, "\t")
}

// decodeAuxiliary will use the bam reader struct to decode binary text to sam auxilary fields.
//In giraf this is what we are calling notes.
//TODO: Look to synchronize auxilary and notes.
func decodeAuxiliary(reader *BamReader) *BamAux {
	aux := &BamAux{}
	var i int
	// number of read bytes
	reader.bytesRead = 0
	// read data
	reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &aux.Tag[0])
	simpleio.ErrorHandle(reader.error)

	reader.bytesRead += 1
	reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &aux.Tag[1])
	simpleio.ErrorHandle(reader.error)

	reader.bytesRead += 1
	reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &aux.Type)
	simpleio.ErrorHandle(reader.error)

	// three bytes read so far
	reader.bytesRead += 1
	// read value

	switch aux.Type {
	case 'A':
		value := byte(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)

		aux.Value = value
		reader.bytesRead += 1
	case 'c':
		value := int8(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)

		aux.Value = value
		aux.Type = 'i'
		reader.bytesRead += 1
	case 'C':
		value := uint8(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)

		aux.Value = value
		aux.Type = 'i'
		reader.bytesRead += 1
	case 's':
		value := int16(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)
		aux.Type = 'i'
		aux.Value = value
		reader.bytesRead += 2
	case 'S':
		value := uint16(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)
		aux.Type = 'i'
		aux.Value = value
		reader.bytesRead += 2
	case 'i':
		value := int32(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)
		aux.Type = 'i'
		aux.Value = value
		reader.bytesRead += 4
	case 'I':
		value := uint32(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)
		aux.Value = value
		reader.bytesRead += 4
	case 'f':
		value := float32(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)

		aux.Value = value
		reader.bytesRead += 4
	case 'd':
		value := float64(0)
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &value)
		simpleio.ErrorHandle(reader.error)

		aux.Value = value
		reader.bytesRead += 8
	case 'Z':
		var b byte
		var buffer bytes.Buffer
		aux.Value = binaryByteToString(reader, b, &buffer)
	case 'H':
		var b byte
		var buffer bytes.Buffer
		for {
			reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &b)
			simpleio.ErrorHandle(reader.error)
			reader.bytesRead += 1
			if b == 0 {
				break
			}
			fmt.Fprintf(&buffer, "%X", b)
		}
		aux.Value = buffer.String()
	case 'B':
		var t byte
		var k int32
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &t)
		simpleio.ErrorHandle(reader.error)

		reader.bytesRead += 1
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &k)
		simpleio.ErrorHandle(reader.error)

		reader.bytesRead += 4
		switch t {
		case 'c':
			data := make([]int32, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 1
			}
			simpleio.ErrorHandle(reader.error)
			aux.Value = data
		case 'C':
			data := make([]uint8, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 1
			}
			simpleio.ErrorHandle(reader.error)
			aux.Value = data
		case 's':
			data := make([]int16, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 2
			}
			simpleio.ErrorHandle(reader.error)
			aux.Value = data
		case 'S':
			data := make([]uint16, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 2
			}
			simpleio.ErrorHandle(reader.error)
			aux.Value = data
		case 'i':
			data := make([]int32, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 4
			}
			simpleio.ErrorHandle(reader.error)
			aux.Value = data
		case 'I':
			data := make([]uint32, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 4
			}
			aux.Value = data
		case 'f':
			data := make([]float32, k)
			for i, reader.error = 0, binary.Read(reader.Gunzip, binary.LittleEndian, &data[i]); i < int(k) && reader.error == nil; i++ {
				reader.bytesRead += 4
			}
			simpleio.ErrorHandle(reader.error)
			aux.Value = data
		default:
			log.Fatalf("Error: encountered unknown auxiliary array value type %c...\n", t)
		}

	default:
		log.Fatalf("Error: Found invalid auxiliary value type %c...\n", aux.Type)
	}
	return aux
}

// ParseCigar will parse and convert a slice of uint32 and return a slice of cigar bytes.
func ParseCigar(bamCigar []uint32) []CigarByte {
	var ans []CigarByte
	var n uint32
	var t byte
	for i := 0; i < len(bamCigar); i++ {
		n = bamCigar[i] >> 4
		t = LookUpCigByte(bamCigar[i] & 0xf)
		ans = append(ans, CigarByte{Op: t, Len: int(n)})
	}
	return ans
}

// binaryByteToString will decode a single binary byte and conver to a string
func binaryByteToString(reader *BamReader, b byte, buf *bytes.Buffer) string {
	for {
		reader.error = binary.Read(reader.Gunzip, binary.LittleEndian, &b)
		simpleio.ErrorHandle(reader.error)
		reader.bytesRead += 1
		if b == 0 {
			return buf.String()
		} else {
			buf.WriteByte(b)
		}
	}
}

/*
// TODO: Plans for the new SAM/BAM record.
type Record struct {
	QName      string
	Ref       string
	Pos       int
	MapQ      byte
	Cigar     []uint32
	Flags     uint16
	MateRef   string
	MatePos   int
	TempLen   int
	Seq       []dna.Base
	Qual      []byte
	AuxFields AuxFields
}*/
