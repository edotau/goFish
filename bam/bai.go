// The following code was forked and modified from grailbio.
// You can follow the source source at: https://github.com/grailbio/bio/blob/d966d878d120dd59fd616ed94b5273c9b16f4309/encoding/bam/index.go
package bam

import (
	"bufio"

	"encoding/binary"
	"fmt"
	"io"
	"sort"

	"github.com/biogo/hts/bgzf"
	htslib "github.com/biogo/hts/bgzf"
	"github.com/edotau/goFish/simpleio"
)

var baiMagic [4]byte = ([4]byte{'B', 'A', 'I', 0x1})

// Bi represents the content of a .bai Bi file (for use with a .bam file).
type Bai struct {
	Magic         [4]byte
	Refs          []Reference
	UnmappedCount *uint64
}

// Reference represents the reference data within a .bai file.
type Reference struct {
	Bins      []bin
	Intervals []htslib.Offset
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
func IndexReader(filename string) *Bai {
	file := simpleio.NewReader(filename)
	return IndexBai(file.Reader)
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

// ReadIndex parses the content of r and returns an Index or nil and an error.
func IndexBai(reader io.Reader) *Bai {
	r := bufio.NewReaderSize(reader, 4<<20)
	i := &Bai{}
	var err error
	if _, err = io.ReadFull(r, i.Magic[0:]); err != nil {
		simpleio.StdError(err)
	}
	if i.Magic != [4]byte{'B', 'A', 'I', 0x1} {
		simpleio.StdError(fmt.Errorf("bam index invalid magic: %v", i.Magic))
	}

	var refCount int32
	if err = binary.Read(r, binary.LittleEndian, &refCount); err != nil {
		simpleio.StdError(err)
	}
	i.Refs = make([]Reference, refCount)

	// Read each Reference
	for refId := 0; int32(refId) < refCount; refId++ {
		// Read each Bin
		var binCount int32
		if err = binary.Read(r, binary.LittleEndian, &binCount); err != nil {
			simpleio.StdError(err)
		}
		ref := Reference{
			Bins: make([]bin, 0, binCount),
		}
		for b := 0; int32(b) < binCount; b++ {
			var binNum uint32
			if err = binary.Read(r, binary.LittleEndian, &binNum); err != nil {
				simpleio.StdError(err)
			}
			var chunkCount int32
			if err = binary.Read(r, binary.LittleEndian, &chunkCount); err != nil {
				simpleio.StdError(err)
			}

			bin := bin{
				binNum: binNum,
				Chunks: make([]Chunk, chunkCount),
			}

			// Read each Chunk
			bin.getNextChunk(r, chunkCount)

			if binNum == 37450 {
				// If we have a metadata chunk, put it in ref.Meta instead of ref.Bins.
				if len(bin.Chunks) != 2 {
					simpleio.StdError(fmt.Errorf("invalid metadata: chunk has %d chunks, should have 2", len(bin.Chunks)))
				}
				ref.Meta = Metadata{
					UnmappedBegin: fromOffset(bin.Chunks[0].Begin),
					UnmappedEnd:   fromOffset(bin.Chunks[0].End),
					MappedCount:   fromOffset(bin.Chunks[1].Begin),
					UnmappedCount: fromOffset(bin.Chunks[1].End),
				}
			} else {
				ref.Bins = append(ref.Bins, bin)
			}
		}

		// Read each Interval.
		var intervalCount int32
		if err = binary.Read(r, binary.LittleEndian, &intervalCount); err != nil {
			simpleio.StdError(err)
		}
		ref.Intervals = make([]bgzf.Offset, intervalCount)
		for inv := 0; int32(inv) < intervalCount; inv++ {
			var ioffset uint64
			if err = binary.Read(r, binary.LittleEndian, &ioffset); err != nil {
				simpleio.StdError(err)
			}
			ref.Intervals[inv] = toOffset(ioffset)
		}
		i.Refs[refId] = ref
	}

	var unmappedCount uint64
	if err = binary.Read(r, binary.LittleEndian, &unmappedCount); err == nil {
		i.UnmappedCount = &unmappedCount
	} else if err != nil && err != io.EOF {
		simpleio.StdError(err)
	}
	return i
}

func (block bin) getNextChunk(reader io.Reader, chunkCount int32) {
	var err error
	// Read each Chunk
	for c := 0; int32(c) < chunkCount; c++ {
		var beginOffset uint64
		if err = binary.Read(reader, binary.LittleEndian, &beginOffset); err != nil {
			simpleio.StdError(err)
		}
		var endOffset uint64
		if err = binary.Read(reader, binary.LittleEndian, &endOffset); err != nil {
			simpleio.StdError(err)
		}
		block.Chunks[c] = Chunk{
			Begin: toOffset(beginOffset),
			End:   toOffset(endOffset),
		}
	}

}

// AllOffsets returns a map of chunk offsets in the index file, it
// includes chunk begin locations, and interval locations.  The Key of
// the map is the Reference ID, and the value is a slice of
// bgzf.Offsets.  The return map will have an entry for every
// reference ID, even if the list of offsets is empty.
func (i *Bai) AllOffsets() map[int][]bgzf.Offset {
	m := make(map[int][]bgzf.Offset)
	for refId, ref := range i.Refs {
		m[refId] = make([]bgzf.Offset, 0)

		// Get the offsets for this ref.
		for _, bin := range ref.Bins {
			for _, chunk := range bin.Chunks {
				if chunk.Begin.File != 0 || chunk.Begin.Block != 0 {
					m[refId] = append(m[refId], chunk.Begin)
				}
			}
		}
		for _, interval := range ref.Intervals {
			if interval.File != 0 || interval.Block != 0 {
				m[refId] = append(m[refId], interval)
			}
		}

		// Sort the offsets
		sort.SliceStable(m[refId], func(i, j int) bool {
			c0 := m[refId][i]
			c1 := m[refId][j]
			if c0.File != c1.File {
				return c0.File < c1.File
			}
			return c0.Block < c1.Block
		})

		// Keep only unique offsets
		uniq := make([]bgzf.Offset, 0)
		previous := bgzf.Offset{-1, 0}
		for _, offset := range m[refId] {
			if offset != previous {
				uniq = append(uniq, offset)
				previous = offset
			}
		}
		m[refId] = uniq
	}
	return m
}
