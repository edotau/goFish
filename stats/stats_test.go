package stats

import (
	"fmt"
	"testing"
)

func TestMeanVariance(t *testing.T) {
	list := []float64{2051, 2053, 2055, 2050, 2051}
	mean, variance := MeanVariance(list)
	if mean != 2052 || variance != 4 {
		t.Errorf("Error: mean variance estimate is not calculating the correct value...\n")
	}
	fmt.Printf("example data: %v\nmean: %v, variance estimator: %v\n", list, mean, variance)
}
