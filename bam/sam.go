package bam

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/simpleio"
)

type Sam struct {
	QName   string
	Flag    uint16
	RName   string
	Pos     int
	MapQ    uint8
	Cigar   Cigar
	MateRef string
	MatePos int
	TmpLen  int
	Seq     Sequence
	Qual    PhredQual
	Aux     string
}

type Cigar []ByteCigar

type Sequence []code.Dna

type PhredQual []byte

type SamReader struct {
	Reader *simpleio.SimpleReader
	record Sam
	line   []string
	done   bool
}

func NewSamReader(filename string) *SamReader {
	var answer SamReader = SamReader{
		Reader: simpleio.NewReader(filename),
		record: Sam{},
	}
	return &answer
}

func ReadSamRecord(filename string) []*Sam {
	reader := NewSamReader(filename)
	ReadSamHeader(reader)
	var answer []*Sam
	for curr, done := UnmarshalSam(reader); !done; curr, done = UnmarshalSam(reader) {
		answer = append(answer, curr)
	}
	return answer
}

func ReadSam(filename string) (*Header, <-chan Sam) {
	reader := NewSamReader(filename)
	header := ReadSamHeader(reader)
	ans := make(chan Sam)
	go func() {
		for curr, done := UnmarshalSam(reader); !done; curr, done = UnmarshalSam(reader) {
			ans <- *curr
		}
		close(ans)
	}()
	return header, ans
}

func ReadSamHeader(file *SamReader) *Header {
	var header Header = Header{
		//Text: &bytes.Buffer{},
	}
	header.ChromSize = make(map[string]int)
	var done bool
	var curr *bytes.Buffer = &bytes.Buffer{}
	for nextBytes, err := file.Reader.Peek(1); err == nil && nextBytes[0] == '@'; nextBytes, err = file.Reader.Peek(1) {
		curr, done = simpleio.ReadLine(file.Reader)
		if !done {
			processHeaderLine(&header, curr.String())
		} else {
			break
		}

	}
	return &header
}

type Header struct {
	Text       bytes.Buffer
	ChromSize  map[string]int
	Chroms     []ChromSize
	bamDecoder *chrInfo
}

type ChromSize struct {
	Name  string
	Size  int
	Order int
}

func processHeaderLine(header *Header, line string) {
	header.Text.WriteString(line)
	header.Text.WriteByte('\n')
	//header.Text.Write(line)
	//header.Text.WriteByte('\n')
	if strings.HasPrefix(line, "@SQ") && strings.Contains(line, "SN:") && strings.Contains(line, "LN:") {
		//if bytes.HasPrefix(line, []byte("@SQ")) && bytes.Contains(line, []byte("SN:")) && bytes.Contains(line, []byte("LN:")) {
		var currName string
		var currLen int
		words := strings.Fields(line)
		for i := 1; i < len(words); i++ {
			elements := strings.Split(words[i], ":")
			switch elements[0] {
			case "SN":
				currName = elements[1]
			case "LN":
				currLen = simpleio.StringToInt(elements[1])
			}
		}
		header.ChromSize[currName] = currLen
	}
}

func UnmarshalSam(file *SamReader) (*Sam, bool) {
	file.Reader.Buffer, file.done = simpleio.ReadLine(file.Reader)
	if !file.done {
		file.line = strings.SplitN(file.Reader.Buffer.String(), "\t", 12)
		if len(file.line) < 11 {
			log.Fatalf("Error: missing sam alignment fields, must contain at least 11...\n")
		}
		ans := &Sam{
			QName:   file.line[0],
			Flag:    simpleio.StringToUInt16(file.line[1]),
			RName:   file.line[2],
			Pos:     simpleio.StringToInt(file.line[3]),
			MapQ:    uint8(simpleio.StringToInt(file.line[4])),
			Cigar:   ReadToBytesCigar([]byte(file.line[5])),
			MateRef: string(file.line[6]),
			MatePos: simpleio.StringToInt(file.line[7]),
			TmpLen:  simpleio.StringToInt(file.line[8]),
			Seq:     code.ToDna([]byte(file.line[9])),
			Qual:    []byte(file.line[10]),
			Aux:     file.line[11],
		}
		if len(ans.Seq) != len(ans.Qual) {
			log.Fatalf("Error: seq and qual lengths should match...\n")
		}
		return ans, false
	} else {
		return nil, true
	}
}

func ToString(record *Sam) string {
	var str strings.Builder
	var err error
	_, err = str.WriteString(record.QName)
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(strconv.Itoa(int(record.Flag)))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(record.RName)
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(strconv.Itoa(record.Pos))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(strconv.Itoa(int(record.MapQ)))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(ByteCigarToString(record.Cigar))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(record.MateRef)
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(strconv.Itoa(record.MatePos))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(strconv.Itoa(record.TmpLen))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(code.ToString(record.Seq))
	simpleio.ErrorHandle(err)
	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	for _, phred := range record.Qual {
		err = str.WriteByte(phred)
		simpleio.ErrorHandle(err)
	}

	err = str.WriteByte('\t')
	simpleio.ErrorHandle(err)
	_, err = str.WriteString(record.Aux)
	simpleio.ErrorHandle(err)
	return str.String()
}
