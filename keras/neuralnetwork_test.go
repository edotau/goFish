package keras

import (
	"testing"
)

func TestNeuralNetworkModel(t *testing.T) {
	var input []float64 = SimpleRandomVector(10)
	model := Sequential([]Layer{
		Conv2D(64, 3, 3, input, Valid), Dense(128, input, ReLU), Dense(32, input, Tanh), Conv2D(32, 3, 1, input, DefaultPadding),
	}, "sequential")

	model.Summary()
}

/*MaxPooling2D(2),
Conv2D(32, 3, 1, DefaultPadding),
MaxPooling2D(2),
Flatten(),
,

Softmax(10),*/
