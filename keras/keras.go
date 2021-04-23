// Package keras implements many functionalities from the popular Keras API
// which is a deep learning API written in Python, running on top of the
// machine learning platform TensorFlow. This package was started mainly for
// fun is still very much under development
package keras

import (
	"time"
)

//Model architecure to implement neural network.
type Model struct {
	//ConvLayers []Layer
	Name string
	//optimizer              Optimizer
	LossFunc   func([]float64, []float64) float64
	LossValues []float64
	Duration   time.Duration
	//Settings     []Metrics
	TrainDataX []float64
	TrainDataY []float64
	//Callbacks    []Callback
	Training     bool
	LearningRate float64
	TrainingLog  TrainingLog
}

//TrainingLog returns model's log
type TrainingLog []string

type Vector []float64
