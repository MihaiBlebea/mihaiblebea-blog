package page

import (
	"sort"
	"strings"

	"github.com/MihaiBlebea/blog/go-broadcast/post"
)

func indexHandler(s *service) (template string, params interface{}, err error) {
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

func articleHandler(s *service, url string) (template string, params interface{}, err error) {
	posts, err := s.postService.GetAllPosts()
	if err != nil {
		return template, params, err
	}

	slug := strings.Replace(url, "/article/", "", -1)

	var p post.Post
	for _, post := range *posts {
		if post.Slug == slug {
			p = post
		}
	}

	var relatedPosts []post.Post
	for i, post := range *posts {
		if i == 3 {
			break
		}
		relatedPosts = append(relatedPosts, post)
	}

	params = struct {
		Articles *[]post.Post
		Article  *post.Post
	}{
		Articles: &relatedPosts,
		Article:  &p,
	}

	return "article", params, err
}

func tagHandler(s *service, url string) (template string, params interface{}, err error) {
	posts, err := s.postService.GetAllPosts()
	if err != nil {
		return template, params, err
	}

	parts := strings.Split(strings.Replace(url, "/tag/", "", -1), "/")
	if len(parts) == 0 {
		return template, params, err
	}

	tag := parts[0]

	var tagPosts []post.Post
	for _, p := range *posts {
		if contains(p.Tags, tag) == true {
			tagPosts = append(tagPosts, p)
		}
	}

	sort.SliceStable(tagPosts, func(i, j int) bool {
		return tagPosts[i].Published.After(tagPosts[j].Published)
	})

	params = struct {
		Articles *[]post.Post
		Tag      string
	}{
		Articles: &tagPosts,
		Tag:      tag,
	}

	return "index", params, err
}
