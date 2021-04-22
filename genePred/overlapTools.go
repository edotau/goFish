package genePred

import (
	"fmt"
	"github.com/goFish/bed"
	"github.com/goFish/simpleio"
	"sort"
	"strings"
)

func QuickSort(file []GenePred) {
	sort.Slice(file, func(i, j int) bool { return bed.Compare(&file[i], &file[j]) == -1 })
}

func RmOverlap(file []GenePred) []GenePred {
	var i, j int
	for i = 0; i < len(file)-1; {
		if !bed.Overlap(&file[i], &file[i+1]) {
			i++
		} else {
			if bed.Len(&file[i]) < bed.Len(&file[i+1]) {
				file[i] = file[i+1]
			}
			//file[i].ChromStart, file[i].ChromEnd, file[i].Score = numbers.MinInt64(file[i].ChromStart, file[i+1].ChromStart), numbers.MaxInt64(file[i].ChromEnd, file[i+1].ChromEnd), file[i].Score+file[i+1].Score
			for j = i + 1; j < len(file)-1; j++ {
				file[j] = file[j+1]
			}
			file = file[:len(file)-1]
		}
	}
	return file
}

func compareExonLen(gp []*GenePred) *GenePred {
	var ans *GenePred = gp[0]
	for i := 1; i < len(gp); i++ {
		if gp[i].ExonCount > ans.ExonCount {
			ans = gp[i]
		}
	}
	return ans
}

func getUniqName(gp *GenePred) string {
	work := strings.Split(gp.GeneName, ".")
	return fmt.Sprintf("%s.%s", work[0], work[1])
}

func ReadToUniqueMap(filename string) map[string][]*GenePred {
	ans := make(map[string][]*GenePred)
	reader := simpleio.NewReader(filename)
	var code string
	for i, err := GenePredLine(reader); !err; i, err = GenePredLine(reader) {
		code = getUniqName(i)
		ans[code] = append(ans[code], i)
	}
	return ans
}

func UniqueGenes(gp map[string][]*GenePred) []GenePred {
	var ans []GenePred
	for _, elements := range gp {
		ans = append(ans, *compareExonLen(elements))
	}
	return ans
}
