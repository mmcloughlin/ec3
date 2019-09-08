package op3

import "github.com/mmcloughlin/ec3/efd/op3/ast"

// AliasCorrect ensures that the program p will compute the same result if the
// given variables are aliased. Aliased variables are pointers to the same
// memory location. If the program is already robust to the given alias sets, no
// changes will be made. Otherwise the program will be modified with additional
// temporaries.
func AliasCorrect(p *ast.Program, aliases [][]ast.Variable, outputs []ast.Variable, vars VariableGenerator) *ast.Program {
	isinput := InputSet(p)
	isoutput := VariableSet(outputs)

	// We'll need the interference graph.
	g := BuildInterferenceGraph(p, outputs)

	// Populate the variable generator.
	vars.MarkUsed(Variables(p)...)

	replacements := map[ast.Variable]ast.Variable{}
	pre := []ast.Assignment{}
	post := []ast.Assignment{}

	for _, set := range aliases {
		for _, v := range set {
			// If v is not written to, we can leave it as is.
			if ReadOnly(p, v) {
				continue
			}

			// Does v interfere with anything in the alias set?
			interfere := false
			for _, u := range set {
				interfere = interfere || g.Interfere(v, u)
			}

			if !interfere {
				continue
			}

			// Replace v with a new variable.
			r := vars.New()
			replacements[v] = r

			if isinput[v] {
				pre = append(pre, ast.Assignment{LHS: r, RHS: v})
			}

			if isoutput[v] {
				post = append(post, ast.Assignment{LHS: v, RHS: r})
			}
		}
	}

	if len(replacements) == 0 {
		return p
	}

	// Build new program.
	q := RenameVariables(p, replacements)
	q.Assignments = append(pre, q.Assignments...)
	q.Assignments = append(q.Assignments, post...)

	return q
}
