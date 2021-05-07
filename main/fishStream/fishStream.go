package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/fasta"
	"github.com/edotau/goFish/geneSeq"
	"github.com/edotau/goFish/reference/stickleback"
	"github.com/edotau/goFish/simpleio"
	"io"
	"log"
	"strings"
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
	var fa *bool = flag.Bool("genome", false, "fetch marine fasta reference genome to stdout")
	var chrom *bool = flag.Bool("chrom", false, "print marine stickleback chrom size info to stdout")
	var wget *bool = flag.Bool("wget", false, "download stickleback genome as a fasta to disk")
	var genes *bool = flag.Bool("gene-info", false, "fetch gene-prediction models to stdout")
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()
	if *wget {
		if *fa {
			fetchHttpStdout()
		}
		if *chrom {
			writer := simpleio.NewWriter("rabsTHREEspine.chromSize.gz")
			writer.Write([]byte(chromTableStdout()))
			writer.Close()
		}
		if *genes {
			stream := simpleio.NewReader(stickleback.GENE_MODEL_RNASEQ)
			defer stream.Close()
			writer := simpleio.NewWriter("rabsTHREEspine.rna-seq.genes.mapped.ensembl.gp.gz")
			io.Copy(writer.Gzip, stream)

			writer.Close()
		}
		return
	}
	if *fa {
		fetchHttpStdout()
	} else if *wget {
		wgetFasta()
	} else if *chrom {
		fmt.Printf("%s\n", chromTableStdout())
	} else if *genes {
		geneModels()
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

func chromTableStdout() string {
	buf := &strings.Builder{}
	for _, i := range stickleback.Chr {
		buf.WriteString(i)
		buf.WriteByte('\t')
		buf.WriteString(simpleio.IntToString(stickleback.GetChrom(i)))
		buf.WriteByte('\n')
	}
	return buf.String()
}

func geneModels() {
	stream := simpleio.NewReader(stickleback.GENE_MODEL_RNASEQ)
	for gene, done := geneSeq.GenePredLine(stream); !done; gene, done = geneSeq.GenePredLine(stream) {
		fmt.Printf("%s\n", gene.ToString())

	}
	stream.Close()

}
