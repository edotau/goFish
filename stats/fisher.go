package stats

import (
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/edotau/goFish/simpleio"
	"github.com/vertgenlab/gonomics/numbers"
)

type PeakStats struct {
	Chr         string
	Start       int
	End         int
	Score       int
	Matrix      *FisherExact
	LeftPvalue  float64
	RightPvalue float64
}

type DiffPeak struct {
	Chr    string
	Start  int
	End    int
	Matrix FisherExact
	Pval   float64
}

type FisherExact struct {
	A int
	B int
	C int
	D int
}

type Results struct {
	Matrix FisherExact
	Pval   float64
}

func RunFishersExactTest(fisherTest FisherExact) Results {
	ans := Results{
		Matrix: fisherTest,
	}
	left := numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, true)
	right := numbers.FisherExact(fisherTest.A, fisherTest.B, fisherTest.C, fisherTest.D, false)
	if left < 0.05 && right < 0.05 {
		log.Fatalf("Error: both left and right fisher's exact tests cannot be significant...\n")
	}
	if left < right {
		ans.Pval = left
	} else {
		ans.Pval = right
	}
	return ans
}

func comparePValue(a DiffPeak, b DiffPeak) int {
	if a.Pval < b.Pval {
		return -1
	}
	if a.Pval > b.Pval {
		return 1
	}
	return 0
}

func SortByPValue(peak []DiffPeak) {
	sort.Slice(peak, func(i, j int) bool { return comparePValue(peak[i], peak[j]) == -1 })
}

func PeakStatsToString(peak PeakStats) string {
	return fmt.Sprintf("%s %d %d\n%s", peak.Chr, peak.Start, peak.End, peak.Matrix.String())
}

func PeakStatsToFile(out io.Writer, peak PeakStats, pValue float64) {
	_, err := fmt.Fprintf(out, "%s,%d,%d,%s,%E\n", peak.Chr, peak.Start, peak.End, StringFmt(peak.Matrix), pValue)
	simpleio.ErrorHandle(err)
}

func PeakStatsSummary(out io.Writer, peak PeakStats, pValue float64) {
	_, err := fmt.Fprintf(out, "%s,%d,%d,%E\n", peak.Chr, peak.Start, peak.End, pValue)
	simpleio.ErrorHandle(err)
}

func WriteDiffPeaks(filename string, peaksDiff []DiffPeak) {
	pvaluePeaks := simpleio.NewWriter(filename)
	SortByPValue(peaksDiff)
	var err error
	for _, peak := range peaksDiff {
		_, err = fmt.Fprintf(pvaluePeaks, "%s,%d,%d,%E\n", peak.Chr, peak.Start, peak.End, peak.Pval)
		simpleio.ErrorHandle(err)
	}
	pvaluePeaks.Close()
}

func GetStats(peak []string, a int, b int, c int, d int) *PeakStats {
	var ans *PeakStats = &PeakStats{
		Chr:   peak[0],
		Start: simpleio.StringToInt(peak[1]),
		End:   simpleio.StringToInt(peak[2]),
		Score: simpleio.StringToInt(peak[3]),
		Matrix: &FisherExact{
			A: a, //fw atac
			B: b, //all fresh, atac + wgs
			C: c, //marine atac
			D: d, //all marine atac + wgs
		},
	}
	return ans
}

func IsHeaderLine(s []string) bool {
	if s[0] == "chr" && s[1] == "start" && s[2] == "end" && s[3] == "score" {
		return true
	} else {
		return false
	}
}

func ParseHeader(header []string) map[int]string {
	hashHeader := make(map[int]string)
	for i := 0; i < len(header); i++ {
		hashHeader[i] = header[i]
	}
	return hashHeader
}

func IsEmpty(fisher *FisherExact) bool {
	if fisher.A == 0 && fisher.B == 0 && fisher.C == 0 && fisher.D == 0 {
		return true
	} else {
		return false
	}
}

func (test *FisherExact) String() string {
	return fmt.Sprintf("%d %d\n%d %d", test.A, test.B, test.C, test.D)
}

func StringFmt(test *FisherExact) string {
	return fmt.Sprintf("%d,%d,%d,%d", test.A, test.B, test.C, test.D)
}
