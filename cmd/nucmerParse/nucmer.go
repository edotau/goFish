// nucmerParse is used parse and extra interesting genomic regions from nucmer suffix tree alignments
package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/simpleio"
	"log"
	"strings"
)

func usage() {
	fmt.Print(
		"nucmerParse - a tool that will parse the output diff file from nucmer and extract SNPs from genomic regions from nucmer suffix tree alignments\n" +
			"Usage:\n" +
			"  ./nucmerParse in.file\n")
	flag.PrintDefaults()
	fmt.Print("\n")
}

func main() {
	var expectedNumArgs int = 1
	//flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()
	flag.Usage = usage
	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	input := flag.Arg(0)
	parseNucmerSNPs(input)
}

func parseNucmerSNPs(filename string) {
	reader := simpleio.NewReader(filename)
	snp := &nucmerSnp{}
	for curr, done := simpleio.ReadLine(reader); !done; curr, done = simpleio.ReadLine(reader) {
		snp = nucmerLine(curr.Bytes())
		fmt.Printf("%s\n", snp.String())

	}
	reader.Close()
}

func (snp *nucmerSnp) String() string {
	return fmt.Sprintf("%s\t%s\t%d\t%d\t%s\t%s\n", snp.RefName, snp.QueryName, snp.RefPos, snp.QueryPos, snp.RefSub, snp.QuerySub)
}
func nucmerLine(b []byte) *nucmerSnp {
	line := strings.Split(string(b), "\t")
	return &nucmerSnp{
		RefName:   line[10],                      //10
		RefPos:    simpleio.StringToInt(line[0]), //0
		RefSub:    line[1],                       //1
		QueryName: line[11],                      //11
		QueryPos:  simpleio.StringToInt(line[4]), //4
		QuerySub:  line[2],
	}
}

type nucmerSnp struct {
	RefName   string //10
	RefPos    int    //0
	RefSub    string //1
	QueryName string //11
	QueryPos  int    //4
	QuerySub  string //2
}
