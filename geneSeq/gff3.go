package geneSeq

import (
	"log"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

type Gff3 struct {
	SeqName    string
	Source     string
	Type       string
	Start      int
	End        int
	Score      *float64
	Strand     *byte
	Phase      *int
	Attributes []Attribute
}

func ReadGff3(filename string) []Gff3 {
	var ans []Gff3
	reader := simpleio.NewReader(filename)
	for i, done := ParseGtff3(reader); !done; i, done = ParseGtff3(reader) {
		ans = append(ans, *i)
	}

	return ans
}

func ParseGtff3(reader *simpleio.SimpleReader) (*Gff3, bool) {
	line, done := simpleio.ReadLine(reader)
	if done {
		return nil, done
	}
	col := strings.Split(line.String(), "\t")
	if len(col) < 9 {
		if strings.HasPrefix(col[0], "#") {
			return ParseGtff3(reader)
		} else {
			return nil, true
		}
	}
	start, end := simpleio.StringToInt(col[3]), simpleio.StringToInt(col[4])
	if start > end {
		log.Fatalf("%s: start (%d) must be < end (%d)", col[0], start, end)
	}
	featFmtThree := Gff3{col[0], col[1], col[2], start, end, nil, nil, nil, nil}
	if col[5] != "." {
		score := simpleio.StringToFloat64(col[5])
		featFmtThree.Score = &score
	}

	if col[6][0] != '.' {
		strand := col[6][0]
		if !(strand == '+' || strand == '-') {
			log.Fatalf("%s: illigal strand: %v", col[0], strand)
		}
		featFmtThree.Strand = &strand
	}

	if col[7] != "." {
		phase := simpleio.StringToInt(col[7])

		if !(phase == 0 || phase == 1 || phase == 2) {
			log.Fatalf("%s: illigal frame: %d", col[0], phase)
		}
		featFmtThree.Phase = &phase
	}

	notes := strings.Split(col[8], ";")
	if len(notes) > 0 {
		featFmtThree.Attributes = []Attribute{}
		var val []string
		for _, each := range notes {
			val = strings.Split(each, "=")
			if len(val) > 1 {

				featFmtThree.Attributes = append(featFmtThree.Attributes, Attribute{val[0], val[1]})
			}
		}

	}
	return &featFmtThree, false
}

func (featFmtThree *Gff3) ToString() string {
	var str strings.Builder
	str.WriteString(featFmtThree.SeqName)
	str.WriteByte('\t')
	str.WriteString(featFmtThree.Source)
	str.WriteByte('\t')
	str.WriteString(featFmtThree.Type)
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(featFmtThree.Start))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(featFmtThree.End))
	str.WriteByte('\t')
	if featFmtThree.Score == nil {
		str.WriteByte('.')
	} else {
		str.WriteString(simpleio.Float64ToString(*featFmtThree.Score))
	}
	str.WriteByte('\t')
	if featFmtThree.Strand == nil {
		str.WriteByte('.')
	} else {
		str.WriteByte(*featFmtThree.Strand)
	}
	str.WriteByte('\t')
	if featFmtThree.Phase == nil {
		str.WriteByte('.')
	} else {
		str.WriteString(simpleio.IntToString(*featFmtThree.Phase))
	}

	str.WriteByte('\t')
	str.WriteString(AttribToStringFmtThree(featFmtThree.Attributes))
	return str.String()
}

func AttribToStringFmtThree(atb []Attribute) string {
	var str strings.Builder
	for i := 0; i < len(atb)-1; i++ {
		str.WriteString(atb[i].Tag)
		str.WriteByte('=')

		str.WriteString(atb[i].Value)

		str.WriteByte(';')

	}
	str.WriteString(atb[len(atb)-1].Tag)
	str.WriteByte('=')

	str.WriteString(atb[len(atb)-1].Value)

	return str.String()

}
