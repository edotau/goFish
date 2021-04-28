package numerical

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/edotau/goFish/stats"
)

type OptimizationMethod string

type LocalLinear struct {
	// If maxIterations is 0, then GradientAscent will run until the algorithm detects convergance.
	// alpha and maxIterations are used only for GradientAscent during learning.
	alpha          float64
	regularization float64
	bandwidth      float64
	maxIterations  int
	// method is the optimization method used when training the model
	method OptimizationMethod
	// trainingSet and expectedResults are the 'x', and 'y' of the data, expressed as vectors, that the model can optimize from
	trainingSet     [][]float64
	expectedResults []float64
	Parameters      []float64 `json:"theta"`

	Output io.Writer
}

type LeastSquares struct {
	// alpha and maxIterations are used only for GradientAscent during learning.
	// If maxIterations is 0, then GradientAscent will run until the algorithm detects convergence.
	alpha          float64
	regularization float64
	maxIterations  int
	method         OptimizationMethod

	// trainingSet and expectedResults are the 'x', and 'y' of the data, expressed as vectors, that the model can optimize from
	trainingSet     [][]float64
	expectedResults []float64

	Parameters []float64 `json:"theta"`
	Output     io.Writer
}

func NewLocalLinear(method OptimizationMethod, alpha, regularization, bandwidth float64, maxIterations int, trainingSet [][]float64, expectedResults []float64) *LocalLinear {
	var params []float64
	if trainingSet == nil || len(trainingSet) == 0 {
		params = []float64{}
	} else {
		params = make([]float64, len(trainingSet[0])+1)
	}

	return &LocalLinear{
		alpha:          alpha,
		regularization: regularization,
		bandwidth:      bandwidth,
		maxIterations:  maxIterations,

		method: method,

		trainingSet:     trainingSet,
		expectedResults: expectedResults,
		// initialize θ as the zero vector (that is, the vector of all zeros)
		Parameters: params,

		Output: os.Stdout,
	}
}

func NewLeastSquares(method OptimizationMethod, alpha, regularization float64, maxIterations int, trainingSet [][]float64, expectedResults []float64, features ...int) *LeastSquares {
	var params []float64
	if len(features) != 0 {
		params = make([]float64, features[0]+1)
	} else if trainingSet == nil || len(trainingSet) == 0 {
		params = []float64{}
	} else {
		params = make([]float64, len(trainingSet[0])+1)
	}

	return &LeastSquares{
		alpha:          alpha,
		regularization: regularization,
		maxIterations:  maxIterations,

		method: method,

		trainingSet:     trainingSet,
		expectedResults: expectedResults,
		Parameters:      params,

		Output: os.Stdout,
	}
}

// LearningRate returns the learning rate α for gradient descent to optimize the model. Could vary as a function of something else later, potentially.
func (l *LocalLinear) LearningRate() float64 {
	return l.alpha
}

// Examples returns the number of training examples (m) that the model currently is training from.
func (l *LocalLinear) Examples() int {
	return len(l.trainingSet)
}

// MaxIterations returns the number of maximum iterations the model will go through in GradientAscent, in the worst case
func (l *LocalLinear) MaxIterations() int {
	return l.maxIterations
}

// Predict takes in a variable x (an array of floats,) and finds the value of the hypothesis function given the current parameter vector θ
// if normalize is given as true, then the input will first be normalized to unit length. Only use this if you trained off of normalized inputs and are feeding an un-normalized input
func (l *LocalLinear) Predict(x []float64, normalize ...bool) ([]float64, error) {
	if len(x)+1 != len(l.Parameters) {
		err := fmt.Errorf("ERROR: Parameter vector should be 1 longer than input vector!\n\tLength of x given: %v\n\tLength of parameters: %v\n", len(x), len(l.Parameters))
		print(err.Error())
		return nil, err
	}

	norm := len(normalize) != 0 && normalize[0]
	if norm {
		stats.NormalizePoint(x)
	}

	if l.trainingSet == nil || l.expectedResults == nil {
		err := fmt.Errorf("ERROR: Attempting to learn with no training examples!\n")
		print(err.Error())
		return nil, err
	}

	examples := len(l.trainingSet)
	if examples == 0 || len(l.trainingSet[0]) == 0 {
		err := fmt.Errorf("ERROR: Attempting to learn with no training examples!\n")
		print(err.Error())
		return nil, err
	}
	if len(l.expectedResults) == 0 {
		err := fmt.Errorf("ERROR: Attempting to learn with no expected results! This isn't an unsupervised model!! You'll need to include data before you learn :)\n")
		print(err.Error())
		return nil, err
	}

	//fmt.Fprintf(l.Output, "Training:\n\tModel: Locally Weighted Linear Regression\n\tOptimization Method: %v\n\tCenter Point: %v\n\tTraining Examples: %v\n\tFeatures: %v\n\tLearning Rate α: %v\n\tRegularization Parameter λ: %v\n...\n\n", l.method, x, examples, len(l.trainingSet[0]), l.alpha, l.regularization)

	var iter int
	features := len(l.Parameters)

	if l.method == BatchGA {
		for ; iter < l.maxIterations; iter++ {
			newTheta := make([]float64, features)
			for j := range l.Parameters {
				dj, err := l.Dj(x, j)
				if err != nil {
					return nil, err
				}

				newTheta[j] = l.Parameters[j] + l.alpha*dj
			}
			// now simultaneously update Theta
			for j := range l.Parameters {
				newθ := newTheta[j]
				if math.IsInf(newθ, 0) || math.IsNaN(newθ) {
					return nil, fmt.Errorf("Sorry! Learning diverged. Some value of the parameter vector theta is ±Inf or NaN")
				}
				l.Parameters[j] = newθ
			}
		}
	} else if l.method == StochasticGA {
		for ; iter < l.maxIterations; iter++ {
			newTheta := make([]float64, features)
			for i := 0; i < examples; i++ {
				for j := range l.Parameters {
					dj, err := l.Dij(x, i, j)
					if err != nil {
						return nil, err
					}

					newTheta[j] = l.Parameters[j] + l.alpha*dj
				}

				// now simultaneously update Theta
				for j := range l.Parameters {
					newθ := newTheta[j]
					if math.IsInf(newθ, 0) || math.IsNaN(newθ) {
						return nil, fmt.Errorf("Sorry! Learning diverged. Some value of the parameter vector theta is ±Inf or NaN")
					}
					l.Parameters[j] = newθ
				}
			}
		}
	} else {
		return nil, fmt.Errorf("Chose a training method not implemented for LocalLinear regression")
	}

	//fmt.Fprintf(l.Output, "Training Completed. Went through %v iterations.\n%v\n\n", iter, l)

	// include constant term in sum
	sum := l.Parameters[0]

	for i := range x {
		sum += x[i] * l.Parameters[i+1]
	}
	return []float64{sum}, nil
}

func (l *LeastSquares) Predict(x []float64, normalize ...bool) ([]float64, error) {
	if len(x)+1 != len(l.Parameters) {
		return nil, fmt.Errorf("Error: Parameter vector should be 1 longer than input vector!\n\tLength of x given: %v\n\tLength of parameters: %v\n", len(x), len(l.Parameters))
	}

	if len(normalize) != 0 && normalize[0] {
		stats.NormalizePoint(x)
	}

	// include constant term in sum
	sum := l.Parameters[0]

	for i := range x {
		sum += x[i] * l.Parameters[i+1]
	}

	return []float64{sum}, nil
}

// LearningRate returns the learning rate α for gradient descent to optimize the model. Could vary as a function of something else later, potentially.
func (l *LeastSquares) LearningRate() float64 {
	return l.alpha
}

// Examples returns the number of training examples (m) that the model currently is training from.
func (l *LeastSquares) Examples() int {
	return len(l.trainingSet)
}

// MaxIterations returns the number of maximum iterations the model will go through in GradientAscent, in the worst case
func (l *LeastSquares) MaxIterations() int {
	return l.maxIterations
}

// Theta returns the parameter vector θ for use in persisting the model, and optimizing the model through gradient descent ( or other methods like Newton's Method)
func (l *LeastSquares) Theta() []float64 {
	return l.Parameters
}

// Learn takes the struct's dataset and expected results and runs batch gradient descent on them, optimizing theta so you can predict based on those results
func (l *LeastSquares) Learn() error {
	if l.trainingSet == nil || l.expectedResults == nil {
		err := fmt.Errorf("ERROR: Attempting to learn with no training examples!\n")
		fmt.Fprintf(l.Output, err.Error())
		return err
	}

	examples := len(l.trainingSet)
	if examples == 0 || len(l.trainingSet[0]) == 0 {
		err := fmt.Errorf("ERROR: Attempting to learn with no training examples!\n")
		fmt.Fprintf(l.Output, err.Error())
		return err
	}
	if len(l.expectedResults) == 0 {
		err := fmt.Errorf("ERROR: Attempting to learn with no expected results! This isn't an unsupervised model!! You'll need to include data before you learn :)\n")
		fmt.Fprintf(l.Output, err.Error())
		return err
	}

	//fmt.Fprintf(l.Output, "Training:\n\tModel: Logistic (Binary) Classification\n\tOptimization Method: %v\n\tTraining Examples: %v\n\tFeatures: %v\n\tLearning Rate α: %v\n\tRegularization Parameter λ: %v\n...\n\n", l.method, examples, len(l.trainingSet[0]), l.alpha, l.regularization)

	var err error
	if l.method == BatchGA {
		err = GradientAscent(l)
	} else if l.method == StochasticGA {
		err = StochasticGradientAscent(l)
	} else {
		err = fmt.Errorf("Chose a training method not implemented for LeastSquares regression")
	}

	if err != nil {
		fmt.Fprintf(l.Output, "\nERROR: Error while learning –\n\t%v\n\n", err)
		return err
	}
	//fmt.Fprintf(l.Output, "Training Completed.\n%v\n\n", l)
	return nil
}

// weight corresponds to the weight given between two datapoints (based on how 'far apart' they are.)
// w[i] = exp(-1 * |x[i] - x|^2 / 2σ^2)
func (l *LocalLinear) weight(X []float64, x []float64) float64 {
	// don't throw error but fail peacefully returning "not at all similar", basically
	if len(X) != len(x) {
		return 0.0
	}

	var diff float64

	for i := range X {
		diff += (X[i] - x[i]) * (X[i] - x[i])
	}

	return math.Exp(-1 * diff / (2 * l.bandwidth * l.bandwidth))
}

// Dj returns the partial derivative of the cost function J(θ) with respect to theta[j] where theta is the parameter vector associated with our hypothesis function Predict (upon which we are optimizing
func (l *LocalLinear) Dj(input []float64, j int) (float64, error) {
	if j > len(l.Parameters)-1 {
		return 0, fmt.Errorf("J (%v) would index out of the bounds of the training set data (len: %v)", j, len(l.Parameters))
	}
	if len(input) != len(l.Parameters)-1 {
		return 0, fmt.Errorf("Length of input x (%v) should be one less than the length of the parameter vector (len: %v)", len(input), len(l.Parameters))
	}

	var sum float64

	for i := range l.trainingSet {
		prediction := l.Parameters[0]
		for k := 1; k < len(l.Parameters); k++ {
			prediction += l.Parameters[k] * input[k-1]
		}
		// account for constant term x is x[i][j] via Andrew Ng's terminology
		var x float64
		if j == 0 {
			x = 1
		} else {
			x = l.trainingSet[i][j-1]
		}

		sum += l.weight(l.trainingSet[i], input) * (l.expectedResults[i] - prediction) * x
	}

	// add in the regularization term λ*θ[j] notice that we don't count the constant term
	if j != 0 {
		sum += l.regularization * l.Parameters[j]
	}

	return sum, nil
}

// Dij returns the derivative of the cost function J(θ) with respect to the j-th parameter of the hypothesis, θ[j], for the training example x[i].
// Used in Stochastic Gradient Descent. assumes that i,j is within the bounds of the data they are looking up! (because this is getting called so much, it needs to be efficient with comparisons)
func (l *LocalLinear) Dij(input []float64, i, j int) (float64, error) {
	if j > len(l.Parameters)-1 || i > len(l.trainingSet)-1 {
		return 0, fmt.Errorf("j (%v) or i (%v) would index out of the bounds of the training set data (len: %v)", j, i, len(l.Parameters))
	}
	if len(input) != len(l.Parameters)-1 {
		return 0, fmt.Errorf("Length of input x (%v) should be one less than the length of the parameter vector (len: %v)", len(input), len(l.Parameters))
	}

	prediction := l.Parameters[0]
	for k := 1; k < len(l.Parameters); k++ {
		prediction += l.Parameters[k] * input[k-1]
	}

	// account for constant term x is x[i][j] via Andrew Ng's terminology
	var x float64
	if j == 0 {
		x = 1
	} else {
		x = l.trainingSet[i][j-1]
	}

	var gradient float64
	gradient = l.weight(l.trainingSet[i], input) * (l.expectedResults[i] - prediction) * x
	// add in the regularization term λ*θ[j] notice that we don't count the constant term
	if j != 0 {
		gradient += l.regularization * l.Parameters[j]
	}

	return gradient, nil
}

// J returns the Least Squares cost function of the given linear model. Could be usefull in testing convergance
func (l *LocalLinear) J() (float64, error) {
	var sum float64

	for i := range l.trainingSet {
		prediction, err := l.Predict(l.trainingSet[i])
		if err != nil {
			return 0, err
		}

		sum += (l.expectedResults[i] - prediction[0]) * (l.expectedResults[i] - prediction[0])
	}

	// add regularization term!
	//
	// notice that the constant term doesn't matter
	for i := 1; i < len(l.Parameters); i++ {
		sum += l.regularization * l.Parameters[i] * l.Parameters[i]
	}

	return sum / float64(2*len(l.trainingSet)), nil
}

// Dj returns the partial derivative of the cost function J(θ) with respect to theta[j] where theta is the parameter vector associated with our hypothesis function Predict (upon which we are optimizing
func (l *LeastSquares) Dj(j int) (float64, error) {
	if j > len(l.Parameters)-1 {
		return 0, fmt.Errorf("J (%v) would index out of the bounds of the training set data (len: %v)", j, len(l.Parameters))
	}
	var sum float64
	for i := range l.trainingSet {
		prediction, err := l.Predict(l.trainingSet[i])
		if err != nil {
			return 0, err
		}
		// account for constant term x is x[i][j] via Andrew Ng's terminology
		var x float64
		if j == 0 {
			x = 1
		} else {
			x = l.trainingSet[i][j-1]
		}
		sum += (l.expectedResults[i] - prediction[0]) * x
	}
	// add in the regularization term λ*θ[j] notice that we don't count the constant term
	if j != 0 {
		sum += l.regularization * l.Parameters[j]
	}
	return sum, nil
}

func (l *LeastSquares) Dij(i int, j int) (float64, error) {
	prediction, err := l.Predict(l.trainingSet[i])
	if err != nil {
		return 0, err
	}
	// account for constant term x is x[i][j] via Andrew Ng's terminology
	var x float64
	if j == 0 {
		x = 1
	} else {
		x = l.trainingSet[i][j-1]
	}
	var gradient float64
	gradient = (l.expectedResults[i] - prediction[0]) * x
	// add in the regularization term λ*θ[j] notice that we don't count the constant term
	if j != 0 {
		gradient += l.regularization * l.Parameters[j]
	}
	return gradient, nil
}
