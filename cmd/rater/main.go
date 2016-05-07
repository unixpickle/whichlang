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

	rating := Rate(classifier, samples)

	fmt.Printf("Success rate: %d/%d or %0.2f%%\n", rating.Correct, rating.Total,
		100*rating.Frac())

	nameLength := len(rating.LongestLangName())
	for _, rating := range rating.LangRatings {
		paddedName := rating.Language
		for len(paddedName) < nameLength {
			paddedName = " " + paddedName
		}
		fmt.Printf("%s  success rate %d/%d or %0.2f%%\n", paddedName,
			rating.Correct, rating.Total, 100*rating.Frac())
	}
}
