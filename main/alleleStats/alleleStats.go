package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/simpleio"
	"github.com/edotau/goFish/stats"
	"github.com/edotau/goFish/vcf"

	"log"
	"strings"
)

func usage() {
	fmt.Print(
		"alleleStats - a software tool to investigate heterozygous snp variance\n" +
			"Usage:\n" +
			"  ./alleleStats [options] input.sam input.vcf\n\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 2
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	var f1Genome *string = flag.String("f1", "", "F1 hybrid sample that appears heterozygous in genotype Vcf``")
	var sampleName *bool = flag.Bool("samples", false, "Get names of samples that appear in Vcf header (Default: /dev/stdout``)")
	var parentOne *string = flag.String("parentOne", "", "Name of first parental genome``")
	var parentTwo *string = flag.String("parentTwo", "", "Name of second parental genome``")
	var counts *bool = flag.Bool("counts", false, "Get allele depth counts for all samples")
	var fisherTest *bool = flag.Bool("fishTest", false, "Perform fisher's exact test for p-value significance between two samples include vcf followed by two sample names that appear in vcf")
	flag.Parse()
	if len(flag.Args()) == 1 {
		file := vcf.NewReader(flag.Arg(0))
		header := vcf.ReadHeader(file)
		if *sampleName {
			if strings.HasSuffix(flag.Arg(0), "vcf.gz") || strings.HasSuffix(flag.Arg(0), ".vcf") {
				fmt.Printf("%s", vcf.PrintSampleNames(header))
			}
		} else if *counts {
			buf := &strings.Builder{}
			//buf.WriteString("Chrom\tPos\tRef,Alt")
			for key := range header.Samples {
				buf.WriteByte('\t')
				buf.WriteString(key)
			}
			fmt.Printf("%s\n", buf.String())
			for i, done := vcf.UnmarshalVcf(file); !done; i, done = vcf.UnmarshalVcf(file) {
				buf.Reset()
				writeLine(i, buf)
				buf.WriteString(strings.Join(vcf.GetAllAlleleDepth(i), "\t"))
				fmt.Printf("%s\n", buf.String())
			}
		} else {
			flag.Usage()
			log.Fatalf("\nExamples:\n./alleleStats -f1 name -parentOne name -parentTwo name input.sam input.vcf\n\nView sample names:\n./alleleSplit -samples file.vcf\n\nRun fisher's exact test:\n./alleleStats -fishTest file.vcf.gz wgs atac/rna\n\n")
		}
	} else if *fisherTest {

		if len(flag.Args()) != 4 {
			log.Fatalf("Error: expecting 4 arguments: vcf, name, name\n./alleleStats -fishTest file.vcf.gz wgs atac/rna\n")
		}

		file := vcf.NewReader(flag.Arg(0))
		header := vcf.ReadHeader(file)
		wgs := header.Samples[flag.Arg(1)]
		//fmt.Printf("wgs:%s=%d\n", flag.Arg(1), header.Samples[flag.Arg(1)])
		hybrid := header.Samples[flag.Arg(2)]
		//fmt.Printf("hybrid:%s=%d\n", flag.Arg(2), header.Samples[flag.Arg(2)])

		results := stats.Results{}
		var sampleOne, sampleTwo alleleDepth
		var significant []stats.DiffPeak

		output := simpleio.NewWriter(flag.Arg(3))
		defer output.Close()
		for v, done := vcf.UnmarshalVcf(file); !done; v, done = vcf.UnmarshalVcf(file) {
			sampleOne, sampleTwo = getGtTwo(v, wgs, hybrid)
			if vcf.IsHeterozygous(sampleOne.Gt) {

				fishingTest := stats.FisherExact{
					A: sampleOne.Depth[sampleOne.Gt.AlleleOne],
					B: sampleOne.Depth[sampleOne.Gt.AlleleTwo],
					C: sampleTwo.Depth[sampleOne.Gt.AlleleOne],
					D: sampleTwo.Depth[sampleOne.Gt.AlleleTwo],
				}
				if equalWithin(fishingTest.A, fishingTest.B) && (fishingTest.C > -1 && fishingTest.D > -1) {
					results = stats.RunFishersExactTest(fishingTest)
					if results.Pval < 0.0001 {
						diff := stats.DiffPeak{
							Chr:    v.Chr,
							Start:  v.Pos,
							Matrix: fishingTest,
							Pval:   results.Pval,
						}
						significant = append(significant, diff)
					}

				}
			}
		}
		stats.SortByPValue(significant)
		//output.WriteString("chr,start,A,B,C,D,Pvalue\n")
		for _, snp := range significant {
			writeDiff(output, snp)
		}
	} else if len(flag.Args()) != expectedNumArgs || (*f1Genome == "" && *parentOne == "" || *parentTwo == "") || !*counts {
		flag.Usage()
		fmt.Printf("\nExamples:\n./alleleStats -f1 name -parentOne name -parentTwo name input.sam input.vcf\n\nView sample names:\n./alleleSplit -samples file.vcf\n\nRun fisher's exact test:\n./alleleStats -fishTest file.vcf.gz wgs atac/rna\n\n")
		log.Fatalf("\n\nError: unexpected number of arguments...\n\n")
	} else {
		SnpSearch(flag.Arg(0), flag.Arg(1), *f1Genome, *parentOne, *parentTwo, *f1Genome)
	}
}

func equalWithin(a, b int) bool {
	if a < b {
		return float64(a)/float64(b) > 0.89
	} else {
		return float64(a)/float64(b) < 1.12
	}
}

func writeDiff(writer *simpleio.SimpleWriter, diff stats.DiffPeak) {
	writer.WriteString(diff.Chr)
	writer.WriteByte('\t')
	writer.WriteString(simpleio.IntToString(diff.Start - 1))
	writer.WriteByte('\t')
	writer.WriteString(simpleio.IntToString(diff.Start))
	writer.WriteByte('\t')
	writer.WriteString(stats.StringFmt(&diff.Matrix))
	writer.WriteByte('\t')
	writer.WriteString(fmt.Sprintf("%E\n", diff.Pval))

}
func GoFishingString(v *vcf.Vcf, results stats.Results) string {
	buf := strings.Builder{}
	buf.WriteString(v.Chr)
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(v.Pos))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(v.Pos + 1))
	buf.WriteByte('\t')
	buf.WriteString(stats.StringFmt(&results.Matrix))
	buf.WriteByte('\t')

	buf.WriteString(fmt.Sprintf("%E\n", results.Pval))
	return buf.String()
}

type alleleDepth struct {
	Gt    vcf.Genotype
	Depth []int
}

func getAdGtAll(v *vcf.Vcf) []alleleDepth {
	format := vcf.MkFormatMap(v)
	ans := make([]alleleDepth, len(v.Genotypes))
	var line []string = make([]string, len(v.Format))
	for i := 0; i < len(ans); i++ {
		line = strings.Split(v.Genotypes[i], ":")
		ans[i].Gt = vcf.ParseGt(line[0])
		ans[i].Depth = simpleio.StringToIntSlice(vcf.FindAlleleDepth(format, line))
	}
	return ans
}

func getGtTwo(v *vcf.Vcf, sampleOne int, sampleTwo int) (alleleDepth, alleleDepth) {
	format := vcf.MkFormatMap(v)
	var line []string = make([]string, len(v.Format))

	wgs := alleleDepth{}
	line = strings.Split(v.Genotypes[sampleOne], ":")
	wgs.Gt = vcf.ParseGt(line[0])
	wgs.Depth = stringToIntSlice(vcf.FindAlleleDepth(format, line))

	hybrid := alleleDepth{}
	line = strings.Split(v.Genotypes[sampleTwo], ":")
	hybrid.Gt = vcf.ParseGt(line[0])
	hybrid.Depth = stringToIntSlice(vcf.FindAlleleDepth(format, line))
	return wgs, hybrid
}

// StringToInts will process strings (usually from column data) and return a slice of []int
func stringToIntSlice(column string) []int {
	work := strings.Split(column, ",")
	var answer []int = make([]int, len(work))
	for i := 0; i < len(work); i++ {
		answer[i] = simpleio.StringToInt(work[i])
	}
	return answer
}

func addTwoGt(v *vcf.Vcf, sampleOne int, sampleTwo int, sampleThree int) (alleleDepth, alleleDepth) {
	format := vcf.MkFormatMap(v)
	var line []string = make([]string, len(v.Format))

	wgs := alleleDepth{}
	line = strings.Split(v.Genotypes[sampleOne], ":")
	wgs.Gt = vcf.ParseGt(line[0])
	wgs.Depth = simpleio.StringToIntSlice(vcf.FindAlleleDepth(format, line))

	hybrid := alleleDepth{}
	line = strings.Split(v.Genotypes[sampleTwo], ":")
	hybrid.Gt = vcf.ParseGt(line[0])
	hybrid.Depth = simpleio.StringToIntSlice(vcf.FindAlleleDepth(format, line))

	hybridTwo := alleleDepth{}
	line = strings.Split(v.Genotypes[sampleThree], ":")
	hybridTwo.Gt = vcf.ParseGt(line[0])
	hybridTwo.Depth = simpleio.StringToIntSlice(vcf.FindAlleleDepth(format, line))

	for i := 0; i < len(hybridTwo.Depth); i++ {
		hybrid.Depth[i] += hybridTwo.Depth[i]
	}
	return wgs, hybrid
}

func writeLine(v *vcf.Vcf, buf *strings.Builder) {
	buf.WriteString(v.Chr)
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(v.Pos))
	buf.WriteByte('\t')
	buf.WriteString(v.Ref)
	buf.WriteByte(',')
	buf.WriteString(v.Alt)
	buf.WriteByte('\t')
}

//func sampleStats(v *vcf.Vcf, buf *strings.Builder, index int) {
//buf.WriteString()
//}
