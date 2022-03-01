package main

import (
	"fmt"
	"testing"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/fastq"
	"github.com/edotau/goFish/stats"
)

func TestMkFastq(t *testing.T) {
	fq := fastq.NewFastq(150)
	fq.Name = "TouchFastq.gz"
	for i := 0; i < len(fq.Seq); i++ {
		fq.Seq[i] = code.NoMaskDnaArray[stats.RandIntInRange(0, len(code.NoMaskDnaArray))]
	}
	if len(fq.Seq) != len(fq.Qual) {
		t.Errorf("Error: fastqs should have one phred qual score for each base...\n")
	}
	fmt.Printf("%s\n", fq.ToString())
}
