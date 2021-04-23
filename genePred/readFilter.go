package genePred

import (
	"fmt"
	"github.com/edotau/goFish/bed"
	"github.com/edotau/goFish/simpleio"
	"strings"
)

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
		fmt.Printf("%s\n", ToString(&curr))
	}
}

func FilterPitx1(filename string) {
	genes := Read(filename)
	pitx1 := bed.Simple{Chr: "chr07", Start: 386418, End: 1051237}
	for _, i := range genes {
		if bed.Overlap(&pitx1, &i) {
			fmt.Printf("%s\n", ToString(&i))
		}
	}
}
