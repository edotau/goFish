package numerical

type Vector struct {
	Cols []float64
	Abs  int
}

func NewVector(length int) (v *Vector) {
	v = new(Vector)
	v.Abs = length
	v.Cols = make([]float64, length)
	return v
}

func (v Vector) Get(i int) float64 {
	return v.Cols[i]
}

func (v Vector) Set(i int, x float64) {
	v.Cols[i] = x
}
