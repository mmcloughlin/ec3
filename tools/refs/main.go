package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/mmcloughlin/ec3/internal/cli"
)

func main() {
	base := cli.NewBaseCommand("refs")
	subcommands.Register(&gen{Command: base}, "")
	subcommands.Register(&bib{Command: base}, "")
	subcommands.Register(&lint{Command: base}, "")
	subcommands.Register(&linkcheck{Command: base}, "")
	subcommands.Register(&format{Command: base}, "")
	subcommands.Register(subcommands.HelpCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

// lint subcommand.
type lint struct {
	cli.Command
}

func (*lint) Name() string     { return "lint" }
func (*lint) Synopsis() string { return "lint a references database" }
func (*lint) Usage() string {
	return `Usage: lint <filename>

Lint a references database.

`
}

func (*lint) Linter() Linter {
	return ConcatLinter(
		ReferenceLinterFunc(RequireURL),
		ReferenceLinterFunc(RequireTitle),
		ReferenceLinterFunc(ValidURL),
		ReferenceLinterFunc(IACRCanonical),
		ReferenceLinterFunc(HALCanonical),
		ReferenceLinterFunc(CheckNewlines),
		DisallowHost("drive.google.com"),
		LinterFunc(DuplicateURLs),
		LinterFunc(CheckSectionTags),
	)
}

func (cmd *lint) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	db, err := LoadDatabaseFile(f.Arg(0))
	if err != nil {
		return cmd.Error(err)
	}

	linter := cmd.Linter()
	errs := linter.Lint(db)
	if len(errs) == 0 {
		return subcommands.ExitSuccess
	}

	for _, err := range errs {
		cmd.Log.Print(err)
	}

	return subcommands.ExitFailure
}

// gen subcommand.
type gen struct {
	cli.Command

	outputtype string
	outputfile string
}

func (*gen) Name() string     { return "gen" }
func (*gen) Synopsis() string { return "generate output from references" }
func (*gen) Usage() string {
	return `Usage: gen [-type <type>] <filename>

Generate output from references.

`
}

func (cmd *gen) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.outputtype, "type", "markdown", "output type")
	f.StringVar(&cmd.outputfile, "out", "", "output file (default to stdout)")
}

func (cmd *gen) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Load database.
	db, err := LoadDatabaseFile(f.Arg(0))
	if err != nil {
		return cmd.Error(err)
	}

	// Load template.
	t, err := LoadOutputTypeTemplate(cmd.outputtype)
	if err != nil {
		return cmd.Error(err)
	}

	// Execute.
	_, w, err := cli.OpenOutput(cmd.outputfile)
	if err != nil {
		return cmd.Error(err)
	}
	defer w.Close()

	err = t.Execute(w, db)
	if err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}

// bib subcommand.
type bib struct {
	cli.Command

	outputfile string
}

func (*bib) Name() string     { return "bib" }
func (*bib) Synopsis() string { return "generate bibtex file from references" }
func (*bib) Usage() string {
	return `Usage: bib [-out <output>] <filename>

Generate bibtex from references.

`
}

func (cmd *bib) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.outputfile, "out", "", "output file (default to stdout)")
}

func (cmd *bib) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Load database.
	db, err := LoadDatabaseFile(f.Arg(0))
	if err != nil {
		return cmd.Error(err)
	}

	// Open output.
	_, w, err := cli.OpenOutput(cmd.outputfile)
	if err != nil {
		return cmd.Error(err)
	}
	defer w.Close()

	// Generate.
	if err := WriteBibTeX(w, db); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}

// linkcheck subcommand.
type linkcheck struct {
	cli.Command
}

func (*linkcheck) Name() string     { return "linkcheck" }
func (*linkcheck) Synopsis() string { return "check whether all urls exist" }
func (*linkcheck) Usage() string {
	return `Usage: linkcheck <filename>

Check whether all URLs in the database exist.

`
}

func (cmd *linkcheck) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Load database.
	db, err := LoadDatabaseFile(f.Arg(0))
	if err != nil {
		return cmd.Error(err)
	}

	// Check all URLs.
	status := subcommands.ExitSuccess
	for _, reference := range db.References {
		if err := CheckLink(reference.URL); err != nil {
			cmd.Log.Printf("url %s failed: %s", reference.URL, err)
			status = subcommands.ExitFailure
		}
	}

	return status
}

// format subcommand.
type format struct {
	cli.Command

	write bool
}

func (*format) Name() string     { return "fmt" }
func (*format) Synopsis() string { return "format database" }
func (*format) Usage() string {
	return `Usage: fmt [-w] <filename>

Format references database.

`
}

func (cmd *format) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.write, "w", false, "write result back to source file")
}

func (cmd *format) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	filename := f.Arg(0)

	// Load database.
	db, err := LoadDatabaseFile(filename)
	if err != nil {
		return cmd.Error(err)
	}

	// Open output.
	output := ""
	if cmd.write {
		output = filename
	}
	_, w, err := cli.OpenOutput(output)
	if err != nil {
		return cmd.Error(err)
	}
	defer w.Close()

	// Write database back.
	if err := StoreDatabase(w, db); err != nil {
		return cmd.Error(err)
	}

	return subcommands.ExitSuccess
}

// LoadDatabaseFile loads the database at filename, falling back to standard
// input if the filename is empty.
func LoadDatabaseFile(filename string) (*Database, error) {
	_, f, err := cli.OpenInput(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadDatabase(f)
}
