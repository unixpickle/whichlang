package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/howeyc/gopass"
)

// A GithubClient uses the Github API
// on behalf of a given user.
type GithubClient struct {
	User string
	Pass string
}

// PromptGithubClient prompts the user for their
// Github account details, then generates a
// *GithubClient based on these details.
func PromptGithubClient() (*GithubClient, error) {
	fmt.Print("Username: ")
	username := ""
	for {
		var ch [1]byte
		_, err := os.Stdin.Read(ch[:])
		if err != nil {
			return nil, err
		} else if ch[0] == '\n' {
			break
		} else if ch[0] == '\r' {
			continue
		}
		username += string(ch[0])
	}

	fmt.Print("Password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}

	return &GithubClient{
		User: strings.TrimSpace(username),
		Pass: string(password),
	}, nil
}

// request accesses an API URL using the
// user's credentials.
//
// It returns an error if the request fails,
// or if Github's API returns an error.
//
// Some requests are naturally paginated, in
// which case the next return argument
// corresponds to the URL of the next page.
func (g *GithubClient) request(u string) (data []byte, next *url.URL, err error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.preview.text-match+json")
	req.SetBasicAuth(g.User, g.Pass)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err == nil {
		if message, ok := raw["message"]; ok {
			if s, ok := message.(string); ok {
				return nil, nil, errors.New(s)
			}
		}
	}

	nextPattern := regexp.MustCompile("\\<(.*?)\\>; rel=\"next\"")
	match := nextPattern.FindStringSubmatch(res.Header.Get("Link"))
	if match != nil {
		u, _ := url.Parse(match[1])
		return body, u, nil
	}

	return body, nil, nil
}
