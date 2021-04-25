package keras

//Network defines the neural network.
type Network struct {
	inputNodes, hiddenNodes, outputNodes int
	weightsIh, weightsHo, biasO, biasH   Matrix
	learningRate                         float64
}

//Package network implements the simple neural network  architecture.

//InitNetwork initializes the network with the number of nodes and the learning rate.
func InitNetwork(inputNodes, hiddenNodes, outputNodes int, lr float64) Network {
	weightsIh := RandomMatrix(hiddenNodes, inputNodes)
	weightsHo := RandomMatrix(outputNodes, hiddenNodes)
	biasO := RandomMatrix(outputNodes, 1)
	biasH := RandomMatrix(hiddenNodes, 1)

	return Network{inputNodes: inputNodes,
		hiddenNodes:  hiddenNodes,
		outputNodes:  outputNodes,
		weightsIh:    weightsIh,
		weightsHo:    weightsHo,
		biasH:        biasH,
		biasO:        biasO,
		learningRate: lr,
	}
}

//Train performs the training.
func (n *Network) Train(inputArray, targetArray []float64) {
	inputs := FromArray(inputArray)
	hidden := n.weightsIh.Multiply(inputs)
	hidden.Add(n.biasH)
	hidden.MapFunc(Sigmoid)

	output := n.weightsHo.Multiply(hidden)
	output.Add(n.biasO)
	output.MapFunc(Sigmoid)

	//Turn targets into
	targets := FromArray(targetArray)

	//Calculate error->still a matrix of values.
	absErrors := targets.Subtract(output)

	//Calculate gradient
	gradients := output.MapFunc(SigmoidPrime)
	gradients.Multiply(absErrors)
	gradients.ScalarAdition(n.learningRate)

	//Derivatives
	devHidden := Transpose(hidden)
	weightsHoDerivative := gradients.Multiply(devHidden)

	// Adjust the weights by deltas
	n.weightsHo.Add(weightsHoDerivative)
	n.biasO.Add(gradients)

	// Calculate the hidden layer errors
	hiddenlayerError := Transpose(n.weightsHo)
	hiddenErrors := hiddenlayerError.Multiply(absErrors)
	hiddenG := hidden.MapFunc(Sigmoid)
	hiddenG.Multiply(hiddenErrors)
	hiddenG.ScalarMultiplication(n.learningRate)

	inputsTranspose := Transpose(inputs)
	weightIHDeltas := hiddenG.Multiply(inputsTranspose)
	n.weightsIh.Add(weightIHDeltas)
	n.biasH.Add(hiddenG)

	output.PrintByRow()
	targets.PrintByRow()
}

//Predict returns the model's prediction based on inputArray
func (n *Network) Predict(inputArray []float64) []float64 {
	inputs := FromArray(inputArray)
	hidden := n.weightsIh.Multiply(inputs)
	hidden.Add(n.biasH)
	hidden.MapFunc(Sigmoid)

	output := n.weightsHo.Multiply(hidden)
	output.Add(n.biasO)
	output.MapFunc(Sigmoid)
	return output.ToArray()
}
