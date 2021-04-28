package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/edotau/goFish/fasta"
	"github.com/edotau/goFish/reference/stickleback"
	"github.com/edotau/goFish/simpleio"
)

func usage() {
	fmt.Print(
		"\nfishStream - fetch stickleback genome reference files into data streams, which can be processed into data structures in golang or downloaded straight onto disk\n\n" +
			"Usage:\n" +
			"  ./fishStream [options] http://stickleback.io\n" +
			"Options:\n")
	flag.PrintDefaults()
	fmt.Print("\n")

}

func main() {
	var fa *string = flag.String("fetch", "", "``provide a filename to download marine stickleback genome or to stdout")
	var chrom *bool = flag.Bool("chrom", false, "print marine stickleback chrom size info to stdout")
	var wget *bool = flag.Bool("wget", false, "download stickleback genome as a fasta to disk")
	var expectedNumArgs int = 0

	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	if strings.Contains(*fa, "stdout") {
		fetchHttpStdout()
	} else if *wget {
		wgetFasta()
	} else if *chrom {
		chromTableStdout()
	} else {
		flag.Usage()
		log.Fatalf("Error: expecting arguments...\n")
	}
	//fetchHttpStdout()
}

func fetchHttpStdout() {
	stream := simpleio.NewReader(stickleback.RabsFasta)
	for i, err := fasta.FastaReader(stream); !err; i, err = fasta.FastaReader(stream) {
		fmt.Printf("%s", i.ToString())
	}
	stream.Close()
}

func wgetFasta() {
	stream := simpleio.NewReader(stickleback.RabsFasta)
	defer stream.Close()
	writer := simpleio.NewWriter("rabsTHREEspine.fa.gz")
	io.Copy(writer.Gzip, stream)

	writer.Close()
}

func chromTableStdout() {
	buf := &strings.Builder{}
	buf.WriteByte('\n')
	for _, i := range stickleback.Chr {
		buf.WriteString(i)
		buf.WriteByte('\t')
		buf.WriteString(simpleio.IntToString(stickleback.GetChrom(i)))
		buf.WriteByte('\n')
	}
	fmt.Printf("%s\n", buf.String())
}
