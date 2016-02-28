package whichlang

import "math"

type ClassifierNode struct {
	Leaf               bool
	LeafClassification string

	Keyword   string
	Threshold float64

	FalseBranch *ClassifierNode
	TrueBranch  *ClassifierNode
}

func (c *ClassifierNode) leafCount() int {
	if c.Leaf {
		return 1
	}
	return c.FalseBranch.leafCount() + c.TrueBranch.leafCount()
}

// A Classifier uses an identification tree to classify a piece of code.
type Classifier struct {
	Keywords []string
	TreeRoot *ClassifierNode
}

func (c *Classifier) Classify(f Frequencies) string {
	normalized := c.normalizeKeywords(f)
	node := c.TreeRoot
	for !node.Leaf {
		if normalized[node.Keyword] > node.Threshold {
			node = node.TrueBranch
		} else {
			node = node.FalseBranch
		}
	}
	return node.LeafClassification
}

func (c *Classifier) LeafCount() int {
	return c.TreeRoot.leafCount()
}

func (c *Classifier) normalizeKeywords(f Frequencies) Frequencies {
	var mag2 float64
	for _, word := range c.Keywords {
		val := f[word]
		mag2 += val * val
	}
	scaler := 1 / math.Sqrt(mag2)
	res := map[string]float64{}
	for _, word := range c.Keywords {
		res[word] *= f[word] * scaler
	}
	return res
}
