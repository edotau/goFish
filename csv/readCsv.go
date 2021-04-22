package csv

/*
import (
	"bufio"
	"encoding/csv"
	"github.com/goFish/simpleio"
	"github.com/goFish/bed"
	"github.com/vertgenlab/gonomics/fileio"
	"fmt"
	"strings"
	"io"
	"os"
)

var overlapFiles = []struct {
	atac string
	rna string
}{

	{atac: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl16_atacseq_FishersExact.csv", rna: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl16_rnaseq_FishersExact.csv"},
	{atac: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl17_atacseq_FishersExact.csv", rna: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl17_rnaseq_FishersExact.csv"},
	{atac: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl18_atacseq_FishersExact.csv", rna: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl18_rnaseq_FishersExact.csv"},
	{atac: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl20_atacseq_FishersExact.csv", rna: "/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/cl20_rnaseq_FishersExact.csv"},
}

type Peak struct {
	Chr    string
	Start  int
	//End    int
	Name   string
	PValue float64
}

func RunReader() {
	output := fileio.EasyCreate("ATAC-Seq_overlappingPeaks_1MB.txt")
	outputTwo := fileio.EasyCreate("RNA-Seq_overlappingPeaks_1MB.txt")
	peakMap := make(map[int]bool)
	peakTwo := make(map[int]bool)
	var str strings.Builder
	var strTwo strings.Builder
	chrMap := ReadChromSize("/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/rabsTHREEspine.sizes")

	//chromInfo.ReadToMap("/Users/bulbasaur/software/golang/src/github.com/edotau/csv/finalData/rabsTHREEspine.sizes")
	for _, each := range overlapFiles {
		atac, rna := readFile(each.atac, chrMap), readFile(each.rna, chrMap)
		for atacPeaks := 0; atacPeaks < len(atac); atacPeaks++ {
			for rnaPeaks := 0; rnaPeaks < len(rna); rnaPeaks++ {
				if bed.Overlap(atac[atacPeaks], rna[rnaPeaks]) {
					if _, ok := peakMap[atacPeaks]; !ok {
						str.Reset()
						str.WriteString(fmt.Sprintf("%s\t%d\t%d\t%s_%d\t%E\n", atac[atacPeaks].Chr, atac[atacPeaks].Start, atac[atacPeaks].Start+1, strings.Split(strings.Split(each.atac, "/")[10], "_FishersExact.csv")[0], atacPeaks, atac[atacPeaks].PValue))
						fmt.Fprintf(output, "%s", str.String())
						//fmt.Printf("%s\n", str.String())
						peakMap[atacPeaks] = true
					} else {
						peakMap[atacPeaks] = false
					}
					if _, ok := peakTwo[rnaPeaks]; !ok {
						strTwo.Reset()
						strTwo.WriteString(fmt.Sprintf("%s\t%d\t%d\t%s_%d\t%E\n", rna[rnaPeaks].Chr, rna[rnaPeaks].Start, rna[rnaPeaks].Start+1, strings.Split(strings.Split(each.rna, "/")[10], "_FishersExact.csv")[0], rnaPeaks, rna[rnaPeaks].PValue))
						fmt.Fprintf(outputTwo, "%s", strTwo.String())
						peakTwo[rnaPeaks] = true
					} else {
						peakTwo[rnaPeaks] = false
					}
				}
			}
		}
	}
	output.Close()
	outputTwo.Close()
}

type ChromPeak struct {
	Len int
	Peaks []Peak
}

func ReadChromSize(filename string) map[string]*ChromPeak {
	reader := simpleio.NewReader(filename)
	ans := make(map[string]*ChromPeak)
	var line []string
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		line = strings.Split(i.String(), "\t")
		curr := ChromPeak{
			Len: simpleio.StringToInt(line[1]),
		}
		ans[line[0]] = &curr
	}
	return ans
}
func readFile(filename string, chrMap map[string]*ChromPeak) map[string]*ChromPeak {
	//var answer []Peak
	f, err := os.Open(filename)
	simpleio.ErrorHandle(err)
	reader := csv.NewReader(bufio.NewReader(f))
	header, hErr := reader.Read()
	simpleio.ErrorHandle(hErr)
	fmt.Printf("%v\n", header)
	for {
		if line, err := reader.Read(); err == nil {
			curr := ToPeak(line)
			if curr.Start > 2000000 && curr.Start+1 < int(chrMap[curr.Chr].Len)-2000000 {
				chrMap[curr.Chr].Peaks = append(chrMap[curr.Chr].Peaks, curr)
				answer = append(answer, *curr)
			}
		} else if err == io.EOF {
			break
		} else {
			simpleio.ErrorHandle(err)
		}
	}
	f.Close()
	return answer
}

func ToPeak(line []string) *Peak {

	start := simpleio.StringToInt(line[1])
	if start < 0 {
		start = 0
	}
	return &Peak{
		Chr: line[0],
		Start: start,
		//End: start+1,
		PValue: simpleio.StringToFloat(line[2]),
	}
}

func (bed *Peak) Chrom() string {
	return bed.Chr
}

func (bed *Peak) ChrStart() int {
	return bed.Start
}

func (bed *Peak) ChrEnd() int {
	return bed.Start+1
}*/
