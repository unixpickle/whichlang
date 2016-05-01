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

func (n *Network) Copy() *Network {
	res := &Network{
		Tokens:        make([]string, len(n.Tokens)),
		Langs:         make([]string, len(n.Langs)),
		HiddenWeights: make([][]float64, len(n.HiddenWeights)),
		OutputWeights: make([][]float64, len(n.OutputWeights)),
	}
	copy(res.Tokens, n.Tokens)
	copy(res.Langs, n.Langs)
	for i, w := range n.HiddenWeights {
		res.HiddenWeights[i] = make([]float64, len(w))
		copy(res.HiddenWeights[i], w)
	}
	for i, w := range n.OutputWeights {
		res.OutputWeights[i] = make([]float64, len(w))
		copy(res.OutputWeights[i], w)
	}
	return res
}

func (n *Network) Classify(f tokens.Freqs) string {
	outputSums := make([]*kahan.Summer64, len(n.OutputWeights))
	for i := range outputSums {
		outputSums[i] = kahan.NewSummer64()
		outputSums[i].Add(n.outputBias(i))
	}

	for hiddenIndex, hiddenWeights := range n.HiddenWeights {
		hiddenSum := kahan.NewSummer64()
		for j, token := range n.Tokens {
			hiddenSum.Add(f[token] * hiddenWeights[j])
		}
		hiddenSum.Add(n.hiddenBias(hiddenIndex))

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

func (n *Network) outputBias(outputIdx int) float64 {
	return n.OutputWeights[outputIdx][len(n.Langs)]
}

func (n *Network) hiddenBias(hiddenIdx int) float64 {
	return n.HiddenWeights[hiddenIdx][len(n.Tokens)]
}

func (n *Network) containsNaN() bool {
	for _, wss := range [][][]float64{n.HiddenWeights, n.OutputWeights} {
		for _, ws := range wss {
			for _, w := range ws {
				if math.IsNaN(w) {
					return true
				}
			}
		}
	}
	return false
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
