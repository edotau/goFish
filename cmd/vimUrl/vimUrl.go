// vimUrl can view http url links and print data stream to stdout
package main

import (
	"flag"
	"fmt"
	"github.com/edotau/goFish/bam"
	"github.com/edotau/goFish/simpleio"
	"log"
	"strings"
)

func usage() {
	fmt.Print(
		"vimUrl - view http url links and print data stream to stdout\n  Usage:\n  ./vimUrl link.com/file.txt\nOptions:\n\nComing soon!\n")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 1
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()
	//var toFile *string = flag.String("out", "", "provide a name to redirect data stream to a `file.txt`")
	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}
	url := flag.Arg(0)
	if strings.HasSuffix(url, ".bam") || strings.HasSuffix(url, ".s am") {
		bam.ViewUrl(url)
	}
	simpleio.VimUrl(url)
}
