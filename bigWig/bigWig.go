package bigWig

import (
	"encoding/binary"
	"github.com/goFish/simpleio"
	"io"
	//"github.com/vertgenlab/gonomics/"
)

const biWigMagic = 0x888FFC26

type BigWigReader struct {
	io.ReadSeeker
	close func() error
}

type file struct {
	blockSize int
	items     int
	levels    []int
}

func OpenBigWig(filename string) *BigWigReader {
	file := simpleio.OpenFile(filename)
	answer := BigWigReader{
		ReadSeeker: file,
		close:      file.Close,
	}
	return &answer
}

func checkMagic(reader *BigWigReader) bool {
	var magic uint32
	err := binary.Read(reader, binary.LittleEndian, &magic)
	simpleio.ErrorHandle(err)
	return magic == biWigMagic
}
