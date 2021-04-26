// Package keras implements many functionalities from the popular Keras API
// which is a deep learning API written in Python, running on top of the
// machine learning platform TensorFlow. This package was started mainly for
// fun is still very much under development
package keras

import (
	"fmt"
	"time"
)

//Model architecure to implement neural network.
type Model struct {
	ConvLayers   []Layer
	Name         string
	Optimizer    Optimizer
	LossFunc     func([]float64, []float64) float64
	LossValues   []float64
	Duration     time.Duration
	Settings     []Metrics
	TrainDataX   []float64
	TrainDataY   []float64
	Callbacks    []Callback
	Training     bool
	LearningRate float64
	TrainingLog  TrainingLog
}

//TrainingLog returns model's log
type TrainingLog []string

type Metrics interface {
	Measure([]float64, []float64) float64
	Name() string
}

//Optimizer interface requires an ApplyGradients function. Pass it to the model compilation.
type Optimizer interface {
	ApplyGradients()
}

//Sequential returns a model given layers and a name.
func Sequential(layers []Layer, name string) *Model {
	return &Model{ConvLayers: layers, Name: name}
}

//Add method adds a layer to the end of the model architecture
func (m *Model) Add(layer Layer) *Model {
	m.ConvLayers[len(m.ConvLayers)] = layer
	return m
}

//GetLayerByIndex returns the ith layer.
func (m *Model) GetLayerByIndex(index int) Layer {
	return m.ConvLayers[index]
}

//GetMetricsByIndex returns the index's model metric
func (m *Model) GetMetricsByIndex(index int) Metrics {
	return m.Settings[index]
}

//GetLayerByName returns the layer given its name.
func (m *Model) GetLayerByName(name string) Layer {
	for i := range m.ConvLayers {
		if m.ConvLayers[i].Name() == name {
			return m.ConvLayers[i]
		}
	}
	return m.ConvLayers[0]
}

//Compile compiles the model given the optimizer, loss and metrics
func (m *Model) Compile(optimizer Optimizer, loss func([]float64, []float64) float64, ms []Metrics) {
	m.Optimizer = optimizer
	m.LossFunc = loss
	m.Settings = ms
}

//Predict does the feed forward magic when fed the inputs.
func (m *Model) Predict(values []float64) []float64 {
	var outputs []float64
	for i := range m.ConvLayers {
		outputs = m.ConvLayers[i].Call()
		m.ConvLayers[i+1].Call()
		if i == len(m.ConvLayers)-1 {
			return outputs
		}
	}
	return outputs
}

// Train trains the model given trainX and  trainY data and the number of epochs. It keeps track of the defined metrics and prints it every epoch. It also prints the training duration.
//It returns a map from strings to floats, where strings represent the metrics name and float the metrics value.
func (m *Model) Train(trainX, trainY []float64, epochs int) map[string]float64 {
	startTime := time.Now()
	metricsValues := make(map[string]float64, len(m.Settings))
	for i := 1; i < epochs; i++ {
		for j := 0; j < len(trainX); j++ {
			lossValue := m.LossFunc(m.Predict(trainX), trainY)
			m.LossValues = append(m.LossValues, lossValue)
			m.Optimizer.ApplyGradients()
		}
		avg := meanValue(m.LossValues)
		for _, met := range m.Settings {
			metricsValues[met.Name()] = met.Measure(m.Predict(trainX), trainY)
		}
		fmt.Printf("Epoch: %d		Loss:%.4f\n", i, avg)
	}
	endTime := time.Now()
	m.Duration = endTime.Sub(startTime)
	fmt.Printf("Training duration: %s\n", m.Duration.String())
	return metricsValues
}

//Summary prints the layer by layer summaary along with trainable parameters.
func (m *Model) Summary() {
	var sum int
	for i := range m.ConvLayers {
		tp := m.ConvLayers[i].TrainableParameters()
		sum += tp
		fmt.Printf("name: %s		trainable parameters: %d\n", m.ConvLayers[i].Name(), tp)
	}
	fmt.Println("Trainable parameters: ", sum)
}
