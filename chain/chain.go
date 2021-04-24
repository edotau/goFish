// Package chain contains data structures and functions to process genome to genome alignments pairwise alignment that allow gaps in both sequences simultaneously. Chain files alignments starts with a header line, contains one or more alignment data lines, and terminates with a blank line
package chain

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

// Chain alignment fields.
type Chain struct {
	Score     int
	TName     string
	TSize     int
	TStrand   byte
	TStart    int
	TEnd      int
	QName     string
	QSize     int
	QStrand   byte
	QStart    int
	QEnd      int
	Alignment []Bases
	Id        int
}

// Bases is a cigar-like info for alignment block: First number is the length/size of bases, then number of target gaps and finally query gaps.
type Bases struct {
	Size   int
	TBases int
	QBases int
}

// NewChain will process text into chain data fields. It will read the first line of the file and assign to header fields and use a reader to read and process the additional lines of the alignment.
func NewChain(text string, reader *simpleio.SimpleReader) *Chain {
	data := strings.Split(text, " ")
	if len(data) == 13 {
		return &Chain{
			Score:     simpleio.StringToInt(data[1]),
			TName:     data[2],
			TSize:     simpleio.StringToInt(data[3]),
			TStrand:   data[4][0],
			TStart:    simpleio.StringToInt(data[5]),
			TEnd:      simpleio.StringToInt(data[6]),
			QName:     data[7],
			QSize:     simpleio.StringToInt(data[8]),
			QStrand:   data[9][0],
			QStart:    simpleio.StringToInt(data[10]),
			QEnd:      simpleio.StringToInt(data[11]),
			Alignment: chainingHelper(reader),
			Id:        simpleio.StringToInt(data[12]),
		}
	} else {
		log.Fatalf("Error: header line needs to contain 13 data fields\n")
		return nil
	}
}

// chainingHelper is the helper function that will process the chain alignment fields and return the alignment stats.
func chainingHelper(reader *simpleio.SimpleReader) []Bases {
	var line *bytes.Buffer
	var data []string
	var answer []Bases
	var curr Bases
	for nextBytes, done := reader.Peek(1); nextBytes[0] != 0 && done == nil; nextBytes, done = reader.Peek(1) {
		line, _ = simpleio.ReadLine(reader)
		data = strings.Split(line.String(), "\t")
		if len(data) == 1 {
			curr = Bases{
				Size:   simpleio.StringToInt(data[0]),
				TBases: 0,
				QBases: 0,
			}
			answer = append(answer, curr)
			//this will advance the reader to the blank line i beliebe the reader will peak at the blank line in the next iteration and exit
			line, _ = simpleio.ReadLine(reader)
			return answer
		} else if len(data) == 3 {
			curr = Bases{
				Size:   simpleio.StringToInt(data[0]),
				TBases: simpleio.StringToInt(data[1]),
				QBases: simpleio.StringToInt(data[2]),
			}
			answer = append(answer, curr)
		} else {
			log.Fatalf("Error: expecting alignment data columns to be 3 or 1 but encountered %d\n", len(data))
		}
	}
	return nil
}

// Chain struct implements the bed interface with the Chrom() method which returns the chromosome name referencing the target sequence.
func (c *Chain) Chrom() string {
	return c.TName
}

// Chain struct implements the bed interface with the ChrStart() method which returns the starting position of the region referencing the target sequence.
func (c *Chain) ChrStart() int {
	return c.TStart
}

// Chain struct implements the bed interface with the ChrEnd() method which returns the starting position of the region referencing the target sequence.
func (c *Chain) ChrEnd() int {
	return c.TEnd
}

// ReadHeaderComments will process header comments that sometimes appear at the beginning of chain file and returns a struct.
func ReadHeaderComments(er *simpleio.SimpleReader) *HeaderComments {
	var line *bytes.Buffer
	var commments HeaderComments
	for nextBytes, done := er.Peek(1); nextBytes[0] == '#' && done == nil; nextBytes, done = er.Peek(1) {
		line, _ = simpleio.ReadLine(er)
		commments.HashTag = append(commments.HashTag, line.String())
	}
	return &commments
}

// HeaderComments stores the comment lines at the beginning of chain alignments into a struct.
type HeaderComments struct {
	HashTag []string
}

// Read is a simple function to read in chain test files and covert into a slice of chain structs
func Read(filename string) []Chain {
	reader := simpleio.NewReader(filename)
	ReadHeaderComments(reader)
	var ans []Chain
	for i, done := ParseChain(reader); !done; i, done = ParseChain(reader) {
		ans = append(ans, *i)
	}
	return ans
}

// ReadAll is a function that will take a list of file names and append all data into a single slice of chain structs
func ReadAll(files []string) []Chain {
	var ans []Chain
	for _, each := range files {
		ans = append(ans, Read(each)...)
	}
	return ans
}

// NextChain will read lines in file and return one chain record at a time and a true false determining the EOF.
func ParseChain(reader *simpleio.SimpleReader) (*Chain, bool) {
	line, done := simpleio.ReadLine(reader)
	if !done {
		return NewChain(line.String(), reader), false
	} else {
		return nil, true
	}
}

// ToString will convert a chain struct to original string format.
func ToString(ch *Chain) string {
	var answer string = fmt.Sprintf("chain %d %s %d %c %d %d %s %d %c %d %d %d\n", ch.Score, ch.TName, ch.TSize, ch.TStrand, ch.TStart, ch.TEnd, ch.QName, ch.QSize, ch.QStrand, ch.QStart, ch.QEnd, ch.Id)
	//minus one in the loop because last line contains 2 zeros and we do not want to print those
	for i := 0; i < len(ch.Alignment)-1; i++ {
		answer += fmt.Sprintf("%d\t%d\t%d\n", ch.Alignment[i].Size, ch.Alignment[i].TBases, ch.Alignment[i].QBases)
	}
	answer = fmt.Sprintf("%s%d\n", answer, ch.Alignment[len(ch.Alignment)-1].Size)
	return answer
}

// PrettyFmt will summarize the chain header lines in a more human readable format
func PrettyFmt(c *Chain) string {
	return fmt.Sprintf("%s\t%s\t%d\t%d\t%s\t%s\t%d\t%d", c.TName, string(c.TStrand), c.TStart, c.TEnd, c.QName, string(c.QStrand), c.QStart, c.QEnd)
}
