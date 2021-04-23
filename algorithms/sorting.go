package algorithms

// QuickSort is a simple implementation to sort a list of ints and return in sorted order
func QuickSort(arr []int) []int {
	newArr := make([]int, len(arr))

	for i, v := range arr {
		newArr[i] = v
	}
	recursiveSort(newArr, 0, len(arr)-1)
	return newArr
}

// recursiveSort is a helper function to iterate sub array to find values
//less than pivot and move them to the beginning of the array keeping
// splitIndex denoting less-value array size
func recursiveSort(arr []int, start, end int) {
	if (end - start) < 1 {
		return
	}

	pivot := arr[end]
	splitIndex := start

	for i := start; i < end; i++ {
		if arr[i] < pivot {
			if splitIndex != i {
				temp := arr[splitIndex]

				arr[splitIndex] = arr[i]
				arr[i] = temp
			}

			splitIndex++
		}
	}

	arr[end] = arr[splitIndex]
	arr[splitIndex] = pivot

	recursiveSort(arr, start, splitIndex-1)
	recursiveSort(arr, splitIndex+1, end)
}

// HeapSort is a simple implementation of the Heap Sorting algorithm
// is space-efficient when represented as an array/slice.
func HeapSort(a []int) {
	a = buildHeap(a)
	size := len(a)
	for i := size - 1; i >= 1; i-- {
		a[0], a[i] = a[i], a[0]
		size--
		heapify(a[:size], 0)
	}
}

// buildHeap is a helper function to build a max heap from our input data
func buildHeap(a []int) []int {
	for i := len(a)/2 - 1; i >= 0; i-- {
		a = heapify(a, i)
	}
	return a
}

// heapify is a helper function to check if the parent node is stored at index i,
// with the left child calculated by 2 * i + 1 and right child by 2 * i + 2
func heapify(a []int, i int) []int {
	l := left(i)
	r := right(i)
	var max int
	if l < len(a) && l >= 0 && a[l] < a[i] {
		max = l
	} else {
		max = i
	}
	if r < len(a) && r >= 0 && a[r] < a[max] {
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
