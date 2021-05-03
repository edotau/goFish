// goToBed is a simplified version of ucsc overlap section to analyze non/overlapping genomic regions in a dataset and filtering var variants of interests
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/edotau/goFish/bed"
	"github.com/edotau/goFish/simpleio"
)

func usage() {
	fmt.Print(
		"goToBed - software toolkit to analyze overlapping genomic regions in a dataset\n" +
			"Usage:\n" +
			"  ./goToBed [options] in.file out.file\n\n" +
			"Default:\n  overlapSelect\n\n" +
			"Options:\n")
	flag.PrintDefaults()
	fmt.Print("\n")
}

func main() {
	var nonoverlap *bool = flag.Bool("nonoverlap", false, "find nonoverlapping genomic regions")
	var concatFiles *bool = flag.Bool("concat", false, "merge input bed files into a unique set that filters out duplicate regions that overlap each other")

	var filterSv *string = flag.String("variant", "", "``filter by a specific structure variant [INS or DEL]")
	var expectedNumArgs int = 2

	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if *concatFiles {
		concat(flag.Args())
		return
	} else if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	keys, selectRegions := bed.SelectGenomeHash(flag.Arg(0), 10)
	reader := simpleio.NewReader(flag.Arg(1))

	if *nonoverlap && *filterSv != "" {
		for i, err := bed.ToGenomeInfo(reader); !err; i, err = bed.ToGenomeInfo(reader) {
			hashKey := bed.GetHashKey(i, keys)
			if !bed.OverlapHashFilterSv(hashKey, i, selectRegions, *filterSv) {
				fmt.Printf("%s\n", bed.GenomeInfoToString(*i))
			}
		}
	} else if *filterSv != "" {
		for i, err := bed.ToGenomeInfo(reader); !err; i, err = bed.ToGenomeInfo(reader) {
			hashKey := bed.GetHashKey(i, keys)
			if bed.OverlapHashFilterSv(hashKey, i, selectRegions, *filterSv) {
				fmt.Printf("%s\n", bed.GenomeInfoToString(*i))
			}
		}
	} else if *nonoverlap {
		for i, err := bed.ToGenomeInfo(reader); !err; i, err = bed.ToGenomeInfo(reader) {
			hashKey := bed.GetHashKey(i, keys)
			if !bed.CheckOverlapHash(hashKey, i, selectRegions) {
				fmt.Printf("%s\n", bed.GenomeInfoToString(*i))
			}
		}
	} else {
		for i, err := bed.ToGenomeInfo(reader); !err; i, err = bed.ToGenomeInfo(reader) {
			hashKey := bed.GetHashKey(i, keys)
			if bed.CheckOverlapHash(hashKey, i, selectRegions) {
				fmt.Printf("%s\n", bed.GenomeInfoToString(*i))
			}
		}
	}

}

func GetSv(b *bed.GenomeInfo) string {

	sv := strings.Split(b.Info.String(), "\t")[0]
	fmt.Printf("%s\n", sv)
	return sv
}

func SvMap(v []bed.StructureVariance) map[string][]*bed.StructureVariance {
	ans := make(map[string][]*bed.StructureVariance)
	for _, i := range v {
		ans[i.TName] = append(ans[i.TName], &i)
	}
	return ans
}

func reading(variance string, filename string) {
	sv := bed.ReadSvTxt(variance)
	svMap := SvMap(sv)
	reader := simpleio.NewReader(filename)
	var hits []*bed.StructureVariance
	var j int
	for i, err := bed.PeakBedReading(reader); !err; i, err = bed.PeakBedReading(reader) {
		//if
		hits = svMap[i.Chr]
		for j = 0; j < len(hits); j++ {
			if bed.Overlap(hits[j], i) {
				fmt.Printf("%s\n", i.String())
				break
			}
		}
	}
}

func concat(input []string) {
	if len(input) < 1 {
		log.Fatalf("Error: must provide more than 1 bed file to use concat feature...\n")
	}

	keys, selectRegions := bed.SelectGenomeHash(input[0], 10)
	for file := 1; file < len(input); file++ {
		reader := simpleio.NewReader(input[file])
		for i, err := bed.ToGenomeInfo(reader); !err; i, err = bed.ToGenomeInfo(reader) {
			hashKey := bed.GetHashKey(i, keys)
			if bed.CheckOverlapHash(hashKey, i, selectRegions) {
			    fmt.Printf("%s\n", bed.GenomeInfoToString(*i))
			}
			//reader.Close()
		}
	}
}
