package main

import (
	"fmt"
	"testing"
)

var faTest Fasta = Fasta{Name: "threeBases", Seq: []Dna{A, A, A}}

func TestKmerValues(t *testing.T) {
	type checkFunc func(int) error

	// Checkes for values of k that are equal to or less than 0, which is equivalent to k < 1
	kLessThanOne := func(k int) error {
		if k < 1 {
			return fmt.Errorf("Error: k=%d cannot be less than or equal to 0...", k)
		} else {
			createKMerHash([]Fasta{faTest}, k)
			return nil
		}
	}
	// Checkes for values of k that are greater than the test fasta:
	kGreaterSeq := func(k int) error {
		if k > len(faTest.Seq) {
			return fmt.Errorf("Error: k cannot be greater than %d...", len(faTest.Seq))
		} else {
			createKMerHash([]Fasta{faTest}, k)
			return nil
		}
	}

	// Test cases that should fail
	tests := [...]struct {
		in    int
		check checkFunc
	}{
		{-10, kLessThanOne},
		{1000, kGreaterSeq},
	}

	// Test logic
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Unit tests to test values of k: %v", tc.in), func(t *testing.T) {
			if err := tc.check(tc.in); err == nil {
				t.Log(err)
			}
		})
	}
}
