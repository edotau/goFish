package keras

import (
	"math/rand"
)

//Vector type
type Vector struct {
	row []float64
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

// RandomVector returns a random valued vector.
func SimpleRandomVector(size int) []float64 {
	slice := make([]float64, size)
	for i := 0; i < size; i++ {
		slice = append(slice, rand.Float64()/0.3)
	}
	return slice
}
