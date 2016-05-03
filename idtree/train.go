package idtree

import (
	"math"
	"runtime"
	"sort"

	"github.com/unixpickle/whichlang/tokens"
)

type splitInfo struct {
	TokenIdx  int
	Threshold float64
	Entropy   float64
}

// Train returns a *Classifier which is the
// result of running ID3 on a set of training
// samples.
func Train(freqs map[string][]tokens.Freqs) *Classifier {
	toks := allTokens(freqs)
	samples := freqsToLinearSamples(toks, freqs)
	return generateClassifier(toks, samples)
}

func allTokens(freqs map[string][]tokens.Freqs) []string {
	words := make([]string, 0)
	seenWords := map[string]bool{}
	for _, freqsList := range freqs {
		for _, freqs := range freqsList {
			for word := range freqs {
				if !seenWords[word] {
					seenWords[word] = true
					words = append(words, word)
				}
			}
		}
	}
	return words
}

// generateClassifier generates a classifier
// for the given set of samples.
func generateClassifier(toks []string, s []linearSample) *Classifier {
	tokIdx, thresh := bestDecision(s)
	if tokIdx == -1 {
		lang := languageMajority(s)
		return &Classifier{
			LeafClassification: &lang,
		}
	}
	res := &Classifier{
		Keyword:   toks[tokIdx],
		Threshold: thresh,
	}
	f, t := splitData(s, tokIdx, thresh)
	res.FalseBranch = generateClassifier(toks, f)
	res.TrueBranch = generateClassifier(toks, t)
	return res
}

func splitData(s []linearSample, tokIdx int, thresh float64) (f, t []linearSample) {
	f = make([]linearSample, 0, len(s))
	t = make([]linearSample, 0, len(s))

	for _, sample := range s {
		if sample.freqs[tokIdx] > thresh {
			t = append(t, sample)
		} else {
			f = append(f, sample)
		}
	}

	return
}

// bestDecision returns the token and threshold
// which split the samples optimally (by the
// criterion of entropy).
// If no split exists, this returns (-1, -1).
func bestDecision(s []linearSample) (tokIdx int, thresh float64) {
	maxProcs := runtime.GOMAXPROCS(0)

	toksPerGo := len(toks) / maxProcs
	splitChan := make(chan *splitInfo, maxProcs)
	for i := 0; i < maxProcs; i++ {
		tokCount := toksPerGo
		tokStart := toksPerGo * i

		// The last set might need to be slightly larger
		// due to division truncation.
		if i == maxProcs-1 {
			tokCount = len(toks) - tokStart
		}

		go bestNodeSubset(tokStart, tokCount, s, splitChan)
	}

	var best *splitInfo
	for i := 0; i < maxProcs; i++ {
		res := <-splitChan
		if res == nil {
			continue
		}
		if best == nil || res.Entropy < best.Entropy {
			best = res
		}
	}

	if best == nil {
		return -1, -1
	}

	return best.TokenIdx, best.Threshold
}

func bestNodeSubset(startIdx, count int, s []linearSample, res chan<- *splitInfo) {
	bestThresh := -1.0
	var bestEntropy float64
	var bestIdx int
	for i := 0; i < count; i++ {
		idx := startIdx + i
		thresh, entropy := bestSplit(s, idx)
		if thresh < 0 {
			continue
		} else if bestThresh < 0 || entropy < bestEntropy {
			bestEntropy = entropy
			bestThresh = thresh
			bestIdx = idx
		}
	}
	if bestThresh == -1 {
		res <- nil
	} else {
		res <- &splitInfo{bestIdx, bestThresh, bestEntropy}
	}
}

// bestSplit finds the ideal threshold for splitting
// samples by a given token (specified by an index).
// This returns the threshold and the resulting entropy.
// The threshold will be -1 if no split is useful.
func bestSplit(s []linearSample, tokenIdx int) (thres float64, entrop float64) {
	sortedArray := make([]linearSample, len(s))
	copy(sortedArray, s)
	sorter := &sampleSorter{sortedGroup, tokenIdx}
	sort.Sort(sorter)

	lowerDistribution := map[string]int{}
	upperDistribution := map[string]int{}

	for _, sample := range sortedArray {
		upperDistribution[sample.lang]++
	}

	if len(upperDistribution) == 1 {
		// Can't split homogeneous data effectively.
		return -1, -1
	}

	thresh = -1
	entrop = -1

	if len(sortedArray) == 0 {
		return
	}

	lastFreq := sortedArray[0].freqs[tokenIdx]
	for i := 1; i < len(s); i++ {
		upperDistribution[s[i-1].freqs[tokenIdx]]--
		lowerDistribution[s[i-1].freqs[tokenIdx]]++

		freq := s[i].freqs[tokenIdx]
		if freq == lastFreq {
			continue
		}

		upperFrac := float64(len(s)-i) / float64(len(s))
		lowerFrac := float64(i) / float64(len(s))
		disorder := upperFrac*distributionEntropy(upperDistribution) +
			lowerFrac*distributionEntropy(lowerDistribution)
		if disorder < entrop || thresh == -1 {
			entrop = disorder
			thresh = (lastFreq + freq) / 2
		}

		lastFreq = freq
	}

	return
}

func distributionEntropy(dist map[string]int) float64 {
	var res float64
	var totalCount int
	for _, count := range dist {
		totalCount += count
	}
	for _, count := range dist {
		fraction := float64(count) / float64(totalCount)
		if fraction != 0 {
			res -= math.Log(fraction) * fraction
		}
	}
	return res
}
