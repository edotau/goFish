// geneSeQueryL processes data tables from gtf, gff3, and genePred gene feture formats, process custom data frames, and sql querys from ensembl and ucsc databases
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/edotau/goFish/geneSeq"
	"github.com/edotau/goFish/simpleio"
)

func usage() {
	fmt.Print(
		"geneSeQuery - processes data tables from gtf, gff3, genePred, and related sql queries to produce custom data frames\n" +
			"Usage:\n" +
			"   ./geneSeQueryL geneSeqFmt.input output.txt\n" +
			"Options:\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 2
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	attribs := flag.Bool("attributes", false, "generates a table containing the attibutes column")
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	geneFeatFmt := geneSeq.ReadGtf(flag.Arg(0))
	output := simpleio.NewWriter(flag.Arg(1))
	if *attribs {
		for i := 0; i < len(geneFeatFmt); i++ {
			io.Copy(output, tabDelimAt(geneFeatFmt[i].Attributes))
		}

	}
	output.Close()
}

func tabDelimAt(a []geneSeq.Attribute) *bytes.Buffer {
	var buf bytes.Buffer = bytes.Buffer{}
	for i := 0; i < len(a)-1; i++ {
		buf.WriteString(a[i].Tag)
		buf.WriteByte('=')
		buf.WriteString(a[i].Value)
		buf.WriteByte('\t')
	}
	buf.WriteString(a[len(a)-1].Tag)
	buf.WriteByte('=')
	buf.WriteString(a[len(a)-1].Value)
	buf.WriteByte('\n')
	return &buf
}

/*
func BuildHashFromEnsembl(filename string) map[string][]geneSeq.Attribute {
	reader := simpleio.NewReader(filename)
	var col []string
	hash := make(map[string][]geneSeq.Attribute)
	for i, done := simpleio.ReadLine(reader); !done; i, done = simpleio.ReadLine(reader) {
		col = strings.Split(i.String(), "\t")
		//transcript_ID
		hash[col[0]] = append(hash[col[0]], geneSeq.Attribute{"gene_id", col[1]})
		hash[col[0]] = append(hash[col[0]], geneSeq.Attribute{"gene_name", col[8]})
		hash[col[0]] = append(hash[col[0]], geneSeq.Attribute{"gene_biotype", col[8]})
	}
}*/
