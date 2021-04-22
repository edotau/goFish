package main

import (
	"flag"
	"fmt"
	"log"
	//"github.com/edotau/simpleio"
	"github.com/vertgenlab/gonomics/chain"
	"github.com/vertgenlab/gonomics/common"
	"github.com/vertgenlab/gonomics/fileio"
)

func usage() {
	fmt.Print(
		"")
	flag.PrintDefaults()
}

func main() {
	var expectedNumArgs int = 2
	flag.Usage = usage
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if len(flag.Args()) != expectedNumArgs {
		flag.Usage()
		log.Fatalf("Error: expecting %d arguments, but got %d\n", expectedNumArgs, len(flag.Args()))
	}

	reader := fileio.EasyOpen(flag.Arg(0))
	writer := fileio.EasyCreate(flag.Arg(1))
	defer reader.Close()
	chain.ReadHeaderComments(reader)
	var err error
	for ch, done := chain.NextChain(reader); !done; ch, done = chain.NextChain(reader) {
		_, err = fmt.Fprintf(writer, "%s,%c,%d,%d,%s,%c,%d,%d\n", ch.TName, common.StrandToRune(ch.TStrand), ch.TStart, ch.TEnd, ch.QName, common.StrandToRune(ch.QStrand), ch.QStart, ch.QEnd)
		common.ExitIfError(err)
	}
	//var answer string = fmt.Sprintf("chain %d %s %d %c %d %d %s %d %c %d %d %d\n", ch.Score, ch.TName, ch.TSize, common.StrandToRune(ch.TStrand), ch.TStart, ch.TEnd, ch.QName, ch.QSize, common.StrandToRune(ch.QStrand), ch.QStart, ch.QEnd, ch.Id)
	writer.Close()
}
