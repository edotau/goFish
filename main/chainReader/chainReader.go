package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/chain"
	"log"
)

func usage() {
	fmt.Print(
		"chainReader - general tool to example chain alignments between two whole genomes\n\n" +
			"Usage:\t" +
			"  chainReader [options] genome.chain, ...\n\n" +
			"Options:\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	var inv = flag.Bool("inversion", false, "find chain regions that span inversions")
	var basic = flag.Bool("default", false, "return a pretty print of all chain regions")

	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
		if !*basic {

		}
	}

	//file := flag.Arg(0)
	ans := chain.ReadAll(flag.Args())
	for _, i := range ans {
		if *inv {
			if Inversion(&i) {
				fmt.Printf("%s\n", chain.PrettyFmt(&i))
			}
		} else {
			fmt.Printf("%s\n", chain.PrettyFmt(&i))
		}
	}
}

func Inversion(c *chain.Chain) bool {
	return c.TStrand != c.QStrand
}
