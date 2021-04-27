package stats

import (
	"math"
	"math/rand"
)

var NextUniform func() float64 = rand.Float64

func Uniform() func() float64 { return NextUniform }

func UniformProbDist() func(x float64) float64 {
	return func(x float64) float64 {
		if 0 <= x && x <= 1 {
			return 1
		}
		return 0
	}
}

func NegativeBinomialProbMassFunc(ρ float64, r int64) func(k int64) float64 {
	return func(k int64) float64 {
		return BinomCoeff(k+r-1, k) * math.Pow(1-ρ, float64(r)) * math.Pow(ρ, float64(k))
	}
}
