// Package bigWig is used to process bigwig files which are binary compressed data
// to provide a visual of aligned sequencing data
package bigWig

import (
	"encoding/binary"
	"io"

	"github.com/edotau/goFish/simpleio"
)

// TODO not a completly finished project
const BigWigMagic = 0x888FFC26

type BigWigReader struct {
	io.ReadSeeker
	close func() error
}

type bwheader struct {
	version           uint16
	nLevels           uint16
	ctOffset          uint64
	dataOffset        uint64
	indexOffset       uint64
	fieldCount        uint16
	definedFieldCount uint16
	sqlOffset         uint64
	summaryOffset     uint64
	bufSize           uint32
	extensionOffset   uint64
	zoomHders         uint64 //bwZoomHdr_t * 	zoomHdrs
	nBasesCovered     uint64
	minVal            float64
	maxVal            float64
	sumData           float64
	sumSquared        float64
}

type zoomHeader struct {
	reductionLevel uint `bin:"len:4"`
	dataOffset     uint `bin:"len:8"`
	indexOffset    uint `bin:"len:8"`
}

type summary struct {
	basesCovered uint    `bin:"len:8"`
	minVal       float64 `bin:"len:8"`
	maxVal       float64 `bin:"len:8"`
	sumData      float64 `bin:"len:8"`
	sumSquares   float64 `bin:"len:8"`
}
type indexWig struct {
	chromId    uint `bin:"len:4"`
	chromStart uint `bin:"len:4"`
	chromEnd   uint `bin:"len:4"`
	itemStep   uint `bin:"len:4"`
	itemSpan   uint `bin:"len:4"`
	Type       uint `bin:"len:1"`
	itemCount  uint `bin:"len:2"`
}

type bgIndex struct {
	NumBlocks       uint64
	BlockSize       uint32
	Entries         uint64
	WidSum          uint64
	ID              uint32
	Start           uint32
	End             uint32
	Span            uint32
	Step            uint32
	Type            uint8
	BufSize         uint32
	NumNode         uint64
	Intervals       uint32
	MaxNumIntervals uint32
}

func NewReader(filename string) *BigWigReader {
	file := simpleio.Vim(filename)
	answer := BigWigReader{
		ReadSeeker: file,
		close:      file.Close,
	}
	return &answer
}

func (reader *BigWigReader) Close() error {
	if reader.close == nil {
		return nil
	} else {
		return reader.close()
	}
}

func MagicBigWig(reader *BigWigReader) bool {
	var magic uint32
	err := binary.Read(reader, binary.LittleEndian, &magic)
	simpleio.StdError(err)
	return magic == BigWigMagic
}

type BwSize struct {
	Block []byte
	Error error
}

/*
func ReadBlocks(reader *BigWigReader) <-chan BwSize {
	// create new channel
	channel := make(chan BwSize, 10)
	// fill channel with blocks
	go func() {
		reader.fillChannel(channel, reader.Reader.Index.Root)
		// close channel and file
		close(channel)
	}()
	return channel
}*/

type file struct {
	blockSize int
	items     int
	levels    []int
}

/*
func BwUnmarshal(reader io.ReadSeeker) (*BigWigReader, error) {
	bwr := new(BigWigReader)
	bwf := new(file)
	if err := bwf.Open(reader); err != nil {
		return nil, err
	}
	bwr.Reader = reader
	bwr.Bwf = *bwf

	seqnames := make([]string, len(bwf.ChromData.Keys))
	lengths := make([]int, len(bwf.ChromData.Keys))

	for i := 0; i < len(bwf.ChromData.Keys); i++ {
		if len(bwf.ChromData.Values[i]) != 8 {
			return nil, fmt.Errorf("invalid chromosome list")
		}
		idx := int(binary.LittleEndian.Uint32(bwf.ChromData.Values[i][0:4]))
		if idx >= len(bwf.ChromData.Keys) {
			return nil, fmt.Errorf("invalid chromosome index")
		}
		seqnames[idx] = strings.TrimRight(string(bwf.ChromData.Keys[i]), "\x00")
		lengths[idx] = int(binary.LittleEndian.Uint32(bwf.ChromData.Values[i][4:8]))
	}
	bwr.Genome = NewGenome(seqnames, lengths)

	return bwr, nil
}*/
