package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/bed"
	"log"
	//"github.com/vertgenlab/gonomics/"
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

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	file := flag.Arg(0)
	ans := bed.ReadSvTxt(file)
	for _, i := range ans {
		fmt.Printf("%s\n", bed.SvToString(i))
	}
}
