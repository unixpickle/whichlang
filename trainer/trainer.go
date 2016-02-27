package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/unixpickle/weakai/svm"
	"github.com/unixpickle/whichlang"
)

var gradientSolver = svm.GradientDescentSolver{
	Steps:    1000,
	StepSize: 0.01,
	Tradeoff: 0.0001,
}

func GenerateClassifiers(s Samples, maxFeatures int) map[string]*whichlang.Classifier {
	res := map[string]*whichlang.Classifier{}
	k := svm.CachedKernel(svm.LinearKernel)
	for lang := range s.Samples {
		fmt.Println(" - generating classifier for:", lang)
		svmClassifier, problem := languageSVMClassifier(lang, s, k)
		classifier := canonicalClassifier(svmClassifier, maxFeatures, problem, s)
		res[lang] = classifier
	}
	return res
}

func languageSVMClassifier(lang string, s Samples,
	k svm.Kernel) (sol *svm.LinearClassifier, p *svm.Problem) {
	p = &svm.Problem{
		Positives: make([]svm.Sample, len(s.Samples[lang])),
		Negatives: make([]svm.Sample, 0),
		Kernel:    k,
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
	features := make(floatIntPairList, len(c.HyperplaneNormal.V))
	for i, f := range c.HyperplaneNormal.V {
		features[i] = floatIntPair{f, i}
	}
	sort.Sort(features)

	for i, sample := range p.Positives {
		p.Positives[i] = takeBestFeatures(features, sample, maxFeatures)
	}
	for i, sample := range p.Negatives {
		p.Negatives[i] = takeBestFeatures(features, sample, maxFeatures)
	}

	// We must use a new cached kernel, since the trimmed Samples have the same UserInfo values as
	// their untrimmed counterparts.
	// This is acceptable, because we probably haven't seen these exact trimmed vectors before, so
	// trying to re-use a cache for them is rather pointless.
	p.Kernel = svm.CachedKernel(svm.LinearKernel)

	solution := gradientSolver.Solve(p).Linearize()

	res := &whichlang.Classifier{
		Keywords:  map[string]float64{},
		Threshold: -solution.Threshold,
	}

	for i, weight := range solution.HyperplaneNormal.V {
		word := s.Words[features[i].i]
		res.Keywords[word] = weight
	}

	return res
}

func takeBestFeatures(f floatIntPairList, s svm.Sample, maxFeatures int) svm.Sample {
	if maxFeatures >= len(s.V) {
		return s
	}
	res := make([]float64, maxFeatures)
	for i := 0; i < maxFeatures; i++ {
		res[i] = s.V[f[i].i]
	}
	return svm.Sample{V: res, UserInfo: s.UserInfo}
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
