package pages

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const templatesPath = "./static/templates"

func parseTemplates() (*template.Template, error) {
	templ := template.New("")

	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
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
