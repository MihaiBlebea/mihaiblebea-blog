package pages

import (
	"sort"

	"github.com/MihaiBlebea/blog/go-broadcast/postrepo"
)

// PageType defines the name of the page
type PageType int

// Page types
const (
	HomePage PageType = iota
)

// Factory is the page builder
type Factory interface {
}

type factory struct {
	postRepo *postrepo.Service
}

// New returns a new Factory interface
func New(postRepo *postrepo.Service) Factory {
	return &factory{postRepo}
}

func (s *factory) Build(pageName PageType) (*Page, error) {
	params = struct {
		Articles *[]post.Post
	}{
		Articles: &p,
	}

	return &Page{}
}

func (s *service) BuildHomePage() (*Page, error) {
	posts, err := s.postService.GetAllPosts()
	if err != nil {
		return template, params, err
	}

	p := *posts

	sort.SliceStable(p, func(i, j int) bool {
		return p[i].Published.After(p[j].Published)
	})

	params = struct {
		Articles *[]post.Post
	}{
		Articles: &p,
	}

	return "index", params, nil
}
