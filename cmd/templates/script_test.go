package main

import (
	"testing"
	// "github.com/edotau/goFish/"
)

var readWriteTests = []struct {
	filename string // input
}{
	{"testdata"},
}

// TODO:
func TestChainReader(t *testing.T) {
	for _, test := range readWriteTests {
		t.Logf("%s\n", test)
	}
}
