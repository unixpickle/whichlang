package tokens

import (
	"math"
	"testing"
)

func TestSampleCountsPruneFreqs(t *testing.T) {
	docs := SampleCounts{
		"A": []Counts{
			{"Foo": 1, "Bar": 3, "Baz": 2, "Once1": 1},
			{"Foo": 1, "Bar": 1, "Once2": 15},
		},
		"B": []Counts{
			{"Baz": 15, "Once3": 17},
		},
	}
	docs.Prune(1)
	actual := docs.SampleFreqs()
	expected := map[string][]Freqs{
		"A": []Freqs{
			{"Foo": 1.0 / 7.0, "Bar": 3.0 / 7.0, "Baz": 2.0 / 7.0},
			{"Foo": 1.0 / 17.0, "Bar": 1.0 / 17.0},
		},
		"B": []Freqs{
			{"Baz": 15.0 / 32.0},
		},
	}
	for lang, freqs := range expected {
		actualFreqs := actual[lang]
		if len(actualFreqs) != len(freqs) {
			t.Error("unexpected document count for", lang)
			continue
		}
		for i, actualFreq := range actualFreqs {
			expFreq := freqs[i]
			if !freqsApproxEqual(actualFreq, expFreq) {
				t.Error("incorrect freq", actualFreq)
			}
		}
	}
	for lang := range actual {
		if expected[lang] == nil {
			t.Error("unexpected language", lang)
		}
	}
}

func freqsApproxEqual(f1, f2 Freqs) bool {
	if len(f1) != len(f2) {
		return false
	}
	for key, val := range f1 {
		if math.Abs(val-f2[key]) > 1e-5 {
			return false
		}
	}
	return true
}
