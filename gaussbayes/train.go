package gaussbayes

import (
	"math"

	"github.com/unixpickle/whichlang/tokens"
)

// Train returns a *Classifier by computing statistical
// properties of the sample data.
func Train(freqs map[string][]tokens.Freqs) *Classifier {
	res := &Classifier{LangGaussians: map[string]map[string]Gaussian{}}
	tokens := allTokens(freqs)
	for lang, samples := range freqs {
		gaussians := computeGaussians(samples)
		addMissing(gaussians, tokens)
		res.LangGaussians[lang] = gaussians
	}
	regularizeVariances(res)
	return res
}

func computeGaussians(samples []tokens.Freqs) map[string]Gaussian {
	res := map[string]Gaussian{}

	computeMeans(samples, res)
	computeVariances(samples, res)

	return res
}

func computeMeans(samples []tokens.Freqs, out map[string]Gaussian) {
	for _, sample := range samples {
		for keyword, freq := range sample {
			outGaussian := out[keyword]
			outGaussian.Mean += freq
			out[keyword] = outGaussian
		}
	}
	meanScaler := 1 / float64(len(samples))
	for key, g := range out {
		g.Mean *= meanScaler
		out[key] = g
	}
}

func computeVariances(samples []tokens.Freqs, out map[string]Gaussian) {
	for _, sample := range samples {
		for keyword, freq := range sample {
			outGaussian := out[keyword]
			outGaussian.Variance += math.Pow(freq-outGaussian.Mean, 2)
			out[keyword] = outGaussian
		}
	}
	varianceScaler := 1 / float64(len(samples))
	for key, g := range out {
		g.Variance *= varianceScaler
		out[key] = g
	}
}

func addMissing(m map[string]Gaussian, tokens []string) {
	for _, token := range tokens {
		if _, ok := m[token]; !ok {
			m[token] = Gaussian{}
		}
	}
}

func allTokens(m map[string][]tokens.Freqs) []string {
	res := map[string]bool{}
	for _, x := range m {
		for _, t := range x {
			for w := range t {
				res[w] = true
			}
		}
	}
	resSlice := make([]string, 0, len(res))
	for w := range res {
		resSlice = append(resSlice, w)
	}
	return resSlice
}

// regularizeVariances ensures that no variances are zero.
func regularizeVariances(c *Classifier) {
	var smallestVariance float64
	for _, m := range c.LangGaussians {
		for _, x := range m {
			if smallestVariance == 0 || (x.Variance < smallestVariance && x.Variance > 0) {
				smallestVariance = x.Variance
			}
		}
	}
	for _, m := range c.LangGaussians {
		for word, x := range m {
			if x.Variance == 0 {
				x.Variance = smallestVariance
				m[word] = x
			}
		}
	}
}
