package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/unixpickle/whichlang"
	"github.com/unixpickle/whichlang/tokens"
)

const HelpColumnSize = 10

func main() {
	if len(os.Args) != 5 {
		dieUsage()
	}

	algorithm := os.Args[1]

	trainer := whichlang.Trainers[algorithm]
	if trainer == nil {
		fmt.Fprintln(os.Stderr, "Unknown algorithm:", algorithm)
		dieUsage()
	}

	ubiquity, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid ubiquity:", ubiquity, "(expected integer)")
		os.Exit(1)
	}

	sampleDir := os.Args[3]
	outputFile := os.Args[4]

	counts, err := tokens.ReadSampleCounts(sampleDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	oldCount := counts.NumTokens()
	fmt.Println("Pruning tokens...")
	counts.Prune(ubiquity)
	newCount := counts.NumTokens()
	fmt.Printf("Pruned %d/%d tokens.\n", (oldCount - newCount), oldCount)

	freqs := counts.SampleFreqs()

	fmt.Println("Training...")
	classifier := trainer(freqs)

	fmt.Println("Saving...")
	data := classifier.Encode()

	if err := ioutil.WriteFile(outputFile, data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing file:", err)
		os.Exit(1)
	}
}

func dieUsage() {
	fmt.Fprintln(os.Stderr, "Usage: trainer <algorithm> <ubiquity> <sample-dir> <output>\n\n"+
		" (ubiquity specifies the number of files in which a\n  keyword should appear.)\n\n"+
		"Available algorithms:")
	for _, name := range whichlang.ClassifierNames {
		spaces := ""
		for i := len(name); i < HelpColumnSize; i++ {
			spaces += " "
		}
		fmt.Fprintln(os.Stderr, " "+name+spaces, whichlang.Descriptions[name])
	}
	fmt.Fprintln(os.Stderr, "")
	os.Exit(1)
}
