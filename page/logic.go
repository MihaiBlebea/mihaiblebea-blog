package page

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/MihaiBlebea/blog/go-broadcast/cache"
	"github.com/MihaiBlebea/blog/go-broadcast/post"
	"github.com/sirupsen/logrus"
)

type service struct {
	postService post.Service
	cache       *cache.Cache
	logger      *logrus.Logger
}

// New returns a new page service
func New(postService post.Service, cache *cache.Cache, logger *logrus.Logger) Service {
	return &service{
		postService: postService,
		cache:       cache,
		logger:      logger,
	}
}

func (s *service) LoadStaticFile(URL string) ([]byte, error) {
	b, err := ioutil.ReadFile(URL)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func (s *service) LoadTemplate(URL string) (*Page, error) {
	var template string
	var params interface{}
	var err error = nil

	if URL == "/" {
		template, params, err = indexHandler(s)
	} else if strings.Contains(URL, "/article") {
		template, params, err = articleHandler(s, URL)
	} else if strings.Contains(URL, "/tag/") {
		template, params, err = tagHandler(s, URL)
	} else {
		template = strings.Split(URL[1:], "/")[0]
		params = nil
	}

	if err != nil {
		return s.LoadErrorPage(err)
	}

	page, err := s.loadPage(template, params)
	if err != nil {
		return s.LoadErrorPage(err)
	}

	return page, nil
}

func (s *service) loadPage(templateName string, params interface{}) (*Page, error) {
	tmpl, err := s.parseTemplates()
	if err != nil {
		return &Page{}, err
	}

	return &Page{
		Params:       params,
		Template:     tmpl,
		TemplateName: templateName,
	}, nil
}

func (s *service) LoadErrorPage(err error) (*Page, error) {
	tmpl, err := s.parseTemplates()
	if err != nil {
		return &Page{}, err
	}

	return &Page{
		Params: struct {
			Err string
		}{
			Err: err.Error(),
		},
		Template:     tmpl,
		TemplateName: "error",
	}, nil
}

func (s *service) parseTemplates() (*template.Template, error) {
	templ := template.New("")
	err := filepath.Walk("./static/templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return templ, nil
}

func contains(hay []string, needle string) bool {
	for _, straw := range hay {
		if straw == needle {
			return true
		}
	}

	return false
}
