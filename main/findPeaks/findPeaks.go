package main

func main() {}

/*
import (
	"flag"
	"github.com/goFish/csv"
	"log"
	"strings"

	"fmt"
	"github.com/goFish/bed"
	"github.com/goFish/simpleio"
	"github.com/vertgenlab/gonomics/fileio"
)

func usage() {
	fmt.Print(
		"")
	flag.PrintDefaults()
}

func mainFIrst() {
	csv.RunReader()
}
func main() {
	//csv.RunReader()
	var expectedNumArgs int = 0
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting at least 2 argument, but got %d\n", len(flag.Args()))
	}
	ans := ReadOverlapPvalues(flag.Arg(0), flag.Arg(1))
	var str strings.Builder
	for i := 0; i < len(ans); i++ {
		str.Reset()
		str.WriteString(fmt.Sprintf("%s\t%d\t%d\t%E\t%s\n", ans[i].Chr, ans[i].Start, ans[i].End, ans[i].PValue, ans[i].Name))
		fmt.Printf("%s", str.String())
	}
}

func mainTwo() {
	//csv.RunReader()
	var expectedNumArgs int = 0
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting at least 1 argument, but got %d\n", len(flag.Args()))
	}
	rna := Read(flag.Arg(0))
	peakMap := make(map[int]bool)
	//reader := simpleio.NewSimpleReader(filename)
	atac := Read(flag.Arg(1))
	output := fileio.EasyCreate("final_overlaps_atac_rna.txt")
	var str strings.Builder
	var i, j int
	for i = 0; i < len(rna); i++ {
		for j = 0; j < len(atac); j++ {
			if bed.Overlap(rna[i], atac[j]) {
				if _, ok := peakMap[i]; !ok && rna[i].Start > 100000 {
					str.Reset()
					str.WriteString(fmt.Sprintf("%s\t%d\t%d\t%E\t%s_%s\n", rna[i].Chr, rna[i].Start, rna[i].End, rna[i].PValue, rna[i].Name, atac[j].Name))
					fmt.Fprintf(output, "%s", str.String())
				}
			}
		}
	}
	output.Close()
}

func Read(filename string) []*bed.Pvalue {
	reader := simpleio.NewSimpleReader(filename)
	var rna []*bed.Pvalue
	//peakMap := make(map[int]bool)
	for line, done := bed.ToBedPValue(reader); !done; line, done = bed.ToBedPValue(reader) {
		rna = append(rna, line)
	}
	return rna
}

func ReadOverlapPvalues(atacFile, rnaFile string) []*bed.Pvalue {
	atacReader := simpleio.NewSimpleReader(atacFile)
	var code string
	atacRnaMap := make(map[string][]*bed.Pvalue)
	for i, done := bed.ToBedPValue(atacReader); !done; i, done = bed.ToBedPValue(atacReader) {
		code = fmt.Sprintf("%s_%d", i.Chr, i.Start)
		atacRnaMap[code] = append(atacRnaMap[code], i)
	}
	rnaReader := simpleio.NewSimpleReader(rnaFile)
	for j, done := bed.ToBedPValue(rnaReader); !done; j, done = bed.ToBedPValue(rnaReader) {
		code = fmt.Sprintf("%s_%d", j.Chr, j.Start)
		atacRnaMap[code] = append(atacRnaMap[code], j)
	}
	var answer []*bed.Pvalue

	var nameBuild string
	var avg float64
	for k := range atacRnaMap {
		nameBuild = ""
		avg = 0
		curr := &bed.Pvalue{}
		curr.Chr = atacRnaMap[k][0].Chr
		curr.Start, curr.End = atacRnaMap[k][0].Start, atacRnaMap[k][0].End
		//curr.Name = fmt.Sprintf("%s_%s_%s", atacRnaMap[k][0].Name, atacRnaMap[k][0].Name, atacRnaMap[k][0].Name)
		if len(atacRnaMap[k]) > 2 {

			for l := 0; l < len(atacRnaMap[k]); l++ {
				nameBuild += atacRnaMap[k][l].Name
				avg += atacRnaMap[k][l].PValue
			}
			avg = avg / 3
			curr.PValue = avg
			curr.Name = nameBuild
			answer = append(answer, curr)
		}
	}
	bed.SortByPValue(answer)
	return answer
}*/
