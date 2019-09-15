package curve

import (
	"crypto/elliptic"
	"strings"

	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/tmpl"
)

//go:generate assets -pkg curve -func loadtemplate -output ztemplates.go tmpl/shortw/*.go

var templates = tmpl.Environment{
	Loader: tmpl.NewBasePath(tmpl.LoaderFunc(loadtemplate), "tmpl/shortw"),
}

type ShortWeierstrass struct {
	PackageName string
	Params      *elliptic.CurveParams
	ShortName   string
}

func (c ShortWeierstrass) Generate() (gen.Files, error) {
	fs := gen.Files{}

	// Build template package.
	pkg, err := templates.Package(
		"curve.go",
		"curve_test.go",
		"recode.go",
		"recode_test.go",
		"util_test.go",
	)
	if err != nil {
		return nil, err
	}

	typename := strings.ToUpper(c.ShortName)
	varname := strings.ToLower(c.ShortName)

	err = pkg.Apply(
		tmpl.GeneratedBy(gen.GeneratedBy),
		tmpl.SetPackageName(c.PackageName),
		tmpl.Rename("CURVENAME", typename),
		tmpl.CommentReplace("CURVENAME", typename),
		tmpl.CommentReplace("CanonicalName", c.Params.Name),

		tmpl.Rename("curvename", varname),

		tmpl.DefineString("ConstCanonicalName", c.Params.Name),
		tmpl.DefineString("ConstPDecimal", c.Params.P.Text(10)),
		tmpl.DefineString("ConstNDecimal", c.Params.N.Text(10)),
		tmpl.DefineString("ConstBHex", c.Params.B.Text(16)),
		tmpl.DefineString("ConstGxHex", c.Params.Gx.Text(16)),
		tmpl.DefineString("ConstGyHex", c.Params.Gy.Text(16)),
		tmpl.DefineIntDecimal("ConstBitSize", c.Params.BitSize),

		tmpl.DefineIntDecimal("ConstW", 6),

		tmpl.DefineIntDecimal("ConstNumTrials", 128),
	)
	if err != nil {
		return nil, err
	}

	for filename, t := range pkg.Templates() {
		src, err := t.Bytes()
		if err != nil {
			return nil, err
		}

		fs.Add(filename, src)
	}

	return fs, nil
}
