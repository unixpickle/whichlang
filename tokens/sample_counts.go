package tokens

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SampleCounts maps programming languages to
// arrays of language sample documents, where
// each sample is represented by a Counts.
type SampleCounts map[string][]Counts

// ReadSampleCounts computes token counts
// for programming language samples in a
// directory.
//
// The directory should contain sub-directories
// for each programming language, and each of
// these languages should contain one or more
// source files.
//
// The returned map maps language names to lists
// of Counts, where each Counts corresponds to
// one source file.
func ReadSampleCounts(sampleDir string) (SampleCounts, error) {
	languages, err := readDirectory(sampleDir, true)
	if err != nil {
		return nil, err
	}

	res := SampleCounts{}
	for _, language := range languages {
		langDir := filepath.Join(sampleDir, language)
		files, err := readDirectory(langDir, false)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			contents, err := ioutil.ReadFile(filepath.Join(langDir, file))
			if err != nil {
				return nil, err
			}
			counts := CountTokens(string(contents))
			res[language] = append(res[language], counts)
		}
	}

	return res, nil
}

// NumTokens returns the number of unique
// tokens in all the documents.
func (s SampleCounts) NumTokens() int {
	toks := map[string]bool{}
	for _, samples := range s {
		for _, sample := range samples {
			for word := range sample {
				toks[word] = true
			}
		}
	}
	return len(toks)
}

// Prune removes tokens which appear in n
// documents or fewer.
func (s SampleCounts) Prune(n int) {
	docCount := map[string]int{}
	for _, samples := range s {
		for _, sample := range samples {
			for word := range sample {
				docCount[word]++
			}
		}
	}

	remove := map[string]bool{}
	for word, count := range docCount {
		if count <= n {
			remove[word] = true
		}
	}

	for _, samples := range s {
		for i, sample := range samples {
			newSample := map[string]int{}
			for word, count := range sample {
				if !remove[word] {
					newSample[word] = count
				}
			}
			samples[i] = newSample
		}
	}
}

// SampleFreqs converts every Counts object
// in s into a Freqs object.
func (s SampleCounts) SampleFreqs() map[string][]Freqs {
	res := map[string][]Freqs{}
	for lang, samples := range s {
		for _, sample := range samples {
			res[lang] = append(res[lang], sample.Freqs())
		}
	}
	return res
}

func readDirectory(dir string, isDir bool) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	contents, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(contents))
	for _, info := range contents {
		if info.IsDir() == isDir && !strings.HasPrefix(info.Name(), ".") {
			res = append(res, info.Name())
		}
	}
	sort.Strings(res)
	return res, nil
}
