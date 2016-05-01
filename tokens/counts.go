package tokens

import (
	"strings"
	"unicode"
)

// Counts records the number of occurrences
// of tokens in a given document.
type Counts map[string]int

// CountTokens counts the tokens in a document.
//
// Two different types of tokens are detected:
//
// - Heterogeneous tokens: any strings which
//   appear surrounded by whitespace.
// - Homogeneous words: strings like "abcd"
//   or "123" which are one type of symbol.
//
// No homogeneous tokens will be counted as
// heterogeneous tokens, or vice versa.
func CountTokens(contents string) Counts {
	res := Counts{}
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
