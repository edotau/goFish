package stats

import (
	"math/rand"
)

//RandIntInRange produces a random value of type int between x and y.
func RandIntInRange(x int, y int) int {
	return int(rand.Float64()*float64(y-x)) + x
}

//RandFloat64InRange returns a a random Float value int between x and y.
func RandFloat64InRange(x float64, y float64) float64 {
	return rand.Float64()*(y-x) + x
}
