package svm

import (
	"encoding/json"
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/whichlang/tokens"
)

// BinaryClassifier stores info for the
// binary classifiers used in a Classifier.
type BinaryClassifier struct {
	// SupportVectors stores indices to
	// elements of Classifier.SampleVectors.
	SupportVectors []int

	// Weights are the corresponding weights
	// for each of the support vectors.
	Weights []float64

	Threshold float64
}

// Classifier uses one-against-all SVMs to
// classify source files.
type Classifier struct {
	Keywords []string
	Kernel   *Kernel

	SampleVectors []linalg.Vector

	// Classifiers maps each language to its
	// corresponding one-against-all binary
	// classifier.
	Classifiers map[string]BinaryClassifier
}

func DecodeClassifier(d []byte) (*Classifier, error) {
	var c Classifier
	if err := json.Unmarshal(d, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Classifier) Classify(sample tokens.Freqs) string {
	products := c.sampleProducts(sample)

	var bestLanguage string
	bestClassification := math.Inf(-1)

	for lang, classifier := range c.Classifiers {
		productSum := kahan.NewSummer64()
		for i, vecIdx := range classifier.SupportVectors {
			productSum.Add(products[vecIdx] * classifier.Weights[i])
		}
		productSum.Add(-classifier.Threshold)
		if productSum.Sum() > bestClassification {
			bestClassification = productSum.Sum()
			bestLanguage = lang
		}
	}

	return bestLanguage
}

func (c *Classifier) Encode() []byte {
	res, _ := json.Marshal(c)
	return res
}

func (c *Classifier) Languages() []string {
	res := make([]string, 0, len(c.Classifiers))
	for lang := range c.Classifiers {
		res = append(res, lang)
	}
	return res
}

func (c *Classifier) sampleProducts(sample tokens.Freqs) []float64 {
	vec := c.sampleVector(sample)
	res := make([]float64, len(c.SampleVectors))
	for i, s := range c.SampleVectors {
		res[i] = c.Kernel.Product(s, vec)
	}
	return res
}

func (c *Classifier) sampleVector(sample tokens.Freqs) linalg.Vector {
	vec := make(linalg.Vector, len(c.Keywords))
	for i, keyword := range c.Keywords {
		vec[i] = sample[keyword]
	}
	return vec
}
