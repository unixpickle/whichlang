package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/whichlang"
)

// RemoveContextualWords removes all words which only occur in one file for a given language.
func RemoveContextualWords(f map[string][]whichlang.Frequencies) {
	for _, samples := range f {
		seenWords := map[string]whichlang.Frequencies{}
		for _, sample := range samples {
			for word := range sample {
				if _, seen := seenWords[word]; !seen {
					seenWords[word] = sample
				} else {
					seenWords[word] = nil
				}
			}
		}
		for seenWord, sample := range seenWords {
			if sample == nil {
				continue
			}
			delete(sample, seenWord)
		}
	}
}

// GetFrequencies processes all the samples in a given sample directory.
func GetFrequencies(d SampleDir) map[string][]whichlang.Frequencies {
	langs, err := d.Languages()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	res := map[string][]whichlang.Frequencies{}
	for _, lang := range langs {
		res[lang] = getLanguageFrequencies(d, lang)
	}
	return res
}

func getLanguageFrequencies(d SampleDir, lang string) []whichlang.Frequencies {
	samples, err := d.SamplesForLanguage(lang)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	res := make([]whichlang.Frequencies, len(samples))
	for i, sample := range samples {
		contents, err := d.ReadSample(lang, sample)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		res[i] = whichlang.ComputeFrequencies(contents)
	}
	return res
}
