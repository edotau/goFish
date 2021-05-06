package geneSeq

import (
	"fmt"
	"log"
	"runtime"
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

// Threads for bread.NewBufferedReader()
var Threads = runtime.NumCPU()

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

	start, end := simpleio.StringToInt(col[3]), simpleio.StringToInt(col[4])
	if start > end {
		simpleio.StdError(fmt.Errorf("%s: start (%d) must be < end (%d)", col[0], start, end))
	}

	var score *float64

	if col[5] != "." {
		s := simpleio.StringToFloat64(col[5])
		score = &s
	}

	var strand *byte
	if col[6][0] != '.' {
		s := col[6][0]
		if !(s == '+' || s == '-') {
			log.Fatalf(fmt.Sprintf("%s: illigal strand: %v", col[0], s))
		}
		strand = &s
	}
	var frame *int
	if col[7] != "." {
		f := simpleio.StringToInt(col[7])
		log.Fatalf(fmt.Sprintf("%s: bad frame: %s", col[0], col[7]))

		if !(f == 0 || f == 1 || f == 2) {
			simpleio.StdError(fmt.Errorf("%s: illigal frame: %d", col[0], f))
		}
		frame = &f
	}

	geneFeature := Gtf{col[0], col[1], col[2], start, end, score, strand, frame, nil}

	tagValues := strings.Split(col[8], "; ")
	if len(tagValues) > 0 {
		geneFeature.Attributes = []Attribute{}
		for _, tagValue := range tagValues[0 : len(tagValues)-1] {
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

/*
type Gtf struct {
	Name       string
	Source     string
	Feature    string
	Start      int
	End        int
	Score      *float64
	Strand     *byte
	Frame      int
	Attributes []Attribute
}

// Attribute are tags contained in the 9th column. Because of the nested formating, information from this column is often lost during file conversions
type Attribute struct {
	Tag   string
	Value string
}

func GtfReader(filename string) []Gtf {
	var ans []Gtf
	reader := simpleio.NewReader(filename)

	for i, done := ParseGtfLine(reader); !done; i, done = ParseGtfLine(reader) {

		ans = append(ans, *i)
	}
	return ans
}

func ParseGtfLine(reader *simpleio.SimpleReader) (*Gtf, bool) {
	i, done := simpleio.ReadLine(reader)
	line := strings.TrimRight(i.String(), "\r\n")
	column := strings.Split(line, "\t")
	if done {
		return nil, false
	} else if len(column) == 0 || column[0] == "#" {
		return nil, false
	} else {

		curr := Gtf{
			Name:    column[0],
			Source:  column[1],
			Feature: column[2],
			//		Start:   simpleio.StringToInt(column[3]),
			//	End:     simpleio.StringToInt(column[4]),
		}

		if column[5] != "." {

			f := simpleio.StringToFloat64(column[5])
			curr.Score = &f
		}
		if column[6][0] != '.' {
			s := column[6][0]
			curr.Strand = &s
		}
		if column[7] != "." {

			curr.Frame = simpleio.StringToInt(column[7])

		}

		fmt.Printf("%s\n", curr.ToString())

		return &curr, false
	}
}

func ParseAttrib(col string) []Attribute {

	tagValues := strings.Split(col, "; ")
	Attributes := []Attribute{}
	if len(tagValues) > 0 {

		for _, tagValue := range tagValues[0 : len(tagValues)-1] {
			col2 := strings.SplitN(tagValue, " ", 2)
			tag := col2[0]
			value := col2[1]
			// if value[len(value)-1] == ';' {
			// 	value = value[0 : len(value)-1]
			// }
			if len(value) > 2 {
				value = value[1 : len(value)-1]
			} else {
				value = ""
			}
			Attributes = append(Attributes, Attribute{tag, value})
		}
	}
	return Attributes
}*/

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
	for _, i := range atb {
		str.WriteString(i.Tag)
		str.WriteByte(' ')
		str.WriteByte('"')
		str.WriteString(i.Value)
		str.WriteByte('"')
		str.WriteByte(';')
		str.WriteByte(' ')
	}
	return str.String()

}
