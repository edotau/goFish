package numerical

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/edotau/goFish/api"
	"github.com/edotau/goFish/stats"
)

var (
	fourClusters  [][]float64
	fourClustersY []float64

	twoClusters  [][]float64
	twoClustersY []float64

	circles [][]float64
	double  [][]float64
)

func init() {
	fourClusters = [][]float64{}
	fourClustersY = []float64{}
	for i := -12.0; i < -8; i += 0.1 {
		for j := -12.0; j < -8; j += 0.1 {
			fourClusters = append(fourClusters, []float64{i, j})
			fourClustersY = append(fourClustersY, 0.0)
		}

		for j := 8.0; j < 12; j += 0.1 {
			fourClusters = append(fourClusters, []float64{i, j})
			fourClustersY = append(fourClustersY, 1.0)
		}
	}

	for i := 8.0; i < 12; i += 0.1 {
		for j := -12.0; j < -8; j += 0.1 {
			fourClusters = append(fourClusters, []float64{i, j})
			fourClustersY = append(fourClustersY, 2.0)
		}

		for j := 8.0; j < 12; j += 0.1 {
			fourClusters = append(fourClusters, []float64{i, j})
			fourClustersY = append(fourClustersY, 3.0)
		}
	}

	twoClusters = [][]float64{}
	twoClustersY = []float64{}
	for i := -10.0; i < -3; i += 0.1 {
		for j := -10.0; j < 10; j += 0.1 {
			twoClusters = append(twoClusters, []float64{i, j})
			twoClustersY = append(twoClustersY, 0.0)
		}
	}

	for i := 3.0; i < 10; i += 0.1 {
		for j := -10.0; j < 10; j += 0.1 {
			twoClusters = append(twoClusters, []float64{i, j})
			twoClustersY = append(twoClustersY, 1.0)
		}
	}
	circles = [][]float64{}
	for i := -12.0; i < -8; i += 0.2 {
		for j := -12.0; j < -8; j += 0.2 {
			circles = append(circles, []float64{i, j})
		}

		for j := 8.0; j < 12; j += 0.2 {
			circles = append(circles, []float64{i, j})
		}
	}

	for i := 8.0; i < 12; i += 0.2 {
		for j := -12.0; j < -8; j += 0.2 {
			circles = append(circles, []float64{i, j})
		}

		for j := 8.0; j < 12; j += 0.2 {
			circles = append(circles, []float64{i, j})
		}
	}

	double = [][]float64{}
	for i := -10.0; i < -3; i += 0.1 {
		for j := -10.0; j < 10; j += 0.1 {
			double = append(double, []float64{i, j})
		}
	}

	for i := 3.0; i < 10; i += 0.1 {
		for j := -10.0; j < 10; j += 0.1 {
			double = append(double, []float64{i, j})
		}
	}

	err := os.MkdirAll("testdata", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("You should be able to create the directory for goml model persistance testing.\n\tError returned: %v\n", err.Error()))
	}
}

func TestKNN(t *testing.T) {
	model := NewKNN(3, fourClusters, fourClustersY, EuclideanDistance)

	var count int
	var wrong int

	duration := time.Duration(0)
	for i := -12.0; i < -8; i += 0.5 {
		for j := -12.0; j < -8; j += 0.5 {
			now := time.Now()
			guess, err := model.Predict([]float64{i, j})
			duration += time.Now().Sub(now)
			api.Nil(t, err, "Prediction error should be nil")

			if 0.0 != guess[0] {
				wrong++
			}
			count++
		}

		for j := 8.0; j < 12; j += 0.5 {
			now := time.Now()
			guess, err := model.Predict([]float64{i, j})
			duration += time.Now().Sub(now)
			api.Nil(t, err, "Prediction error should be nil")

			if 1.0 != guess[0] {
				wrong++
			}
			count++
		}
	}

	for i := 8.0; i < 12; i += 0.5 {
		for j := -12.0; j < -8; j += 0.5 {
			now := time.Now()
			guess, err := model.Predict([]float64{i, j})
			duration += time.Now().Sub(now)
			api.Nil(t, err, "Prediction error should be nil")

			if 2.0 != guess[0] {
				wrong++
			}
			count++
		}

		for j := 8.0; j < 12; j += 0.5 {
			now := time.Now()
			guess, err := model.Predict([]float64{i, j})
			duration += time.Now().Sub(now)
			api.Nil(t, err, "Prediction error should be nil")

			if 3.0 != guess[0] {
				wrong++
			}
			count++
		}
	}

	accuracy := 100 * (1 - float64(wrong)/float64(count))
	api.True(t, accuracy > 95, "Accuracy (%v) should be greater than 95 percent", accuracy)
	fmt.Printf("Accuracy: %v percent\n\tPoints Tested: %v\n\tMisclassifications: %v\n\tAverage Prediction Time: %v\n", accuracy, count, wrong, duration/time.Duration(count))
}

func TestKMeans(t *testing.T) {
	norm := append([][]float64{}, circles...)
	stats.Normalize(norm)
	model := NewKMeans(4, 2, norm)

	api.Nil(t, model.Learn(), "Learning error should be nil")

	c1, err := model.Predict([]float64{-10, -10}, true)
	api.Nil(t, err, "Prediction error should be nil")

	c2, err := model.Predict([]float64{-10, 10}, true)
	api.Nil(t, err, "Prediction error should be nil")

	c3, err := model.Predict([]float64{10, -10}, true)
	api.Nil(t, err, "Prediction error should be nil")

	c4, err := model.Predict([]float64{10, 10}, true)
	api.Nil(t, err, "Prediction error should be nil")

	var count int
	var wrong int

	for i := -12.0; i < -8; i += 0.2 {
		for j := -12.0; j < -8; j += 0.2 {
			guess, err := model.Predict([]float64{i, j}, true)
			api.Nil(t, err, "Prediction error should be nil")

			if c1[0] != guess[0] {
				wrong++
			}
			count++
		}

		for j := 8.0; j < 12; j += 0.2 {
			guess, err := model.Predict([]float64{i, j}, true)
			api.Nil(t, err, "Prediction error should be nil")

			if c2[0] != guess[0] {
				wrong++
			}
			count++
		}
	}

	for i := 8.0; i < 12; i += 0.2 {
		for j := -12.0; j < -8; j += 0.2 {
			guess, err := model.Predict([]float64{i, j}, true)
			api.Nil(t, err, "Prediction error should be nil")

			if c3[0] != guess[0] {
				wrong++
			}
			count++
		}

		for j := 8.0; j < 12; j += 0.2 {
			guess, err := model.Predict([]float64{i, j}, true)
			api.Nil(t, err, "Prediction error should be nil")

			if c4[0] != guess[0] {
				wrong++
			}
			count++
		}
	}

	accuracy := 100 * (1 - float64(wrong)/float64(count))
	api.True(t, accuracy > 87, "Accuracy (%v) should be greater than 87 percent", accuracy)
	fmt.Printf("Accuracy: %v percent\n\tPoints Tested: %v\n\tMisclassifications: %v\n\tClasses: %v\n", accuracy, count, wrong, []float64{c1[0], c2[0], c3[0], c4[0]})
}

/*
func TestKMeansGaussian(t *testing.T) {
	gaussian := [][]float64{}
	for i := 0; i < 40; i++ {
		x := rand.NormFloat64() + 4
		y := rand.NormFloat64()*0.25 + 5
		gaussian = append(gaussian, []float64{x, y})
	}
	for i := 0; i < 66; i++ {
		x := rand.NormFloat64()
		y := rand.NormFloat64() + 10
		gaussian = append(gaussian, []float64{x, y})
	}
	for i := 0; i < 100; i++ {
		x := rand.NormFloat64()*3 - 10
		y := rand.NormFloat64()*0.25 - 7
		gaussian = append(gaussian, []float64{x, y})
	}
	for i := 0; i < 23; i++ {
		x := rand.NormFloat64() * 2
		y := rand.NormFloat64() - 1.25
		gaussian = append(gaussian, []float64{x, y})
	}

	model := NewKMeans(4, 15, gaussian)

	if model.Learn() != nil {
		panic("Oh NO!!! There was an error learning!!")
	}

	// now you can predict like normal!
	guess, err := model.Predict([]float64{-3, 6})
	if err != nil {
		panic("prediction error")
	}
	// or if you want to get the clustering results from the data
	results := model.Guesses()
}*/
