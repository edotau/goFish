package main

import (
	"flag"
	"fmt"
	"github.com/goFish/simpleio"
	"github.com/vertgenlab/gonomics/common"
	"log"
	"strings"
)

func main() {
	var expectedNumArgs int = 1
	//flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

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
		RefName:   line[10],                    //10
		RefPos:    common.StringToInt(line[0]), //0
		RefSub:    line[1],                     //1
		QueryName: line[11],                    //11
		QueryPos:  common.StringToInt(line[4]), //4
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
