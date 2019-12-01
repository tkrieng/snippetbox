package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/tkrieng/snippetbox/pkg/forms"
	"github.com/tkrieng/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for // any dynamic data that we want to pass to our HTML templates. // At the moment it only contains one field, but we'll add more // to it as the build progresses.
type templateData struct {
	CurrentYear int
	Flash       string
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
