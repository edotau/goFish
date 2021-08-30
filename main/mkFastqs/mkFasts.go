package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/fastq"
	"github.com/edotau/goFish/simpleio"
)

func usage() {
	fmt.Print(
		"")
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

	reads := readsNNN(1000)
	writer := simpleio.NewWriter(flag.Arg(0))
	defer writer.Close()
	for i := 0; i < len(reads); i++ {
		fastq.WriteFastq(&reads[i], writer)
	}
	writer.WriteByte('\n')

}

func readsNNN(length int) []fastq.Fastq {
	var ans []fastq.Fastq = make([]fastq.Fastq, length)
	for i := 0; i < length; i++ {
		ans[i] = newFastq(150)
		ans[i].Name = "NNNNNs"
	}
	return ans
}

func newFastq(length int) fastq.Fastq {
	var ans fastq.Fastq
	ans.Seq, ans.Qual = make([]code.Dna, length), make([]byte, length)
	for i := 0; i < length; i++ {
		ans.Seq[i], ans.Qual[i] = code.N, 'J'
	}
	return ans
}
