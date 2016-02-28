package whichlang

import (
	"math"
	"sort"
)

type Sample struct {
	// Vector is a normalized vector of word frequencies, corresponding to words in the classifier's
	// Keywords list
	Vector []float64

	Language string
}

func (s Sample) distance(v []float64) float64 {
	var dot float64
	for i, x := range s.Vector {
		dot += x * v[i]
	}
	return 1 - dot
}

// A Classifier uses k-nearest-neighbors to figure out which language novel samples are most likely
// written in.
type Classifier struct {
	Keywords     []string
	Samples      []Sample
	NumNeighbors int
}

func (c *Classifier) Classify(f Frequencies) string {
	vec := c.computeVector(f)
	pairs := make(langDistPairList, len(c.Samples))
	for i, sample := range c.Samples {
		pairs[i] = langDistPair{
			distance: sample.distance(vec),
			language: sample.Language,
		}
	}
	sort.Sort(pairs)

	languageWeights := map[string]float64{}
	for i := 0; i < len(pairs) && i < c.NumNeighbors; i++ {
		pair := pairs[i]
		dist := pair.distance
		if dist == 0 {
			return pair.language
		}
		languageWeights[pair.language] += 1 / dist
	}

	var bestWeight float64
	var bestLang string

	for lang, weight := range languageWeights {
		if weight > bestWeight || bestLang == "" {
			bestLang = lang
			bestWeight = weight
		}
	}

	return bestLang
}

func (c *Classifier) computeVector(f Frequencies) []float64 {
	res := make([]float64, len(c.Keywords))
	var mag2 float64
	for i, word := range c.Keywords {
		val := f[word]
		mag2 += val * val
		res[i] = val
	}
	normalizer := math.Sqrt(mag2)
	for i, x := range res {
		res[i] = x / normalizer
	}
	return res
}

type langDistPair struct {
	distance float64
	language string
}

type langDistPairList []langDistPair

func (l langDistPairList) Len() int {
	return len(l)
}

func (l langDistPairList) Less(i, j int) bool {
	return l[i].distance < l[j].distance
}

func (l langDistPairList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
