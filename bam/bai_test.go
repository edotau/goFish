package bam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/edotau/goFish/simpleio"
	"github.com/grailbio/testutil/expect"
	"github.com/stretchr/testify/assert"
)

func TestReadIndex(t *testing.T) {
	reader := IndexReader("testdata/tenXbarcodeTest.bam.bai")
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		fmt.Printf("%s\n", i.String())
	}

}

func toInt(t *testing.T, s string) int {
	i, err := strconv.Atoi(s)
	assert.Nil(t, err)
	return i
}

func writeBin(t *testing.T, w io.Writer, s string) {
	bins := strings.Split(s, ":")
	// Write the number of bins
	err := binary.Write(w, binary.LittleEndian, int32(len(bins)))
	expect.Nil(t, err)

	for _, bin := range bins {
		binContent := strings.Split(bin, ",")

		// Write the bin number
		err = binary.Write(w, binary.LittleEndian, assert.Nil(t, uint32(simpleio.StringToInt(binContent[0]))))
		expect.Nil(t, err)
		binContent = binContent[1:]

		// Write the number of chunks
		err = binary.Write(w, binary.LittleEndian, int32(len(binContent)/2))
		expect.Nil(t, err)

		// Write the chunks
		for _, voffset := range binContent {
			err = binary.Write(w, binary.LittleEndian, assert.Nil(t, uint32(simpleio.StringToInt(voffset))))
			expect.Nil(t, err)
		}
	}
}

func writeIntervals(t *testing.T, w io.Writer, s string) {
	intervals := strings.Split(s, ",")

	// Write the number of intervals
	err := binary.Write(w, binary.LittleEndian, int32(len(intervals)))
	expect.Nil(t, err)

	for _, voffset := range intervals {
		err = binary.Write(w, binary.LittleEndian, assert.Nil(t, uint64(simpleio.StringToInt(voffset))))
		expect.Nil(t, err)
	}
}

func writeIndex(t *testing.T, bins, intervals []string, unmapped int) *bytes.Buffer {
	var buf bytes.Buffer
	magic := []byte{'B', 'A', 'I', 0x1}
	_, err := buf.Write(magic)
	expect.Nil(t, err)

	// Two references
	err = binary.Write(&buf, binary.LittleEndian, int32(len(bins)))
	expect.Nil(t, err)

	for i := range bins {
		writeBin(t, &buf, bins[i])
		writeIntervals(t, &buf, intervals[i])
	}

	// Write unmapped count
	if unmapped >= 0 {
		err = binary.Write(&buf, binary.LittleEndian, uint64(unmapped))
		expect.Nil(t, err)
	}

	return &buf
}
