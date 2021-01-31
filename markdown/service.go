package markdown

import (
	"bytes"
	"errors"
	"html/template"

	"github.com/MihaiBlebea/blog/go-broadcast/postrepo"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	rendererHTML "github.com/yuin/goldmark/renderer/html"
)

// Errors
var (
	ErrInvalidType      = errors.New("Invalid type while converting from interface")
	ErrPostNotPublished = errors.New("Post was not published yet")
)

// Service is a service that transforms string content into Post models
type Service interface {
	ParsePost(content string) (*postrepo.Post, error)
}

type markdown struct {
	parser goldmark.Markdown
}

// New returns a new markdown service
func New() Service {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			rendererHTML.WithHardWraps(),
			rendererHTML.WithXHTML(),
			rendererHTML.WithUnsafe(),
		),
	)

	return &markdown{md}
}

func (m *markdown) ParsePost(content string) (*postrepo.Post, error) {
	context := parser.NewContext()

	var buf bytes.Buffer
	err := m.parser.Convert([]byte(content), &buf, parser.WithContext(context))
	if err != nil {
		return &postrepo.Post{}, err
	}

	params := meta.Get(context)

	var published string
	if _, ok := params["Published"]; ok != false {
		published, ok = params["Published"].(string)
		if ok != true {
			return &postrepo.Post{}, ErrInvalidType
		}
	}

	title, ok := params["Title"].(string)
	if ok != true {
		return &postrepo.Post{}, ErrInvalidType
	}

	slug, ok := params["Slug"].(string)
	if ok != true {
		return &postrepo.Post{}, ErrInvalidType
	}

	var image string
	if _, ok := params["Image"]; ok != false {
		image, ok = params["Image"].(string)
		if ok != true {
			return &postrepo.Post{}, ErrInvalidType
		}
	}

	var summary string
	if _, ok := params["Summary"]; ok != false {
		summary, ok = params["Summary"].(string)
		if ok != true {
			return &postrepo.Post{}, ErrInvalidType
		}
	}

	tags, err := castTagsToString(params)
	if err != nil {
		return &postrepo.Post{}, err
	}

	p := &postrepo.Post{
		Title:   title,
		Slug:    slug,
		Summary: summary,
		Image:   image,
		HTML:    template.HTML(buf.String()),
		Tags:    tags,
	}

	p.SetPublished(published)

	return p, nil
}

func castTagsToString(params map[string]interface{}) ([]string, error) {
	if _, ok := params["Tags"]; ok == false {
		return []string{}, nil
	}

	tags, ok := params["Tags"].([]interface{})
	if ok == false {
		return []string{}, ErrInvalidType
	}

	strTags := make([]string, 0, len(tags))

	for _, tag := range tags {
		t, ok := tag.(string)
		if ok == false {
			continue
		}

		strTags = append(strTags, t)
	}

	return strTags, nil
}
