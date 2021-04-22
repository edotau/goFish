package main

import (
	"flag"
	"fmt"
	"github.com/goFish/genePred"
	"github.com/goFish/simpleio"
	"log"
)

func usage() {
	fmt.Print(
		"geneSym - a small program to switch out prediction id's with known gene symbols.\nUsage:\n  ./geneSym genePred.gp ensemblToGenes.txt output.gp")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 3
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	gp := genePred.Read(flag.Arg(0))
	geneNames := genePred.ReadBioMart(flag.Arg(1))
	writer := simpleio.NewWriter(flag.Arg(2))
	for _, i := range gp {

		curr := i
		gene, ok := geneNames[i.GeneName]
		if ok {
			i.GeneName = geneNames[i.GeneName]
			curr.GeneName = gene
			writer.WriteLine(genePred.ToString(&curr))
		} else {
			writer.WriteLine(genePred.ToString(&i))
		}

	}

}
