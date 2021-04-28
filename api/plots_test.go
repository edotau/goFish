package api

import (
	"fmt"
	"os"
	"testing"
)

func TestPlots(t *testing.T) {
	err := os.MkdirAll("testdata", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("You should be able to create the directory for goml model persistance testing.\n\tError returned: %v\n", err.Error()))
	}
	filename := "testdata/xscatter2"
	x := []float64{-4, -3.3, -1.8, -1, 0.2, 0.8, 1.8, 3.1, 4, 5.3, 6, 7, 8, 9}
	y := []float64{22, 18, -3, 0, 0.5, 2, 45, 12, 16.5, 24, 30, 55, 60, 70}
	ScatterChart(filename, x, y)
}
