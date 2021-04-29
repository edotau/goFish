package genePred

import (
	"strings"

	"github.com/edotau/goFish/simpleio"
)

type GeneSym struct {
	Symbol  string
	Ensembl string
}

func ReadBioMart(filename string) map[string]string {
	//var ans
	reader := simpleio.NewReader(filename)
	var work []string
	hash := make(map[string]string)
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		work = strings.SplitN(i.String(), "\t", 2)
		hash[work[0]] = strings.ReplaceAll(work[1], " ", "_")
	}

	return hash
}
