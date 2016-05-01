package neuralnet

import (
	"math/rand"

	"github.com/unixpickle/whichlang/tokens"
)

const InitialIterationCount = 100
const DefaultMaxIterations = 6400

// HiddenLayerScale specifies how much larger
// the hidden layer is than the output layer.
const HiddenLayerScale = 2.0

var StepSizes = []float64{1e-8, 1e-7, 1e-6, 1e-5, 1e-4, 1e-3, 1e-2, 1e-1, 1, 1e1, 1e2}

func Train(data map[string][]tokens.Freqs) *Network {
	ds := NewDataSet(data)

	var best *Network
	var bestCrossScore float64
	var bestTrainScore float64

	for _, stepSize := range StepSizes {
		t := NewTrainer(ds, stepSize)
		t.Train(DefaultMaxIterations)

		n := t.Network()
		crossScore := ds.CrossScore(n)
		trainScore := ds.TrainingScore(n)
		if (crossScore == bestCrossScore && trainScore >= bestTrainScore) ||
			best == nil || (crossScore > bestCrossScore) {
			bestCrossScore = crossScore
			bestTrainScore = trainScore
			best = n
		}
	}

	return best
}

type Trainer struct {
	n *Network
	d *DataSet
	g *gradientCalc

	stepSize float64
}

func NewTrainer(d *DataSet, stepSize float64) *Trainer {
	hiddenCount := int(HiddenLayerScale * float64(len(d.TrainingSamples)))
	n := &Network{
		Tokens:        d.Tokens(),
		Langs:         d.Langs(),
		HiddenWeights: make([][]float64, hiddenCount),
		OutputWeights: make([][]float64, len(d.TrainingSamples)),
	}
	for i := range n.OutputWeights {
		n.OutputWeights[i] = make([]float64, hiddenCount+1)
		for j := range n.OutputWeights {
			n.OutputWeights[i][j] = rand.Float64()
		}
	}
	for i := range n.HiddenWeights {
		n.HiddenWeights[i] = make([]float64, len(n.Tokens)+1)
		for j := range n.HiddenWeights[i] {
			n.HiddenWeights[i][j] = rand.Float64()
		}
	}
	return &Trainer{
		n:        n,
		d:        d,
		g:        newGradientCalc(n),
		stepSize: stepSize,
	}
}

func (t *Trainer) Train(maxIters int) {
	iters := InitialIterationCount
	if iters > maxIters {
		iters = maxIters
	}
	for i := 0; i < iters; i++ {
		t.runAllSamples()
	}
	if iters == maxIters {
		return
	}

	// Use cross-validation to find the best
	// number of iterations.
	crossScore := t.d.CrossScore(t.n)
	trainScore := t.d.TrainingScore(t.n)
	for {
		nextAmount := iters
		if nextAmount+iters > maxIters {
			nextAmount = maxIters - iters
		}
		for i := 0; i < nextAmount; i++ {
			t.runAllSamples()
		}
		iters += nextAmount

		if iters == maxIters {
			break
		}

		newCrossScore := t.d.CrossScore(t.n)
		newTrainScore := t.d.TrainingScore(t.n)
		if (newCrossScore == crossScore && newTrainScore == trainScore) ||
			newCrossScore < crossScore {
			break
		}

		crossScore = newCrossScore
		trainScore = newTrainScore
	}
}

func (t *Trainer) Network() *Network {
	return t.n
}

func (t *Trainer) runAllSamples() {
	for i, lang := range t.n.Langs {
		samples := t.d.TrainingSamples[lang]
		for _, sample := range samples {
			t.descendSample(sample, i)
		}
	}
}

// descendSample performs gradient descent to
// reduce the output error for a given sample.
func (t *Trainer) descendSample(f tokens.Freqs, langIdx int) {
	t.g.Compute(f, langIdx)
	for i, partials := range t.g.HiddenPartials {
		for j, partial := range partials {
			t.n.HiddenWeights[i][j] -= partial * t.stepSize
		}
	}
	for i, partials := range t.g.OutputPartials {
		for j, partial := range partials {
			t.n.OutputWeights[i][j] -= partial * t.stepSize
		}
	}
}
