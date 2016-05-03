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
// Four different types of tokens are detected:
//
// - Heterogeneous tokens: any strings which
//   appear surrounded by whitespace.
// - Homogeneous words: strings like "abcd"
//   or "123" which are one type of symbol.
// - Line-initial words: both heterogeneous
//   and homogeneous words which begin a line.
//   These tokens start with "\n".
// - Line-final words: both heterogeneous and
//   homogeneous words which end a line.
//   These tokens end with "\n".
//
// No homogeneous tokens will be counted as
// heterogeneous tokens, or vice versa.
// All line-boundary words are counted twice,
// once with the newline and once without it.
// If one token makes up an entire line, it is
// counted as both line-initial and line-final.
func CountTokens(contents string) Counts {
	res := Counts{}
	for _, t := range heterogeneousTokens(contents) {
		res[t] += 1
	}
	for _, t := range homogeneousTokens(contents) {
		res[t] += 1
	}
	for _, t := range lineBoundaryTokens(contents) {
		res[t] += 1
	}
	return res
}

func lineBoundaryTokens(contents string) []string {
	var res []string

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		for i, field := range []string{fields[0], fields[len(fields)-1]} {
			homog := homogeneousTokens(field)
			hetero := heterogeneousTokens(field)
			for _, tokList := range [][]string{homog, hetero} {
				if len(tokList) > 0 {
					if i == 0 {
						res = append(res, "\n"+tokList[0])
					} else {
						res = append(res, tokList[len(tokList)-1]+"\n")
					}
				}
			}
		}
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
