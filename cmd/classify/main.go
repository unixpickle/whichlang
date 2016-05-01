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
		fmt.Fprintln(os.Stderr, "Usage: classify <algorithm> <classifier-file> <file>")
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

	contents, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	counts := tokens.CountTokens(string(contents))
	freqs := counts.Freqs()
	language := classifier.Classify(freqs)
	fmt.Println("Classification:", language)
}
