package main

import (
	"text/template"

	"golang.org/x/xerrors"
)

//go:generate assets -pkg main -func loadtemplate -output ztemplates.go *.tmpl

// LoadOutputTypeTemplate returns a template for the given output type.
func LoadOutputTypeTemplate(name string) (*template.Template, error) {
	tmplname := name + ".tmpl"
	data, err := loadtemplate(tmplname)
	if err != nil {
		return nil, xerrors.Errorf("unknown output type %q", name)
	}

	t, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, xerrors.Errorf("template parse: %w", err)
	}

	return t, nil
}
