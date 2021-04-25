package keras

import (
	"fmt"
	"log"
	"testing"
)

var mat [][]float64 = [][]float64{
	{2, 4, 5, 6},
	{1, 2, 3, 6},
}

var forTrans [][]float64 = [][]float64{
	{2, 4},
	{1, 2},
}

func TestColNum(t *testing.T) {
	test := Matrix{Matrix: mat}
	if ColNum(test) != 4 {
		log.Fatalf("Error: number of columns should equal 4...\n")
	}
}

func TestRowNum(t *testing.T) {
	test := Matrix{Matrix: mat}
	if RowNum(test) != 2 {
		log.Fatalf("Error: number of rows should equal 2...\n")
	}
}

func TestNewMatrix(t *testing.T) {
	fmt.Printf("Testing new matrix function...\n")
	ans := NewMatrix(3, 2)
	PrintfMatrix(ans)
}

func TestTransposeMatrix(t *testing.T) {
	before := RandomMatrix(3, 2)
	fmt.Printf("Before: \n")
	PrintByRows(before)
	fmt.Printf("Transpose: \n")
	after := Transpose(before)
	PrintByRows(after)
}
