package genePred

import (
	"github.com/edotau/goFish/simpleio"
	"strings"
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
		work = strings.Split(i.String(), "\t")
		hash[work[0]] = work[1]
	}
	return hash
}
