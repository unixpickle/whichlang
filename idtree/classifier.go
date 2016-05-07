package idtree

import (
	"encoding/json"

	"github.com/unixpickle/whichlang/tokens"
)

type Classifier struct {
	LeafClassification *string

	Keyword   string
	Threshold float64

	FalseBranch *Classifier
	TrueBranch  *Classifier
}

func DecodeClassifier(d []byte) (*Classifier, error) {
	var res Classifier
	if err := json.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Classifier) Classify(f tokens.Freqs) string {
	if c.LeafClassification == nil {
		if f[c.Keyword] > c.Threshold {
			return c.TrueBranch.Classify(f)
		} else {
			return c.FalseBranch.Classify(f)
		}
	} else {
		return *c.LeafClassification
	}
}

func (c *Classifier) Encode() []byte {
	res, _ := json.Marshal(c)
	return res
}

func (c *Classifier) Languages() []string {
	if c.LeafClassification != nil {
		return []string{*c.LeafClassification}
	}

	seen := map[string]bool{}
	for _, lang := range c.FalseBranch.Languages() {
		seen[lang] = true
	}
	for _, lang := range c.TrueBranch.Languages() {
		seen[lang] = true
	}

	res := make([]string, 0, len(seen))
	for lang := range seen {
		res = append(res, lang)
	}
	return res
}
