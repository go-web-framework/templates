package templates

import (
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const partialsDir = "partials"

// Set is a collection of templates.
type Set struct {
	FuncMap     template.FuncMap
	DefaultArgs map[string]interface{}
	Templates   map[string]*template.Template
}

var ErrNoSuchTemplate = errors.New("templates: no matching template for name")

func (s *Set) execute(name string, w io.Writer, args interface{}) error {
	t, ok := s.Templates[name]
	if !ok {
		return ErrNoSuchTemplate
	}

	a := args
	if m, ok := args.(map[string]interface{}); ok {
		a = normalize(s.DefaultArgs, m)
	}

	return t.Execute(w, a)
}

func (s *Set) Execute(name string, w io.Writer, args interface{}) error {
	return s.execute(name, w, args)
}

func normalize(def, new map[string]interface{}) map[string]interface{} {
	var ret map[string]interface{}

	for k, v := range def {
		if ret == nil {
			ret = make(map[string]interface{})
		}
		ret[k] = v
	}

	for k, v := range new {
		if ret == nil {
			ret = make(map[string]interface{})
		}
		ret[k] = v
	}

	return ret
}

func (s *Set) Parse(path string) error {
	m, err := readDir(path)
	if err != nil {
		return err
	}

	// Standardize path separators,
	for k, v := range m {
		delete(m, k)
		m[filepath.ToSlash(k)] = v
	}

	var partials []string
	for k, v := range m {
		if strings.HasPrefix(k, partialsDir+"/") {
			partials = append(partials, string(v))
		}
	}

	// The partials should be parsed with each main template
	// to be available in the main template.
	for k, v := range m {
		if strings.HasPrefix(k, partialsDir+"/") {
			continue
		}

		templ, err := template.New(k).Funcs(s.FuncMap).Parse(string(v))
		if err != nil {
			return err
		}
		for _, contents := range partials {
			if _, err := templ.Parse(contents); err != nil {
				return err
			}
		}

		if s.Templates == nil {
			s.Templates = make(map[string]*template.Template)
		}
		s.Templates[k] = templ
	}

	return nil
}

func readDir(root string) (map[string][]byte, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, err
	}

	var m map[string][]byte // Lazily initialized in Walk.

	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		relp, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}

		if m == nil {
			m = make(map[string][]byte)
		}
		m[relp] = b

		return nil
	})

	return m, err
}