// goIndels - is a quick script to fild indels in sam/bam alignment
package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/bam"
	"log"
)

func usage() {
	fmt.Print(
		"goIndels - software to get indels in sam/bam alignment\nUsage:\n  ./goIndels .sam/.bam\n\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	_, reader := bam.Read(flag.Arg(0))
	var currPos int
	var j int
	for i := range reader {
		currPos = i.Pos
		for j = 0; j < len(i.Cigar); j++ {
			if bam.ConsumesReference(i.Cigar[j].Op) {
				currPos += int(i.Cigar[j].RunLen)
			}
			if i.Cigar[j].Op == bam.Insertion || i.Cigar[j].Op == bam.Deletion {
				fmt.Printf("%s\t%d\t%s\n", i.RName, currPos, string(i.Cigar[j].Op))
			}
		}
	}
}
