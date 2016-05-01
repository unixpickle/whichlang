package neuralnet

import (
	"math/rand"

	"github.com/unixpickle/whichlang/tokens"
)

const ValidationFraction = 0.3

// A DataSet is a set of data split into training
// samples and validation samples.
type DataSet struct {
	ValidationSamples map[string][]tokens.Freqs
	TrainingSamples   map[string][]tokens.Freqs
}

// NewDataSet creates a DataSet by randomly
// partitioning some data samples into
// validation and training samples.
func NewDataSet(samples map[string][]tokens.Freqs) *DataSet {
	res := &DataSet{
		ValidationSamples: map[string][]tokens.Freqs{},
		TrainingSamples:   map[string][]tokens.Freqs{},
	}
	for lang, langSamples := range samples {
		shuffled := make([]tokens.Freqs, len(langSamples))
		perm := rand.Perm(len(shuffled))
		for i, x := range perm {
			shuffled[i] = langSamples[x]
		}

		numValid := int(float64(len(langSamples)) * ValidationFraction)
		res.ValidationSamples[lang] = shuffled[:numValid]
		res.TrainingSamples[lang] = shuffled[numValid:]
	}
	return res
}

// CrossScore returns the fraction of withheld
// samples the Network worked for.
func (c *DataSet) CrossScore(n *Network) float64 {
	return scoreNetwork(n, c.ValidationSamples)
}

// TrainingScore returns the fraction of
// training samples the Network worked for.
func (c *DataSet) TrainingScore(n *Network) float64 {
	return scoreNetwork(n, c.TrainingSamples)
}

// Tokens returns all of the tokens from all
// of the training samples.
func (c *DataSet) Tokens() []string {
	toks := map[string]bool{}
	for _, samples := range c.TrainingSamples {
		for _, sample := range samples {
			for tok := range sample {
				toks[tok] = true
			}
		}
	}

	res := make([]string, 0, len(toks))
	for tok := range toks {
		res = append(res, tok)
	}
	return res
}

// Langs returns all of the languages represented
// by the training samples.
func (c *DataSet) Langs() []string {
	res := make([]string, 0, len(c.TrainingSamples))
	for lang := range c.TrainingSamples {
		res = append(res, lang)
	}
	return res
}

func scoreNetwork(n *Network, samples map[string][]tokens.Freqs) float64 {
	var totalRight int
	var total int
	for lang, langSamples := range samples {
		for _, sample := range langSamples {
			if n.Classify(sample) == lang {
				totalRight++
			}
			total++
		}
	}
	return float64(totalRight) / float64(total)
}
