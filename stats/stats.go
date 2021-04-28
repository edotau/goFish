// Package stats contains functions to perform statistical tests
package stats

import "math"

const lnSqrt2Pi = 0.918938533204672741780329736406 // log(sqrt(2*pi))

// Sample mean variance estimates for a data vector with Bessel correction or n - 1.
// https://en.wikipedia.org/wiki/Bessel%27s_correction
func MeanVariance(x []float64) (float64, float64) {
	var n int
	var m, m2 float64 = 0.0, 0.0
	var mean float64 = 0.0
	var varEst float64 = 0.0

	for _, val := range x {
		n += 1
		mean += val
		delta := val - m
		m += delta / float64(n)
		m2 += delta * (val - m)
	}

	varEst = m2 / float64(n-1)
	mean /= float64(len(x))
	return mean, varEst
}

func RejectionSample(targetDensity func(float64) float64, sourceDensity func(float64) float64, source func() float64, K float64) float64 {
	x := source()
	for ; NextUniform() >= targetDensity(x)/(K*sourceDensity(x)); x = source() {

	}
	return x
}

// x[i][j] := x[i][j] / |x[i]|
func Normalize(x [][]float64) {
	for i := range x {
		NormalizePoint(x[i])
	}
}

// NormalizePoint is the same as Normalize, but it only operates on one singular datapoint, normalizing it's value to unit length.
func NormalizePoint(x []float64) {
	var sum float64
	for i := range x {
		sum += x[i] * x[i]
	}
	mag := math.Sqrt(sum)

	for i := range x {
		if math.IsInf(x[i]/mag, 0) || math.IsNaN(x[i]/mag) {
			x[i] = 0
			continue
		}
		x[i] /= mag
	}
}
