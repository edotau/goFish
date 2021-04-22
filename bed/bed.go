// Package bed provides an interface along with functions to compute different Browser Extinsible Data (BED) formated files.
// Other data types can implement the bed interface by implementing the Chrom(), ChrStart(), ChrEnd() methods.
package bed

import (
	"bytes"
	"github.com/goFish/simpleio"
	"github.com/vertgenlab/gonomics/numbers"
	"log"
	"sort"
	"strconv"
	"strings"
)

// Bed interface is a data structure that implements the following methods:
// Chrom() []byte returns the chromosome name
// ChrStart() int returns starting the position
// ChrEnd() int returns ending the position
type Bed interface {
	Chrom() string
	ChrStart() int
	ChrEnd() int
	//ToString() string
}

// Simple is the most basic Bed implementation and only records the current chromosome, start and end coordinates.
type Simple struct {
	Chr   string
	Start int
	End   int
}

// Overlap will compare two bed regions and check if there is any overlap.
func Overlap(alpha Bed, beta Bed) bool {
	if (numbers.Max(alpha.ChrStart(), beta.ChrStart()) < numbers.Min(alpha.ChrEnd(), beta.ChrEnd())) && alpha.Chrom() == beta.Chrom() {
		return true
	} else {
		return false
	}
}

// overlap is the helper function used as an attempt to reduce method calls to the interface amd attempts to reduce duplication of code for different overlap settings.
func Overlapping(alphaChr string, betaChr string, alphaStart int, alphaEnd int, betaStart int, betaEnd int) bool {
	if (numbers.Max(alphaStart, betaStart) < numbers.Min(alphaEnd, betaEnd)) && alphaChr == betaChr {
		return true
	} else {
		return false
	}
}

type GenomeInfo struct {
	Chr   string
	Start int
	End   int
	Info  bytes.Buffer
}

// GenomeInfo bed struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (b *GenomeInfo) Chrom() string {
	return b.Chr
}

// GenomeInfo bed struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (b *GenomeInfo) ChrStart() int {
	return b.Start
}

// GenomeInfo bed struct implements the bed interface with the ChrEnd() method which returns the starting position of the region.
func (b *GenomeInfo) ChrEnd() int {
	return b.End
}

func GenomeInfoToString(b GenomeInfo) string {
	buf := strings.Builder{}
	buf.WriteString(b.Chr)
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(b.ChrStart()))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(b.ChrEnd()))
	if b.Info.String() != "" {
		buf.WriteByte('\t')
		buf.Write(b.Info.Bytes())
	}
	return buf.String()
}

func ToGenomeInfo(reader *simpleio.SimpleReader) (*GenomeInfo, bool) {
	var err bool
	reader.Buffer, err = simpleio.ReadLine(reader)

	if !err {
		columns := strings.SplitN(reader.Buffer.String(), "\t", 4)
		ans := GenomeInfo{}

		ans.Chr = columns[0]
		ans.Start = simpleio.StringToInt(columns[1])
		ans.End = simpleio.StringToInt(columns[2])
		if len(columns) > 3 {
			ans.Info.WriteString(columns[3])
		}
		return &ans, false
	} else {
		return nil, true
	}
}

func ReadHeader(reader *simpleio.SimpleReader) *strings.Builder {
	header := &strings.Builder{}
	for i, done := ParseComments(reader); !done; i, done = ParseComments(reader) {
		header.Write(i.Bytes())
	}
	return header
}

func ParseComments(reader *simpleio.SimpleReader) (*bytes.Buffer, bool) {
	if b, err := reader.Peek(1); err == nil && b[0] == byte('#') {
		return simpleio.ReadLine(reader)
	} else {
		return nil, true
	}
}

// HeadOverlapByLen checks for overlap while modifying the starting coordinates.
func HeadOverlapByLen(alpha Bed, beta Bed, length int) bool {
	var alphaStart, betaStart int = alpha.ChrStart() - length, beta.ChrStart() - length

	if alphaStart < 0 {
		alphaStart = 0
	}

	if betaStart < 0 {
		betaStart = 0
	}

	return Overlapping(alpha.Chrom(), beta.Chrom(), alphaStart, alpha.ChrEnd(), betaStart, beta.ChrEnd())
}

// TailOverlapByLen checks for overlap while modifying the ending coordinates.
func TailOverlapByLen(alpha Bed, beta Bed, alphaSize int, betaSize int, length int) bool {
	var alphaEnd, betaEnd int = alpha.ChrStart() - length, beta.ChrStart() - length

	if alphaEnd > alphaSize {
		alphaEnd = alphaSize
	}

	if betaEnd > betaSize {
		betaEnd = betaSize
	}

	return Overlapping(alpha.Chrom(), beta.Chrom(), alpha.ChrStart(), alphaEnd, beta.ChrStart(), betaEnd)
}

// ToSimpleBed will take a simpleReader and return a simpleBed struct
func ToSimpleBed(reader *simpleio.SimpleReader) (*Simple, bool) {
	curr, done := simpleio.ReadLine(reader)
	if !done {
		columns := bytes.Split(curr.Bytes(), []byte{'\t'})
		answer := Simple{
			Chr:   string(columns[0]),
			Start: simpleio.StringToInt(string(columns[1])),
			End:   simpleio.StringToInt(string(columns[2])),
		}
		return &answer, false
	} else {
		return nil, true
	}
}

// Simple bed struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (bed *Simple) Chrom() string {
	return bed.Chr
}

// Simple bed struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (bed *Simple) ChrStart() int {
	return bed.Start
}

// Simple bed struct implements the bed interface with the ChrEnd() method which returns the starting position of the region.
func (bed *Simple) ChrEnd() int {
	return bed.End
}

// Five is a slightly more detailed than the basic Bed implementation and records the current chromosome, start and end coordinates
// a name, and a bed score
type Five struct {
	Chr   string
	Start int
	End   int
	Name  string
	Score int
}

// Five bed struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (bed *Five) Chrom() string {
	return bed.Chr
}

// Five bed struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (bed *Five) ChrStart() int {
	return bed.Start
}

// Five bed struct implements the bed interface with the ChrEnd() method which returns the ending position of the region.
func (bed *Five) ChrEnd() int {
	return bed.End
}

// Six is a slightly more detailed than the basic Bed implementation and records the current chromosome, start and end coordinates
// Similar to five, but includes the strand information.
type Six struct {
	Chr    string
	Start  int
	End    int
	Name   string
	Score  int
	Strand bool
}

// Six bed struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (bed *Six) Chrom() string {
	return bed.Chr
}

// Six bed struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (bed *Six) ChrStart() int {
	return bed.Start
}

// Six bed struct implements the bed interface with the ChrEnd() method which returns the ending position of the region.
func (bed *Six) ChrEnd() int {
	return bed.End
}

// BedPlus is the most detailed of the bed formats and will store all information found in bed formatted files.
// TODO: Consider bed 12 plus format (very useful for working with gtf and gff files)
type BedPlus struct {
	Chr   string
	Start int
	End   int
	Name  string
	Score int
	Info  string
}

// BedPlus struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (bed *BedPlus) Chrom() string {
	return bed.Chr
}

// BedPlus struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (bed *BedPlus) ChrStart() int {
	return bed.Start
}

// BedPlus struct implements the bed interface with the ChrEnd() method which returns the ending position of the region.
func (bed *BedPlus) ChrEnd() int {
	return bed.End
}

// ToBedPlus will take a simpleReader as an input and returns a BedPlus struct.
func ToBedPlus(reader *simpleio.SimpleReader) (*BedPlus, bool) {
	curr, done := simpleio.ReadLine(reader)
	if !done {
		columns := bytes.SplitN(curr.Bytes(), []byte{'\t'}, 6)
		answer := BedPlus{
			Chr:   string(columns[0]),
			Start: simpleio.StringToInt(string(columns[1])),
			End:   simpleio.StringToInt(string(columns[2])),
			Name:  string(columns[3]),
			Score: simpleio.StringToInt(string(columns[4])),
			Info:  string(columns[5]),
		}
		return &answer, false
	} else {
		return nil, true
	}
}

func (b *BedPlus) String() string {
	str := strings.Builder{}
	str.WriteString(b.Chrom())
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrStart()))
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrEnd()))
	str.WriteByte('\t')
	str.WriteString(b.Name)
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.Score))
	str.WriteByte('\t')
	str.WriteString(b.Info)
	return str.String()
}

// Pvalue struct is a simple bed interface with an added PValue field
type Pvalue struct {
	Chr    string
	Start  int
	End    int
	Name   string
	PValue float64
}

// Pvalue struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (bed *Pvalue) Chrom() string {
	return bed.Chr
}

// Pvalue struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (bed *Pvalue) ChrStart() int {
	return bed.Start
}

// Pvalue struct implements the bed interface with the ChrEnd() method which returns the starting position of the region.
func (bed *Pvalue) ChrEnd() int {
	return bed.End
}

// comaarePValue is a helper function to compare pvalues between two Pvalue beds used in sorting slices.
func comparePValue(a *Pvalue, b *Pvalue) int {
	if a.PValue < b.PValue {
		return -1
	}
	if a.PValue > b.PValue {
		return 1
	}
	return 0
}

// SortByPValue performs a Pvalue sort to find the smallest and most significant.
func SortByPValue(peak []*Pvalue) {
	sort.Slice(peak, func(i, j int) bool { return comparePValue(peak[i], peak[j]) == -1 })
}

// ToBedPValue will convert bytes read from simpleReader and return a bed Pvalue.
func ToBedPValue(reader *simpleio.SimpleReader) (*Pvalue, bool) {
	curr, done := simpleio.ReadLine(reader)
	if !done {
		columns := bytes.SplitN(curr.Bytes(), []byte{'\t'}, 6)
		answer := Pvalue{
			Chr:    string(columns[0]),
			Start:  simpleio.StringToInt(string(columns[1])),
			End:    simpleio.StringToInt(string(columns[2])),
			Name:   string(columns[3]),
			PValue: float64(simpleio.StringToFloat(string(columns[6]))),
		}
		return &answer, false
	} else {
		return nil, true
	}
}

//To string will cast a bed interface and return a simple string with 3 fields.
func ToString(b Bed) string {
	var str strings.Builder
	str.WriteString(b.Chrom())
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrStart()))
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrEnd()))
	return str.String()
}

func ReadBed(filename string) []Simple {
	reader := simpleio.NewReader(filename)
	var ans []Simple
	for i, done := SimpleLine(reader); !done; i, done = SimpleLine(reader) {
		ans = append(ans, *i)
	}
	return ans
}

func SimpleLine(reader *simpleio.SimpleReader) (*Simple, bool) {
	buffer, done := simpleio.ReadLine(reader)
	if !done {
		fields := strings.Split(buffer.String(), "\t")
		return &Simple{
			Chr:   fields[0],
			Start: simpleio.StringToInt(fields[1]),
			End:   simpleio.StringToInt(fields[2]),
		}, false
	} else {
		return nil, true
	}

}

func PeakBedReading(reader *simpleio.SimpleReader) (*BedPlus, bool) {
	var err bool
	reader.Buffer, err = simpleio.ReadLine(reader)
	if !err {
		columns := strings.SplitN(reader.Buffer.String(), "\t", 6)
		ans := BedPlus{}

		ans.Chr = columns[0]
		ans.Start = simpleio.StringToInt(columns[1])
		ans.End = simpleio.StringToInt(columns[2])
		ans.Name = columns[3]
		ans.Score = simpleio.StringToInt(columns[4])
		if ans.Score > 1000 {
			log.Fatalf("Error: bed scores should not be greater than 1000...\n")
		}

		ans.Info = columns[5]
		return &ans, false
	} else {
		return nil, true
	}
}

//func WriteBed(b Bed)
