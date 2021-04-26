package bed

import (
	"log"
	"strconv"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

func NewAnnoPeaks(filename string) *simpleio.SimpleReader {

	reader := simpleio.NewReader(filename)
	i, err := simpleio.ReadLine(reader)
	if !err {
		if !strings.HasPrefix(i.String(), "PeakID") {
			log.Fatalf("Error: file input may not be an annotated peak file...\n")
		}
	}
	return reader
}

// Five is a slightly more detailed than the basic Bed implementation and records the current chromosome, start and end coordinates
// a name, and a bed score
type Annoatation struct {
	Chr             string
	Start           int
	End             int
	Strand          byte
	Score           int
	Annotation      string
	Detailed        string
	DistanceToTSS   int
	NearestPromoter string
	GeneID          string
	NearestUnigene  string
}

// Five bed struct implements the bed interface with the Chrom() method which returns the chromosome name.
func (bed *Annoatation) Chrom() string {
	return bed.Chr
}

// Five bed struct implements the bed interface with the ChrStart() method which returns the starting position of the region.
func (bed *Annoatation) ChrStart() int {
	return bed.Start
}

// Five bed struct implements the bed interface with the ChrEnd() method which returns the ending position of the region.
func (bed *Annoatation) ChrEnd() int {
	return bed.End
}

// ToBedPlus will take a simpleReader as an input and returns a BedPlus struct.
func ToBedAnnotation(reader *simpleio.SimpleReader) (*Annoatation, bool) {
	curr, done := simpleio.ReadLine(reader)
	if !done {
		columns := strings.Split(curr.String(), "\t")
		answer := Annoatation{
			Chr:             columns[1],
			Start:           simpleio.StringToInt(columns[2]),
			End:             simpleio.StringToInt(columns[3]),
			Strand:          byte(columns[4][0]),
			Score:           simpleio.StringToInt(columns[5]),
			Annotation:      strings.ReplaceAll(columns[7], " ", "_"),
			DistanceToTSS:   simpleio.StringToInt(columns[9]),
			NearestPromoter: columns[10],
			GeneID:          columns[12],
			NearestUnigene:  columns[13],
		}
		return &answer, false
	} else {
		return nil, true
	}
}

func (b *Annoatation) ToString() string {
	str := strings.Builder{}
	str.WriteString(b.Chrom())
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrStart()))
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrEnd()))
	str.WriteByte('\t')
	str.WriteByte(b.Strand)
	str.WriteByte('\t')
	str.WriteString(b.Annotation)
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.DistanceToTSS))
	str.WriteByte('\t')
	str.WriteString(b.NearestPromoter)
	str.WriteByte('\t')
	str.WriteString(b.GeneID)
	str.WriteByte('\t')
	str.WriteString(b.NearestUnigene)
	return str.String()
}
