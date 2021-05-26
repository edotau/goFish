package main

import ()

type Kmer struct {
	Seq   string
	Count int
}

// HeapSort is a simple implementation of the Heap Sorting algorithm
// is space-efficient when represented as an array/slice.
func HeapSort(a map[string]int) []Kmer {
	var ans []Kmer
	for key := range a {
		ans = append(ans, Kmer{Seq: key, Count: a[key]})
	}
	ans = buildHeap(ans)
	size := len(a)
	for i := size - 1; i >= 1; i-- {
		ans[0], ans[i] = ans[i], ans[0]
		size--
		heapify(ans[:size], 0)
	}
	return ans
}

// buildHeap is a helper function to build a max heap from our input data
func buildHeap(a []Kmer) []Kmer {
	for i := len(a)/2 - 1; i >= 0; i-- {
		a = heapify(a, i)
	}
	return a
}

// heapify is a helper function to check if the parent node is stored at index i,
// with the left child calculated by 2 * i + 1 and right child by 2 * i + 2
func heapify(a []Kmer, i int) []Kmer {
	l := left(i)
	r := right(i)
	var max int
	if l < len(a) && l >= 0 && a[l].Count < a[i].Count {
		max = l
	} else {
		max = i
	}
	if r < len(a) && r >= 0 && a[r].Count < a[max].Count {
		max = r
	}
	if max != i {
		a[i], a[max] = a[max], a[i]
		a = heapify(a, max)
	}
	return a
}

// left is a helper function to calculate the left child, 2 * i + 1
func left(i int) int {
	return 2*i + 1
}

// right is a helper function to calculate the right child, 2 * i + 2
func right(i int) int {
	return 2*i + 2
}
