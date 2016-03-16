package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// A GithubClient uses the Github API on behalf of a given user.
type GithubClient struct {
	User string
	Pass string
}

// LanguageRepositories returns a list of repository names in the form "username/repository".
func (g *GithubClient) LanguageRepositories(lang string) ([]string, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "search/repositories",
		RawQuery: url.Values{
			"order": []string{"desc"},
			"q":     []string{"language:" + lang},
		}.Encode(),
	}
	body, err := g.request(u.String())
	if err != nil {
		return nil, err
	}

	var obj langRepoList
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}

	res := make([]string, len(obj.Items))
	for i, x := range obj.Items {
		res[i] = x.FullName
	}
	return res, nil
}

// FirstFile scans the given repository and returns the contents of the first file which meet the
// given criteria.
//
// This may return (nil, nil) if no file meeting the criteria was found, but no network error
// occurred.
func (g *GithubClient) FirstFile(repo string, minFileSize int,
	extensions []string) ([]byte, error) {
	return g.firstFileSearch(repo, minFileSize, extensions, "/")
}

func (g *GithubClient) firstFileSearch(repo string, minFileSize int, extensions []string,
	searchPath string) (match []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   path.Join("/repos", repo, "/contents", searchPath),
	}
	body, err := g.request(u.String())
	if err != nil {
		return nil, err
	}

	var result []entity
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	for _, ent := range result {
		if ent.meetsSearchCriterion(minFileSize, extensions) {
			return g.readFile(repo, ent.Path)
		}
	}

	sourceDirectoryHeuristic(result, repo)

	for _, ent := range result {
		if ent.Dir() {
			match, err = g.firstFileSearch(repo, minFileSize, extensions, ent.Path)
			if match != nil || err != nil {
				return
			}
		}
	}

	return nil, nil
}

func (g *GithubClient) readFile(repo, filePath string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   path.Join("/repos", repo, "/contents", filePath),
	}
	body, err := g.request(u.String())
	if err != nil {
		return nil, err
	}

	var result fileEntity
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Encoding == "base64" {
		return base64.StdEncoding.DecodeString(result.Content)
	} else {
		return nil, errors.New("unknown encoding: " + result.Encoding)
	}
}

func (g *GithubClient) request(u string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.preview.text-match+json")
	req.SetBasicAuth(g.User, g.Pass)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err == nil {
		if message, ok := raw["message"]; ok {
			if s, ok := message.(string); ok {
				return nil, errors.New(s)
			}
		}
	}

	return body, nil
}

type langRepoList struct {
	Items []langRepoItem `json:"items"`
}

type langRepoItem struct {
	FullName string `json:"full_name"`
}

type entity struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (e entity) Dir() bool {
	return e.Type == "dir"
}

type fileEntity struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func (e *entity) meetsSearchCriterion(minSize int, exts []string) bool {
	if e.Type != "file" {
		return false
	}
	if e.Size < minSize {
		return false
	}
	for _, ext := range exts {
		if strings.HasSuffix(e.Name, "."+ext) {
			return true
		}
	}
	return false
}

// sourceDirectoryHeuristic puts directories which are likely to contain source code at the
// beginning of a list of entities.
func sourceDirectoryHeuristic(results []entity, repoName string) {
	sourceDirs := []string{"src", repoName, "lib", "com", "org", "net"}
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
