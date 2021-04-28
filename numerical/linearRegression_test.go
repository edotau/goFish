package numerical

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/edotau/goFish/api"
	"github.com/edotau/goFish/csv"
	"github.com/edotau/goFish/stats"
)

var flatX [][]float64
var flatY []float64

var increasingX [][]float64
var increasingY []float64

var threeDLineX [][]float64
var threeDLineY []float64

var normX [][]float64
var normY []float64

var noisyX [][]float64
var noisyY []float64

func init() {
	err := os.MkdirAll("testdata", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Error: %v\n", err.Error()))
	}
	// the line y=3
	flatX = [][]float64{}
	flatY = []float64{}
	for i := -10; i < 10; i++ {
		for j := -10; j < 10; j++ {
			for k := -10; k < 10; k++ {
				flatX = append(flatX, []float64{float64(i), float64(j), float64(k)})
				flatY = append(flatY, 3.0)
			}
		}
	}
	// the line y=x
	increasingX = [][]float64{}
	increasingY = []float64{}
	for i := -10; i < 10; i++ {
		increasingX = append(increasingX, []float64{float64(i)})
		increasingY = append(increasingY, float64(i))
	}

	threeDLineX = [][]float64{}
	threeDLineY = []float64{}

	normX = [][]float64{}
	normY = []float64{}
	// the line z = 10 + (x/10) + (y/5)
	for i := -10; i < 10; i++ {
		for j := -10; j < 10; j++ {
			threeDLineX = append(threeDLineX, []float64{float64(i), float64(j)})
			threeDLineY = append(threeDLineY, 10+float64(i)/10+float64(j)/5)

			normX = append(normX, []float64{float64(i), float64(j)})
		}
	}

	stats.Normalize(normX)
	for i := range normX {
		normY = append(normY, 10+float64(normX[i][0])/10+float64(normX[i][1])/5)
	}

	// noisy x has random noise embedded
	rand.Seed(42)
	noisyX = [][]float64{}
	noisyY = []float64{}
	for i := 256.0; i < 1024; i += 2 {
		noisyX = append(noisyX, []float64{i + (rand.Float64()-0.5)*3})
		noisyY = append(noisyY, 0.5*i+rand.NormFloat64()*25)
	}
	// save the random data to make some nice plots!
	csv.Write("testdata/noisy_linear.csv", noisyX, noisyY, true)
}

func TestYEqualsThreeFlatLine(t *testing.T) {
	fmt.Printf("Testing y = 3 linear regression model...\n")
	var err error

	model := NewLeastSquares(BatchGA, .000001, 0, 800, flatX, flatY)

	err = model.Learn()
	api.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				api.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				api.InDelta(t, 3, guess[0], 1e-2, "Guess should be really close to 3 (within 1e-2) for y=3")
				api.Nil(t, err, "Prediction error should be nil")
			}
		}
	}
}

func TestLocalLinearShouldPass1(t *testing.T) {
	x := [][]float64{}
	y := []float64{}

	// throw in some junk points which
	// should be more-or-less ignored
	// by the weighting
	for i := -70.0; i < -65; i += 2 {
		for j := -70.0; j < -65; j += 2 {
			x = append(x, []float64{i, j})
			y = append(y, 20*(rand.Float64()-0.5))
		}
	}
	for i := 65.0; i < 70; i += 2 {
		for j := 65.0; j < 70; j += 2 {
			x = append(x, []float64{i, j})
			y = append(y, 20*(rand.Float64()-0.5))
		}
	}

	// put in some linear points
	for i := -20.0; i < 20; i++ {
		for j := -20.0; j < 20; j++ {
			x = append(x, []float64{i, j})
			y = append(y, 5*i-5*j-10)
		}
	}

	model := NewLocalLinear(BatchGA, 1e-4, 0, 0.75, 500, x, y)

	var count int
	var err float64
	for i := -15.0; i < 15; i += 7 {
		for j := -15.0; j < 15; j += 7 {
			guess, predErr := model.Predict([]float64{i, j})
			api.Nil(t, predErr, "learning/prediction error should be nil")
			count++

			err += math.Abs(guess[0] - (5*i - 5*j - 10))
		}
	}

	avgError := err / float64(count)

	api.True(t, avgError < 0.4, "Average error should be less than 0.4 from the expected value of the linear data (currently %v)", avgError)
	fmt.Printf("Average Error: %v\n\tPoints Tested: %v\n\tTotal Error: %v\n", avgError, count, err)
}
