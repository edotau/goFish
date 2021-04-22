package main

/*
import (
	"flag"
	"fmt"
	"github.com/vertgenlab/gonomics/fastq"
	"github.com/vertgenlab/gonomics/simpleGraph"
	"log"
)

func usage() {
	fmt.Print(
		"randomFq - simulates random paired end fastq records and outputs to a file\n\n" +
			"Usage:\n" +
			"  randomFq [options] in.fa prefix_R1.fastq.gz prefix_R2.fastq.gz\n\n" +
			"Options:\n\n")
	flag.PrintDefaults()

}
*/
func main() {

}
	/*

    var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	var readLength *int = flag.Int("length", 150, "specify a read length.")
	var numberOfReads *int = flag.Int("number", 10000, "specify how many reads to generate")
	var prefix *string = flag.String("prefix", "simulated_reads", "provide a prefix for file names")
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n\n", expectedNumArgs, len(flag.Args()))
	} else {
		genome := simpleGraph.Read(flag.Arg(0))
		simReads := simpleGraph.RandomPairedReads(genome, *readLength, *numberOfReads, 0)
		fastq.WritePair(*prefix+"_R1.fastq.gz", *prefix+"_R2.fastq.gz", simReads)
	}
}
*/
