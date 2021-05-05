package geneSeq

import (
	"log"
	"testing"
)

func TestGtfReading(t *testing.T) {
	gf, _ := ParseGtfLine("testdata/denovo.reference.transcrtipt.assembly.annotated.gtf")
	for i := 0; i < 10; i++ {
		log.Printf("%v\n", gf[i])
	}

}
