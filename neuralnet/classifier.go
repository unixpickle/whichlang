package neuralnet

import (
	"encoding/json"
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/whichlang/tokens"
)

// A Network is a feedforward neural network with
// a single hidden layer.
type Network struct {
	Tokens []string
	Langs  []string

	// HiddenCount is the number of hidden neurons
	// in this network.
	HiddenCount int

	// In the following weights, the last weight for
	// each neuron corresponds to a constant shift,
	// and is not multiplied by an input's value.
	HiddenWeights [][]float64
	OutputWeights [][]float64
}

func DecodeNetwork(data []byte) (*Network, error) {
	var n Network
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func (n *Network) Classify(f tokens.Freqs) string {
	outputSums := make([]*kahan.Summer64, len(n.OutputWeights))
	for i := range outputSums {
		outputSums[i] = kahan.NewSummer64()

		// The last weight is a constant offset.
		outputSums[i].Add(n.OutputWeights[i][len(n.Langs)])
	}

	for hiddenIndex, hiddenWeights := range n.HiddenWeights {
		hiddenSum := kahan.NewSummer64()
		for j, token := range n.Tokens {
			hiddenSum.Add(f[token] * hiddenWeights[j])
		}

		// The last weight is a constant offset.
		hiddenSum.Add(hiddenWeights[len(n.Tokens)])

		hiddenOut := sigmoid(hiddenSum.Sum())
		for j, outSum := range outputSums {
			weight := n.OutputWeights[j][hiddenIndex]
			outSum.Add(weight * hiddenOut)
		}
	}

	maxSum := outputSums[0].Sum()
	maxIdx := 0
	for i, x := range outputSums {
		if x.Sum() > maxSum {
			maxSum = x.Sum()
			maxIdx = i
		}
	}
	return n.Langs[maxIdx]
}

func (n *Network) Encode() []byte {
	enc, _ := json.Marshal(n)
	return enc
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
