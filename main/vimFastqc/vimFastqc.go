package main

import (
	"flag"
	"fmt"
	"log"
	//"github.com/goFish/simpleio"
	"github.com/edotau/goFish/fastq"
)

func usage() {
	fmt.Print(
		"vimFastqc - quickly access fastq quality control metrics...\n\n")
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

	fqs := fastq.FastqReader(flag.Arg(0))
	var numReads int = 0
	//log.Printf("Average read length is: %f...\n", fastq.FindAveReadLength(fqs))
	for each := range fqs {
		fmt.Printf("%s\n", fastq.MetricsTable(&each))
		numReads++
	}
	fmt.Printf("Total Reads Process: %d...\n", numReads)
	//data := simpleio.NewWriter(flag.Arg(1))
	//fastq.MetricsTable(fqs)
}
