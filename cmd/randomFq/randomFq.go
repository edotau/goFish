package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/edotau/goFish/fasta"
	"github.com/edotau/goFish/fastq"
	"github.com/edotau/goFish/simpleio"

	"github.com/edotau/goFish/reference/stickleback"
)

// TODO: program is not finished
func usage() {
	fmt.Print(
		"randomFq - provide a fasta reference to simulate random paired end fastq.gz records and outputs to a file\n\n" +
			"Usage:\n" +
			"  randomFq [options] -fasta in.fa prefix\n\n" +
			"Options:\n\n")
	flag.PrintDefaults()

}

func main() {

	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	var ref *string = flag.String("fasta", stickleback.RabsFasta, "fastq reads are generated from an input fasta reference. If you omit this field, the marine stickleback genome will be used")
	var readLength *int = flag.Int("length", 150, "specify a read length.")
	var numberOfReads *int = flag.Int("number", 10000, "specify how many reads to generate")
	//var prefix *string = flag.String("prefix", "simulated_reads", "provide a prefix for file names")
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n\n", expectedNumArgs, len(flag.Args()))
	} else {
		genome := fasta.Read(*ref)
		readOne := simpleio.NewWriter(fmt.Sprintf("%s_R1.fastq.gz", flag.Arg(0)))
		defer readOne.Close()

		readTwo := simpleio.NewWriter(fmt.Sprintf("%s_R2.fastq.gz", flag.Arg(0)))
		defer readTwo.Close()
		fqs := fastq.RandomPairedReads(genome, *readLength, *numberOfReads, 500)
		for _, i := range fqs {
			readOne.Write(fastq.ToBytes(&i.ReadOne))
			readTwo.Write(fastq.ToBytes(&i.ReadTwo))
		}

		/*
			randomRead := simpleio.NewWriter(fmt.Sprintf("%s_R1.fastq.gz", flag.Arg(0)))
			for i := 0; i < *numberOfReads; i++ {
				fq := fastq.NewFastq(*readLength)
				fq.Name = fmt.Sprintf("%s_GOFISH_%c_R:1", flag.Arg(0), '+')
				for i := 0; i < len(fq.Seq); i++ {
					fq.Seq[i] = code.NoMaskDnaArray[stats.RandIntInRange(0, len(code.NoMaskDnaArray))]

				}
				randomRead.Write(fastq.ToBytes(&fq))
			}
			randomRead.Close()
		*/
		/*

		 */
	}

}
