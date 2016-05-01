package neuralnet

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
)

// A gradientCalc can compute gradients of the
// error function 0.5*||Actual - Expected||^2
// for a neural network on a given input.
//
// A gradientCalc demands a lot of scratch
// memory, so it is a good idea to create one
// gradientCalc and then reuse it over and over.
type gradientCalc struct {
	n *Network

	hiddenOutputs []float64
	outputs       []float64
	inputs        []float64
	expectedOut   []float64

	hiddenOutPartials []float64

	OutputPartials [][]float64
	HiddenPartials [][]float64
}

func newGradientCalc(n *Network) *gradientCalc {
	res := &gradientCalc{
		n:                 n,
		hiddenOutputs:     make([]float64, len(n.HiddenWeights)),
		outputs:           make([]float64, len(n.OutputWeights)),
		expectedOut:       make([]float64, len(n.OutputWeights)),
		hiddenOutPartials: make([]float64, len(n.HiddenWeights)),
		OutputPartials:    make([][]float64, len(n.OutputWeights)),
		HiddenPartials:    make([][]float64, len(n.HiddenWeights)),
	}

	for i := range res.OutputPartials {
		res.OutputPartials[i] = make([]float64, len(res.hiddenOutputs)+1)
	}
	for i := range res.HiddenPartials {
		res.HiddenPartials[i] = make([]float64, len(n.Tokens)+1)
	}

	return res
}

func (g *gradientCalc) Compute(inputs []float64, langIdx int) {
	g.inputs = inputs
	for j := range g.expectedOut {
		if j == langIdx {
			g.expectedOut[j] = 1
		} else {
			g.expectedOut[j] = 0
		}
	}

	g.computeOutputs()
	g.computeGradients()
}

// Normalize normalizes the gradient using
// the Euclidean norm.
func (g *gradientCalc) Normalize() {
	sum := kahan.NewSummer64()
	for _, xss := range [][][]float64{g.HiddenPartials, g.OutputPartials} {
		for _, xs := range xss {
			for _, x := range xs {
				sum.Add(x * x)
			}
		}
	}

	normalizer := 1.0 / math.Sqrt(sum.Sum())

	for _, xss := range [][][]float64{g.HiddenPartials, g.OutputPartials} {
		for _, xs := range xss {
			for i, x := range xs {
				xs[i] = x * normalizer
			}
		}
	}
}

func (g *gradientCalc) computeOutputs() {
	outputSums := make([]*kahan.Summer64, len(g.outputs))

	for i := range outputSums {
		outputSums[i] = kahan.NewSummer64()
		outputSums[i].Add(g.n.outputBias(i))
	}

	for hiddenIndex, hiddenWeights := range g.n.HiddenWeights {
		hiddenSum := kahan.NewSummer64()
		for j, input := range g.inputs {
			hiddenSum.Add(input * hiddenWeights[j])
		}
		hiddenSum.Add(g.n.hiddenBias(hiddenIndex))

		hiddenOut := sigmoid(hiddenSum.Sum())
		g.hiddenOutputs[hiddenIndex] = hiddenOut
		for j, outSum := range outputSums {
			weight := g.n.OutputWeights[j][hiddenIndex]
			outSum.Add(weight * hiddenOut)
		}
	}

	for i, sum := range outputSums {
		g.outputs[i] = sigmoid(sum.Sum())
	}
}

func (g *gradientCalc) computeGradients() {
	for i, output := range g.outputs {
		gradient := g.OutputPartials[i]
		diff := output - g.expectedOut[i]
		sumPartial := (1 - output) * output * diff
		for j, input := range g.hiddenOutputs {
			gradient[j] = input * sumPartial
			g.hiddenOutPartials[j] = g.n.OutputWeights[i][j] * sumPartial
		}
		gradient[len(g.hiddenOutputs)] = sumPartial
	}
	for i, output := range g.hiddenOutputs {
		gradient := g.HiddenPartials[i]
		sumPartial := (1 - output) * output * g.hiddenOutPartials[i]
		for j, input := range g.inputs {
			gradient[j] = input * sumPartial
		}
		gradient[len(g.inputs)] = sumPartial
	}
}
