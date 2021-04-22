package bed

import (
	"fmt"
	"github.com/goFish/simpleio"
	"sort"
	"strings"
)

type StructureVariance struct {
	Name string
	Len  int

	TName  string
	TStart int
	TEnd   int
	QName  string
	QStart int
	QEnd   int
}

func (sv *StructureVariance) Chrom() string {
	return sv.TName
}

func (sv *StructureVariance) ChrStart() int {
	return sv.TStart
}

func (sv *StructureVariance) ChrEnd() int {
	return sv.TEnd
}

func SvToString(sv StructureVariance) string {
	return fmt.Sprintf("%s\t%d\t%d\t%s\t%d\t%d\t%d\t%s", sv.TName, sv.TStart, sv.TEnd, sv.QName, sv.QStart, sv.QEnd, sv.Len, sv.Name)
}

func ReadSvTxt(filename string) []StructureVariance {
	reader := simpleio.NewReader(filename)
	var data []string
	var curr StructureVariance
	line, done := simpleio.ReadLine(reader)
	var ans []StructureVariance
	for line, done = simpleio.ReadLine(reader); !done; line, done = simpleio.ReadLine(reader) {
		//fmt.Printf("%s\n", line.String())
		data = strings.Split(line.String(), "\t")
		curr = StructureVariance{
			Name:   data[3],
			TName:  data[0],
			TStart: simpleio.StringToInt(data[1]),
			TEnd:   simpleio.StringToInt(data[2]),
			QName:  data[4],
			QStart: simpleio.StringToInt(data[5]),
			QEnd:   simpleio.StringToInt(data[6]),
			Len:    simpleio.StringToInt(data[8]),
		}
		if curr.Len > 1000 {
			ans = append(ans, curr)
			//fmt.Printf("%s\n", SvToString(&curr))
		}
	}
	SortBySvLen(ans)
	return ans
}

func compareSv(a StructureVariance, b StructureVariance) int {
	if a.Len < b.Len {
		return -1
	}
	if a.Len > b.Len {
		return 1
	}
	return 0
}

// SortByPValue performs a Pvalue sort to find the smallest and most significant.
func SortBySvLen(sv []StructureVariance) {
	sort.Slice(sv, func(i, j int) bool { return compareSv(sv[i], sv[j]) == -1 })
}
