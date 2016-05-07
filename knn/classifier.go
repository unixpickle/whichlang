package knn

import (
	"encoding/json"
	"math"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/whichlang/tokens"
)

type Sample struct {
	Language string
	Vector   linalg.Vector
}

type Classifier struct {
	Tokens  []string
	Samples []Sample

	NeighborCount int
}

func DecodeClassifier(d []byte) (*Classifier, error) {
	var res Classifier
	if err := json.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Classifier) Classify(f tokens.Freqs) string {
	vec := make(linalg.Vector, len(c.Tokens))
	for i, keyw := range c.Tokens {
		vec[i] = f[keyw]
	}

	vecMag := vec.Dot(vec)
	if vecMag == 0 {
		return c.Samples[0].Language
	}
	vec.Scale(1 / math.Sqrt(vecMag))

	return c.classifyVector(vec)
}

func (c *Classifier) Encode() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *Classifier) Languages() []string {
	seenLangs := map[string]bool{}
	for _, sample := range c.Samples {
		seenLangs[sample.Language] = true
	}
	res := make([]string, 0, len(seenLangs))
	for lang := range seenLangs {
		res = append(res, lang)
	}
	return res
}

func (c *Classifier) classifyVector(vec linalg.Vector) string {
	matches := make([]match, 0, c.NeighborCount)
	for _, sample := range c.Samples {
		correlation := sample.Vector.Dot(vec)
		insertIdx := matchInsertionIndex(matches, correlation)
		if insertIdx >= c.NeighborCount {
			continue
		}
		if len(matches) < c.NeighborCount {
			matches = append(matches, match{})
		}
		copy(matches[insertIdx+1:], matches[insertIdx:])
		matches[insertIdx] = match{
			Language:    sample.Language,
			Correlation: correlation,
		}
	}

	return dominantClassification(matches)
}

func dominantClassification(matches []match) string {
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
