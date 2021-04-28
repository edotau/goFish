// Modified from https://github.com/cdipaolo/goml for the sole purpose of teaching myself how to implement machine learning models from scratch
// Original code and documentation can be found at: https://github.com/cdipaolo/goml
package numerical

import (
	"fmt"
	"math"
)

const (
	BatchGA      OptimizationMethod = "Batch Gradient Ascent"
	StochasticGA                    = "Stochastic Gradient Descent"
)

// Ascendable is an interface that can be used with batch gradient descent where the parameter vector theta is in one dimension only (so softmax regression would need it's own model, for example)
type Ascendable interface {
	// LearningRate returns the learning rate α to be used in Gradient Descent as the modifier term
	LearningRate() float64
	// Dj returns the derivative of the cost function J(θ) with respect to the j-th parameter of the hypothesis, θ[j]. Called as Dj(j)
	Dj(int) (float64, error)
	// Theta returns a pointer to the parameter vector theta, which is 1D vector of floats
	Theta() []float64
	// MaxIterations returns the maximum number of iterations to try using gradient ascent. Might return after less if strong convergance is detected, but it'll let the user set a cap.
	MaxIterations() int
}

// StochasticAscendable is an interface that can be used with stochastic gradient descent where the parameter vector theta is in one dimension only (so softmax regression would need it's own model)
type StochasticAscendable interface {
	// LearningRate returns the learning rate α to be used in Gradient Descent as the modifier term
	LearningRate() float64
	// Examples returns the number of examples in the training set the model is using
	Examples() int
	// Dj returns the derivative of the cost function J(θ) with respect to the j-th parameter of the hypothesis, θ[j], for the training example x[i]. Called as Dij(i,j)
	Dij(int, int) (float64, error)
	// Theta returns a pointer to the parameter vector theta, which is 1D vector of floats
	Theta() []float64
	// MaxIterations returns the maximum number of iterations to try using gradient ascent. Might return after less if strong convergance is detected, but it'll let the user set a cap.
	MaxIterations() int
}

// Gradient Ascent follows the following algorithm:
// θ[j] := θ[j] + α·∇J(θ) where J(θ) is the cost function, α is the learning rate, and θ[j] is the j-th value in the parameter vector
func GradientAscent(d Ascendable) error {
	Theta := d.Theta()
	Alpha := d.LearningRate()
	MaxIterations := d.MaxIterations()
	// if the iterations given is 0, set it to be 250 (seems reasonable base value)
	if MaxIterations == 0 {
		MaxIterations = 250
	}

	var iter int
	features := len(Theta)

	// Stop iterating if the number of iterations exceeds the limit
	for ; iter < MaxIterations; iter++ {
		newTheta := make([]float64, features)
		for j := range Theta {
			dj, err := d.Dj(j)
			if err != nil {
				return err
			}
			newTheta[j] = Theta[j] + Alpha*dj
		}
		// now simultaneously update Theta
		for j := range Theta {
			newθ := newTheta[j]
			if math.IsInf(newθ, 0) || math.IsNaN(newθ) {
				return fmt.Errorf("Sorry! Learning diverged. Some value of the parameter vector theta is ±Inf or NaN")
			}
			Theta[j] = newθ
		}
	}
	return nil
}

// StochasticGradientAscent operates on a StochasticAscendable model and further optimizes the parameter vector Theta of the model, which is then used within the Predict function.
// Stochastic gradient descent updates the parameter vector after looking at each individual training example, which can result in never converging to the absolute minimum; even raising the cost function potentially, but it will typically converge faster than batch gradient descent
// Gradient Ascent follows the following algorithm:
// θ[j] := θ[j] + α·∇J(θ) where J(θ) is the cost function, α is the learningrate, and θ[j] is the j-th value in the parameter vector
func StochasticGradientAscent(d StochasticAscendable) error {
	Theta := d.Theta()
	Alpha := d.LearningRate()
	MaxIterations := d.MaxIterations()
	Examples := d.Examples()

	// if the iterations given is 0, set it to be 250 (seems reasonable base value)
	if MaxIterations == 0 {
		MaxIterations = 250
	}

	var iter int
	features := len(Theta)

	// Stop iterating if the number of iterations exceeds the limit
	for ; iter < MaxIterations; iter++ {
		newTheta := make([]float64, features)
		for i := 0; i < Examples; i++ {
			for j := range Theta {
				dj, err := d.Dij(i, j)
				if err != nil {
					return err
				}

				newTheta[j] = Theta[j] + Alpha*dj
			}
			// now simultaneously update Theta
			for j := range Theta {
				newθ := newTheta[j]
				if math.IsInf(newθ, 0) || math.IsNaN(newθ) {
					return fmt.Errorf("Sorry! Learning diverged. Some value of the parameter vector theta is ±Inf or NaN")
				}
				Theta[j] = newθ
			}
		}
	}
	return nil
}
