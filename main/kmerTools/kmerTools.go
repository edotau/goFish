package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Print(
		"kmerTools - software to find kmer matching strings\n" +
			"Usage:\n" +
			"  ./kmerTools [options] path/diectory\n\n")
	flag.PrintDefaults()
	fmt.Print("\n")
}

func main() {

	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	var kmerLength *int = flag.Int("kmer", 4, "provide a value k for the size of kmer to evaluate``")
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	// Takes a directory as input
	dir := flag.Arg(0)
	// Creates a new folder for challenge1 results
	results := os.MkdirAll(fmt.Sprintf("results_kmer_%d_challenge1", *kmerLength), 0755)
	if results != nil {
		panic(results)
	}

	err := filepath.Walk(dir, func(dirFile string, info os.FileInfo, err error) error {
		if strings.HasSuffix(dirFile, ".fasta") || strings.HasSuffix(dirFile, ".fa") {
			kmer := createKMerHash(ReadFasta(dirFile), *kmerLength)
			sorted := HeapSort(kmer)
			output, writeErr := os.Create(fmt.Sprintf("results_kmer_%d_challenge1/%s_kmer.txt", *kmerLength, strings.Trim(path.Base(dirFile), ".fasta")))
			if writeErr != nil {
				panic(writeErr)
			}
			for _, i := range sorted {
				output.Write([]byte(fmt.Sprintf("%s\t%d\n", i.Seq, i.Count)))
			}
			output.Close()
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// createKMerHash is the main function that takes a list of fasta records and a value k to output a map of kmers n kmers occurrences
func createKMerHash(fa []Fasta, k int) map[string]int {
	// If k is equal to or less than 0, ie k < 1 program will exit and print the error below
	if k < 1 {
		log.Fatalf("Error: kmer length must be at least 1...\n")
	}
	hash := make(map[string]int)
	for _, scaffold := range fa {
		// If k is greater than
		if k > len(scaffold.Seq) {
			log.Fatalf("Error: kmer length must be less than %d...\n", len(scaffold.Seq))
		}
		for i := 0; i < len(scaffold.Seq)-k; i++ {
			hash[DnaToString(scaffold.Seq[i:i+k])]++
		}
	}
	return hash
}
