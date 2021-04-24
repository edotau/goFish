// genePredDict is a software tool used to operate (concat, find, query or and colapse dups) on ucsc genePred format
package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/bed"
	"github.com/edotau/goFish/genePred"
	"github.com/edotau/goFish/simpleio"
	"log"
)

func usage() {
	fmt.Print(
		"genePredDict - software tools to operate on ucsc genePred format\n\t  Usage: ./genePredDict target.gp query.gp\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	remove := flag.Bool("remove", false, "collapse all genePred records")
	find := flag.Bool("find", false, "find query nonoverlaping genePred with target")
	cat := flag.Bool("cat", false, "concat genePred Files")
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	if *find && *remove {
		log.Fatalf("Error: remove and find flags cannot be given at the same time...\n")
	} else if *find {
		findNonOverlap(flag.Arg(0), flag.Arg(1))
	} else if *remove {
		reduceCapacity(flag.Arg(0))
	} else if *cat {
		concat(flag.Args())
	} else {
		flag.Usage()
	}

}

func concat(files []string) {
	var ans []genePred.GenePred
	for i := 0; i < len(files); i++ {
		ans = append(ans, genePred.Read(files[i])...)
	}
	genePred.QuickSort(ans)
	ans = genePred.RmOverlap(ans)

	for j := 0; j < len(ans); j++ {
		fmt.Printf("%s\n", genePred.ToString(&ans[j]))
	}
}

func reduceCapacity(filename string) {
	geneModel := genePred.Read(filename)
	genePred.QuickSort(geneModel)
	geneModel = genePred.RmOverlap(geneModel)
	for _, i := range geneModel {
		fmt.Printf("%s\n", genePred.ToString(&i))
	}
	//genePred.WriteGenePred(output, geneModel)
}

func findNonOverlap(t string, q string) {
	target := genePred.ReadToMap(t)
	var curr []*genePred.GenePred

	reader := simpleio.NewReader(q)
	var i int
	var nonOverlap bool = false
	for query, err := genePred.GenePredLine(reader); !err; query, err = genePred.GenePredLine(reader) {
		curr = target[query.Chr]
		nonOverlap = false
		for i = 0; i < len(curr); i++ {
			if bed.Overlap(curr[i], query) {
				nonOverlap = true
			}
		}
		if !nonOverlap {
			fmt.Printf("%s\n", genePred.ToString(query))
		}

	}
}
