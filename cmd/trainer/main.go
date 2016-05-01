package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: trainer <sample dir> <output.json>")
		os.Exit(1)
	}

	freqs := GetFrequencies(SampleDir(os.Args[1]))
	NormalizeFrequencies(freqs)
	RemoveContextualWords(freqs, 5)
	classifier := GenerateClassifier(freqs)

	fmt.Println("Generated classifier with", classifier.LeafCount(), "leaves")

	output, err := json.Marshal(classifier)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode JSON:", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(os.Args[2], output, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing file:", err)
		os.Exit(1)
	}
}
