// Package templates contain golang template scripts for quick and easy golang programing
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/fasta"
	"github.com/edotau/goFish/simpleio"
)

func usage() {
	fmt.Print(
		"")
	flag.PrintDefaults()
}

func main() {

	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	fa := fasta.Read(flag.Arg(0))
	for i := 0; i < len(fa); i++ {
		var curr []code.Dna
		for j := 0; j < len(fa[i].Seq); j++ {
			if fa[i].Seq[j] != code.Dot {
				curr = append(curr, fa[i].Seq[j])
			}
		}
		fa[i].Name = fmt.Sprintf("AlphaBeta_%d", i)
		fa[i].Seq = curr

	}
	writer := simpleio.NewWriter("vdjTR.fa")
	defer writer.Close()
	for _, each := range fa {
		writer.Write([]byte(each.ToString()))
	}

}
