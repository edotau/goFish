package geneSeq

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/edotau/goFish/simpleio"
	"github.com/shenwei356/breader"
)

// Version is the GTF version
const Version = 2.2

// Feature is the gff feature struct
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

// ReadFeatures returns gtf features of a file
func ParseGtfLine(file string) ([]Gtf, error) {
	return ReadFilteredFeatures(file, []string{}, []string{}, []string{})
}

// ReadFilteredFeatures returns gtf features of specific chrs in a file
func ReadFilteredFeatures(file string, chrs []string, feats []string, attrs []string) ([]Gtf, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}
	chrsMap := make(map[string]struct{}, len(chrs))
	for _, chr := range chrs {
		chrsMap[strings.ToLower(chr)] = struct{}{}
	}

	featsMap := make(map[string]struct{}, len(feats))
	for _, f := range feats {
		featsMap[strings.ToLower(f)] = struct{}{}
	}

	attrsMap := make(map[string]struct{}, len(attrs))
	for _, f := range attrs {
		attrsMap[strings.ToLower(f)] = struct{}{}
	}

	fn := func(line string) (interface{}, bool, error) {
		if len(line) == 0 || line[0] == '#' {
			return nil, false, nil
		}
		line = strings.TrimRight(line, "\r\n")
		items := strings.Split(line, "\t")

		if len(items) != 9 {
			return nil, false, nil
		}

		if len(chrs) > 0 { // selected chrs
			if _, ok := chrsMap[strings.ToLower(items[0])]; !ok {
				return nil, false, nil
			}
		}

		if len(feats) > 0 { // selected features
			if _, ok := featsMap[strings.ToLower(items[2])]; !ok {
				return nil, false, nil
			}
		}
		var err error

		start, err := strconv.Atoi(items[3])
		if err != nil {
			return nil, false, fmt.Errorf("%s: bad start: %s", items[0], items[3])
		}

		end, err := strconv.Atoi(items[4])
		if err != nil {
			return nil, false, fmt.Errorf("%s: bad end: %s", items[0], items[4])
		}

		if start > end {
			return nil, false, fmt.Errorf("%s: start (%d) must be < end (%d)", items[0], start, end)
		}

		var score *float64
		if items[5] != "." {
			s, err := strconv.ParseFloat(items[5], 64)
			if err != nil {
				return nil, false, fmt.Errorf("%s: bad score: %s", items[0], items[5])
			}
			score = &s
		}

		var strand *byte
		if items[6][0] != '.' {
			s := items[6][0]
			if !(s == '+' || s == '-') {
				return nil, false, fmt.Errorf("%s: illigal strand: %v", items[0], s)
			}

			strand = &s
		}

		var frame *int
		if items[7] != "." {
			f, err := strconv.Atoi(items[7])
			if err != nil {
				return nil, false, fmt.Errorf("%s: bad frame: %s", items[0], items[7])
			}
			if !(f == 0 || f == 1 || f == 2) {
				return nil, false, fmt.Errorf("%s: illigal frame: %d", items[0], f)
			}
			frame = &f
		}

		feature := Gtf{items[0], items[1], items[2], start, end, score, strand, frame, nil}

		tagValues := strings.Split(items[8], "; ")
		if len(tagValues) > 0 {
			var ok bool
			feature.Attributes = []Attribute{}
			for _, tagValue := range tagValues[0 : len(tagValues)-1] {
				items2 := strings.SplitN(tagValue, " ", 2)
				tag := items2[0]
				if _, ok = attrsMap[tag]; !ok {
					continue
				}
				value := items2[1]
				// if value[len(value)-1] == ';' {
				// 	value = value[0 : len(value)-1]
				// }
				if len(value) > 2 {
					value = value[1 : len(value)-1]
				} else {
					value = ""
				}
				feature.Attributes = append(feature.Attributes, Attribute{tag, value})
			}
		}
		return feature, true, nil
	}
	reader, err := breader.NewBufferedReader(file, Threads, 100, fn)
	if err != nil {
		return nil, err
	}
	features := []Gtf{}
	for chunk := range reader.Ch {
		if chunk.Err != nil {
			return nil, chunk.Err
		}
		for _, data := range chunk.Data {
			features = append(features, data.(Gtf))
		}
	}
	return features, nil
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
			items2 := strings.SplitN(tagValue, " ", 2)
			tag := items2[0]
			value := items2[1]
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
		str.WriteByte('\t')
	} else {
		str.WriteString(simpleio.Float64ToString(*gf.Score))
	}
	if gf.Strand == nil {
		str.WriteByte('.')
		str.WriteByte('\t')
	} else {
		str.WriteByte(*gf.Strand)
	}
	if gf.Frame == nil {
		str.WriteByte('.')
		str.WriteByte('\t')
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
