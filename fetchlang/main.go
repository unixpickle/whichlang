package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const MinFileSize = 100
const MaxFileSize = 500000
const MaxRequestsPerRepo = 20

type Language struct {
	Name       string
	Extensions []string
}

var Languages = []Language{
	{"ActionScript", []string{"as"}},
	{"C", []string{"c"}},
	{"C#", []string{"cs"}},
	{"C++", []string{"cpp", "c++", "C", "cc"}},
	{"Clojure", []string{"clj"}},
	{"CoffeeScript", []string{"coffee"}},
	{"CSS", []string{"css"}},
	{"Go", []string{"go"}},
	{"Haskell", []string{"hs"}},
	{"HTML", []string{"html", "htm"}},
	{"Java", []string{"java"}},
	{"JavaScript", []string{"js"}},
	{"Lua", []string{"lua"}},
	{"Matlab", []string{"m"}},
	{"Objective-C", []string{"m"}},
	{"Perl", []string{"pl"}},
	{"PHP", []string{"php"}},
	{"Python", []string{"py"}},
	{"R", []string{"r", "R"}},
	{"Ruby", []string{"rb"}},
	{"Scala", []string{"scala"}},
	{"Shell", []string{"sh", "bash"}},
	{"Swift", []string{"swift"}},
	{"TeX", []string{"tex", "TeX"}},
	{"VimL", []string{"vim"}},
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: fetchlang <sample dir> <sample count>")
		os.Exit(1)
	}
	sampleDir := os.Args[1]
	sampleCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid sample count:", os.Args[2])
		os.Exit(1)
	}

	client, err := PromptGithubClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read credentials:", err)
		os.Exit(1)
	}

	for _, lang := range Languages {
		fmt.Println("Fetching samples for", lang.Name, "...")
		err := fetchLanguage(client, sampleDir, sampleCount, lang)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func fetchLanguage(github *GithubClient, sampleDir string, count int, lang Language) error {
	if err := os.Mkdir(filepath.Join(sampleDir, lang.Name), 0755); err != nil {
		return err
	}

	doneChan := make(chan struct{}, 1)
	repoChan, errChan := github.Search(lang.Name, doneChan)

	defer func() {
		close(doneChan)
	}()

	resCount := 0
	for repo := range repoChan {
		file, err := github.SearchFile(FileSearch{
			Repository:  repo,
			MinFileSize: MinFileSize,
			MaxFileSize: MaxFileSize,
			Extensions:  lang.Extensions,
			MaxRequests: MaxRequestsPerRepo,
		})
		if err == ErrNoResults {
			fmt.Println("No results for:", repo)
			continue
		} else if err == ErrMaxRequests {
			fmt.Println("Max requests exceeded:", repo)
			continue
		} else if err != nil {
			return err
		} else if file == nil {
			continue
		}

		fileName := fmt.Sprintf("%d.%s", resCount, lang.Extensions[0])
		targetFile := filepath.Join(sampleDir, lang.Name, fileName)
		if err := ioutil.WriteFile(targetFile, file, 0755); err != nil {
			return err
		}

		resCount++
		if resCount == count {
			break
		}
	}

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
