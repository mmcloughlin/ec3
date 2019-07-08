package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/printer"

	"github.com/mmcloughlin/ec3/asm"
)

// GeneratedBy is the name used in code generation warnings.
const GeneratedBy = "ec3"

type File struct {
	Path   string
	Source []byte
}

type Files []File

func (f *Files) Add(path string, src []byte) {
	*f = append(*f, File{
		Path:   path,
		Source: src,
	})
}

func (f *Files) CompileAsm(pkg, prefix string, ctx *build.Context) error {
	s, err := asm.Compile(ctx)
	if err != nil {
		return err
	}

	cfg := printer.Config{
		Pkg:  pkg,
		Name: GeneratedBy,
	}

	// Stubs.
	stubs, err := printer.NewStubs(cfg).Print(s)
	if err != nil {
		return err
	}
	f.Add(prefix+".go", stubs)

	// Assembly.
	goasm, err := printer.NewGoAsm(cfg).Print(s)
	if err != nil {
		return err
	}
	f.Add(prefix+".s", goasm)

	return nil
}

// Output writes bundled files to disk rooted at path. Directories are created as necessary.
func (f Files) Output(path string) error {
	for _, file := range f {
		// Ensure directory exists.
		filename := filepath.Join(path, file.Path)
		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}

		// Write file.
		if err := ioutil.WriteFile(filename, file.Source, 0644); err != nil {
			return err
		}
	}
	return nil
}
