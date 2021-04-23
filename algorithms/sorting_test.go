package algorithms

import (
	"testing"
)

var ans []int = []int{-1, 2, 3, 5, 6, 7, 8}
var maxHeap []int = []int{8, 7, 6, 5, 3, 2, -1}

func TestQuickSort(t *testing.T) {
	input := []int{6, 5, 3, 7, 2, 8, -1}
	output := QuickSort(input)
	if !equalSort(ans, output) {
		t.Errorf("Error: QuickSort did not sort the correct values...\nanswer:\t%v\noutput:\t%v", ans, output)
	}
}

func TestHeapSort(t *testing.T) {
	input := []int{6, 5, 3, 7, 2, 8, -1}
	output := input
	HeapSort(output)
	if !equalSort(maxHeap, output) {
		t.Errorf("Error: HeapSort did not sort the correct values...\nanswer:\t%v\noutput:\t%v", ans, output)
	}
}

func BenchmarkHeapSort(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		input := []int{6, 5, 3, 7, 2, 8, -1}
		output := input
		HeapSort(output)
	}
	b.Name()
}

func BenchmarkQuickSort(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		input := []int{6, 5, 3, 7, 2, 8, -1}
		QuickSort(input)
	}
	b.Name()
}

func equalSort(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
