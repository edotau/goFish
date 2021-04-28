package numerical

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

type Matrix struct {
	Rows int
	Cols int
	Data []float64
}

func NewMatrix(rows, cols int) *Matrix {
	m := Matrix{}
	m.Rows = rows
	m.Cols = cols
	m.Data = make([]float64, rows*cols)
	return &m
}

func (m Matrix) Get(i int, j int) float64 {
	return m.Data[i*m.Cols+j]
}

func (m Matrix) Set(i int, j int, x float64) {
	//	m.A[i+j*m.C] = x
	m.Data[i*m.Cols+j] = x
}

func RandomMatrix(rows, cols int) *Matrix {
	var i, j int
	m := NewMatrix(rows, cols)
	for i = 0; i < m.Rows; i++ {
		for j = 0; j < m.Cols; j++ {
			m.Set(i, j, rand.Float64())
		}

	}
	return m
}

func (m *Matrix) ToString() string {
	var i, j int
	str := &strings.Builder{}
	for i = 0; i < m.Rows; i++ {
		for j = 0; j < m.Cols; j++ {
			str.WriteString(simpleio.Float64ToString(m.Get(i, j)))
			str.WriteByte('\t')
		}
		str.WriteByte('\n')
	}
	return str.String()
}

func (m *Matrix) Print() {
	var i, j int
	for i = 0; i < m.Rows; i++ {
		for j = 0; j < m.Cols; j++ {
			fmt.Printf("%f ", m.Get(i, j))
		}
		fmt.Print("\n")
	}
}

// Transpose will tranpose a matrix and modify a given matrix.
func Transpose(m *Matrix) *Matrix {
	ans := NewMatrix(m.Cols, m.Rows)
	var i, j int
	for i = 0; i < m.Cols; i++ {
		for j = 0; j < m.Rows; j++ {
			ans.Set(i, j, m.Get(j, i))
		}
	}
	return ans
}

func Equal(a, b *Matrix) bool {
	if len(a.Data) != len(b.Data) {
		return false
	} else {
		for i := 0; i < len(a.Data); i++ {
			if a.Data[i] != b.Data[i] {
				return false
			}
		}
		return true
	}
}
