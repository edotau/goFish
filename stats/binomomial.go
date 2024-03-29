package stats

import (
	"math"
)

func Binomial(ρ float64, n int64) func() int64 {
	return func() int64 { return NextBinomial(ρ, n) }
}

func NextBinomial(ρ float64, n int64) (result int64) {
	for i := int64(0); i <= n; i++ {
		result += NextBernoulli(ρ)
	}
	return
}

//  Binomial coefficient calculates the number of ways k-element subsets (or k-combinations) of an n-element set, disregarding order
func BinomCoeff(n, k int64) float64 {
	if k == 0 {
		return 1
	}
	if n == 0 {
		return 0
	}
	if n < 10 && k < 10 {
		return BinomCoeff(n-1, k-1) + BinomCoeff(n-1, k)
	}
	return Round(math.Exp(LnFactBig(float64(n)) - LnFactBig(float64(k)) - LnFactBig(float64(n-k))))
}

func Bernoulli(ρ float64) func() int64 { return func() int64 { return NextBernoulli(ρ) } }

func NextBernoulli(ρ float64) int64 {
	if NextUniform() < ρ {
		return 1
	}
	return 0
}

//NegativeBinomial(ρ, r) => number of NextBernoulli(ρ) failures before r successes
func NextNegativeBinomial(ρ float64, r int64) int64 {
	k := int64(0)
	for r >= 0 {
		i := NextBernoulli(ρ)
		r -= i
		k += (1 - i)
	}
	return k
}
func NegativeBinomial(ρ float64, r int64) func() int64 {
	return func() int64 {
		return NextNegativeBinomial(ρ, r)
	}
}

// Round to nearest integer
func Round(x float64) float64 {
	var i float64
	f := math.Floor(x)
	c := math.Ceil(x)
	if x-f < c-x {
		i = f
	} else {
		i = c
	}
	return i
}

// cLnFactBig(n) = Gamma(n+1)
func LnFactBig(n float64) float64 {
	n = math.Trunc(n)
	return LnGamma(n + 1)
}

// Natural logarithm of the Gamma function
func LnGamma(x float64) (res float64) {
	res = (x-0.5)*math.Log(x+4.5) - (x + 4.5)
	res += lnSqrt2Pi
	res += math.Log1p(
		76.1800917300/(x+0) - 86.5053203300/(x+1) +
			24.0140982200/(x+2) - 1.23173951600/(x+3) +
			0.00120858003/(x+4) - 0.00000536382/(x+5))

	return
}
