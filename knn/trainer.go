package knn

import (
	"math/rand"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/whichlang/tokens"
)

// crossCorrelationFrac specifies the fraction of
// samples which are used for cross-validation
// when determining the optimal k-value.
const crossCorrelationFrac = 0.3

func Train(f map[string][]tokens.Freqs) *Classifier {
	seenToks := map[string]bool{}
	sampleCount := 0
	for _, samples := range f {
		for _, sample := range samples {
			for tok := range sample {
				seenToks[tok] = true
			}
			sampleCount++
		}
	}

	toks := make([]string, 0, len(seenToks))
	for tok := range seenToks {
		toks = append(toks, tok)
	}

	samples := make([]Sample, 0, sampleCount)
	for lang, freqSamples := range f {
		for _, freqs := range freqSamples {
			vec := make(linalg.Vector, len(toks))
			for i, token := range toks {
				vec[i] = freqs[token]
			}
			if mag := vec.Dot(vec); mag != 0 {
				vec.Scale(1 / mag)
			}
			samples = append(samples, Sample{
				Language: lang,
				Vector:   vec,
			})
		}
	}

	kValue := optimalKValue(samples)
	return &Classifier{
		Tokens:        toks,
		Samples:       samples,
		NeighborCount: kValue,
	}
}

func optimalKValue(s []Sample) int {
	crossCount := int(crossCorrelationFrac * float64(len(s)))
	if crossCount == 0 {
		return 1
	}
	samples := shuffleSamples(s)
	crossSamples := samples[0:crossCount]
	trainingSamples := samples[crossCount:]

	c := &Classifier{
		Samples: trainingSamples,
	}

	bestK := 1
	bestCorrect := 0
	for k := 1; k <= len(trainingSamples); k++ {
		crossCorrect := 0
		c.NeighborCount = k
		for _, cross := range crossSamples {
			if c.classifyVector(cross.Vector) == cross.Language {
				crossCorrect++
			}
		}
		if crossCorrect > bestCorrect {
			bestK = k
			bestCorrect = crossCorrect
		}
	}

	return bestK
}

func shuffleSamples(s []Sample) []Sample {
	res := make([]Sample, len(s))

	p := rand.Perm(len(s))
	for i, x := range p {
		res[i] = s[x]
	}

	return res
}
