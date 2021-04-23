package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/genePred"
	"log"
)

func usage() {
	fmt.Print(
		"uniqGenePred - locate and/or filter the gene prediction with the most number of exons.\n  Usage: ./uniqGenePred input.gp > output.gp\n\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	geneModels := genePred.FilterStrand(genePred.UniqueGenes(genePred.ReadToUniqueMap(flag.Arg(0))))
	genePred.QuickSort(geneModels)
	for _, i := range geneModels {
		fmt.Print(genePred.ToString(&i))
	}
}
