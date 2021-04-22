package main

import (
	"flag"
	"fmt"
	"github.com/goFish/bed"
	"github.com/goFish/simpleio"
	"github.com/vertgenlab/gonomics/fileio"
	"log"
	"strconv"
	"strings"
)

func usage() {
	fmt.Print(
		"./overlappingAtacRna atac.txt rna.txt\n  find overlapping regions between atac-seq and rna-seq data\n\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 3
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	log.Printf("Reading atac seq file to process...\n")
	rnaSeq := ReadToMap(flag.Arg(0))
	log.Printf("Reading rna seq file to map...\n")
	atacSeq := Read(flag.Arg(1))
	output := fileio.EasyCreate(flag.Arg(2))
	var extend int = 50000

	var regions []*bed.Pvalue
	var j int
	var start, end int
	var err error
	//var atacStart int
	for _, i := range atacSeq {
		regions = rnaSeq[i.Chr]
		for j = 0; j < len(regions); j++ {
			//atacStart = simpleio.StringToInt()
			start = i.Start - extend
			if start < 0 {
				start = 0
			}
			end = i.Start + extend
			if end > getChrom(i.Chr) {
				end = getChrom(i.Chr)
			}
			if bed.Overlapping(regions[j].Chr, i.Chr, regions[j].Start, regions[j].End, start, end) {
				_, err = fmt.Fprintf(output, "%s\n", ToString(*i))
				simpleio.ErrorHandle(err)
				break
			}
		}
	}
	//atacFile, rnaFile := flag.Arg(0), flag.Arg(1)

}
func ToString(b bed.Pvalue) string {
	var str strings.Builder
	str.WriteString(b.Chrom())
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrStart()))
	str.WriteByte('\t')
	str.WriteString(strconv.Itoa(b.ChrEnd()))
	//str.WriteByte('\t')
	//str.WriteString(b.Name)
	str.WriteString(fmt.Sprintf("\t%E", b.PValue))
	return str.String()
}

func Read(filename string) []*bed.Pvalue {
	reader := simpleio.NewReader(filename)
	var ans []*bed.Pvalue
	name := strings.Split(filename, ".")[0]
	//var ans []bed.PValue
	simpleio.ReadLine(reader)
	for i, err := getData(reader, name); !err; i, err = getData(reader, name) {
		ans = append(ans, i)
	}
	return ans
}

func ReadToMap(filename string) map[string][]*bed.Pvalue {
	reader := simpleio.NewReader(filename)
	ans := make(map[string][]*bed.Pvalue)
	name := strings.Split(filename, ".")[0]
	//var ans []bed.PValue
	simpleio.ReadLine(reader)
	for i, err := getData(reader, name); !err; i, err = getData(reader, name) {
		ans[i.Chr] = append(ans[i.Chr], i)
	}
	return ans
}

func getData(reader *simpleio.SimpleReader, name string) (*bed.Pvalue, bool) {
	line, err := simpleio.ReadLine(reader)
	data := strings.Split(line.String(), ",")
	if !err {
		return &bed.Pvalue{
			Chr:    data[0],
			Start:  simpleio.StringToInt(data[1]),
			End:    simpleio.StringToInt(data[1]) + 1,
			Name:   name,
			PValue: float64(simpleio.StringToFloat(data[6])),
		}, false
	} else {
		return nil, true
	}
}
