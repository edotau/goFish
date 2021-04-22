package main

/*
import (
	"github.com/goFish/vcf"
	//"bufio"
	//"compress/gzip"
	//"encoding/csv"
	"flag"
	"sort"
	"fmt"
	"github.com/vertgenlab/gonomics/common"
	"github.com/vertgenlab/gonomics/fileio"
	"github.com/vertgenlab/gonomics/numbers"
	"log"
	//"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Print(
		"stickleback atac rna allele analysis\nUsage:\n  ./atacRnaFish vcf.gz\n\n")
	flag.PrintDefaults()
}


func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	reader := vcf.NewVcfReader(flag.Arg(0))
	vcf.ReadHeader(reader)
	//leftTest, rightTest := make([]*Diff, 0), make([]*Diff, 0)
	//for v, done := vcf.UnmarshalVcf(reader); !done; v, done = vcf.UnmarshalVcf(reader) {
	//	fmt.Printf("%v\n", v)
	//}

	//fmt.Printf("CL16_wgs.merged\tCL16_rnaseq\n")
	//leftTest, rightTest := make([]*Diff, 0), make([]*Diff, 0)
	for line, done := ParseLine(reader); !done; line, done = ParseLine(reader) {
	//	if curr := atac(line, 25,22); curr != nil {
	//		if curr.LeftPvalue < 0.05 {
	//			leftTest = append(leftTest, &Diff{Chr: curr.Chr, Start: curr.Start, Pval: curr.LeftPvalue})
	//		}
	//		if curr.RightPvalue < 0.05 {
	//			rightTest = append(rightTest, &Diff{Chr: curr.Chr, Start: curr.Start, Pval: curr.RightPvalue})
	//		}
	//	}
	}
	//WriteDiff("cl20_atacseq_FishersExactLeft.csv", leftTest)
	//WriteDiff("cl20_atacseq_FishersExactRight.csv", rightTest)
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
	Chr         string
	Start       string
	Pval 		float64
}

func (test *FisherExact) String() string {
	return fmt.Sprintf("%d %d\n%d %d", test.A, test.B, test.C, test.D)
}

func comparePValue(a *Diff, b *Diff) int {
	if a.Pval < b.Pval {
		return -1
	}
	if a.Pval > b.Pval {
		return 1
	}
	return 0
}

func SortByPValue(peak []*Diff) {
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

func WriteDiff(filename string, peaksDiff []*Diff) {
	pvaluePeaks := fileio.EasyCreate(filename)
	SortByPValue(peaksDiff)
	var err error
	_, err = fmt.Fprintf(pvaluePeaks, "%s\n", "chr,start,Pvalue")
	for _, peak := range peaksDiff {
		_, err = fmt.Fprintf(pvaluePeaks, "%s,%s,%E\n", peak.Chr, peak.Start, peak.Pval)
		common.ExitIfError(err)
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
	Matrix      *FisherExact
	LeftPvalue  float64
	RightPvalue float64
}
func atac(line []string, wgsIdx int, colNum int) *Stats {
	//CL16: wgs=13, atac=10 XXXX
	//CL17: wgs=17, atac=14 XXXXX
	//CL18: wgs=21 atac=18 XXXXXX
	//CL20: wgs=25, atac=22 XXXXX
	wholeGenome := wgs(line[wgsIdx])
	atac := wgs(line[colNum])

	if wholeGenome != nil && atac != nil {
		var fisherTest FisherExact = FisherExact{}
		//if wholeGenome.Allele[0] == rna.Allele[0] {
		fisherTest.A = wholeGenome.Depth[0]
		fisherTest.B = wholeGenome.Depth[0] + atac.Depth[0]
		//}
		//if wholeGenome.Allele[1] == atac.Allele[1] {
		fisherTest.C = wholeGenome.Depth[1]
		fisherTest.D = wholeGenome.Depth[1] + atac.Depth[1]
		//}
		var ans Stats = Stats{
			Chr: line[0],
			Start: line[1],
			Matrix: &fisherTest,
			LeftPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func cl16(line []string) *Stats {
	wholeGenome := wgs(line[13])
	rna := rnaSeq(line[11], line[12])

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
			Chr: line[0],
			Start: line[1],
			Matrix: &fisherTest,
			LeftPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func cl17(line []string) *Stats {
	wholeGenome := wgs(line[17])
	rna := rnaSeq(line[15], line[16])

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
			Chr: line[0],
			Start: line[1],
			Matrix: &fisherTest,
			LeftPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func cl18(line []string) *Stats {
	wholeGenome := wgs(line[21])
	rna := rnaSeq(line[19], line[20])

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
			Chr: line[0],
			Start: line[1],
			Matrix: &fisherTest,
			LeftPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
			RightPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false),
		}
		return &ans
	} else {
		return nil
	}
}

func cl20(line []string) *Stats {
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
			Chr: line[0],
			Start: line[1],
			Matrix: &fisherTest,
			LeftPvalue: numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true),
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
}

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
}
*/
