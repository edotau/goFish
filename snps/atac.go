package main

import (
	"flag"
	"fmt"
	"github.com/goFish/bed"
	"github.com/goFish/simpleio"
	"log"
	"strings"
	//"github.com/goFish/vcf"
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

	Vim(flag.Arg(0))
}

func Vim(filename string) {
	reader := simpleio.NewReader(filename)
	var work []string
	var atac, rna *bed.Simple
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		work = strings.Split(i.String(), "\t")
		atac, rna = textToSimpleBed(work)
		fmt.Printf("%s\n%s\n", bed.ToString(atac), bed.ToString(rna))
	}
}

func textToSimpleBed(line []string) (*bed.Simple, *bed.Simple) {
	atac := bed.Simple{}
	atac.Chr = line[0]
	atac.Start = simpleio.StringToInt(line[1])
	//atac.End = atac.Start  + simpleio.StringToInt(line[5])
	atac.End = simpleio.StringToInt(line[2])

	rna := bed.Simple{}
	rna.Chr = line[0]
	rna.Start = atac.Start + simpleio.StringToInt(line[5])
	//rna.End = rna.Start  + )
	rna.End = rna.Start + 1
	return &atac, &rna
}

/*
func textToVcf(data []string) *bed.Simple {
	ans := vcf.Vcf{
		Chr: file.data[0],
		Pos: simpleio.StringToInt(data[1]),
		Id: data[2],
		Ref: data[3],
		Alt: data[4],
		Filter: data[6],
		Info: data[7],
		Format: data[8]}
	return &ans
}*/
