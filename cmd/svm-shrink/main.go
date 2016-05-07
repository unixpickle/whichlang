package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/whichlang/svm"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: svm-shrink <svm_in.json> <svm_out.json>")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		die(err)
	}

	classifier, err := svm.DecodeClassifier(data)
	if err != nil {
		die(err)
	} else if classifier.Kernel.Type != svm.LinearKernel {
		die(errors.New("can only shrink linear classifiers"))
	}

	langs := classifier.Languages()

	newClassifier := &svm.Classifier{
		Keywords:      classifier.Keywords,
		Kernel:        classifier.Kernel,
		SampleVectors: make([]linalg.Vector, len(langs)),
		Classifiers:   map[string]svm.BinaryClassifier{},
	}

	for i, lang := range langs {
		newClassifier.SampleVectors[i] = combineLanguageVecs(classifier, lang)
		bc := svm.BinaryClassifier{
			SupportVectors: []int{i},
			Weights:        []float64{1},
			Threshold:      classifier.Classifiers[lang].Threshold,
		}
		newClassifier.Classifiers[lang] = bc
	}

	encoded := newClassifier.Encode()
	if err := ioutil.WriteFile(os.Args[2], encoded, 0755); err != nil {
		die(err)
	}
}

func combineLanguageVecs(c *svm.Classifier, lang string) linalg.Vector {
	sum := make(linalg.Vector, len(c.Keywords))
	bc := c.Classifiers[lang]
	for i, idx := range bc.SupportVectors {
		sum.Add(c.SampleVectors[idx].Copy().Scale(bc.Weights[i]))
	}
	return sum
}

func die(e error) {
	fmt.Fprintln(os.Stderr, e)
	os.Exit(1)
}
