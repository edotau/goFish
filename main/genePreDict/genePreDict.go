// genePredDict is a software tool used to operate (concat, find, query or and colapse dups) on ucsc geneSeq format
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/edotau/goFish/bed"
	"github.com/edotau/goFish/geneSeq"
	"github.com/edotau/goFish/reference/stickleback"
	"github.com/edotau/goFish/simpleio"
)

func usage() {
	fmt.Print(
		"genePreDict - software tools to operate on ucsc genePred format\n  Usage: ./genePredDict target.gp query.gp\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	remove := flag.Bool("remove", false, "collapse all genePred records")
	find := flag.Bool("find", false, "find query nonoverlaping genePred with target")
	cat := flag.Bool("cat", false, "concat genePred Files")
	rename := flag.Bool("gene-names", false, "switch out prediction id's with known gene symbols")
	flag.Parse()

	if len(flag.Args()) < expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	if *find && *remove {
		log.Fatalf("Error: remove and find flags cannot be given at the same time...\n")
	} else if *rename {
		toGeneNames(flag.Arg(0), flag.Arg(1))
	} else if *find {
		findNonOverlap(flag.Arg(0), flag.Arg(1))
	} else if *remove {
		reduceCapacity(flag.Arg(0))
	} else if *cat {
		concat(flag.Args())
	} else {
		flag.Usage()
	}

}

func concat(files []string) {
	var ans []geneSeq.GenePred
	for i := 0; i < len(files); i++ {
		ans = append(ans, geneSeq.Read(files[i])...)
	}
	geneSeq.QuickSort(ans)
	ans = geneSeq.RmOverlap(ans)

	for j := 0; j < len(ans); j++ {
		fmt.Printf("%s\n", ans[j].ToString())
	}
}

func reduceCapacity(filename string) {
	geneModel := geneSeq.Read(filename)
	geneSeq.QuickSort(geneModel)
	geneModel = geneSeq.RmOverlap(geneModel)
	for _, i := range geneModel {
		fmt.Printf("%s\n", i.ToString())
	}
	//geneSeq.WriteGenePred(output, geneModel)
}

func findNonOverlap(t string, q string) {
	target := geneSeq.ReadToMap(t)
	var curr []*geneSeq.GenePred

	reader := simpleio.NewReader(q)
	var i int
	var nonOverlap bool = false
	for query, err := geneSeq.GenePredLine(reader); !err; query, err = geneSeq.GenePredLine(reader) {
		curr = target[query.Chr]
		nonOverlap = false
		for i = 0; i < len(curr); i++ {
			if bed.Overlap(curr[i], query) {
				nonOverlap = true
			}
		}
		if !nonOverlap {
			fmt.Printf("%s\n", query.ToString())
		}

	}
}

func toGeneNames(input, output string) {
	gp := geneSeq.Read(input)

	writer := simpleio.NewWriter(output)

	geneNames := geneSeq.ReadBioMart(stickleback.ENSEMBL_GENE_NAMES)

	for _, i := range gp {

		curr := i
		gene, ok := geneNames[i.GeneName]
		if ok {
			i.GeneName = strings.ReplaceAll(geneNames[i.GeneName], " ", "_")
			curr.GeneName = strings.ReplaceAll(gene, " ", "_")
			writer.Writer.WriteString(curr.ToString())
			writer.Writer.WriteByte('\n')
		} else {
			writer.Writer.WriteString(i.ToString())
			writer.Writer.WriteByte('\n')
		}

	}
	writer.Close()
}

/*
	geneModels := geneSeq.FilterStrand(geneSeq.UniqueGenes(geneSeq.ReadToUniqueMap(flag.Arg(0))))
	geneSeq.QuickSort(geneModels)
	for _, i := range geneModels {
		fmt.Print(geneSeq.ToString(&i))
	}*/
