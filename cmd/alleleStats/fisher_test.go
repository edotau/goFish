package main

import (
	"testing"

	"github.com/vertgenlab/gonomics/numbers"
)

// fisherExactGreaterTests examples were taken stright from gonomics
var fisherExactGreaterTests = []struct {
	a      int     // top left in matrix
	b      int     // top right in matrix
	c      int     // bottom left in matrix
	d      int     // bottom right in matrix
	pvalue float64 // pvalue of a being greater
}{
	{108, 1432, 742, 70208, 2.95e-50},
	{76, 542, 774, 71098, 1.04e-52},
	{47, 1253, 50, 84636, 7.12e-59},
}

func TestFisherTest(t *testing.T) {

	for _, test := range fisherExactGreaterTests {
		calculated := numbers.FisherExact(test.a, test.b, test.c, test.d, false)
		if calculated > test.pvalue*1.01 || calculated < test.pvalue*0.99 {
			t.Errorf("For a fisher test (greater) on (%d, %d, %d, %d): expected %e, but got %e", test.a, test.b, test.c, test.d, test.pvalue, calculated)
		}
	}

}
