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
		fmt.Fprintln(os.Stderr, "Usage: trainer <sample dir> <num neighbors> <output.json>")
		os.Exit(1)
	}

	numNeighbors, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid neighbor count:", os.Args[2])
		os.Exit(1)
	}

	freqs := GetFrequencies(SampleDir(os.Args[1]))
	RemoveContextualWords(freqs, numNeighbors)
	classifier := GenerateClassifier(freqs)
	classifier.NumNeighbors = numNeighbors

	fmt.Println("Generated classifier with", len(classifier.Keywords), "keywords.")

	output, err := json.Marshal(classifier)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode JSON:", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(os.Args[3], output, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing file:", err)
		os.Exit(1)
	}
}
