package svm

import (
	"log"
	"math/rand"
	"time"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/weakai/svm"
	"github.com/unixpickle/whichlang/tokens"
)

const crossValidationFraction = 0.3
const farAwayTimeout = time.Hour * 24 * 365

func Train(data map[string][]tokens.Freqs) *Classifier {
	params, err := EnvTrainerParams()
	if err != nil {
		panic(err)
	}
	return TrainParams(data, params)
}

func TrainParams(data map[string][]tokens.Freqs, p *TrainerParams) *Classifier {
	crossFreqs, trainingFreqs := partitionSamples(data, crossValidationFraction)
	tokens, samples := vectorizeSamples(trainingFreqs)

	solver := svm.GradientDescentSolver{
		Timeout:  farAwayTimeout,
		Tradeoff: p.Tradeoff,
	}

	var bestClassifier *Classifier
	var bestValidationScore float64

	for _, kernel := range p.Kernels {
		if p.Verbose {
			log.Println("Trying kernel:", kernel)
		}
		solverKernel := cachedKernel(kernel)
		classifier := &Classifier{
			Keywords:    tokens,
			Kernel:      kernel,
			Classifiers: map[string]BinaryClassifier{},
		}

		usedSamples := map[int]linalg.Vector{}
		for lang := range samples {
			if p.Verbose {
				log.Println("Training classifier for language:", lang)
			}
			problem := svmProblem(samples, lang, solverKernel)
			solution := solver.Solve(problem)
			binClass := BinaryClassifier{
				SupportVectors: make([]int, len(solution.SupportVectors)),
				Weights:        make([]float64, len(solution.Coefficients)),
				Threshold:      -solution.Threshold,
			}
			copy(binClass.Weights, solution.Coefficients)
			for i, v := range solution.SupportVectors {
				// v.UserInfo will be turned into a support
				// vector index by makeSampleVectorList().
				binClass.SupportVectors[i] = v.UserInfo
				usedSamples[v.UserInfo] = linalg.Vector(v.V)
			}
			classifier.Classifiers[lang] = binClass
		}

		makeSampleVectorList(classifier, usedSamples)

		score := correctFraction(classifier, crossFreqs)
		if p.Verbose {
			trainingScore := correctFraction(classifier, trainingFreqs)
			log.Printf("Results: cross=%f training=%f support=%d/%d", score,
				trainingScore, len(classifier.SampleVectors), countSamples(samples))
		}
		if score > bestValidationScore || bestClassifier == nil {
			bestClassifier = classifier
		}
	}

	return bestClassifier
}

func partitionSamples(data map[string][]tokens.Freqs, crossFrac float64) (cross,
	training map[string][]tokens.Freqs) {

	cross = map[string][]tokens.Freqs{}
	training = map[string][]tokens.Freqs{}

	for lang, samples := range data {
		p := rand.Perm(len(samples))
		newSamples := make([]tokens.Freqs, len(samples))
		for i, x := range p {
			newSamples[i] = samples[x]
		}
		crossCount := int(crossFrac * float64(len(samples)))
		cross[lang] = newSamples[:crossCount]
		training[lang] = newSamples[crossCount:]
	}

	return
}

func vectorizeSamples(data map[string][]tokens.Freqs) ([]string, map[string][]svm.Sample) {
	seenToks := map[string]bool{}
	for _, samples := range data {
		for _, sample := range samples {
			for tok := range sample {
				seenToks[tok] = true
			}
		}
	}
	toks := make([]string, 0, len(seenToks))
	for tok := range seenToks {
		toks = append(toks, tok)
	}

	sampleMap := map[string][]svm.Sample{}
	sampleID := 1
	for lang, samples := range data {
		vecSamples := make([]svm.Sample, 0, len(samples))
		for _, sample := range samples {
			vec := make([]float64, len(toks))
			for i, tok := range toks {
				vec[i] = sample[tok]
			}
			svmSample := svm.Sample{
				V:        vec,
				UserInfo: sampleID,
			}
			sampleID++
			vecSamples = append(vecSamples, svmSample)
		}
		sampleMap[lang] = vecSamples
	}

	return toks, sampleMap
}

func countSamples(s map[string][]svm.Sample) int {
	var count int
	for _, samples := range s {
		count += len(samples)
	}
	return count
}

func cachedKernel(k *Kernel) svm.Kernel {
	return svm.CachedKernel(func(s1, s2 svm.Sample) float64 {
		return k.Product(linalg.Vector(s1.V), linalg.Vector(s2.V))
	})
}

func svmProblem(data map[string][]svm.Sample, posLang string, k svm.Kernel) *svm.Problem {
	var positives, negatives []svm.Sample
	for lang, samples := range data {
		if lang == posLang {
			positives = append(positives, samples...)
		} else {
			negatives = append(negatives, samples...)
		}
	}
	return &svm.Problem{
		Positives: positives,
		Negatives: negatives,
		Kernel:    k,
	}
}

func correctFraction(c *Classifier, data map[string][]tokens.Freqs) float64 {
	var correct, total int
	for lang, samples := range data {
		for _, sample := range samples {
			total++
			if c.Classify(sample) == lang {
				correct++
			}
		}
	}
	return float64(correct) / float64(total)
}

func makeSampleVectorList(c *Classifier, used map[int]linalg.Vector) {
	userInfoToVecIdx := map[int]int{}

	for userInfo, sample := range used {
		userInfoToVecIdx[userInfo] = len(c.SampleVectors)
		c.SampleVectors = append(c.SampleVectors, sample)
	}

	for _, binClass := range c.Classifiers {
		for i, userInfo := range binClass.SupportVectors {
			binClass.SupportVectors[i] = userInfoToVecIdx[userInfo]
		}
	}
}
