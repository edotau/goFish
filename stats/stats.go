// Package stats contains functions to perform statistical tests
package stats

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
