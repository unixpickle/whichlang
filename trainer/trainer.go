package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/unixpickle/weakai/svm"
	"github.com/unixpickle/whichlang"
)

var gradientSolver = svm.GradientDescentSolver{
	Steps:    10,
	StepSize: 0.01,
	Tradeoff: 0.1,
}

func GenerateClassifiers(s Samples, maxFeatures int) map[string]*whichlang.Classifier {
	res := map[string]*whichlang.Classifier{}
	for lang := range s.Samples {
		fmt.Println(" - generating classifier for:", lang)
		svmClassifier, problem := languageSVMClassifier(lang, s)
		classifier := canonicalClassifier(svmClassifier, maxFeatures, problem, s)
		res[lang] = classifier
	}
	return res
}

func languageSVMClassifier(lang string, s Samples) (sol *svm.LinearClassifier, p *svm.Problem) {
	p = &svm.Problem{
		Positives: make([]svm.Sample, len(s.Samples[lang])),
		Negatives: make([]svm.Sample, 0),
		Kernel:    svm.LinearKernel,
	}

	copy(p.Positives, s.Samples[lang])

	for l, samples := range s.Samples {
		if l == lang {
			continue
		}
		p.Negatives = append(p.Negatives, samples...)
	}

	sol = gradientSolver.Solve(p).Linearize()
	return
}

func canonicalClassifier(c *svm.LinearClassifier, maxFeatures int,
	p *svm.Problem, s Samples) *whichlang.Classifier {
	features := make(floatIntPairList, len(c.HyperplaneNormal))
	for i, f := range c.HyperplaneNormal {
		features[i] = floatIntPair{f, i}
	}
	sort.Sort(features)

	for i, sample := range p.Positives {
		p.Positives[i] = takeBestFeatures(features, sample, maxFeatures)
	}
	for i, sample := range p.Negatives {
		p.Negatives[i] = takeBestFeatures(features, sample, maxFeatures)
	}

	solution := gradientSolver.Solve(p).Linearize()

	res := &whichlang.Classifier{
		Keywords:  map[string]float64{},
		Threshold: -solution.Threshold,
	}

	for i, weight := range solution.HyperplaneNormal {
		word := s.Words[features[i].i]
		res.Keywords[word] = weight
	}

	return res
}

func takeBestFeatures(f floatIntPairList, s svm.Sample, maxFeatures int) svm.Sample {
	if maxFeatures >= len(s) {
		return s
	}
	res := make(svm.Sample, maxFeatures)
	for i := 0; i < maxFeatures; i++ {
		res[i] = s[f[i].i]
	}
	return res
}

type floatIntPair struct {
	f float64
	i int
}

type floatIntPairList []floatIntPair

func (f floatIntPairList) Len() int {
	return len(f)
}

func (f floatIntPairList) Less(i, j int) bool {
	return math.Abs(f[i].f) > math.Abs(f[j].f)
}

func (f floatIntPairList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
