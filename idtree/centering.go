package idtree

import "github.com/unixpickle/whichlang/tokens"

// centerThresholdsRoot ensures that every
// threshold lies directly between the max
// and min permissible threshold.
func centerThresholdsRoot(c *Classifier, f map[string][]tokens.Freqs) {
	vecs := []tokens.Freqs{}
	for _, list := range f {
		for _, wordMap := range list {
			vecs = append(vecs, wordMap)
		}
	}
	centerThresholds(vecs, c)
}

func centerThresholds(vecs []tokens.Freqs, node *Classifier) {
	if node.LeafClassification != nil {
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
	centerThresholds(t, node.TrueBranch)
	centerThresholds(f, node.FalseBranch)
}

func splitOnNode(vecs []tokens.Freqs, node *Classifier) (t, f []tokens.Freqs) {
	t = make([]tokens.Freqs, 0, len(vecs))
	f = make([]tokens.Freqs, 0, len(vecs))
	for _, v := range vecs {
		if v[node.Keyword] > node.Threshold {
			t = append(t, v)
		} else {
			f = append(f, v)
		}
	}
	return
}
