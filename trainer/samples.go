package main

import (
	"github.com/unixpickle/weakai/svm"
	"github.com/unixpickle/whichlang"
)

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

	for lang, list := range frequencies {
		sampleList := make([]svm.Sample, len(list))
		for i, freqs := range list {
			sample := make(svm.Sample, len(wordList))
			for word, freq := range freqs {
				sample[wordIndices[word]] = freq
			}
			sampleList[i] = sample
		}
		res.Samples[lang] = sampleList
	}

	return res
}
