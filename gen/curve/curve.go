package curve

import (
	"crypto/elliptic"
	"strings"

	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/tmpl"
)

//go:generate assets -pkg curve -func loadtemplate -output ztemplates.go tmpl/shortw/curve.go

var templates = tmpl.Environment{
	Loader: tmpl.LoaderFunc(loadtemplate),
}

type ShortWeierstrass struct {
	PackageName string
	Params      *elliptic.CurveParams
	ShortName   string
}

func (c ShortWeierstrass) Generate() (gen.Files, error) {
	fs := gen.Files{}

	t, err := templates.Load("tmpl/shortw/curve.go")
	if err != nil {
		return nil, err
	}

	typename := strings.ToUpper(c.ShortName)
	varname := strings.ToLower(c.ShortName)

	err = t.Apply(
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
	)
	if err != nil {
		return nil, err
	}

	src, err := t.Bytes()
	if err != nil {
		return nil, err
	}

	fs.Add("curve.go", src)

	return fs, nil
}
