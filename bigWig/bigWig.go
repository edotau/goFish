// Package bigWig is used to process bigwig files which are binary compressed data
// to provide a visual of aligned sequencing data
package bigWig

import (
	"encoding/binary"
	"io"

	"github.com/edotau/goFish/simpleio"
)

const BigWigMagic = 0x888FFC26

type BigWigReader struct {
	io.ReadSeeker
	close func() error
}

func NewBigWigReader(filename string) *BigWigReader {
	file := simpleio.OpenFile(filename)
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
	simpleio.ErrorHandle(err)
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
