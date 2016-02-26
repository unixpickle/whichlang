package whichlang

import (
	"strings"
	"unicode"
)

// A Frequencies object represents a histogram of word occurrences in a document.
// Each word is mapped to a number indicating the number of times it occurred.
//
// Frequency counts are floating points rather than integers to allow Frequencies objects to be
// averaged and manipulated in various ways.
type Frequencies map[string]float64

// ComputeFrequencies returns a frequency map of all the words which appear in a string.
// Two different types of words are detected:
//
// - Heterogeneous words: any strings which appear surrounded by whitespace.
// - Homogeneous words: strings of one particular character type (e.g. letter) which may appear
//   inside a larger heterogenous keyword.
//
// Both types of words are weighted equally in the result.
// No re-counting will occur (e.g. "this is a test" has four homogeneous words and 0 heterogeneous
// ones).
func ComputeFrequencies(contents string) Frequencies {
	res := Frequencies{}
	for _, t := range heterogeneousTokens(contents) {
		res[t] += 1
	}
	for _, t := range homogeneousTokens(contents) {
		res[t] += 1
	}
	return res
}

func heterogeneousTokens(contents string) []string {
	fields := strings.Fields(contents)
	res := make([]string, 0, len(fields))
	for _, f := range fields {
		if !isHeterogeneous(f) {
			res = append(res, f)
		}
	}
	return res
}

func homogeneousTokens(contents string) []string {
	tokens := []string{}
	res := ""
	lastClass := charClassSpace
	for _, ch := range contents {
		c := classForRune(ch)
		if c == lastClass {
			res += string(ch)
			continue
		}
		if lastClass != charClassSpace && len(res) > 0 {
			tokens = append(tokens, res)
		}
		res = string(ch)
		lastClass = c
	}
	if lastClass != charClassSpace && len(res) > 0 {
		tokens = append(tokens, res)
	}
	return tokens
}

type charClass int

const (
	charClassLetter charClass = iota
	charClassNumber
	charClassSpace
	charClassSymbol
)

func classForRune(r rune) charClass {
	if unicode.IsLetter(r) {
		return charClassLetter
	} else if unicode.IsDigit(r) {
		return charClassNumber
	} else if unicode.IsSpace(r) {
		return charClassSpace
	}
	return charClassSymbol
}

func isHeterogeneous(s string) bool {
	if len(s) == 0 {
		return true
	}
	c := classForRune([]rune(s)[0])
	for _, r := range s {
		if classForRune(r) != c {
			return false
		}
	}
	return true
}
