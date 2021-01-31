package postrepo

import (
	"errors"
	"time"
)

type service struct {
	githubToken string
	articles    []string
}

// Service represents the interface for the service
type Service interface {
	RefillCache() error
	GetAllContent() []string
}

// New returns a new repo service
func New(githubToken string) Service {
	return &service{githubToken: githubToken}
}

// RefillCache makes a request to fetch all the data from the github repo, then it saves all posts into a slice
func (s *service) RefillCache() error {
	url, err := getTreeSHA(s.githubToken)
	if err != nil {
		return err
	}

	urls, err := getFileURLs(url, s.githubToken)
	if err != nil {
		return err
	}

	type Result struct {
		content string
		err     error
	}
	res := make(chan Result, len(urls))

	for _, url := range urls {
		go func(url string) {
			content, err := getFileContent(url, s.githubToken)
			if err != nil {
				res <- Result{err: err}
			}
			res <- Result{content: content}
		}(url)
	}

	for {
		select {
		case r := <-res:
			if r.err != nil {
				continue
			}

			s.articles = append(s.articles, r.content)
		case <-time.After(5 * time.Second):
			return errors.New("Github request has timedout")
		}
	}
}

// GetAllContent returns all the posts saved as strings into the posts slice
func (s *service) GetAllContent() []string {
	return s.articles
}
