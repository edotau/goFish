package axt

import (
	"github.com/goFish/alphabet"
	"github.com/goFish/bed"
	"github.com/goFish/simpleio"
	"log"
	"strings"
)

// Axt struct: Naming convention is hard here because UCSC website does not
// match the UCSC Kent source tree.
type Axt struct {
	RName      string
	RStart     int
	REnd       int
	QName      string
	QStart     int
	QEnd       int
	QStrandPos byte // true is positive strand, false is negative strand
	Score      int
	RSeq       []alphabet.Dna
	QSeq       []alphabet.Dna
}

func FindIndels(filename string, output string) {
	writer := simpleio.NewWriter(output)

	defer writer.Close()
	reader := simpleio.NewReader(filename)
	var j bed.GenomeInfo
	for i, done := AxtRecord(reader); !done; i, done = AxtRecord(reader) {
		indels := AxtToGenomeInfo(i)
		for _, j = range indels {
			writer.WriteString(bed.GenomeInfoToString(j))
			writer.WriteByte('\n')
		}
	}
}

//NextAxt processes the next Axt alignment in the provided input.
func AxtRecord(reader *simpleio.SimpleReader) (*Axt, bool) {
	header, hDone := simpleio.ReadLine(reader)
	if hDone {
		return nil, true
	}
	var words []string = strings.Split(header.String(), " ")
	if len(words) != 9 {
		log.Fatalf("Error: missing fields in header or sequences\n")
	}
	var answer *Axt = &Axt{
		RName:      words[1],
		RStart:     simpleio.StringToInt(words[2]),
		REnd:       simpleio.StringToInt(words[3]),
		QName:      words[4],
		QStart:     simpleio.StringToInt(words[5]),
		QEnd:       simpleio.StringToInt(words[6]),
		QStrandPos: words[7][0],
		Score:      simpleio.StringToInt(words[8]),
	}
	target, tDone := simpleio.ReadLine(reader)

	answer.RSeq = make([]alphabet.Dna, len(target.Bytes()))
	copy(answer.RSeq, alphabet.ToDna(target.Bytes()))

	query, qDone := simpleio.ReadLine(reader)

	answer.QSeq = make([]alphabet.Dna, len(query.Bytes()))

	copy(answer.QSeq, alphabet.ToDna(query.Bytes()))

	blank, bDone := simpleio.ReadLine(reader)
	if blank.String() != "" {
		log.Fatalf("Error: every fourth line should be blank %s\n", blank.String())
	}
	if tDone || qDone || bDone {
		return nil, true
	}
	return answer, false
}
