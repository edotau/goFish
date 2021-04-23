package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/chain"
	"log"
)

func usage() {
	fmt.Print(
		"")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	//file := flag.Arg(0)
	ans := chain.ReadAll(flag.Args())
	for _, i := range ans {
		if Inversion(&i) {
			fmt.Printf("%s\n", chain.PrettyFmt(&i))
		}

	}
}

func Inversion(c *chain.Chain) bool {
	return c.TStrand != c.QStrand
}
