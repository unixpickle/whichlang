package main

import (
	"encoding/json"
	"net/url"
)

// Search asynchronously lists repositories which
// are written in the given programming language.
// Repository names are of the form "user/repo".
//
// The caller should close the done argument when
// they do not need any more results.
// When no results are left, or when done is closed,
// or on error, both returned channels are closed.
//
// If search results cannot be obtained, an error
// is sent on the error channel.
func (g *GithubClient) Search(lang string, done <-chan struct{}) (<-chan string, <-chan error) {
	nameChan := make(chan string, 0)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			close(nameChan)
			close(errChan)
		}()
		u := repositorySearchURL(lang)
		for u != nil {
			select {
			case <-done:
				return
			default:
			}

			body, next, err := g.request(u.String())
			if err != nil {
				errChan <- err
				return
			}
			u = next

			var obj struct {
				Items []struct {
					FullName string `json:"full_name"`
				} `json:"items"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				errChan <- err
				return
			}

			for _, x := range obj.Items {
				select {
				case nameChan <- x.FullName:
				case <-done:
					return
				}
			}
		}
	}()
	return nameChan, errChan
}

func repositorySearchURL(lang string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "search/repositories",
		RawQuery: url.Values{
			"order": []string{"desc"},
			"q":     []string{"language:" + lang},
		}.Encode(),
	}
}
