package main

import (
	"math"

	"github.com/unixpickle/whichlang"
)

func GenerateClassifier(freqs map[string][]whichlang.Frequencies) *whichlang.Classifier {
	words := map[string]bool{}
	sampleCount := 0
	for _, list := range freqs {
		for _, wordMap := range list {
			for word := range wordMap {
				words[word] = true
				sampleCount++
			}
		}
	}

	res := &whichlang.Classifier{
		Keywords:     make([]string, 0, len(words)),
		Samples:      make([]whichlang.Sample, 0, sampleCount),
		NumNeighbors: 1,
	}

	for word := range words {
		res.Keywords = append(res.Keywords, word)
	}

	for lang, list := range freqs {
		for _, wordMap := range list {
			vec := make([]float64, len(words))
			var mag float64
			for i, word := range res.Keywords {
				vec[i] += wordMap[word]
				mag += vec[i] * vec[i]
			}
			scaler := 1 / math.Sqrt(mag)
			for i, x := range vec {
				vec[i] = x * scaler
			}
			res.Samples = append(res.Samples, whichlang.Sample{
				Language: lang,
				Vector:   vec,
			})
		}
	}

	return res
}
