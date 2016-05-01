package main

import "github.com/unixpickle/whichlang"

// ComputeConfidences sets the LeafConfidence value on each leaf.
func ComputeConfidences(c *whichlang.Classifier, f map[string][]whichlang.Frequencies) {
	for lang, list := range f {
		computeConfidenceForLang(list, len(list), lang, c.TreeRoot)
	}
}

func computeConfidenceForLang(vecs []whichlang.Frequencies, totalLen int, lang string,
	node *whichlang.ClassifierNode) {
	if node.Leaf {
		if node.LeafClassification == lang {
			node.LeafConfidence = float64(len(vecs)) / float64(totalLen)
		}
		return
	}
	t, f := splitOnNode(vecs, node)
	computeConfidenceForLang(t, totalLen, lang, node.TrueBranch)
	computeConfidenceForLang(f, totalLen, lang, node.FalseBranch)
}

// CenterThresholds makes sure that every threshold lies directly between the maximum and minimum
// possible thresholds that would still split the samples the exact same way.
func CenterThresholds(c *whichlang.Classifier, f map[string][]whichlang.Frequencies) {
	vecs := []whichlang.Frequencies{}
	for _, list := range f {
		for _, wordMap := range list {
			vecs = append(vecs, wordMap)
		}
	}
	centerThresholdsForNode(vecs, c.TreeRoot)
}

func centerThresholdsForNode(vecs []whichlang.Frequencies, node *whichlang.ClassifierNode) {
	if node.Leaf {
		return
	}
	var lowerSide float64
	var upperSide float64
	for i, vec := range vecs {
		if vec[node.Keyword] <= node.Threshold {
			if i == 0 || vec[node.Keyword] > lowerSide {
				lowerSide = node.Threshold
			}
		} else {
			if i == 0 || vec[node.Keyword] < upperSide {
				upperSide = node.Threshold
			}
		}
	}
	node.Threshold = (lowerSide + upperSide) / 2

	t, f := splitOnNode(vecs, node)
	centerThresholdsForNode(t, node.TrueBranch)
	centerThresholdsForNode(f, node.FalseBranch)
}

func splitOnNode(vecs []whichlang.Frequencies,
	node *whichlang.ClassifierNode) (t, f []whichlang.Frequencies) {
	t = make([]whichlang.Frequencies, 0, len(vecs))
	f = make([]whichlang.Frequencies, 0, len(vecs))
	for _, v := range vecs {
		if v[node.Keyword] > node.Threshold {
			t = append(t, v)
		} else {
			f = append(f, v)
		}
	}
	return
}
