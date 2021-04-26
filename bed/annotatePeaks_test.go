package bed

import (
	"fmt"
	"testing"
)

func TestReadAnnotationBed(t *testing.T) {
	annotatedPeaks := "testdata/annotatePeaks.txt"
	reader := NewAnnoPeaks(annotatedPeaks)
	var lines int = 0
	for i, done := ToBedAnnotation(reader); !done; i, done = ToBedAnnotation(reader) {
		lines++
		fmt.Printf("%s\n", i.ToString())
	}

	if lines != 9 {
		t.Errorf("Error: reader did not read all lines in file...\n")
	}
}
