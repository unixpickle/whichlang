package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: trainer <sample dir> <max keywords> <output.json>")
		os.Exit(1)
	}
	freqs := getFrequencies(SampleDir(os.Args[1]))
	fmt.Println("got frequencies for", len(freqs), "languages.")
}

