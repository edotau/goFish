package algorithms

import (
	"bytes"
	"log"
	"sort"
)

func BurrowsWheelerTransform(s []byte, idx byte) []byte {
	sa := GetSuffixArray(s)
	return ComputeSuffixArray(s, sa, idx)
}

// TODO: InverseBurrowsWheelerTransform(){}

func GetSuffixArray(s []byte) []int {
	sa := make([]int, len(s)+1)
	sa[0] = len(s)

	for i := 0; i < len(s); i++ {
		sa[i+1] = i
	}
	sort.Slice(sa[1:], func(i, j int) bool {
		return bytes.Compare(s[sa[i+1]:], s[sa[j+1]:]) < 0
	})
	return sa
}

// ComputeSuffixArray computes BWT using a generated suffix array as input
func ComputeSuffixArray(s []byte, sa []int, idx byte) []byte {
	if len(s)+1 != len(sa) || sa[0] != len(s) {
		log.Fatalf("Error: the length of suffix array is invalid...\n")
	}
	bwt := make([]byte, len(sa))
	bwt[0] = s[len(s)-1]
	for i := 1; i < len(sa); i++ {
		if sa[i] == 0 {
			bwt[i] = idx
		} else {
			bwt[i] = s[sa[i]-1]
		}
	}
	return bwt
}
