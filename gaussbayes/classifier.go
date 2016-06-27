// Package gaussbayes implements naive Bayesian
// classification using the assumption that token
// frequencies follow Gaussian distributions.
package gaussbayes

import (
	"encoding/json"
	"math"

	"github.com/unixpickle/whichlang/tokens"
)

// Gaussian is a Gaussian probability distribution.
type Gaussian struct {
	Mean     float64
	Variance float64
}

// EvalLog evaluates the natural logarithm of the
// density function at a given x value.
func (g Gaussian) EvalLog(x float64) float64 {
	coeff := 1 / math.Sqrt(2*g.Variance*math.Pi)
	exp := -(x - g.Mean) / (2 * g.Variance)
	return math.Log(coeff) + exp
}

type Classifier struct {
	LangGaussians map[string]map[string]Gaussian
}

func DecodeClassifier(d []byte) (*Classifier, error) {
	var c Classifier
	if err := json.Unmarshal(d, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Classifier) Classify(f tokens.Freqs) string {
	var bestLanguage string
	var bestLogProbability float64
	for lang, dists := range c.LangGaussians {
		var probLog float64
		for token, gaussian := range dists {
			probLog += gaussian.EvalLog(f[token])
		}
		if bestLanguage == "" || probLog > bestLogProbability {
			bestLanguage = lang
			bestLogProbability = probLog
		}
	}
	return bestLanguage
}

func (c *Classifier) Encode() []byte {
	res, _ := json.Marshal(c)
	return res
}

func (c *Classifier) Languages() []string {
	var languages []string
	for lang := range c.LangGaussians {
		languages = append(languages, lang)
	}
	return languages
}
