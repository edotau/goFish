package bam

import (
	"testing"
)

func TestReadIndex(t *testing.T) {
	idx := IndexReader("testdata/tenXbarcodeTest.bam.bai")
	t.Logf("%v\n", idx.Refs)
}
