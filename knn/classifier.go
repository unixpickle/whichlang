package knn

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/whichlang/tokens"
)

type Classifier struct {
	Keywords []string
	Samples  map[string][]linalg.Vector

	NeighborCount int
}

func (c *Classifier) Classify(f tokens.Freqs) string {
	vec := make(linalg.Vector, len(c.Keywords))
	for i, keyw := range c.Keywords {
		vec[i] = f[keyw]
	}

	vecMag := vec.Dot(vec)
	if vecMag == 0 {
		for lang := range c.Samples {
			return lang
		}
	}
	vec.Scale(1 / math.Sqrt(vecMag))

	matches := make([]match, 0, c.NeighborCount)
	for lang, samples := range c.Samples {
		for _, sample := range samples {
			correlation := sample.Dot(vec)
			insertIdx := matchInsertionIndex(matches, correlation)
			if insertIdx >= c.NeighborCount {
				continue
			}
			if len(matches) < c.NeighborCount {
				matches = append(matches, match{})
			}
			copy(matches[insertIdx+1:], matches[insertIdx:])
			matches[insertIdx] = match{
				Language:    lang,
				Correlation: correlation,
			}
		}
	}

	scores := map[string]float64{}
	for _, m := range matches {
		scores[m.Language] += m.Correlation
	}

	var bestLang string
	bestScore := math.Inf(-1)
	for lang, score := range scores {
		if score > bestScore {
			bestScore = score
			bestLang = lang
		}
	}

	return bestLang
}

type match struct {
	Language    string
	Correlation float64
}

func matchInsertionIndex(m []match, corr float64) int {
	for i, x := range m {
		if x.Correlation < corr {
			return i
		}
	}
	return len(m)
}
