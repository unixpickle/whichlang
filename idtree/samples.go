package idtree

import "github.com/unixpickle/whichlang/tokens"

type linearSample struct {
	freqs []float64
	lang  string
}

func freqsToLinearSamples(toks []string, freqs map[string][]tokens.Freqs) []linearSample {
	var res []linearSample
	for lang, freqsList := range freqs {
		for _, freqs := range freqsList {
			s := linearSample{
				lang:  lang,
				freqs: make([]float64, len(toks)),
			}
			for i, tok := range toks {
				s.freqs[i] = freqs[tok]
			}
			res = append(res, s)
		}
	}
	return res
}

func languageMajority(samples []linearSample) string {
	counts := map[string]int{}
	for _, sample := range samples {
		counts[sample.lang]++
	}

	var maxCount int
	var maxLang string
	for lang, count := range counts {
		if count > maxCount {
			maxCount = count
			maxLang = lang
		}
	}

	return maxLang
}

// A sampleSorter implements sort.Interface
// and facilitates sorting linear samples
// by the frequency of a given token.
type sampleSorter struct {
	samples  []linearSample
	tokenIdx int
}

func (s *sampleSorter) Len() int {
	return len(s.samples)
}

func (s *sampleSorter) Swap(i, j int) {
	s.samples[i], s.samples[j] = s.samples[j], s.samples[i]
}

func (s *sampleSorter) Less(i, j int) bool {
	f1 := s.samples[i].freqs[s.tokenIdx]
	f2 := s.samples[j].freqs[s.tokenIdx]
	return f1 < f2
}
