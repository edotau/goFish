package main

/*
import(
	"github.com/vertgenlab/gonomics/vcf"
	"github.com/vertgenlab/gonomics/fileio"
	"flag"
	"fmt"
	"log"
)

func usage() {
	fmt.Print(
		"vcfSort - sort vcf file by coordinate position\n\n" +
			"Usage:\n" +
			"  vcfSort [options] in.vcf out.vcf\n\n" +
			"Options:\n\n")
	flag.PrintDefaults()
}

func main() {
	//var expectedNumArgs int = 2
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()
	vcfs, header := vcf.GoReadToChan(flag.Arg(0))
	var ans []*vcf.Vcf
	for i := range vcfs {
		ans = append(ans, i)
	}
	vcf.Sort(ans)
	output := fileio.EasyCreate(flag.Arg(1))
	vcf.NewWriteHeader(output, header)
	for _, j := range ans {
		vcf.WriteVcf(output, j)
	}
}*/
