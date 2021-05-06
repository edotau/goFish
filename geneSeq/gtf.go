package geneSeq

import (
	"fmt"
	"log"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

// Version is the GTF version
const Version = 2.2

type Gtf struct {
	SeqName    string
	Source     string
	Feature    string
	Start      int
	End        int
	Score      *float64
	Strand     *byte
	Frame      *int
	Attributes []Attribute
}

// Attribute is the attribute
type Attribute struct {
	Tag   string
	Value string
}

func ReadGtf(filename string) []Gtf {
	reader := simpleio.NewReader(filename)
	var ans []Gtf
	for i, done := ParseNextGtf(reader); !done; i, done = ParseNextGtf(reader) {
		ans = append(ans, *i)
	}
	return ans
}

func ParseNextGtf(reader *simpleio.SimpleReader) (*Gtf, bool) {
	line, done := simpleio.ReadLine(reader)
	if done {
		return nil, done
	}
	col := strings.Split(line.String(), "\t")
	if len(col) < 9 {
		if strings.HasPrefix(col[0], "#") {
			return ParseNextGtf(reader)
		}
		return nil, false
	}
	start, end := simpleio.StringToInt(col[3]), simpleio.StringToInt(col[4])
	if start > end {
		simpleio.StdError(fmt.Errorf("%s: start (%d) must be < end (%d)", col[0], start, end))
	}
	geneFeature := Gtf{col[0], col[1], col[2], start, end, nil, nil, nil, nil}

	//var score *float64

	if col[5] != "." {
		score := simpleio.StringToFloat64(col[5])
		geneFeature.Score = &score
	}

	if col[6][0] != '.' {
		strand := col[6][0]
		if !(strand == '+' || strand == '-') {
			log.Fatalf("%s: illigal strand: %v", col[0], strand)
		}
		geneFeature.Strand = &strand
	}

	if col[7] != "." {
		frame := simpleio.StringToInt(col[7])

		if !(frame == 0 || frame == 1 || frame == 2) {
			log.Fatalf("%s: illigal frame: %d", col[0], frame)
		}
		geneFeature.Frame = &frame
	}

	tagValues := strings.Split(col[8], "; ")
	if len(tagValues) > 0 {
		geneFeature.Attributes = []Attribute{}
		for _, tagValue := range tagValues {
			col2 := strings.SplitN(tagValue, " ", 2)
			tag := col2[0]

			value := col2[1]

			if len(value) > 2 {
				value = value[1 : len(value)-1]
			} else {
				value = ""
			}
			geneFeature.Attributes = append(geneFeature.Attributes, Attribute{tag, value})
		}
	}
	return &geneFeature, false
}

func (gf *Gtf) ToString() string {
	var str strings.Builder
	str.WriteString(gf.SeqName)
	str.WriteByte('\t')
	str.WriteString(gf.Source)
	str.WriteByte('\t')
	str.WriteString(gf.Feature)
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gf.Start))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gf.End))
	str.WriteByte('\t')
	if gf.Score == nil {
		str.WriteByte('.')
	} else {
		str.WriteString(simpleio.Float64ToString(*gf.Score))
	}
	str.WriteByte('\t')
	if gf.Strand == nil {
		str.WriteByte('.')
	} else {
		str.WriteByte(*gf.Strand)
	}
	str.WriteByte('\t')
	if gf.Frame == nil {
		str.WriteByte('.')
	} else {
		str.WriteString(simpleio.IntToString(*gf.Frame))
	}

	str.WriteByte('\t')
	str.WriteString(AttribToString(gf.Attributes))
	return str.String()
}

func AttribToString(atb []Attribute) string {
	var str strings.Builder
	for i := 0; i < len(atb)-1; i++ {
		str.WriteString(atb[i].Tag)
		str.WriteByte(' ')
		str.WriteByte('"')
		str.WriteString(atb[i].Value)
		str.WriteByte('"')
		str.WriteByte(';')
		str.WriteByte(' ')
	}
	str.WriteString(atb[len(atb)-1].Tag)
	str.WriteByte(' ')
	str.WriteByte('"')
	str.WriteString(atb[len(atb)-1].Value)
	str.WriteByte(';')

	return str.String()

}
