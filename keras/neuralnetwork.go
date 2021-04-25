package keras

//Network defines a simple neural network architecture.
type Network struct {
	InputNodes, HiddenNodes, OutputNodes int
	WeightsIh, WeightsHo, BiasO, BiasH   Matrix
	LearningRate                         float64
}

//Package network implements the

// InitNetwork initializes the network with the number of nodes and the learning rate.
func InitNetwork(InputNodes, HiddenNodes, OutputNodes int, lr float64) Network {
	WeightsIh := RandomMatrix(HiddenNodes, InputNodes)
	WeightsHo := RandomMatrix(OutputNodes, HiddenNodes)
	BiasO := RandomMatrix(OutputNodes, 1)
	BiasH := RandomMatrix(HiddenNodes, 1)

	return Network{InputNodes: InputNodes,
		HiddenNodes:  HiddenNodes,
		OutputNodes:  OutputNodes,
		WeightsIh:    WeightsIh,
		WeightsHo:    WeightsHo,
		BiasH:        BiasH,
		BiasO:        BiasO,
		LearningRate: lr,
	}
}

// Train performs the training.
func (n *Network) Train(inputArray, targetArray []float64) {
	inputs := FromArray(inputArray)
	hidden := n.WeightsIh.Multiply(inputs)
	hidden.Add(n.BiasH)
	hidden.MapFunc(Sigmoid)

	output := n.WeightsHo.Multiply(hidden)
	output.Add(n.BiasO)
	output.MapFunc(Sigmoid)

	//Turn targets into
	targets := FromArray(targetArray)

	//Calculate error->still a matrix of values.
	absErrors := targets.Subtract(output)

	//Calculate gradient
	gradients := output.MapFunc(SigmoidPrime)
	gradients.Multiply(absErrors)
	gradients.ScalarAdition(n.LearningRate)

	//Derivatives
	devHidden := Transpose(hidden)
	WeightsHoDerivative := gradients.Multiply(devHidden)

	// Adjust the weights by deltas
	n.WeightsHo.Add(WeightsHoDerivative)
	n.BiasO.Add(gradients)

	// Calculate the hidden layer errors
	hiddenlayerError := Transpose(n.WeightsHo)
	hiddenErrors := hiddenlayerError.Multiply(absErrors)
	hiddenG := hidden.MapFunc(Sigmoid)
	hiddenG.Multiply(hiddenErrors)
	hiddenG.ScalarMultiplication(n.LearningRate)

	inputsTranspose := Transpose(inputs)
	weightIHDeltas := hiddenG.Multiply(inputsTranspose)
	n.WeightsIh.Add(weightIHDeltas)
	n.BiasH.Add(hiddenG)

	PrintByRows(output)
	PrintByRows(targets)
}

//Predict returns the model's prediction based on inputArray
func (n *Network) Predict(inputArray []float64) []float64 {
	inputs := FromArray(inputArray)
	hidden := n.WeightsIh.Multiply(inputs)
	hidden.Add(n.BiasH)
	hidden.MapFunc(Sigmoid)

	output := n.WeightsHo.Multiply(hidden)
	output.Add(n.BiasO)
	output.MapFunc(Sigmoid)
	return ToArray(output)
}
