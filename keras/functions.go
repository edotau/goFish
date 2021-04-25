package keras

import (
	"math"
)

//ReLU or rectified linear activation function is a piecewise linear function that will output the input directly if it is positive, otherwise, it will output zero
func ReLU(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}

// Sigmoid activation function is commonly known as the logistic function
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

//SigmoidPrime calculates the derivative of the Sigmoid
func SigmoidPrime(x float64) float64 {
	return Sigmoid(x) * (1 - Sigmoid(x))
}

// Tanh calculates tanh(x) of a number
func Tanh(x float64) float64 {
	return math.Tanh(x)
}

// ELUExponential Linear Unit or its widely known name ELU is a function that tend to converge cost to zero faster and produce more accurate results
func ELU(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0.7 * (math.Exp(x) - 1)
}

// Swish is an activation function proposed byGoogle Brain Team, which is simply f(x) = x · sigmoid(x)
func Swish(x float64) float64 {
	return x * Sigmoid(x)
}

// CrossEntropy returns the cross entropy loss
func CrossEntropy(prediction, truth []float64) float64 {
	var loss float64
	for i := range prediction {
		loss += prediction[i]*math.Log(truth[i]) + (1-prediction[i])*math.Log(1-truth[i])
	}
	return loss
}

// Sensitivity returns the sensitivity
func Sensitivity(predicted, actual []float64) int {
	tp := TruePositivies(predicted, actual)
	fn := FalseNegatives(predicted, actual)
	return tp / (tp + fn)
}

// Specificity returns the specificity
func Specificity(predicted, actual []float64) int {
	fp := FalsePositives(predicted, actual)
	tn := TrueNegatives(predicted, actual)
	return fp / (fp + tn)
}

// Precision returns the precision.
func Precision(predicted, actual []float64) int {
	tp := TruePositivies(predicted, actual)
	fp := FalsePositives(predicted, actual)
	return tp / (tp + fp)
}

// Recall returns the recall.
func Recall(predicted, actual []float64) int {
	tp := TruePositivies(predicted, actual)
	fn := FalseNegatives(predicted, actual)
	return tp / (tp + fn)
}

// TrueNegative calculates the true negative predicted values
func TrueNegatives(predicted, actual []float64) int {
	var sum int
	for i := range predicted {
		if predicted[i] == actual[i] {
			sum++
		}
	}
	return sum
}

// TruePositivies returns the number of true positive predicted values.
func TruePositivies(predicted, actual []float64) int {
	var sum int
	for i := range predicted {
		if predicted[i] == actual[i] {
			sum++
		}
	}
	return sum
}

// FalsePositives estimates the false positive predicted values
func FalsePositives(predicted, actual []float64) int {
	var sum int
	for i := range predicted {
		if predicted[i] == actual[i] {
			sum++
		}
	}
	return sum
}

// FalseNegatives returns the number of false negative predicted values.
func FalseNegatives(predicted, actual []float64) int {
	var sum int
	for i := range predicted {
		if predicted[i] == actual[i] {
			sum++
		}
	}
	return sum
}

// MeanSqRoot returns the mean squared error between prediction and truth arrays.
func MeanSqRoot(prediction, truth []float64) float64 {
	loss := 0.0
	for i := range prediction {
		loss += math.Pow(truth[i]-prediction[i], 2)
	}
	return loss
}

//Rmse returns the root mean squared error between prediction and truth arrays.
func RootMeanSq(prediction, truth []float64) float64 {
	return math.Sqrt(MeanSqRoot(prediction, truth))
}
