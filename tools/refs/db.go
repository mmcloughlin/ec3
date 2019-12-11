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
	Title       string            `yaml:"title,omitempty"`
	URL         string            `yaml:"url,omitempty"`
	Author      string            `yaml:"author,omitempty"`
	Note        string            `yaml:"note,omitempty"`
	Section     string            `yaml:"section,omitempty"`
	Highlight   bool              `yaml:"highlight,omitempty"`
	Supplements []Supplement      `yaml:"supplements,omitempty"`
	Queued      bool              `yaml:"queued,omitempty"`
	ID          string            `yaml:"id,omitempty"`
	Type        string            `yaml:"type,omitempty"`
	Fields      map[string]string `yaml:"fields,omitempty"`
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

// StoreDatabase writes the database to w.
func StoreDatabase(w io.Writer, db *Database) error {
	e := yaml.NewEncoder(w)
	if err := e.Encode(db); err != nil {
		return err
	}
	return e.Close()
}
