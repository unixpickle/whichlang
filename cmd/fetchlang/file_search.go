package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strings"
)

var (
	ErrNoResults   = errors.New("no results")
	ErrMaxRequests = errors.New("too many requests")
)

// A FileSearch defines parameters for
// searching a repository for files.
type FileSearch struct {
	// Repository is the repository name,
	// formatted as "user/repo".
	Repository string

	MinFileSize int
	MaxFileSize int
	Extensions  []string

	// MaxRequests is the maximum number of
	// API requests to be performed by the
	// search before giving up.
	MaxRequests int
}

// SearchFile runs a FileSearch.
// It returns ErrMaxRequests if more than
// s.MaxRequests requests are used.
// It returns ErrNoResults if no results
// are found.
func (g *GithubClient) SearchFile(s FileSearch) (contents []byte, err error) {
	return g.firstFileSearch(&s, "/")
}

func (g *GithubClient) firstFileSearch(s *FileSearch, dir string) (match []byte, err error) {
	if s.MaxRequests == 0 {
		return nil, ErrMaxRequests
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   path.Join("/repos", s.Repository, "/contents", dir),
	}

	body, _, err := g.request(u.String())
	if err != nil {
		return nil, err
	}

	s.MaxRequests--

	var result []entity
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	for _, ent := range result {
		if ent.Match(s) {
			return g.readFile(s.Repository, ent.Path)
		}
	}

	sourceDirectoryHeuristic(result, s.Repository)

	for _, ent := range result {
		if ent.Dir() {
			match, err = g.firstFileSearch(s, ent.Path)
			if match != nil || (err != nil && err != ErrNoResults) {
				return
			}
		}
	}

	return nil, ErrNoResults
}

func (g *GithubClient) readFile(repo, filePath string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   path.Join("/repos", repo, "/contents", filePath),
	}
	body, _, err := g.request(u.String())
	if err != nil {
		return nil, err
	}

	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Encoding == "base64" {
		return base64.StdEncoding.DecodeString(result.Content)
	} else {
		return nil, errors.New("unknown encoding: " + result.Encoding)
	}
}

type entity struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (e *entity) Dir() bool {
	return e.Type == "dir"
}

func (e *entity) Match(s *FileSearch) bool {
	if e.Type != "file" {
		return false
	}
	if e.Size < s.MinFileSize || e.Size > s.MaxFileSize {
		return false
	}
	for _, ext := range s.Extensions {
		if strings.HasSuffix(e.Name, "."+ext) {
			return true
		}
	}
	return false
}

// sourceDirectoryHeuristic sorts a list of
// entities so that the first ones are more
// likely to contain source code.
func sourceDirectoryHeuristic(results []entity, repoName string) {
	sourceDirs := []string{"src", repoName, "lib", "com", "org", "net", "css", "assets"}
	numFound := 0
	for _, sourceDir := range sourceDirs {
		for i, ent := range results[numFound:] {
			if ent.Dir() && ent.Name == sourceDir {
				results[numFound], results[i] = results[i], results[numFound]
				numFound++
				break
			}
		}
	}
}
