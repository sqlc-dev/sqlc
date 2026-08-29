package dolphin

import (
	"fmt"
	"strings"

	pcast "github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/parser"
)

// deparen unwraps parenthesized expressions so a redundant pair of
// parentheses added or dropped by formatting does not change the
// fingerprint. Everything else about the tree survives untouched.
type deparen struct{}

func (v deparen) Enter(n pcast.Node) (pcast.Node, bool) { return n, false }
func (v deparen) Leave(n pcast.Node) (pcast.Node, bool) {
	if p, ok := n.(*pcast.ParenthesesExpr); ok {
		return p.Expr, true
	}
	return n, true
}

// Fingerprint reduces a statement to a canonical form that survives changes
// in whitespace, keyword case, quoting style and redundant parentheses —
// and nothing else. Identifier case is preserved: MySQL table names are
// case-sensitive on most servers, so a formatting pass that changes one is
// a semantic change and must not fingerprint equal. fmt uses this as its
// proof that a formatted statement still means what the author wrote; a
// statement whose fingerprint it cannot match is left exactly as written.
func (p *Parser) Fingerprint(sql string) (string, error) {
	// A fresh parser: p.pingcap carries per-parse comment state that
	// ParseFile is still using.
	stmts, _, err := parser.New().Parse(sql, "", "")
	if err != nil {
		return "", normalizeErr(err)
	}
	flags := format.RestoreStringSingleQuotes |
		format.RestoreKeyWordUppercase |
		format.RestoreNameBackQuotes
	parts := make([]string, 0, len(stmts))
	for _, s := range stmts {
		n, ok := s.Accept(deparen{})
		if !ok {
			return "", fmt.Errorf("fingerprint: rewrite failed")
		}
		var sb strings.Builder
		if err := n.Restore(format.NewRestoreCtx(flags, &sb)); err != nil {
			return "", fmt.Errorf("fingerprint: %w", err)
		}
		parts = append(parts, sb.String())
	}
	return strings.Join(parts, "; "), nil
}
