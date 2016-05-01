package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixpickle/whichlang"
	"github.com/unixpickle/whichlang/tokens"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: rater <algorithm> <classifier-file> <dir>")
		os.Exit(1)
	}

	decoder := whichlang.Decoders[os.Args[1]]
	if decoder == nil {
		fmt.Fprintln(os.Stderr, "Unknown algorithm:", os.Args[1])
		os.Exit(1)
	}

	classifierData, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	classifier, err := decoder(classifierData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to decode classifier:", err)
		os.Exit(1)
	}

	samples, err := tokens.ReadSampleCounts(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read samples:", err)
		os.Exit(1)
	}

	var totalSamples int
	var totalSuccesses int
	langSuccesses := map[string]int{}

	for lang, langSamples := range samples {
		for _, sample := range langSamples {
			totalSamples++
			if classifier.Classify(sample.Freqs()) == lang {
				totalSuccesses++
				langSuccesses[lang]++
			}
		}
	}

	fmt.Printf("Success rate: %d/%d or %0.2f%%\n", totalSuccesses, totalSamples,
		100*float64(totalSuccesses)/float64(totalSamples))

	for lang, s := range samples {
		fmt.Printf("%s - success rate %d/%d or %0.2f%%\n", lang, langSuccesses[lang],
			len(s), 100*float64(langSuccesses[lang])/float64(len(s)))
	}
}
