package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixpickle/whichlang"
	"github.com/unixpickle/whichlang/tokens"
)

const PruneCount = 1
const HelpColumnSize = 10

func main() {
	if len(os.Args) != 4 {
		dieUsage()
	}

	algorithm := os.Args[1]

	trainer := whichlang.Trainers[algorithm]
	if trainer == nil {
		fmt.Fprintln(os.Stderr, "Unknown algorithm:", algorithm)
		dieUsage()
	}

	sampleDir := os.Args[2]
	outputFile := os.Args[3]

	counts, err := tokens.ReadSampleCounts(sampleDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Pruning tokens...")
	counts.Prune(PruneCount)
	fmt.Println("Training...")

	freqs := counts.SampleFreqs()
	classifier := trainer(freqs)
	data := classifier.Encode()

	if err := ioutil.WriteFile(outputFile, data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing file:", err)
		os.Exit(1)
	}
}

func dieUsage() {
	fmt.Fprintln(os.Stderr, "Usage: trainer <algorithm> <sample-dir> <output>\n\n"+
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
