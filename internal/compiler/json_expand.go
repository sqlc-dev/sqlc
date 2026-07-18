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

// expandJSON rewrites sqlc's JSON directives into real Postgres functions,
// replacing only the qualified name span and leaving the argument list and
// any enclosing ARRAY(...) untouched:
//
//   - sqlc.jsonb_build_object."Name"(k, v, ...) -> jsonb_build_object(k, v, ...)
//   - sqlc.embed.jsonb(table)                   -> to_jsonb(table)
func (c *Compiler) expandJSON(raw *ast.RawStmt, query string) ([]source.Edit, error) {
	calls := astutils.Search(raw, func(node ast.Node) bool {
		return rewrite.IsJSONCall(node) || rewrite.IsEmbedJSONCall(node)
	})
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

		replacement := "jsonb_build_object"
		if rewrite.IsEmbedJSONCall(call) {
			replacement = "to_jsonb"
		}

		edits = append(edits, source.Edit{
			Location: loc,
			OldFunc:  func(string) int { return length },
			New:      replacement,
		})
	}

	return edits, nil
}
