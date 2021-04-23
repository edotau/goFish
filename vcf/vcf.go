// Package vcf declares vcf struct data fields and contains methods and functions that operate on the Vcf struct
package vcf

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/goFish/simpleio"
)

// Vcf struct is the declaration of data fields.
type Vcf struct {
	Chr       string
	Pos       int
	Id        string
	Ref       string
	Alt       string
	Qual      float32
	Filter    string
	Info      string
	Format    []string
	Genotypes []string
}

//ReadToChan is a helper function.
func ReadToChan(file *VcfReader, data chan<- Vcf) {
	for curr, done := UnmarshalVcf(file); !done; curr, done = UnmarshalVcf(file) {
		data <- *curr
	}
	close(data)
}

/*
type Genotype struct {
	One    int16
	Two    int16
	Phased bool
	Depth  []byte
}*/

// VcfReader struct contains a simple reader and additional fields to help reduce memory allocation when processing lines in a vcf file.
type VcfReader struct {
	Reader *simpleio.SimpleReader
	record *Vcf
	data   []string
	done   bool
	Info   map[string]int
}

// NewHeader will allocate initial memory for a new Vcf Header and create maps used for genotypes.
func NewHeader() Header {
	answer := Header{
		Text: strings.Builder{},
	}
	answer.Ref = make(map[string]int)
	answer.Samples = make(map[string]int)
	return answer
}

// NewReader will open a text file and return a pointer to a VcfReader which contains the capibility of processing vcf text.
func NewReader(filename string) *VcfReader {
	var finish bool
	var v *Vcf
	return &VcfReader{
		Reader: simpleio.NewReader(filename),
		record: v,
		data:   make([]string, 0, 10),
		done:   finish,
	}
}

// UnmarshalVcf takes a VcfReader as an input and assign data fields and return a Vcf struct and a bool if the conversion was successful.
func UnmarshalVcf(file *VcfReader) (*Vcf, bool) {
	file.Reader.Buffer, file.done = simpleio.ReadLine(file.Reader)
	if !file.done {
		file.data = strings.SplitN(file.Reader.Buffer.String(), "\t", 10)
		if len(file.data) < 9 {
			log.Fatalf("Error when reading this vcf line:\n%s\nExpecting at least 9 columns", file.data)
		}
		file.record = &Vcf{Chr: file.data[0], Pos: simpleio.StringToInt(file.data[1]), Id: file.data[2], Ref: file.data[3], Alt: file.data[4], Filter: file.data[6], Info: file.data[7], Format: strings.Split(file.data[8], ":")}
		if strings.Compare(file.data[5], ".") == 0 {
			file.record.Qual = 255
		} else {
			file.record.Qual = simpleio.StringToFloat(file.data[5])
		}
		if len(file.data) > 9 {
			file.record.Genotypes = strings.Split(file.data[9], "\t")
			//file.record.Genotypes = GetAlleleGenotype(file.data[9])
		}
		return file.record, file.done
	} else {
		return nil, file.done
	}
}

// ReadHeader will allocate inital memory needed to process and store Vcf Header data.
func ReadHeader(file *VcfReader) *Header {
	var line *bytes.Buffer
	var done bool
	var err error
	var nextBytes []byte
	header := NewHeader()
	for nextBytes, err = file.Reader.Peek(1); err == nil && nextBytes[0] == '#'; nextBytes, err = file.Reader.Peek(1) {
		line, done = simpleio.ReadLine(file.Reader)
		if !done {
			parseHeaderData(&header, line)
		}
	}
	return &header
}

// Chrom is a method that returns the chromosome name of the record. The Chrom() method is used to implement the bed interface.
func (v *Vcf) Chrom() string {
	return v.Chr
}

// ChrStart is a method that returns the starting position of the vcf record. The ChrStart() method is used to implement the bed interface.
func (v *Vcf) ChrStart() int {
	return v.Pos
}

// ChrEnd is a method that returns the ending positing of the vcf record. The ChrEnd() method is used to implement the bed interface.
func (v *Vcf) ChrEnd() int {
	return v.Pos
}

// ToString will take a Vcf pointer and build a string.
func ToString(v *Vcf) string {
	var str strings.Builder
	str.WriteString(v.Chr)
	str.WriteByte('\t')
	str.WriteString(v.Id)
	str.WriteByte('\t')
	str.WriteString(v.Ref)
	str.WriteByte('\t')
	str.WriteString(v.Alt)
	str.WriteByte('\t')
	str.WriteString(fmt.Sprintf("%f", v.Qual))
	str.WriteByte('\t')
	str.WriteString(v.Filter)
	str.WriteByte('\t')
	str.WriteString(v.Info)
	str.WriteByte('\t')
	str.WriteString(strings.Join(v.Format, ":"))
	str.WriteByte('\t')
	str.WriteString(strings.Join(v.Genotypes, "\t"))
	return str.String()
}

func ReadVcfs(filename string) []Vcf {
	file := NewReader(filename)
	var ans []Vcf
	ReadHeader(file)
	for file.record, file.done = UnmarshalVcf(file); !file.done; file.record, file.done = UnmarshalVcf(file) {
		ans = append(ans, *file.record)
	}
	return ans
}

// GenomeKeyUint64 will combine chromosome number with a genome coordinate position and generate a key used as a hash lookup in the form of an uint64.
func GenomeKeyUint64(chrom int, start int) uint64 {
	var chromCode uint64 = uint64(chrom)
	chromCode = chromCode << 32
	var answer uint64 = chromCode | uint64(start)
	return answer
}

// Header contains a string builder to hold header info from a vcf file.
type Header struct {
	Text strings.Builder
	//Format map[string]int
	Ref        map[string]int
	ChromSizes []ChromSize
	Samples    map[string]int
}

type ChromSize struct {
	Name string
	Size int
}

// processHeader is a helper function that takes a preallocated header struct and a bytes buffer as input and will store vcf header information as a string.Builder
func parseHeaderData(header *Header, line *bytes.Buffer) {
	var hapIdx int
	words := strings.Split(line.String(), ",")
	if strings.HasPrefix(line.String(), "#") {
		header.Text.Write(line.Bytes())
		header.Text.WriteByte('\n')
		if strings.HasPrefix(line.String(), "##contig") {
			chrom := ChromSize{Name: strings.Split(words[0], "=")[2], Size: simpleio.StringToInt(strings.Split(words[1], "=")[1])}

			_, ok := header.Ref[chrom.Name]
			if !ok {
				header.Ref[chrom.Name] = len(header.ChromSizes)
				header.ChromSizes = append(header.ChromSizes, chrom)
			}
		}
		if strings.HasPrefix(line.String(), "#CHROM") {
			words := strings.Split(line.String(), "\t")[9:]
			for hapIdx = 0; hapIdx < len(words); hapIdx++ {
				header.Samples[words[hapIdx]] = hapIdx
			}
		}
	} else {
		log.Fatal("There was an error reading the header line")
	}
}

func GetGenotypes(samples string) []Genotype {
	text := strings.Split(samples, "\t")
	var hap string
	var alleles []string
	var err error
	var n int64
	var answer []Genotype = make([]Genotype, len(text))
	for i := 0; i < len(text); i++ {
		hap = strings.Split(text[i], ":")[0]
		if strings.Compare(hap, "./.") == 0 || strings.Compare(hap, ".|.") == 0 {
			//answer[i] = nil
			answer[i] = Genotype{AlleleOne: -1, AlleleTwo: -1, Phased: false}
		} else if strings.Contains(hap, "|") {
			alleles = strings.SplitN(hap, "|", 2)
			answer[i] = Genotype{AlleleOne: int16(simpleio.StringToInt(alleles[0])), AlleleTwo: int16(simpleio.StringToInt(alleles[1])), Phased: true}
		} else if strings.Contains(hap, "/") {
			alleles = strings.SplitN(hap, "/", 2)
			answer[i] = Genotype{AlleleOne: int16(simpleio.StringToInt(alleles[0])), AlleleTwo: int16(simpleio.StringToInt(alleles[1])), Phased: false}
		} else {
			//Deal with single haps. There might be a better soltuion, but I think this should work.
			n, err = strconv.ParseInt(alleles[0], 10, 16)
			if err != nil && n < int64(len(text)) {
				answer[i] = Genotype{AlleleOne: int16(n), AlleleTwo: -1, Phased: false}
			} else {
				log.Fatalf("Error: Unexpected parsing error...\n")
			}
		}
	}
	return answer

}

//Parse Vcf header to quickly print sample names that appear inside Vcf
func PrintSampleNames(header *Header) string {
	var buffer strings.Builder

	for keys := range header.Samples {
		buffer.WriteString(keys)
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func PrintSampleNamesLine(header *Header) string {
	var buffer strings.Builder

	for keys := range header.Samples {
		buffer.WriteString(keys)
		buffer.WriteByte('\t')
	}
	return buffer.String()
}

type Genotype struct {
	AlleleOne int16
	AlleleTwo int16
	Phased    bool
}

func ParseGt(data string) Genotype {
	gt := []byte(data)
	ans := Genotype{}
	if gt[0] != '.' && gt[2] != '.' {
		ans = Genotype{AlleleOne: int16(simpleio.StringToInt(string(gt[0]))), AlleleTwo: int16(simpleio.StringToInt(string(gt[2])))}
		switch gt[1] {
		case '/':
			ans.Phased = false
		case '|':
			ans.Phased = true
		default:
			log.Fatalf("Error: unexpected parasing occured, expecting '/' or '|', but found %v\n", gt[1])
		}
	} else if gt[0] != '.' && gt[2] == '.' {
		ans = Genotype{AlleleOne: int16(simpleio.StringToInt(string(gt[0]))), AlleleTwo: -1, Phased: false}
	} else {
		ans = Genotype{AlleleOne: -1, AlleleTwo: -1, Phased: false}
	}
	return ans
}

func GtToString(gt Genotype) string {
	buffer := strings.Builder{}
	buffer.Grow(3)
	if gt.AlleleOne == -1 && gt.AlleleTwo == -1 {
		buffer.WriteByte('.')
		buffer.WriteByte('/')
	} else {
		buffer.WriteString(simpleio.Int16ToString((gt.AlleleOne)))
		if !gt.Phased {
			buffer.WriteByte('/')
		} else {
			buffer.WriteByte('|')
		}
		buffer.WriteString(simpleio.Int16ToString((gt.AlleleTwo)))
	}
	return buffer.String()
}

func MkFormatMap(v *Vcf) map[string]int {
	ans := make(map[string]int)
	for i := 0; i < len(v.Format); i++ {
		ans[v.Format[i]] = i
	}
	return ans
}

func GetAllAlleleDepth(v *Vcf) []string {
	format := MkFormatMap(v)
	var ans []string = make([]string, len(v.Genotypes))
	var line []string = make([]string, len(v.Format))
	for i := 0; i < len(ans); i++ {
		line = strings.Split(v.Genotypes[i], ":")
		ans[i] = FindAlleleDepth(format, line)
	}
	return ans
}

func FindAlleleDepth(format map[string]int, line []string) string {
	return line[format["AD"]]
}

type AlleleDepth struct {
	Gt    Genotype
	Depth []int
}

func GetGenotypeDepth(v *Vcf) []string {
	format := MkFormatMap(v)
	var ans []string = make([]string, len(v.Genotypes))
	var line []string = make([]string, len(v.Format))
	for i := 0; i < len(ans); i++ {
		line = strings.Split(v.Genotypes[i], ":")
		ans[i] = line[0] + ":"
		ans[i] += FindAlleleDepth(format, line)
	}
	return ans
}
