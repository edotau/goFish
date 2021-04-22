package keras

import (
	"fmt"
	"math"
	"math/rand"
)

// Matrix type is a 2D slice of float64 Values.
type Matrix struct {
	Matrix [][]float64
}

// ColNum returns the number of columns in a matrix.
func ColNum(m Matrix) int {
	return len(m.Matrix[len(m.Matrix)-1])
}

// RowNum returns the number of rows in a matrix.
func RowNum(m Matrix) int {
	return len(m.Matrix)
}

// PrintfByRows will print a matrix by rows.
func PrintfByRows(m Matrix) {
	for r := 0; r < len(m.Matrix); r++ {
		fmt.Println(m.Matrix[r])
	}
}

// NewMatrix allocates the appropriate memory for an m x n matrix.
func NewMatrix(m, n int) Matrix {
	ans := Matrix{}
	ans.Matrix = make([][]float64, n)
	for each := range ans.Matrix {
		ans.Matrix[each] = make([]float64, m)
	}
	return ans
}

// PrintfMatrix is a helper function that will print a matrix to stdout.
func PrintfMatrix(m Matrix) {
	var i, j int
	for i = 0; i < len(m.Matrix); i++ {
		for j = 0; j < len(m.Matrix[i])-1; j++ {
			fmt.Printf("%f, ", m.Matrix[i][j])
		}
		fmt.Printf("%f\n", m.Matrix[i][len(m.Matrix[i])-1])

	}
}

// RandomMatrix will create a new matrix and randomize float64 values.
func RandomMatrix(m, n int) Matrix {
	ans := NewMatrix(m, n)
	var i, j int
	for i = 0; i < len(ans.Matrix); i++ {
		for j = 0; j < len(ans.Matrix[0]); j++ {
			ans.Matrix[i][j] = rand.Float64()
		}
	}
	return ans
}

// Transpose will tranpose a matrix and modify a given matrix.
func Transpose(m Matrix) Matrix {
	ans := NewMatrix(len(m.Matrix), len(m.Matrix[0]))
	var i, j int
	for i = 0; i < len(m.Matrix[0]); i++ {
		for j = 0; j < len(m.Matrix); j++ {
			ans.Matrix[i][j] = m.Matrix[j][i]
		}
	}
	return ans
}

// Sigmoid returns sigmoid of x.
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// SigmoidPrime returns the sigmoid derivative of x.
func SigmoidPrime(x float64) float64 {
	return Sigmoid(x) * (1 - Sigmoid(x))
}

/*
func DotProduct(one, two Matrix) float64 {
	var ans float64
	var i, j int
	for i = 0; i < len
}
*/
