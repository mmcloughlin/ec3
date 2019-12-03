package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Database is a collection of references categorized into sections.
type Database struct {
	Sections   []*Section   `yaml:"sections"`
	References []*Reference `yaml:"references"`
}

// Section is a reference category.
type Section struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

// Reference describes a referenced resource.
type Reference struct {
	Title       string       `yaml:"title"`
	URL         string       `yaml:"url"`
	Author      string       `yaml:"author"`
	Note        string       `yaml:"note"`
	Section     string       `yaml:"section"`
	Highlight   bool         `yaml:"highlight"`
	Supplements []Supplement `yaml:"supplements"`
	Queued      bool         `yaml:"queued"`
}

// Supplement is another resource associated with a reference. For example, in
// the case of a paper, supplements might be slides for an associated talk, or a
// link to the code.
type Supplement struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

// LoadDatabase loads a reference database in YAML format.
func LoadDatabase(r io.Reader) (*Database, error) {
	d := yaml.NewDecoder(r)
	db := &Database{}
	if err := d.Decode(&db); err != nil {
		return nil, err
	}
	return db, nil
}
