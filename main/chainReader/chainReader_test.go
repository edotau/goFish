package main

import(
    "testing"
    "github.com/edotau/goFish/chain"
)

var readWriteTests = []struct {
    filename string // input
}{
    {"testdata/axtTest.chain.gz"},
}

func TestChainReader(t *testing.T) {
    for _, test := range readWriteTests {
        chain.Read(test.filename)
    }
    t.Logf("Chain file read without errors\n")

}
