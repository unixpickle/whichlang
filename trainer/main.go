package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: trainer <sample dir> <max keywords> <output.json>")
		os.Exit(1)
	}

	maxKeywords, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid max keywords:", os.Args[2])
		os.Exit(1)
	}

	fmt.Println("Computing frequencies...")
	freqs := GetFrequencies(SampleDir(os.Args[1]))
	NormalizeFrequencies(freqs)
	RemoveContextualWords(freqs)
	samples := NewSamples(freqs)

	fmt.Println("Generating classifiers (dimensionality is " + strconv.Itoa(len(samples.Words)) +
		")...")
	classifiers := GenerateClassifiers(samples, maxKeywords)

	output, err := json.Marshal(classifiers)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode JSON:", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(os.Args[3], output, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing file:", err)
		os.Exit(1)
	}
}
