package stats

import (
	"fmt"
	"testing"
)

func TestBinomomial(t *testing.T) {
	fmt.Printf("NextUniform => %f\n", NextUniform())
	fmt.Printf("NextBernoulli => %d\n", NextBernoulli(.5))
	fmt.Printf("NextBinomial => %d\n", NextBinomial(.5, 10))
	fmt.Printf("NextPoisson => %d\n", NextPoisson(1.5))
	fmt.Printf("NextNegativeBinomial => %d\n", NextNegativeBinomial(.5, 10))
}
