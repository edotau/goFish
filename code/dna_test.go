package code

import (
	"testing"
)

var dnaStrings = []string{
	"ACGTacgtNn",
	"AAAAAACCCCCGGGGTTTTT",
	"aaaaaCCCCTTTaaaaa",
	"NNNNNAAAaaaTTT",
}

func TestDnaToFromString(t *testing.T) {
	for _, input := range dnaStrings {
		bases := ToDna([]byte(input))
		answer := ToString(bases)
		if input != answer {
			t.Errorf("Converting %s to bases and back gave %s", input, answer)
		}
	}
}
