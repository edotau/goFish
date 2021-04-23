package main

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/edotau/goFish/simpleio"
	"github.com/vertgenlab/gonomics/common"
	"github.com/vertgenlab/gonomics/fileio"
	"github.com/vertgenlab/gonomics/numbers"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func usage() {
	fmt.Print(
		"./vcfReader $vcf.gz\n performs FisherExact test\n\n")
	flag.PrintDefaults()
}

type VcfReader struct {
	*csv.Reader
	file *os.File
	line []string
	//buffer *bytes.Buffer
}

func NewVcfReader(filename string) *VcfReader {
	var answer VcfReader = VcfReader{
		file: fileio.MustOpen(filename),
		line: make([]string, 100),
	}
	switch true {
	case strings.HasSuffix(filename, ".gz"):
		gzipReader, err := gzip.NewReader(answer.file)
		common.ExitIfError(err)
		answer.Reader = csv.NewReader(bufio.NewReader(gzipReader))
	default:
		answer.Reader = csv.NewReader(bufio.NewReader(answer.file))
	}
	answer.Reader.Comma = '\t'
	answer.Reader.Comment = '#'
	return &answer
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	reader := NewVcfReader(flag.Arg(0))
	answer := make([]Diff, 0)
	for line, done := ParseLine(reader); !done; line, done = ParseLine(reader) {
		// new index
		// CL16: wgs=19, atac=16
		// CL17: wgs=23 , atac=20

		if curr := CL16rna(line); curr != nil {
			if curr.LeftPvalue < 0.05 && curr.RightPvalue < 0.05 {
				log.Fatalf("Error: both left and right sides can't be significant...\n")
			} else if curr.LeftPvalue < 0.05 && simpleio.StringToInt(curr.Start) > 5000000 && simpleio.StringToInt(curr.Start) < getChrom(curr.Chr)-5000000 {
				//fmt.Printf("wgs=%s, atac=%s, %v\n", line[19], line[16], fisherTest)
				answer = append(answer, Diff{Chr: curr.Chr, Start: curr.Start, Pval: curr.LeftPvalue, FS: curr.Matrix})
			} else if curr.RightPvalue < 0.05 && simpleio.StringToInt(curr.Start) > 5000000 && simpleio.StringToInt(curr.Start) < getChrom(curr.Chr)-5000000 {
				answer = append(answer, Diff{Chr: curr.Chr, Start: curr.Start, Pval: curr.RightPvalue, FS: curr.Matrix})
			}
		}
	}
	WriteDiff("updatedAnswers/CL16_RNASEQ.FS.csv", answer)
}

func CL16atac(line []string) *Stats {
	return atac(line, 19, 16)
}

func CL17atac(line []string) *Stats {
	return atac(line, 23, 20)
}

func CL18atac(line []string) *Stats {
	return atac(line, 27, 24)
}

func CL20atac(line []string) *Stats {
	return atac(line, 31, 28)
}

//wgs=19, 1=17,2=18

func CL16rna(line []string) *Stats {
	return rna(line, 19, 17, 18)
}

//wgs=23

func CL17rna(line []string) *Stats {
	return rna(line, 23, 21, 22)
}

func CL18rna(line []string) *Stats {
	return rna(line, 27, 25, 26)
}

func CL20rna(line []string) *Stats {
	return rna(line, 31, 29, 30)
}

type FisherExact struct {
	A int
	B int
	C int
	D int
}

type AlleleDepth struct {
	Allele []int
	Depth  []int
}

type Diff struct {
	Chr   string
	Start string
	Pval  float64
	FS    FisherExact
	//RawWgs string
	//RawAS string
}

func (test *FisherExact) String() string {
	return fmt.Sprintf("%d\t%d\t%d\t%d", test.A, test.B, test.C, test.D)
}

func comparePValue(a Diff, b Diff) int {
	if a.Pval < b.Pval {
		return -1
	}
	if a.Pval > b.Pval {
		return 1
	}
	return 0
}

func SortByPValue(peak []Diff) {
	sort.Slice(peak, func(i, j int) bool { return comparePValue(peak[i], peak[j]) == -1 })
}

func ToStringAD(ad *AlleleDepth) string {
	var str strings.Builder
	str.Grow(len(ad.Allele) + len(ad.Depth))
	var i int
	//str.WriteString("Genotypes:\n")
	for i = 0; i < len(ad.Allele); i++ {
		str.WriteString(strconv.Itoa(ad.Allele[i]))
		str.WriteByte('\t')
	}
	str.WriteByte('\n')
	for i = 0; i < len(ad.Depth); i++ {
		str.WriteString(strconv.Itoa(ad.Depth[i]))
		str.WriteByte('\t')
	}
	return str.String()
}

func WriteDiff(filename string, peaksDiff []Diff) {
	pvaluePeaks := fileio.EasyCreate(filename)
	SortByPValue(peaksDiff)
	var err error
	_, err = fmt.Fprintf(pvaluePeaks, "%s\n", "chr,start,A,B,C,D,Pvalue")
	for _, peak := range peaksDiff {
		if peak.Chr != "chrM" {

			_, err = fmt.Fprintf(pvaluePeaks, "%s,%s,%d,%d,%d,%d,%E\n", peak.Chr, peak.Start, peak.FS.A, peak.FS.B, peak.FS.C, peak.FS.D, peak.Pval)
			common.ExitIfError(err)
		}
	}
	pvaluePeaks.Close()
}

func getCounts(s string, a []int) []int {
	var ans []int
	genotype := strings.Split(s, ",")
	for i := 0; i < len(a); i++ {
		ans = append(ans, common.StringToInt(genotype[a[i]]))
	}
	return ans
}

type Stats struct {
	Chr         string
	Start       string
	Matrix      FisherExact
	LeftPvalue  float64
	RightPvalue float64
}

func withinPercent(a int, b int, percent float64) bool {
	//errRate := percent*(float64(a)+float64(b))/2
	if math.Abs(float64(a-b)) < percent*(float64(a)+float64(b))/2 {
		return true
	} else {
		return false
	}
}

// new index
// CL16: wgs=19, atac=16

// CL17: wgs=23 , atac=20
func atac(line []string, wgsIdx int, colNum int) *Stats {
	//CL16: wgs=13, atac=10 XXXX
	//CL17: wgs=17, atac=14 XXXXX
	//CL18: wgs=21 atac=18 XXXXXX
	//CL20: wgs=25, atac=22 XXXXX

	atacOne, atacTwo := parseGenotype(line[colNum])

	//fmt.Printf("%s\n", )
	if atacOne != nil && atacTwo != nil {
		wholeGenomeOne, wholeGenomeTwo := parseWgsDepth(line[wgsIdx], atacOne.Id, atacTwo.Id)
		if wholeGenomeOne != nil && wholeGenomeTwo != nil {
			var fisherTest FisherExact = FisherExact{}
			//if atacOne.Id == atacTwo.Id {
			//	return nil
			//}
			//if wholeGenomeOne.Id == wholeGenomeTwo.Id {
			//	return nil
			//}
			if wholeGenomeOne.Depth == 0 || wholeGenomeTwo.Depth == 0 {
				return nil
			}

			////if !withinPercent(wholeGenomeOne.Depth, wholeGenomeTwo.Depth, .3) {
			//	return nil
			//}
			//if math.Abs(float64(wholeGenomeOne.Depth-wholeGenomeTwo.Depth)) > 5 {
			//	return nil
			//}
			//if wholeGenomeOne.Depth < 10 && wholeGenomeTwo.Depth < 10 {
			//	return nil
			//}
			if wholeGenomeOne.Id == atacOne.Id && wholeGenomeTwo.Id == atacTwo.Id && wholeGenomeTwo.Id != 0 {
				fisherTest.A = wholeGenomeOne.Depth
				fisherTest.B = wholeGenomeTwo.Depth
				fisherTest.C = atacOne.Depth
				fisherTest.D = atacTwo.Depth
				//fmt.Printf("wgs=%s, atac=%s, answer=%v\n", line[19], line[16], fisherTest)
				return &Stats{
					Chr:         line[0],
					Start:       line[1],
					Matrix:      fisherTest,
					LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
					RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
				}
			}
			//if wholeGenomeTwo.Id != atacTwo.Id {
			//	log.Fatalf("Error: a and c or b and d positions in matrix are not equal %d != %d or %d != %d...\n", wholeGenomeOne.Id, atacOne.Id, wholeGenomeTwo.Id, atacTwo.Id)
			//}

			//if wholeGenomeOne.Id != wholeGenomeTwo.Id

			//if wholeGenome.Allele[0] == rna.Allele[0] {
			fisherTest.A = wholeGenomeOne.Depth
			if wholeGenomeTwo.Id != 0 {
				fisherTest.B = wholeGenomeTwo.Depth
			} else {
				fisherTest.B = 0
			}

			//}
			//if wholeGenome.Allele[1] == atac.Allele[1] {
			fisherTest.C = atacOne.Depth
			if atacOne.Id == atacTwo.Id && atacOne.Id != 0 {
				fisherTest.C = 0
				fisherTest.D = atacTwo.Depth
			} else if atacOne.Id == atacTwo.Id && atacOne.Id == 0 {
				fisherTest.D = 0
			} else {
				fisherTest.D = atacTwo.Depth
			}

			//}
			var ans Stats = Stats{
				Chr:         line[0],
				Start:       line[1],
				Matrix:      fisherTest,
				LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
				RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
			}
			return &ans
		}
	} //else {
	return nil
	//}
}

type genotype struct {
	Id    int
	Depth int
}

func parseWgsDepth(input string, alleleOne int, alleleTwo int) (*genotype, *genotype) {
	gt := strings.Split(input, ":")
	if strings.Contains(gt[0], "./.") {
		return nil, nil
	}
	if strings.Contains(gt[0], "/") {
		info := strings.Split(gt[0], "/")
		work := strings.Split(gt[1], ",")
		one := &genotype{}
		two := &genotype{}
		if strings.Contains(info[0], ".") {
			one.Id = 0
		} else {
			one.Id = alleleOne
		}
		if strings.Contains(info[1], ".") {
			two.Id = 0
		} else {
			two.Id = alleleTwo
		}

		if work[one.Id] == "." {
			one.Depth = 0
		} else {
			one.Depth = simpleio.StringToInt(work[one.Id])
		}
		if work[two.Id] == "." {
			two.Depth = 0
		} else {
			two.Depth = simpleio.StringToInt(work[two.Id])
		}
		return one, two
	}

	if strings.Contains(gt[0], "|") {
		info := strings.Split(gt[0], "|")
		work := strings.Split(gt[1], ",")
		one := &genotype{}
		two := &genotype{}
		if strings.Contains(info[0], ".") {
			one.Id = 0
		} else {
			one.Id = alleleOne
		}
		if strings.Contains(info[1], ".") {
			two.Id = 0
		} else {
			two.Id = alleleTwo
		}

		if work[one.Id] == "." {
			one.Depth = 0
		} else {
			one.Depth = simpleio.StringToInt(work[one.Id])
		}
		if work[two.Id] == "." {
			two.Depth = 0
		} else {
			two.Depth = simpleio.StringToInt(work[two.Id])
		}
		return one, two
	}
	return nil, nil
}

func parseGenotype(input string) (*genotype, *genotype) {
	gt := strings.Split(input, ":")
	if strings.Contains(gt[0], "./.") {
		return nil, nil
	}
	if strings.Contains(gt[0], "/") {
		info := strings.Split(gt[0], "/")
		work := strings.Split(gt[1], ",")
		one := &genotype{}
		two := &genotype{}
		if strings.Contains(info[0], ".") {
			one.Id = 0
		} else {
			one.Id = simpleio.StringToInt(info[0])
		}
		if strings.Contains(info[1], ".") {
			two.Id = 0
		} else {
			two.Id = simpleio.StringToInt(info[1])
		}

		if work[one.Id] == "." {
			one.Depth = 0
		} else {
			one.Depth = simpleio.StringToInt(work[one.Id])
		}
		if work[two.Id] == "." {
			two.Depth = 0
		} else {
			two.Depth = simpleio.StringToInt(work[two.Id])
		}
		return one, two
	}

	if strings.Contains(gt[0], "|") {
		info := strings.Split(gt[0], "|")
		work := strings.Split(gt[1], ",")
		one := &genotype{}
		two := &genotype{}
		if strings.Contains(info[0], ".") {
			one.Id = 0
		} else {
			one.Id = simpleio.StringToInt(info[0])
		}
		if strings.Contains(info[1], ".") {
			two.Id = 0
		} else {
			two.Id = simpleio.StringToInt(info[1])
		}

		if work[one.Id] == "." {
			one.Depth = 0
		} else {
			one.Depth = simpleio.StringToInt(work[one.Id])
		}
		if work[two.Id] == "." {
			two.Depth = 0
		} else {
			two.Depth = simpleio.StringToInt(work[two.Id])
		}
		return one, two
	}
	return nil, nil
}

/*
func rna(line []string, wI int, r1 int, r2 int) *Stats {
	wholeGenome := wgs(line[wI])
	rna := rnaSeq(line[r1], line[r2])

	if wholeGenome != nil && rna != nil {
		var fisherTest FisherExact = FisherExact{}
		//if wholeGenome.Allele[0] == rna.Allele[0] {
		fisherTest.A = wholeGenome.Depth[0]
		fisherTest.B = rna.Depth[0]
		//}
		//if wholeGenome.Allele[1] == rna.Allele[1] {
		fisherTest.C = wholeGenome.Depth[1]
		fisherTest.D = rna.Depth[1]
		//}
		var ans Stats = Stats{
			Chr:         line[0],
			Start:       line[1],
			Matrix:      fisherTest,
			LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}*/

func rna(line []string, wgsIdx int, colNumOne int, colNumTwo int) *Stats {
	rnaOne, rnaTwo := parseGenotype(line[colNumOne])
	rnaOneTmp, rnaTwoTmp := parseGenotype(line[colNumTwo])

	//fmt.Printf("%s\n", )
	if rnaOne != nil && rnaTwo != nil && rnaOneTmp != nil && rnaTwoTmp != nil {
		if rnaOne.Id != rnaOneTmp.Id || rnaTwo.Id != rnaTwoTmp.Id {
			return nil
			//log.Fatalf("Error: genotypes for both rnaseq sampels should be the same...%s != %s\n", line[colNumOne], line[colNumTwo])
		}
		rnaOne.Depth += rnaOneTmp.Depth
		rnaTwo.Depth += rnaTwoTmp.Depth
		wholeGenomeOne, wholeGenomeTwo := parseWgsDepth(line[wgsIdx], rnaOne.Id, rnaTwo.Id)
		if wholeGenomeOne != nil && wholeGenomeTwo != nil {
			var fisherTest FisherExact = FisherExact{}
			//if rnaOne.Id == rnaTwo.Id {
			//	return nil
			//}
			//if wholeGenomeOne.Id == wholeGenomeTwo.Id {
			//	return nil
			//}
			if wholeGenomeOne.Depth == 0 || wholeGenomeTwo.Depth == 0 {
				return nil
			}
			//if !withinPercent(wholeGenomeOne.Depth, wholeGenomeTwo.Depth, .3) {
			//	return nil
			//}
			//if math.Abs(float64(wholeGenomeOne.Depth-wholeGenomeTwo.Depth)) > 5 {
			//	return nil
			//}
			//if wholeGenomeOne.Depth < 10 && wholeGenomeTwo.Depth < 10 {
			//	return nil
			//}
			if wholeGenomeOne.Id == rnaOne.Id && wholeGenomeTwo.Id == rnaTwo.Id && wholeGenomeTwo.Id != 0 {
				fisherTest.A = wholeGenomeOne.Depth
				fisherTest.B = wholeGenomeTwo.Depth
				fisherTest.C = rnaOne.Depth
				fisherTest.D = rnaTwo.Depth
				//fmt.Printf("wgs=%s, atac=%s, answer=%v\n", line[19], line[16], fisherTest)
				return &Stats{
					Chr:         line[0],
					Start:       line[1],
					Matrix:      fisherTest,
					LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
					RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
				}
			}
			//if wholeGenomeTwo.Id != rnaTwo.Id {
			//	log.Fatalf("Error: a and c or b and d positions in matrix are not equal %d != %d or %d != %d...\n", wholeGenomeOne.Id, rnaOne.Id, wholeGenomeTwo.Id, rnaTwo.Id)
			//}

			//if wholeGenomeOne.Id != wholeGenomeTwo.Id

			//if wholeGenome.Allele[0] == rna.Allele[0] {
			fisherTest.A = wholeGenomeOne.Depth
			if wholeGenomeTwo.Id != 0 {
				fisherTest.B = wholeGenomeTwo.Depth
			} else {
				fisherTest.B = 0
			}

			//}
			//if wholeGenome.Allele[1] == atac.Allele[1] {
			fisherTest.C = rnaOne.Depth
			if rnaOne.Id == rnaTwo.Id && rnaOne.Id != 0 {
				fisherTest.C = 0
				fisherTest.D = rnaTwo.Depth
			} else if rnaOne.Id == rnaTwo.Id && rnaOne.Id == 0 {
				fisherTest.D = 0
			} else {
				fisherTest.D = rnaTwo.Depth
			}

			//}
			var ans Stats = Stats{
				Chr:         line[0],
				Start:       line[1],
				Matrix:      fisherTest,
				LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
				RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
			}
			return &ans
		}
	} //else {
	return nil
	//}
}

// RNA-Seq CL16 wgs=19, 1=17,2=18
func rnaCl16(line []string) *Stats {
	wholeGenome := wgs(line[19])
	rna := rnaSeq(line[17], line[18])

	if wholeGenome != nil && rna != nil {
		var fisherTest FisherExact = FisherExact{}
		//if wholeGenome.Allele[0] == rna.Allele[0] {
		fisherTest.A = wholeGenome.Depth[0]
		fisherTest.B = rna.Depth[0]
		//}
		//if wholeGenome.Allele[1] == rna.Allele[1] {
		fisherTest.C = wholeGenome.Depth[1]
		fisherTest.D = rna.Depth[1]
		//}
		var ans Stats = Stats{
			Chr:         line[0],
			Start:       line[1],
			Matrix:      fisherTest,
			LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func rnaCl17(line []string) *Stats {
	wholeGenome := wgs(line[17])
	rna := rnaSeq(line[15], line[16])

	if wholeGenome != nil && rna != nil {
		var fisherTest FisherExact = FisherExact{}
		//if wholeGenome.Allele[0] == rna.Allele[0] {
		fisherTest.A = wholeGenome.Depth[0]
		fisherTest.B = rna.Depth[0]
		//}
		//if wholeGenome.Allele[1] == rna.Allele[1] {
		fisherTest.C = wholeGenome.Depth[1]
		fisherTest.D = rna.Depth[1]
		//}
		var ans Stats = Stats{
			Chr:         line[0],
			Start:       line[1],
			Matrix:      fisherTest,
			LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func rnaCl18(line []string) *Stats {
	wholeGenome := wgs(line[21])
	rna := rnaSeq(line[19], line[20])

	if wholeGenome != nil && rna != nil {
		var fisherTest FisherExact = FisherExact{}
		//if wholeGenome.Allele[0] == rna.Allele[0] {
		fisherTest.A = wholeGenome.Depth[0]
		fisherTest.B = rna.Depth[0]
		//}
		//if wholeGenome.Allele[1] == rna.Allele[1] {
		fisherTest.C = wholeGenome.Depth[1]
		fisherTest.D = rna.Depth[1]
		//}
		var ans Stats = Stats{
			Chr:         line[0],
			Start:       line[1],
			Matrix:      fisherTest,
			LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func rnaCl20(line []string) *Stats {
	wholeGenome := wgs(line[25])
	rna := rnaSeq(line[23], line[24])

	if wholeGenome != nil && rna != nil {
		var fisherTest FisherExact = FisherExact{}
		//if wholeGenome.Allele[0] == rna.Allele[0] {
		fisherTest.A = wholeGenome.Depth[0]
		fisherTest.B = wholeGenome.Depth[0] + rna.Depth[0]
		//}
		//if wholeGenome.Allele[1] == rna.Allele[1] {
		fisherTest.C = wholeGenome.Depth[1]
		fisherTest.D = wholeGenome.Depth[1] + rna.Depth[1]
		//}
		var ans Stats = Stats{
			Chr:         line[0],
			Start:       line[1],
			Matrix:      fisherTest,
			LeftPvalue:  numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func rnaSeq(info string, infoTwo string) *AlleleDepth {
	if strings.Contains(info, "./.") || strings.Contains(info, ".|.") {
		//log.Printf("%s", info)
		return nil
	}
	words := strings.Split(info, ":")
	var ans *AlleleDepth = &AlleleDepth{
		Allele: GetGenotype(words[0]),
	}
	if len(ans.Allele) == 1 {
		return nil
	}
	ans.Depth = getCounts(words[1], ans.Allele)
	if strings.Contains(infoTwo, "./.") || strings.Contains(infoTwo, ".|.") {
		//log.Printf("%s", info)
		return nil
	}
	words = strings.Split(infoTwo, ":")
	var tmp *AlleleDepth = &AlleleDepth{
		Allele: GetGenotype(words[0]),
	}
	tmp.Depth = getCounts(words[1], tmp.Allele)
	if len(tmp.Depth) == len(ans.Depth) {
		for i := 0; i < len(ans.Depth); i++ {
			ans.Depth[i] += tmp.Depth[i]
		}
	}
	return ans
}

func wgs(info string) *AlleleDepth {
	if strings.Contains(info, "./.") || strings.Contains(info, ".|.") {
		//log.Printf("%s", info)
		return nil
	}
	words := strings.Split(info, ":")
	var ans *AlleleDepth = &AlleleDepth{
		Allele: GetGenotype(words[0]),
	}
	if len(ans.Allele) == 1 {
		return nil
	}
	ans.Depth = getCounts(words[1], ans.Allele)
	return ans
}

/*
func cl16(line []string) *FisherExact {

	var ans *FisherExact = &FisherExact{}
	info := strings.Split(line[13], ":")
	wgs := GetGenotype(info[0])
	if wgs[0] == wgs[1] {
		return nil
	}
	rawAD := WgsAlleleDepth(info[1])

	ans.A = rawAD[0]


	//
	//if rawAD == nil {
	//	return nil
	//}
	ans[0], ans[1] = rawAD[0], rawAD[1]

	rawAD = RnaSeq(line[11])
	ans[2], ans[3] = rawAD[0], rawAD[1]

	rawAD = RnaSeq(line[12])
	ans[2]+= rawAD[0]
	ans[3]+= rawAD[1]
	//ans[] = processAlleleDepth(line[11])
	//ans[2] = processAlleleDepth(line[12])
	return &FisherExact{}
}*/

func ParseLine(reader *VcfReader) ([]string, bool) {
	var err error
	reader.line, err = reader.Read()
	if err == nil {
		return reader.line, false
	} else {
		fileio.CatchErrThrowEOF(err)
		reader.Close()
		return nil, true
	}
}

func GetGenotype(s string) []int {
	var ans []int
	if strings.Contains(s, "/") {
		a := strings.Split(s, "/")
		ans = append(ans, common.StringToInt(a[0]))
		for i := 1; i < len(a); i++ {
			if next := common.StringToInt(a[i]); next != ans[i-1] {
				ans = append(ans, next)
			}
		}
		return ans
	} else if strings.Contains(s, "|") {
		a := strings.Split(s, "|")
		ans = append(ans, common.StringToInt(a[0]))
		for i := 1; i < len(a); i++ {
			if next := common.StringToInt(a[i]); next != ans[i-1] {
				ans = append(ans, next)
			}
		}
		return ans
	} else {
		return nil
	}
}

/*
func RnaSeq(info string) []int {
	var ans []int = make([]int, 2)
	text := strings.Split(info, ":")
	genotypes := GetGenotype(text[0])
	words := strings.Split(text[1], ",")
	ans[0] = common.StringToInt(words[genotypes[0]])
	if genotypes[0] != genotypes[1] {
		ans[1] = common.StringToInt(words[genotypes[1]])
	} else {
		ans[1] = 0
	}

	return ans
}*/

func (reader *VcfReader) Close() {
	if reader != nil {
		err := reader.file.Close()
		common.ExitIfError(err)
	}
}
