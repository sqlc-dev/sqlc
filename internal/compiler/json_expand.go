package compiler

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
)

func findOpenParen(s string, from int) (int, error) {
	idx := strings.IndexByte(s[from:], '(')
	if idx < 0 {
		return 0, fmt.Errorf("no opening parenthesis found from location %d", from)
	}
	return from + idx, nil
}

// expandJSON rewrites sqlc.jsonb_build_object."Name"(k, v, ...) into
// jsonb_build_object(k, v, ...), replacing only the qualified name span and
// leaving the argument list and any enclosing ARRAY(...) untouched.
func (c *Compiler) expandJSON(raw *ast.RawStmt, query string) ([]source.Edit, error) {
	calls := astutils.Search(raw, rewrite.IsJSONCall)
	if len(calls.Items) == 0 {
		return nil, nil
	}

	var edits []source.Edit
	for _, item := range calls.Items {
		call := item.(*ast.FuncCall)
		loc := call.Location - raw.StmtLocation

		openParen, err := findOpenParen(query, loc)
		if err != nil {
			return nil, err
		}
		length := openParen - loc

		edits = append(edits, source.Edit{
			Location: loc,
			OldFunc:  func(string) int { return length },
			New:      "jsonb_build_object",
		})
	}

	return edits, nil
}
