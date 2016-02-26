package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SampleDir string

func (s SampleDir) Languages() ([]string, error) {
	return readDirectory(string(s), true)
}

func (s SampleDir) SamplesForLanguage(lang string) ([]string, error) {
	return readDirectory(filepath.Join(string(s), lang), false)
}

func (s SampleDir) ReadSample(lang, sample string) (string, error) {
	p := filepath.Join(string(s), lang, sample)
	contents, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(contents), nil
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
