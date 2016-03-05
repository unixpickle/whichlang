package whichlang

type ClassifierNode struct {
	Leaf               bool
	LeafClassification string
	LeafConfidence     float64

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
	TreeRoot *ClassifierNode
}

func (c *Classifier) Classify(f Frequencies) (lang string, confidence float64) {
	normalized := c.normalizeKeywords(f)
	node := c.TreeRoot
	for !node.Leaf {
		if normalized[node.Keyword] > node.Threshold {
			node = node.TrueBranch
		} else {
			node = node.FalseBranch
		}
	}
	return node.LeafClassification, node.LeafConfidence
}

func (c *Classifier) LeafCount() int {
	return c.TreeRoot.leafCount()
}

func (c *Classifier) normalizeKeywords(f Frequencies) Frequencies {
	var totalSum float64
	for _, val := range f {
		totalSum += val
	}
	if totalSum == 0 {
		totalSum = 1
	}
	scaler := 1 / totalSum
	res := map[string]float64{}
	for key, val := range f {
		res[key] = val * scaler
	}
	return res
}
