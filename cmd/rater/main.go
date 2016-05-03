package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

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

	var totalSamples int
	var totalSuccesses int
	langSuccesses := map[string]int{}

	for lang, langSamples := range samples {
		for _, sample := range langSamples {
			totalSamples++
			if classifier.Classify(sample.Freqs()) == lang {
				totalSuccesses++
				langSuccesses[lang]++
			}
		}
	}

	fmt.Printf("Success rate: %d/%d or %0.2f%%\n", totalSuccesses, totalSamples,
		100*float64(totalSuccesses)/float64(totalSamples))

	sorter := RatingSorter{Ratings: make([]Rating, 0, len(samples))}
	for lang, s := range samples {
		r := Rating{
			Language: lang,
			Correct:  langSuccesses[lang],
			Total:    len(s),
		}
		sorter.Ratings = append(sorter.Ratings, r)
	}
	sort.Sort(sorter)

	for _, rating := range sorter.Ratings {
		fmt.Printf("%s - success rate %d/%d or %0.2f%%\n", rating.Language,
			rating.Correct, rating.Total, 100*rating.Frac())
	}
}

type Rating struct {
	Language string
	Correct  int
	Total    int
}

func (r Rating) Frac() float64 {
	return float64(r.Correct) / float64(r.Total)
}

type RatingSorter struct {
	Ratings []Rating
}

func (r RatingSorter) Len() int {
	return len(r.Ratings)
}

func (r RatingSorter) Less(i, j int) bool {
	return r.Ratings[i].Frac() > r.Ratings[j].Frac()
}

func (r RatingSorter) Swap(i, j int) {
	r.Ratings[i], r.Ratings[j] = r.Ratings[j], r.Ratings[i]
}
