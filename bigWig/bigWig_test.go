package bigWig

import (
	"testing"
)

func TestIsBigWig(t *testing.T) {
	bwig := NewReader("testdata/track.bw")
	if !MagicBigWig(bwig) {
		t.Errorf("Error: there is a bug in the function to check bigwig header if BIGWIG_MAGIC = 0x888FFC26\n")
	}
}

func TestBigWigReaderSource(t *testing.T) {
	ImportBigWigReader("testdata/track.bw")
}
