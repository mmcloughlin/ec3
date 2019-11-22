package main

import (
	"context"
	"flag"
	"os"

	"github.com/mmcloughlin/ec3/internal/cli"

	"github.com/google/subcommands"
)

func main() {
	base := cli.NewBaseCommand("refs")
	subcommands.Register(&lint{Command: base}, "")
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
		ReferenceLinterFunc(AuthorPeriod),
		ReferenceLinterFunc(IACRCanonical),
		ReferenceLinterFunc(CheckNewlines),
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

func LoadDatabaseFile(filename string) (*Database, error) {
	_, f, err := cli.OpenInput(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadDatabase(f)
}
