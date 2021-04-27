package stats

import (
	"math"
)

func Poisson(lamda float64) func() int64 {
	return func() int64 {
		return NextPoisson(lamda)
	}
}

func NextPoisson(lamda float64) int64 {
	// this can be improved upon
	i := int64(0)
	t := math.Exp(-lamda)
	p := float64(1.0)
	for ; p > t; p *= NextUniform() {
		i++
	}
	return i
}

func PoissonProbDistFunc(lamda float64) func(k int64) float64 {
	pmf := Poisson_LnPMF(lamda)
	return func(k int64) float64 {
		p := math.Exp(pmf(k))
		return p
	}
}

func PoissonCumDistFunc(lamda float64) func(k int64) float64 {
	return func(k int64) float64 {
		var p float64 = 0
		var i int64
		pmf := PoissonProbDistFunc(lamda)
		for i = 0; i <= k; i++ {
			p += pmf(i)
		}
		return p
	}
}

func Poisson_LnPMF(λ float64) func(k int64) float64 {
	return func(k int64) (p float64) {
		i := float64(k)
		a := math.Log(λ) * i
		b := math.Log(math.Gamma((i + 1)))
		p = a - b - λ
		return p
	}
}
