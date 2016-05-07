package main

import (
	"runtime"
	"sync"

	"github.com/unixpickle/whichlang"
	"github.com/unixpickle/whichlang/tokens"
)

type Challenge struct {
	Language string
	Sample   tokens.Freqs
}

type Result struct {
	Language string
	Correct  bool
}

func Rate(c whichlang.Classifier, s map[string][]tokens.Counts) *OverallRating {
	var wg sync.WaitGroup
	challengeChan := make(chan Challenge, 0)
	resultChan := make(chan Result, 0)

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for challenge := range challengeChan {
				correct := (c.Classify(challenge.Sample) == challenge.Language)
				resultChan <- Result{
					Language: challenge.Language,
					Correct:  correct,
				}
			}
		}()
	}

	go func() {
		for lang, langSamples := range s {
			for _, sample := range langSamples {
				challengeChan <- Challenge{lang, sample.Freqs()}
			}
		}
		close(challengeChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var total, successes int
	langSuccesses := map[string]int{}
	langTotals := map[string]int{}

	for result := range resultChan {
		total++
		langTotals[result.Language]++
		if result.Correct {
			successes++
			langSuccesses[result.Language]++
		}
	}

	langRatings := makeLangRatings(langSuccesses, langTotals)
	return NewOverallRating(successes, total, langRatings)
}

func makeLangRatings(succ, total map[string]int) []*LangRating {
	res := make([]*LangRating, 0, len(total))
	for lang, totalCount := range total {
		lr := &LangRating{
			Rating: Rating{
				Correct: succ[lang],
				Total:   totalCount,
			},
			Language: lang,
		}
		res = append(res, lr)
	}
	return res
}
