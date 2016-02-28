package main

import (
	"math"

	"github.com/unixpickle/weakai/svm"
	"github.com/unixpickle/whichlang"
)

// Samples stores the vector representations of whichlang.Frequencies, since the vector
// representations can be handed off to the SVM solver.
type Samples struct {
	Samples map[string][]svm.Sample
	Words   []string
}

func NewSamples(frequencies map[string][]whichlang.Frequencies) Samples {
	words := map[string]bool{}
	for _, list := range frequencies {
		for _, freqs := range list {
			for word := range freqs {
				words[word] = true
			}
		}
	}

	wordIndices := map[string]int{}
	wordList := make([]string, 0, len(words))
	for word := range words {
		wordIndices[word] = len(wordList)
		wordList = append(wordList, word)
	}

	res := Samples{
		Samples: map[string][]svm.Sample{},
		Words:   wordList,
	}

	sampleIndex := 1
	for lang, list := range frequencies {
		sampleList := make([]svm.Sample, len(list))
		for i, freqs := range list {
			sampleVec := make([]float64, len(wordList))
			var magSquared float64
			for word, freq := range freqs {
				sampleVec[wordIndices[word]] = freq
				magSquared += freq * freq
			}
			scaler := 1 / math.Sqrt(magSquared)
			for i, x := range sampleVec {
				sampleVec[i] = x * scaler
			}
			sampleList[i] = svm.Sample{V: sampleVec, UserInfo: sampleIndex}
			sampleIndex++
		}
		res.Samples[lang] = sampleList
	}

	return res
}
