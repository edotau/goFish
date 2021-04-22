package bed

import (
	"sort"
	"strings"
)

type BedSlice []Bed

func HeapSort(a BedSlice) {
	a = buildHeap(a)
	size := len(a)
	for i := size - 1; i >= 1; i-- {
		a[0], a[i] = a[i], a[0]
		size--
		heapify(a[:size], 0)
	}
}

func QuickSort(beds []Bed) {
	sort.Slice(beds, func(i, j int) bool { return Compare(beds[i], beds[j]) == -1 })
}

func heapify(beds []Bed, i int) []Bed {
	var l int = 2*i + 1
	var r int = 2*i + 2
	var max int
	if l < len(beds) && l >= 0 && Compare(beds[l], beds[i]) < 0 {
		max = l
	} else {
		max = i
	}
	if r < len(beds) && r >= 0 && Compare(beds[r], beds[max]) < 0 {
		max = r
	}
	if max != i {
		beds[i], beds[max] = beds[max], beds[i]
		beds = heapify(beds, max)
	}
	return beds
}

func buildHeap(a []Bed) []Bed {
	for i := len(a)/2 - 1; i >= 0; i-- {
		a = heapify(a, i)
	}
	return a
}

func Compare(a Bed, b Bed) int {
	if chrName := strings.Compare(a.Chrom(), b.Chrom()); chrName != 0 {
		return chrName
	}
	if a.ChrStart() < b.ChrStart() {
		return -1
	}
	if a.ChrStart() > b.ChrStart() {
		return 1
	}
	if a.ChrEnd() < b.ChrEnd() {
		return -1
	}
	if a.ChrEnd() > b.ChrEnd() {
		return 1
	}
	return 0
}

func Len(b Bed) int {
	return b.ChrEnd() - b.ChrStart()
}
