package db

import (
	"bufio"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/parse"
)

type Key struct {
	Path           string
	Collection     string
	Class          string
	Section        string
	Shape          string
	Representation string
	Operation      string
	Name           string
	Ext            string
}

func KeyFromFilename(filename string) Key {
	k := Key{Path: filename}
	k.Ext = filepath.Ext(filename)
	path := strings.TrimSuffix(filename, k.Ext)
	parts := strings.Split(path, "/")
	n := len(parts)
	dst := []*string{&k.Collection, &k.Class, &k.Section, &k.Shape, &k.Representation, &k.Operation}
	for i := 0; i < n-1 && i < len(dst); i++ {
		*dst[i] = parts[i]
	}
	k.Name = parts[n-1]
	return k
}

func (k Key) IsOP3() bool { return k.Ext == ".op3" }

func (k Key) IsFormula() bool { return k.Operation != "" }

func (k Key) IsShape() bool { return k.Shape != "" && k.Name == "coordinates" }

func (k Key) IsRepresentation() bool { return k.Representation != "" && k.Name == "variables" }

func (k Key) ShapeID() string { return path.Join(k.Class, k.Shape) }

func (k Key) RepresentationID() string { return path.Join(k.ShapeID(), k.Representation) }

func (k Key) OperationID() string { return path.Join(k.RepresentationID(), k.Operation) }

func (k Key) FormulaID() string { return path.Join(k.OperationID(), k.Name) }

// Collection specifies information about a collection of data included in the
// database.
type Collection struct {
	FormulaURL string
}

// DefaultCollection specifies information about the original EFD database.
var DefaultCollection = Collection{
	// Example: https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
	FormulaURL: "https://hyperelliptic.org/EFD/CLASS/auto-SHAPE-REPRESENTATION.html#OPERATION-TAG",
}

type Database struct {
	Collections     map[string]*Collection
	Shapes          map[string]*efd.Shape
	Representations map[string]*efd.Representation
	Formulae        map[string]*efd.Formula
}

func New() *Database {
	return &Database{
		Collections:     map[string]*Collection{},
		Shapes:          map[string]*efd.Shape{},
		Representations: map[string]*efd.Representation{},
		Formulae:        map[string]*efd.Formula{},
	}
}

func (d Database) collection(k string) *Collection {
	if _, ok := d.Collections[k]; !ok {
		return &DefaultCollection
	}
	return d.Collections[k]
}

func (d Database) shape(k string) *efd.Shape {
	if _, ok := d.Shapes[k]; !ok {
		d.Shapes[k] = &efd.Shape{
			ID: k,
		}
	}
	return d.Shapes[k]
}

func (d Database) representation(k string) *efd.Representation {
	if _, ok := d.Representations[k]; !ok {
		d.Representations[k] = &efd.Representation{
			ID: k,
		}
	}
	return d.Representations[k]
}

func (d Database) formula(k string) *efd.Formula {
	if _, ok := d.Formulae[k]; !ok {
		d.Formulae[k] = &efd.Formula{
			ID: k,
		}
	}
	return d.Formulae[k]
}

func (d Database) finalize() {
	// Set formula URLs.
	for _, f := range d.Formulae {
		if f.URL != "" {
			continue
		}
		r := strings.NewReplacer(
			"CLASS", f.Class,
			"SHAPE", f.Shape.Tag,
			"REPRESENTATION", f.Representation.Tag,
			"OPERATION", f.Operation,
			"TAG", f.Tag,
		)
		tmpl := d.collection(f.Collection).FormulaURL
		f.URL = r.Replace(tmpl)
	}
}

func Read(s Store) (*Database, error) {
	p := parser{
		DB: New(),
	}
	if err := Walk(s, p); err != nil {
		return nil, err
	}
	p.DB.finalize()
	return p.DB, nil
}

type parser struct {
	DB *Database
}

func (p parser) Visit(f File) error {
	k := KeyFromFilename(f.Path())
	switch {
	case k.IsOP3():
		return p.op3(k, f)
	case k.IsShape():
		return p.shape(k, f)
	case k.IsRepresentation():
		return p.representation(k, f)
	case k.IsFormula():
		return p.formula(k, f)
	case k.Name == "metadata":
		return p.metadata(k, f)
	}
	return nil
}

func (p parser) metadata(k Key, r io.Reader) error {
	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	c := &Collection{}
	for prop, vs := range props {
		switch prop {
		case "formulaurl":
			c.FormulaURL, err = atmostone(prop, vs)
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	p.DB.Collections[k.Collection] = c
	return nil
}

func (p parser) formula(k Key, r io.Reader) error {
	f := p.DB.formula(k.FormulaID())
	f.Collection = k.Collection
	f.Tag = k.Name
	f.Class = k.Class
	f.Shape = p.DB.shape(k.ShapeID())
	f.Representation = p.DB.representation(k.RepresentationID())
	f.Operation = k.Operation

	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	for prop, vs := range props {
		switch prop {
		case "source":
			f.Source, err = atmostone(prop, vs)
		case "url":
			f.URL, err = atmostone(prop, vs)
		case "compute":
			f.Compute = vs
		case "assume":
			f.Assume = vs
		case "parameter":
			f.Parameters = vs
		case "appliesto":
			f.AppliesTo, err = atmostone(prop, vs)
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p parser) representation(k Key, r io.Reader) error {
	repr := p.DB.representation(k.RepresentationID())
	repr.Collection = k.Collection
	repr.Tag = k.Representation
	repr.Class = k.Class
	repr.Shape = p.DB.shape(k.ShapeID())

	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	for prop, vs := range props {
		switch prop {
		case "name":
			repr.Name, err = exactlyone(prop, vs)
		case "assume":
			repr.Assume = vs
		case "parameter":
			repr.Parameters = vs
		case "variable":
			repr.Variables = vs
		case "satisfying":
			repr.Satisfying = vs
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p parser) shape(k Key, r io.Reader) error {
	s := p.DB.shape(k.ShapeID())
	s.Collection = k.Collection
	s.Tag = k.Shape
	s.Class = k.Class

	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	for prop, vs := range props {
		switch prop {
		case "name":
			s.Name, err = exactlyone(prop, vs)
		case "coordinate":
			s.Coordinates = vs
		case "a0", "a1", "a2", "a3", "a4", "a6":
			i := int(prop[1] - '0')
			s.A[i], err = exactlyone(prop, vs)
		case "satisfying":
			s.Satisfying = vs
		case "parameter":
			s.Parameters = vs
		case "addition":
			s.Addition = vs
		case "doubling":
			s.Doubling = vs
		case "negation":
			s.Negation = vs
		case "neutral":
			s.Neutral = vs
		case "fromweierstrass":
			s.FromWeierstrass = vs
		case "toweierstrass":
			s.ToWeierstrass = vs
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p parser) op3(k Key, r io.Reader) error {
	prog, err := parse.Reader(k.Path, r)
	// Note we expect some files to fail parsing, so we suppress errors here.
	if err != nil {
		return nil
	}

	f := p.DB.formula(k.FormulaID())
	f.Program = prog
	return nil
}

var whitespace = regexp.MustCompile("[[:space:]]+")

func ParseProperties(r io.Reader) (map[string][]string, error) {
	properties := map[string][]string{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		parts := whitespace.Split(line, 2)
		if len(parts) != 2 || parts[0] == "" {
			continue
		}
		key, value := parts[0], parts[1]
		properties[key] = append(properties[key], value)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return properties, nil
}

func exactlyone(prop string, vs []string) (string, error) {
	if len(vs) != 1 {
		return "", xerrors.Errorf("expected exactly one value for %q", prop)
	}
	return vs[0], nil
}

func atmostone(prop string, vs []string) (string, error) {
	if len(vs) > 1 {
		return "", xerrors.Errorf("expected at most one value for %q", prop)
	}
	if len(vs) == 1 {
		return vs[0], nil
	}
	return "", nil
}
