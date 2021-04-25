package keras

import (
	"fmt"
)

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

// PrintByRow prints the matrix by row.
func PrintByRows(m Matrix) {
	for r := range m.Matrix {
		fmt.Println(m.Matrix[r])
	}
}
