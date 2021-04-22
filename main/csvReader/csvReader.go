package main

/*
import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/vertgenlab/gonomics/common"
	"github.com/vertgenlab/gonomics/fileio"
	"github.com/vertgenlab/gonomics/numbers"
	"io"
	"log"
	"strings"
	"github.com/goFish/stats"
)

func usage() {
	fmt.Print(
		"Gvcf to FisherExact test\n")
	flag.PrintDefaults()
}
*/
func main() {

}
/*
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	input := flag.Arg(0)
	csvFile := fileio.EasyOpen(input)
	reader := csv.NewReader(csvFile.File)

	var i []string
	var err error
	i, err = reader.Read()
	common.ExitIfError(err)

	if !stats.IsHeaderLine(i) {
		log.Fatalf("Error: check header line of input csv file...\n")
	}
	a, b, c, d := atacRnaMatrix(i)
	prefix := strings.TrimSuffix(input, ".csv")
	leftTest, rightTest := make([]*stats.DiffPeak, 0), make([]*stats.DiffPeak, 0)

	var curr *stats.PeakStats = &stats.PeakStats{}
	for {
		i, err = reader.Read()
		if err == nil {
			curr = stats.GetStats(i, common.StringToInt(i[a]), common.StringToInt(i[b]), common.StringToInt(i[c]), common.StringToInt(i[d]))
			if !stats.IsEmpty(curr.Matrix) {
				curr.LeftPvalue = numbers.FisherExact(curr.Matrix.A, curr.Matrix.B, curr.Matrix.C, curr.Matrix.D, true)
				curr.RightPvalue = numbers.FisherExact(curr.Matrix.A, curr.Matrix.B, curr.Matrix.C, curr.Matrix.D, false)
				if curr.LeftPvalue < 0.05 && curr.RightPvalue < 0.05 {
					log.Fatalf("Error: both left and right cannot be significant...\n")
				} else {
					if curr.LeftPvalue < 0.05 {
						leftTest = append(leftTest, &stats.DiffPeak{Chr: curr.Chr, Start: curr.Start, End: curr.End, Pval: curr.LeftPvalue})
					} else if curr.RightPvalue < 0.05 {
						rightTest = append(rightTest, &stats.DiffPeak{Chr: curr.Chr, Start: curr.Start, End: curr.End, Pval: curr.RightPvalue})
					}
				}
			}
		} else if err == io.EOF {
			break
		} else {
			log.Fatal(err)
		}
	}
	csvFile.Close()
	stats.WriteDiffPeaks(prefix+"_FishersExactLeft.csv", leftTest)
	stats.WriteDiffPeaks(prefix+"_FishersExactRight.csv", rightTest)
}

func isMarine(name string) bool {
	if strings.Contains(name, "RABS") || strings.Contains(name, "LITC") {
		return true
	} else {
		return false
	}
}

func isFreshwater(name string) bool {
	if strings.Contains(name, "BEPA") || strings.Contains(name, "MATA") {
		return true
	} else {
		return false
	}
}

func atacRnaMatrix(csvLine []string) (int, int, int, int) {
	var a, b, c, d int
	for i := 4; i < len(csvLine); i++ {
		if isFreshwater(csvLine[i]) {
			if strings.Contains(csvLine[i], "atac") || strings.Contains(csvLine[i], "rnaseq") {
				a = i
			} else if strings.Contains(csvLine[i], "wgs") {
				b = i
			} else {
				log.Fatalf("Error: did not find freshwater atac or wgs sample...\n")
			}
		} else if isMarine(csvLine[i]) {
			if strings.Contains(csvLine[i], "atac") || strings.Contains(csvLine[i], "rnaseq") {
				c = i
			} else if strings.Contains(csvLine[i], "wgs") {
				d = i
			} else {
				log.Fatalf("Error: did not find marine atac or wgs sample...\n")
			}
		} else {
			log.Fatalf("Error: unknown freshwater or marine name...\n")
		}
	}
	return a, b, c, d
}*/
