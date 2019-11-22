package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Database struct {
	Sections   []*Section   `yaml:"sections"`
	References []*Reference `yaml:"references"`
}

type Section struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type Reference struct {
	Title       string       `yaml:"title"`
	URL         string       `yaml:"url"`
	Author      string       `yaml:"author"`
	Note        string       `yaml:"note"`
	Section     string       `yaml:"section"`
	Supplements []Supplement `yaml:"supplement"`
}

type Supplement struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

func LoadDatabase(r io.Reader) (*Database, error) {
	d := yaml.NewDecoder(r)
	db := &Database{}
	if err := d.Decode(&db); err != nil {
		return nil, err
	}
	return db, nil
}
