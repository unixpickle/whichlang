package neuralnet

import (
	"log"
	"math/rand"

	"github.com/unixpickle/whichlang/tokens"
)

const InitialIterationCount = 200

// HiddenLayerScale specifies how much larger
// the hidden layer is than the output layer.
const HiddenLayerScale = 2.0

func Train(data map[string][]tokens.Freqs) *Network {
	ds := NewDataSet(data)

	var best *Network
	var bestCrossScore float64
	var bestTrainScore float64

	verbose := verboseFlag()

	for _, stepSize := range stepSizes() {
		if verbose {
			log.Printf("trying step size %f", stepSize)
		}

		t := NewTrainer(ds, stepSize, verbose)
		t.Train(maxIterations())

		n := t.Network()
		if n.containsNaN() {
			if verbose {
				log.Printf("got NaN for step size %f", stepSize)
			}
			continue
		}
		crossScore := ds.CrossScore(n)
		trainScore := ds.TrainingScore(n)
		if verbose {
			log.Printf("stepSize=%f crossScore=%f trainScore=%f", stepSize,
				crossScore, trainScore)
		}
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
	verbose  bool
}

func NewTrainer(d *DataSet, stepSize float64, verbose bool) *Trainer {
	hiddenCount := int(HiddenLayerScale * float64(len(d.TrainingSamples)))
	n := &Network{
		Tokens:        d.Tokens(),
		Langs:         d.Langs(),
		HiddenWeights: make([][]float64, hiddenCount),
		OutputWeights: make([][]float64, len(d.TrainingSamples)),
		InputShift:    -d.MeanFrequency,
		InputScale:    1 / d.FrequencyStddev,
	}
	for i := range n.OutputWeights {
		n.OutputWeights[i] = make([]float64, hiddenCount+1)
		for j := range n.OutputWeights {
			n.OutputWeights[i][j] = rand.Float64()*2 - 1
		}
	}
	for i := range n.HiddenWeights {
		n.HiddenWeights[i] = make([]float64, len(n.Tokens)+1)
		for j := range n.HiddenWeights[i] {
			n.HiddenWeights[i][j] = rand.Float64()*2 - 1
		}
	}
	return &Trainer{
		n:        n,
		d:        d,
		g:        newGradientCalc(n),
		stepSize: stepSize,
		verbose:  verbose,
	}
}

func (t *Trainer) Train(maxIters int) {
	iters := InitialIterationCount
	if iters > maxIters {
		iters = maxIters
	}
	for i := 0; i < iters; i++ {
		if verboseStepsFlag() {
			log.Printf("done %d iterations, cross=%f training=%f",
				i, t.d.CrossScore(t.n), t.d.TrainingScore(t.n))
		}
		t.runAllSamples()
	}
	if iters == maxIters {
		return
	}

	if t.n.containsNaN() {
		return
	}

	// Use cross-validation to find the best
	// number of iterations.
	crossScore := t.d.CrossScore(t.n)
	trainScore := t.d.TrainingScore(t.n)
	lastNet := t.n.Copy()

	for {
		if t.verbose {
			log.Printf("current scores: cross=%f train=%f iters=%d",
				crossScore, trainScore, iters)
		}

		nextAmount := iters
		if nextAmount+iters > maxIters {
			nextAmount = maxIters - iters
		}
		for i := 0; i < nextAmount; i++ {
			if verboseStepsFlag() {
				log.Printf("done %d iterations, cross=%f training=%f",
					i+iters, t.d.CrossScore(t.n), t.d.TrainingScore(t.n))
			}
			t.runAllSamples()
			if t.n.containsNaN() {
				break
			}
		}
		iters += nextAmount

		if t.n.containsNaN() {
			t.n = lastNet
			break
		}

		newCrossScore := t.d.CrossScore(t.n)
		newTrainScore := t.d.TrainingScore(t.n)
		if (newCrossScore == crossScore && newTrainScore == trainScore) ||
			newCrossScore < crossScore {
			t.n = lastNet
			return
		}

		crossScore = newCrossScore
		trainScore = newTrainScore

		if iters == maxIters {
			return
		}
		lastNet = t.n.Copy()
	}
}

func (t *Trainer) Network() *Network {
	return t.n
}

func (t *Trainer) runAllSamples() {
	for i, lang := range t.n.Langs {
		samples := t.d.NormalTrainingSamples[lang]
		for _, sample := range samples {
			t.descendSample(sample, i)
		}
	}
}

// descendSample performs gradient descent to
// reduce the output error for a given sample.
func (t *Trainer) descendSample(inputs []float64, langIdx int) {
	t.g.Compute(inputs, langIdx)
	t.g.Normalize()

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
