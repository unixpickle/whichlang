package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/whichlang"
)

func getFrequencies(d SampleDir) map[string][]whichlang.Frequencies {
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

