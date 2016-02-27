package whichlang

type Classifier struct {
	// Keywords is a weighted list of keywords for this classifier, with potentially negative
	// weights.
	Keywords map[string]float64

	// Threshold is a value above which the weighted sum of keywords should trigger a positive for
	// this classifier.
	Threshold float64
}

func (c *Classifier) Classify(f Frequencies) bool {
	return c.WeightedSum(f) > c.Threshold
}

func (c *Classifier) WeightedSum(f Frequencies) float64 {
	var sum float64
	for keyword, weight := range c.Keywords {
		sum += weight * f[keyword]
	}
	return sum
}
