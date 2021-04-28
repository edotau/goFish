package numerical

import (
	"fmt"
	"testing"
)

var matrix *Matrix = &Matrix{
	Rows: 2,
	Cols: 4,
	Data: []float64{2, 4, 5, 6, 1, 2, 3, 6},
}

var transpose *Matrix = &Matrix{
	Rows: 4,
	Cols: 2,
	Data: []float64{2, 1, 4, 2, 5, 3, 6, 6},
}

func TestNewMatrix(t *testing.T) {
	fmt.Printf("Testing new matrix function...\n")
	ans := RandomMatrix(3, 2)
	fmt.Printf("%s\n", ans.ToString())
}

func TestTransposeMatrix(t *testing.T) {
	before := matrix
	fmt.Printf("Before: \n")
	fmt.Printf("%s\n", before.ToString())

	fmt.Printf("Transpose: \n")
	after := Transpose(before)
	fmt.Printf("%s\n", after.ToString())

	if !Equal(after, transpose) {
	}
}
