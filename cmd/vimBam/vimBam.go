// vimBam
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/edotau/goFish/bam"
	//"github.com/edotau/goFish/"
)

func usage() {
	fmt.Print(
		"vimBam - view option of samtools integrated with other features from goFish\n" +
			"  Usage:\n" +
			"./vimBam [options] align.bam\n" +
			"options:\n")
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

	header, alignments := bam.Read(flag.Arg(0))
	fmt.Printf("%s\n", header.Text.String())
	for i := range alignments {
		fmt.Printf("%s\n", bam.ToString(&i))
	}
}
