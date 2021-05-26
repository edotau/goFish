package main

import (
	"testing"
)

var files []string = []string{
	"challenge1/experiment1.fasta",
	"challenge1/experiment2.fasta",
	"challenge1/experiment3.fasta",
	"challenge1/experiment4.fasta",
}

func TestFastaReader(t *testing.T) {
	for _, test := range files {
		ReadFasta(test)
	}
}
