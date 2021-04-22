package main

func main() {}

/*
import(
	"strings"
	"github.com/vertgenlab/gonomics/vcf"
	"github.com/vertgenlab/gonomics/common"
	"encoding/csv"
	"io"
	"log"
)

func ReadVcf(filename string) []*vcf.Vcf {
	var ans []*vcf.Vcf
	data := NewSimpleReader(filename)
	reader := csv.NewReader(data)
	reader.Comma = '\t'
	reader.Comment = '#'
	var curr *vcf.Vcf
	var i []string
	var err error
	for {
		i, err = reader.Read()
		if err == nil {
			if len(i) > 9 {
				curr = &vcf.Vcf{
					Chr: i[0],
					Pos: common.StringToInt64(i[1]),
					Id: i[2],
					Ref: i[3],
					Alt: i[4],
					Filter: i[6],
					Info: i[7],
					Format: i[8],
					Notes: strings.Join(i[9:], "\n"),
				}
				ans = append(ans, curr)
			} else {
				log.Fatalf("Error: line is less than 10: %v", i)
			}
		} else if err == io.EOF {
			break
		} else {
			log.Fatal(err)
		}
	}
	data.Close()
	return ans
}*/
