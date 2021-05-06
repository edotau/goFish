package geneSeq

import (
	"fmt"
	"sort"
	"strings"

	"github.com/edotau/goFish/bed"
	"github.com/edotau/goFish/simpleio"
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

func ReadFilter(mapFile string, gpFile string) {
	reader := simpleio.NewReader(mapFile)

	var column []string
	mapOfGenes := make(map[string]string)

	for line, done := simpleio.ReadLine(reader); !done; line, done = simpleio.ReadLine(reader) {
		column = strings.Split(line.String(), "\t")
		mapOfGenes[column[10]] = column[0]
	}

	for _, gp := range Read(gpFile) {
		curr := gp
		curr.GeneName = mapOfGenes[gp.GeneName]
		fmt.Printf("%s\n", curr.ToString())
	}
}

func FilterPitx1(filename string) {
	genes := Read(filename)
	pitx1 := bed.Simple{Chr: "chr07", Start: 386418, End: 1051237}
	for _, i := range genes {
		if bed.Overlap(&pitx1, &i) {
			fmt.Printf("%s\n", i.ToString())
		}
	}
}
