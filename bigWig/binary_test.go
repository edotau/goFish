package bigWig

import (
	"encoding/binary"
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/edotau/goFish/simpleio"
	"github.com/ghostiam/binstruct"
)

var smallWig Wig = Wig{
	StepType: "fixedStep",
	Chrom:    "test1",
	Start:    1,
	Step:     10,
	Span:     10,
	Val:      []float64{0.1, 1.2, 2.3, 3.4, 4.5, 5.6, 6.7, 7.8, 8.9, 9},
}

type bw struct {
	chromId    uint `bin:"len:4"`
	chromStart uint `bin:"len:4"`
	chromEnd   uint `bin:"len:4"`
	itemStep   uint `bin:"len:4"`
	itemSpan   uint `bin:"len:4"`
	stepType   uint `bin:"len:1"`
	reserve    uint `bin:"-"`
	itemCount  uint `bin:"len:2"`
	//start     uint   `bin:"len:4"`
	//end       uint   `bin:"len:4"`
	//itemStep  uint   `bin:"len:4"`
	//itemSpan  uint   `bin:"len:4"` //Number of bases in item in fixedStep and varStep sections.
	//stepType  uint   `bin:"len:1"` //Section type. 1 for bedGraph, 2 for varStep, 3 for fixedStep.
	//itemCount uint   `bin:"len:2"`
}

func TestBigWigReadAll(t *testing.T) {

	file, err := os.Open("testdata/track.bw")
	var magic uint32
	err = binary.Read(file, binary.LittleEndian, &magic)
	simpleio.ErrorHandle(err)
	if magic != BigWigMagic {
		t.Errorf("Error: this is not a bigwig file...\n")
	}

	if err != nil {
		log.Fatal(err)
	}
	var wig bwheader
	decoder := binstruct.NewDecoder(file, binary.LittleEndian)
	err = decoder.Decode(&wig)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(wig)
}
