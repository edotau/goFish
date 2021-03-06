package keras

import (
	"math"
	"math/rand"
	"time"
)

//Layer interface given these 5 functions which every layer must have.
type Layer interface {
	Call() []float64
	GetWeights() Matrix
	GetBiases() Vector
	Name() string
	TrainableParameters() int
}

//DenseLayer defines a fully connected layer.
type DenseLayer struct {
	units             int
	inputs, outputs   []float64
	weights           Weights
	biases            Biases
	trainable         bool
	name              string
	kernelRegularizer func([]float64) []float64
	biasRegularizer   func([]float64) []float64
	Activation        func(float64) float64
	KernelInit        func(float64) float64
	BiasInit          func(float64) float64
}

//DenseLayer defines a fully connected layer.
type Conv2DLayer struct {
	Filters           int
	Inputs, Outputs   []float64
	Weights           Weights
	Biases            Biases
	trainable         bool
	name              string
	kernelRegularizer func([]float64) []float64
	biasRegularizer   func([]float64) []float64
	Activation        func(float64) float64
	KernelInit        func(float64) float64
	BiasInit          func(float64) float64
}

//Dense fully connected layer initializer
func Dense(units int, inputs []float64, activation func(float64) float64) DenseLayer {
	weights := WeightInit(units, len(inputs), HeUniform)
	biases := BiasInit(units, ZeroInitializer)
	return DenseLayer{units: units,
		inputs:     inputs,
		Activation: activation,
		weights:    weights,
		biases:     biases,
	}
}

//Conv2D initializes a Conv2DLayer
func Conv2D(numFilter int, x int, y int, inputs []float64, padding func(x, y int) Weights) Conv2DLayer {
	weights := padding(x, y)
	biases := BiasInit(numFilter, ZeroInitializer)
	return Conv2DLayer{Filters: numFilter,
		Inputs:     inputs,
		Activation: ReLU,
		Weights:    weights,
		Biases:     biases,
	}
}

func Valid(x, y int) Weights {
	weights := WeightInit(x, y, HeUniform)
	return weights
}

func DefaultPadding(x, y int) Weights {
	w := RandomMatrix(x, y)
	padding := Zeros(Max(x, y), Max(x, y))
	var i, j int
	for i = 0; i < RowNum(w); i++ {
		for j = 0; j < ColNum(w); j++ {
			padding.Matrix[i][j] = w.Matrix[i][j] + padding.Matrix[i][j]
		}
	}
	padding.MapFunc(HeUniform)
	return Weights{Kernels: padding, KernelInit: HeUniform}
}

func Max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func (cd Conv2DLayer) Call() []float64 {
	vec := NewVector(cd.Inputs).ApplyMatrix(cd.Weights.Kernels).Add(cd.Biases.bs)
	return vec.Map(cd.Activation).Slice()
}

//Name of the dense layer
func (cd Conv2DLayer) Name() string {
	return cd.name
}

//GetWeights returns the layer's weights.
func (cd Conv2DLayer) GetWeights() Matrix {
	return cd.Weights.Kernels
}

//GetBiases returns the layer's biases.
func (cd Conv2DLayer) GetBiases() Vector {
	return cd.Biases.bs
}

//TrainableParameters returns the count of trainable parameters.
func (cd Conv2DLayer) TrainableParameters() int {
	return NumberOfElements(cd.Weights.Kernels) + cd.Biases.bs.NumberOfElements()
}

/*
tf.keras.layers.Conv2D(
    filters,
    kernel_size,
    strides=(1, 1),
    padding="valid",
    data_format=None,
    dilation_rate=(1, 1),
    groups=1,
    activation=None,
    use_bias=True,
    kernel_initializer="glorot_uniform",
    bias_initializer="zeros",
    kernel_regularizer=None,
    bias_regularizer=None,
    activity_regularizer=None,
    kernel_constraint=None,
    bias_constraint=None,
    **kwargs
)*/

// Weights struct with the actual Kernels and the kernel initializer function.
type Weights struct {
	Kernels    Matrix
	KernelInit func(float64) float64
}

//Biases struct with the actual biases and the bias initializer function.
type Biases struct {
	bs       Vector
	BiasInit func(float64) float64
}

type Shape []float64

//WeightInit used for weight initialization. Already defined at the initialization of the dense layer.
func WeightInit(a, b int, kernelInit func(float64) float64) Weights {
	w := RandomMatrix(a, b).MapFunc(kernelInit)
	return Weights{Kernels: w, KernelInit: kernelInit}
}

//BiasInit used for bias initialization. Already defined at the initialization of the dense layer.
func BiasInit(a int, biasInit func(float64) float64) Biases {
	bs := RandomVector(a).Map(biasInit)
	return Biases{bs: bs, BiasInit: biasInit}
}

//Call of the dense layer.Outputs the next tensors.
func (d DenseLayer) Call() []float64 {
	vec := NewVector(d.inputs).ApplyMatrix(d.weights.Kernels).Add(d.biases.bs)
	return vec.Map(d.Activation).Slice()
}

//Name of the dense layer
func (d DenseLayer) Name() string {
	return d.name
}

//GetWeights returns the layer's weights.
func (d DenseLayer) GetWeights() Matrix {
	return d.weights.Kernels
}

//GetBiases returns the layer's biases.
func (d DenseLayer) GetBiases() Vector {
	return d.biases.bs
}

//TrainableParameters returns the count of trainable parameters.
func (d DenseLayer) TrainableParameters() int {
	return NumberOfElements(d.weights.Kernels) + d.biases.bs.NumberOfElements()
}

//SetWeights is used for manually defining the weight
func (d *DenseLayer) SetWeights(Kernels Matrix) {
	d.weights.Kernels = Kernels
}

//SetBiases is used for manually defining the bias vector.
func (d *DenseLayer) SetBiases(bs Vector) {
	d.biases.bs = bs
}

//InputLayer layer, much like the keras one.
type InputLayer struct {
	Inputs, outputs []float64
	Weights         Weights
	Biases          Biases
	Trainable       bool
	Name            string
}

//Input
func Input(inputs []float64) InputLayer {
	weights := WeightInit(len(inputs), 1, HeUniform)
	biases := BiasInit(len(inputs), ZeroInitializer)
	return InputLayer{
		Inputs:  inputs,
		Weights: weights,
		Biases:  biases,
	}
}

//Call of the input layer
func (i *InputLayer) Call() []float64 {
	vec := NewVector(i.Inputs).ApplyMatrix(i.Weights.Kernels).Add(i.Biases.bs)
	i.outputs = vec.Slice()
	return vec.Slice()
}

//BatchNormLayer layer
type BatchNormLayer struct {
	inputs, outputs      []float64
	beta, epsilon, alpha float64
	trainable            bool
	name                 string
}

//BatchNorm init
func BatchNorm(inputs []float64) BatchNormLayer {
	return BatchNormLayer{inputs: inputs}
}

//Call for the batch normalization layer
func (bn *BatchNormLayer) Call() []float64 {
	outputs := make([]float64, len(bn.inputs))
	variance := Variance(bn.inputs)
	mean := meanValue(bn.inputs)
	for _, x := range bn.inputs {
		newX := (x - mean) / math.Sqrt(variance+bn.epsilon)
		outputs = append(outputs, bn.alpha*newX+bn.beta)
	}
	bn.outputs = outputs
	return outputs
}

//Variance returns the variance
func Variance(fls []float64) float64 {
	var sum float64
	for _, f := range fls {
		sum += math.Pow(f-meanValue(fls), 2)
	}
	return sum / float64(len(fls))
}

func meanValue(fls []float64) float64 {
	mean := sum(fls) / float64(len(fls))
	return mean
}

func sum(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

//DropoutLayer layer
type DropoutLayer struct {
	inputs []float64
	rate   float64
}

//Dropout init
func Dropout(inputs []float64, rate float64) DropoutLayer {
	return DropoutLayer{inputs: inputs, rate: rate}

}

//Call for the dropout layer
func (dr *DropoutLayer) Call() []float64 {
	weightCount := dr.rate * float64(len(dr.inputs))
	for i := int(weightCount); i > 0; i-- {
		if len(dr.inputs)%int(weightCount) == 0 {
			dr.inputs[i] = 0
		}
	}
	return dr.inputs
}

//SoftmaxLayer layer
type SoftmaxLayer struct {
	inputs, outputs []float64
	classes         int
}

// Softmax returns the softmax layer based on values.
func Softmax(inputs []float64, classes int) SoftmaxLayer {
	return SoftmaxLayer{inputs: inputs, classes: classes}
}

// Call of the softmax
func (s *SoftmaxLayer) Call() []float64 {
	sum := 0.0
	preds := make([]float64, len(s.inputs))
	for i, n := range s.inputs {
		preds[i] -= math.Exp(n - findMax(s.inputs))
		sum += preds[i]
	}
	for k := range preds {
		preds[k] /= sum
	}
	outputs := preds[:s.classes]
	s.outputs = outputs
	return outputs
}

func findMax(fls []float64) float64 {
	max := -10000.0
	for _, k := range fls {
		if k > max {
			max = k
		}
	}
	return max
}

// FlattenLayer layer
type FlattenLayer struct {
	inputs, outputs []float64
	name            string
	trainable       bool
}

// Call of the FlattenLayer
func (f *FlattenLayer) Call() []float64 {
	return f.outputs
}

// Flatten init.
func Flatten(m Matrix) FlattenLayer {
	return FlattenLayer{outputs: ToArray(m)}
}

// HeUniform stands for He Initialization or the glorot_unifom for kernel_initialization.
func HeUniform(x float64) float64 {
	rand.Seed(time.Now().UnixNano())
	down, upper := x-0.4, x+0.4
	return down + rand.Float64()*(upper-down)
}

// ZeroInitializer returns the zeros initializer for the bias initialization
func ZeroInitializer(x float64) float64 {
	return 0
}

// OnesInitializer returns the ones initializer for the bias initialization
func OnesInitializer(x float64) float64 {
	return 1
}
