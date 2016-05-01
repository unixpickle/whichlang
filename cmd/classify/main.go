package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixpickle/whichlang"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: classify <classifier.json> <file>")
		os.Exit(1)
	}

	classifierData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var classifier whichlang.Classifier
	if err := json.Unmarshal(classifierData, &classifier); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to decode classifier:", err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	freqs := whichlang.ComputeFrequencies(string(contents))
	language := classifier.Classify(freqs)
	fmt.Println("Code appears to be", language)
}
