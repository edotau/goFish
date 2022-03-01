package main

import (
	"github.com/edotau/goFish/chain"
	"testing"
)

var readWriteTests = []struct {
	filename string // input
}{
	{"testdata/axtTest.chain.gz"},
}

// TODO: finish writing unit test with a little bit more logic
func TestChainReader(t *testing.T) {
	for _, test := range readWriteTests {
		chain.Read(test.filename)
	}
	t.Logf("Chain file read without errors\n")

}
