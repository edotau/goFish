package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/fastq"
	"log"
)

func usage() {
	fmt.Print(
		"fastqReader - software to process fastq files io\n  Usage:\n./fastqReader read.fastq.gz\n")
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
	fq := fastq.Read(flag.Arg(0))
	for _, i := range fq {
		if !isValid(&i) {
			log.Printf("Error: fastq record %s is not valid...\n", i.Name)
		}
	}
}

func isValid(read *fastq.Fastq) bool {
	if len(read.Seq) < 1 || len(read.Qual) < 1 {
		log.Printf("Warning: Fastq file does not contain sequence or quality scores...\n")
		return false
	} else {
		return true
	}
}

/*
func compareTrim(raw, trim *fastq.Fastq) bool {

}*/
