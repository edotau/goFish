// GeneSeq contains functions and other data structures for processing gene models and annotations
// from difference public data bases like ENSEMBL and UCSC
// Currently supports genePred and GTF formats
package geneSeq

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/edotau/goFish/concurrency"
	"github.com/edotau/goFish/simpleio"
)

type GenePred struct {
	GeneName  string
	Chr       string
	Strand    byte
	TxStart   int
	TxEnd     int
	CdsStart  int
	CdsEnd    int
	ExonCount int
	ExonStart []int
	ExonEnd   []int
	Ext       string
}

type GeneModels []GenePred

type GeneSym struct {
	Symbol  string
	Ensembl string
}

func Read(filename string) []GenePred {
	reader := simpleio.NewReader(filename)
	var ans []GenePred
	for i, err := GenePredLine(reader); !err; i, err = GenePredLine(reader) {
		ans = append(ans, *i)
	}
	return ans
}

func GenePredLine(reader *simpleio.SimpleReader) (*GenePred, bool) {
	var err bool
	reader.Buffer, err = simpleio.ReadLine(reader)
	if !err {
		columns := strings.SplitN(reader.Buffer.String(), "\t", 11)
		if len(columns) >= 10 {
			ans := &GenePred{
				GeneName:  columns[0],
				Chr:       columns[1],
				Strand:    columns[2][0],
				TxStart:   simpleio.StringToInt(columns[3]),
				TxEnd:     simpleio.StringToInt(columns[4]),
				CdsStart:  simpleio.StringToInt(columns[5]),
				CdsEnd:    simpleio.StringToInt(columns[6]),
				ExonCount: simpleio.StringToInt(columns[7]),
				ExonStart: simpleio.StringToIntSlice(columns[8]),
				ExonEnd:   simpleio.StringToIntSlice(columns[9]),
			}
			if len(columns) == 11 {
				ans.Ext = columns[10]
			}
			if len(ans.ExonStart) == ans.ExonCount && len(ans.ExonEnd) == ans.ExonCount {
				return ans, false
			} else {
				log.Fatalf("Error: ExonCount must equal length of ExonStart and ExonEnd...\n")
			}
		} else {
			log.Fatalf("Error: line must contains %d, must be at least 10 columns in gene Prediction format...\n", len(columns))
		}
	}
	return nil, true
}

// ReadBioPark processes a two column table and produces a hash table look up used to map gene names from difference data bases to one another
func ReadBioMart(filename string) map[string]string {
	//var ans
	reader := simpleio.NewReader(filename)
	var work []string
	hash := make(map[string]string)
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		work = strings.SplitN(i.String(), "\t", 2)
		hash[work[0]] = strings.ReplaceAll(work[1], " ", "_")
	}

	return hash
}

func ReadPipe(filename string) {
	reader := simpleio.NewReader(filename)
	outC := concurrency.Work(func(inC chan interface{}) {
		defer close(inC)
		//for i, err := GenePredLine(reader); !err; i, err = GenePredLine(reader) {

		//}
		for i := 0; i < 1; i++ {
			inC <- reader
		}
	}).
		Pipe(func(in interface{}) interface{} {
			var ans GeneModels
			for i, err := GenePredLine(in.(*simpleio.SimpleReader)); !err; i, err = GenePredLine(in.(*simpleio.SimpleReader)) {
				ans = append(ans, *i)
			}
			return ans
		}).
		//Pipe(func(in interface{}) (interface{}) {
		//WriteTest("testdata/rewrite.vcf", in.(VcfSlice))
		//return nil
		//}).
		Merge()

	for range outC {
		//fmt.Printf("%v\n", j)
		// Do nothing, just for  drain out channel
	}
}

// FilterSrand will filter out genePred lines that do not have strand information
func FilterStrand(gp []GenePred) []GenePred {
	var ans []GenePred
	for _, i := range gp {
		if i.Strand != '.' {
			ans = append(ans, i)
		}
	}
	return ans
}

func (gp *GenePred) Chrom() string {
	return gp.Chr
}

func (gp *GenePred) ChrStart() int {
	return gp.TxStart
}

func (gp *GenePred) ChrEnd() int {
	return gp.TxEnd
}

func WriteGenePred(filename string, geneModels []GenePred) {
	output := simpleio.NewWriter(filename)
	var err error
	for _, i := range geneModels {
		_, err = fmt.Fprintf(output, "%s\n", i.ToString())
		simpleio.StdError(err)
	}
	output.Close()
}

func ReadToMap(filename string) map[string][]*GenePred {
	ans := make(map[string][]*GenePred)
	reader := simpleio.NewReader(filename)
	for i, err := GenePredLine(reader); !err; i, err = GenePredLine(reader) {
		ans[i.Chr] = append(ans[i.Chr], i)
	}
	return ans
}

func (gp *GenePred) ToString() string {
	var str strings.Builder
	str.WriteString(gp.GeneName)
	str.WriteByte('\t')
	str.WriteString(gp.Chr)
	str.WriteByte('\t')
	str.WriteByte(gp.Strand)
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gp.TxStart))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gp.TxEnd))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gp.CdsStart))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gp.CdsEnd))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntToString(gp.ExonCount))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntSliceToString(gp.ExonStart))
	str.WriteByte('\t')
	str.WriteString(simpleio.IntSliceToString(gp.ExonEnd))
	if gp.Ext != "" {
		str.WriteByte('\t')
		str.WriteString(gp.Ext)
	}
	return str.String()
}

func (gp *GenePred) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteString(gp.GeneName)
	buf.WriteByte('\t')
	buf.WriteString(gp.Chr)
	buf.WriteByte('\t')
	buf.WriteByte(gp.Strand)
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(gp.TxStart))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(gp.TxEnd))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(gp.CdsStart))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(gp.CdsEnd))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntToString(gp.ExonCount))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntSliceToString(gp.ExonStart))
	buf.WriteByte('\t')
	buf.WriteString(simpleio.IntSliceToString(gp.ExonEnd))
	if gp.Ext != "" {
		buf.WriteByte('\t')
		buf.WriteString(gp.Ext)
	}
	buf.WriteByte('\n')
	return buf.Bytes()
}
