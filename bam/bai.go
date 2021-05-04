// The following code was forked and modified from grailbio.
// You can follow the source source at: https://github.com/grailbio/bio/blob/d966d878d120dd59fd616ed94b5273c9b16f4309/encoding/bam/index.go
package bam

import (
	"github.com/edotau/goFish/simpleio"
	htslib "github.com/grailbio/bio/encoding/bam"
	"github.com/grailbio/hts/bgzf"
)

// Bi represents the content of a .bai Bi file (for use with a .bam file).
type Bai struct {
	Magic         [4]byte
	Refs          []Reference
	UnmappedCount *uint64
}

// Reference represents the reference data within a .bai file.
type Reference struct {
	Bins      []bin
	Intervals []bgzf.Offset
	Meta      Metadata
}

// Bin represents the bin data within a .bai file.
type bin struct {
	binNum uint32
	Chunks []Chunk
}

// Chunk represents the Chunk data within a .bai file.
type Chunk struct {
	Begin bgzf.Offset
	End   bgzf.Offset
}

// Metadata represents the Metadata data within a .bai file.
type Metadata struct {
	UnmappedBegin uint64
	UnmappedEnd   uint64
	MappedCount   uint64
	UnmappedCount uint64
}

// IndexReader decompress a binary index of r and returns an Index or nil and an error.
func IndexReader(filename string) *htslib.Bai {
	file := simpleio.NewReader(filename)
	reader, err := htslib.IndexReader(file)
	simpleio.ErrorHandle(err)
	return reader
}

func toOffset(voffset uint64) bgzf.Offset {
	return bgzf.Offset{
		File:  int64(voffset >> 16),
		Block: uint16(voffset),
	}
}

func fromOffset(offset bgzf.Offset) uint64 {
	return uint64(offset.File<<16) | uint64(offset.Block)
}
