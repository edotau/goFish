package keras

import (
	"fmt"
	"math/rand"
)

// Matrix type is a 2D slice of float64 Values.
type Matrix struct {
	Matrix [][]float64
}

//Vector type
type Vector struct {
	row []float64
}

// ColNum returns the number of columns in a matrix.
func ColNum(m Matrix) int {
	return len(m.Matrix[len(m.Matrix)-1])
}

// RowNum returns the number of rows in a matrix.
func RowNum(m Matrix) int {
	return len(m.Matrix)
}

//Dimensions returns the number of rows and columns of m.
func (m Matrix) Dimensions() (int, int) {
	return RowNum(m), ColNum(m)
}

//NumberOfElements returns the number of elements.
func (m Matrix) NumberOfElements() int {
	return RowNum(m) * ColNum(m)
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

//PrintByRow prints the matrix by row.
func (m Matrix) PrintByRow() {
	for r := range m.Matrix {
		fmt.Println(m.Matrix[r])
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

//MapFunc applies f to every element
func (m Matrix) MapFunc(f func(x float64) float64) Matrix {
	for i := 0; i < RowNum(m); i++ {
		for j := 0; j < ColNum(m); j++ {
			m.Matrix[i][j] = f(m.Matrix[i][j])
		}
	}
	return m
}

//ToArray returns the matrix in array form.
func (m Matrix) ToArray() []float64 {
	var arr []float64
	for i := 0; i < RowNum(m); i++ {
		for j := 0; j < ColNum(m); j++ {
			arr = append(arr, m.Matrix[i][j])
		}
	}
	return arr
}

// ApplyMatrix returns the vector through a matrix transformation.
func (v Vector) ApplyMatrix(matrix Matrix) Vector {
	var product Vector
	for _, r := range matrix.Matrix {
		for i := 0; i < len(r); i++ {
			product.row[i] = r[i] * v.row[i]
		}
	}
	return product
}

//Add performs elementary matrix addition
func (m Matrix) Add(mat Matrix) Matrix {
	var product Matrix
	for i := 0; i < RowNum(m); i++ {
		for j := 0; j < ColNum(m); j++ {
			product.Matrix[i][j] = m.Matrix[i][j] + mat.Matrix[i][j]
		}
	}
	return product
}

//Subtract performs elementary matrix subtraction
func (m Matrix) Subtract(mat Matrix) Matrix {
	var product Matrix
	for i := 0; i < RowNum(m); i++ {
		for j := 0; j < ColNum(m); j++ {
			product.Matrix[i][j] = m.Matrix[i][j] - mat.Matrix[i][j]
		}
	}
	return product
}

//Multiply performs elementary matrix multiplication
func (m Matrix) Multiply(mat Matrix) Matrix {
	var product Matrix
	for i := 0; i < RowNum(m); i++ {
		for j := 0; j < ColNum(m); j++ {
			product.Matrix[i][j] = m.Matrix[i][j] * mat.Matrix[i][j]
		}
	}
	return product
}

// NumberOfElements returns the number of elements.
func (v Vector) NumberOfElements() int {
	return len(v.row)
}

// NewVector returns a vector type
func NewVector(slice []float64) Vector {
	return Vector{row: slice}
}

// RandomVector returns a random valued vector.
func RandomVector(size int) Vector {
	slice := make([]float64, size)
	for i := 0; i < size; i++ {
		slice = append(slice, rand.Float64()/0.3)
	}
	return NewVector(slice)
}

//Map maps the vector by with the function
func (v Vector) Map(f func(float64) float64) Vector {
	for i := range v.row {
		v.row[i] = f(v.row[i])
	}
	return v
}

// Add returns an elementary operation on two vectors.
func (v Vector) Add(v2 Vector) Vector {
	var resultVector Vector
	for i := 0; i < len(v.row); i++ {
		resultVector.row[i] = v.row[i] + v2.row[i]
	}
	return resultVector
}

// Slice returns vector.slice.
// You can perform indexing with this method.
func (v Vector) Slice() []float64 {
	return v.row
}

//ScalarMultiplication multiplies every element with a scalar
func (m Matrix) ScalarMultiplication(scalar float64) Matrix {
	for _, r := range m.Matrix {
		for i := range r {
			r[i] = r[i] * scalar
		}
	}
	return m
}

//FromArray returns a matrix from array
func FromArray(arr []float64) Matrix {
	m := Zeros(len(arr), 1)
	for i := 0; i < len(arr); i++ {
		m.Matrix[i][0] = arr[0]
	}
	return m
}

//Zeros returns a matrix of zeros.
func Zeros(row, column int) Matrix {
	b := make([][]float64, row)
	v := make([]float64, column)
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			v[j] = 0
			b[i] = v
		}
	}

	return Matrix{Matrix: b}
}

//ScalarAdition adds a scalar to every elements
func (m Matrix) ScalarAdition(scalar float64) Matrix {
	for _, r := range m.Matrix {
		for i := range r {
			r[i] = r[i] + scalar
		}
	}
	return m
}

//NewVector returns a vector type

/*
func DotProduct(one, two Matrix) float64 {
	var ans float64
	var i, j int
	for i = 0; i < len
}
*/
