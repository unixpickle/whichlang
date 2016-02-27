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
		fmt.Fprintln(os.Stderr, "Usage: classify <classifiers.json> <file>")
		os.Exit(1)
	}

	classifierData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var classifiers map[string]*whichlang.Classifier
	if err := json.Unmarshal(classifierData, &classifiers); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to decode classifiers:", err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	freqs := whichlang.ComputeFrequencies(string(contents))
	matched := false
	for lang, classifier := range classifiers {
		if classifier.Classify(freqs) {
			fmt.Println("Match for", lang)
			matched = true
		}
	}
	if !matched {
		fmt.Println("No matches found.")
	}
}
