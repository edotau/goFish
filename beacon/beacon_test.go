package beacon

import (
	"fmt"
	"testing"
)

func TestBeacon(t *testing.T) {
	ans := Read("testdata/D82698_OptoSeqBCRSummaryFile.csv")
	for _, i := range ans {
		fmt.Printf("%s\n", replaceRef(&i).String())
	}
}
