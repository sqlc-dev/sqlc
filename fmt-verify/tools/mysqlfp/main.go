// mysqlfp parses two MySQL query files with marino (the parser dolphin
// wraps) and compares their statements' canonical restore forms, with
// identifier case preserved. Redundant parentheses are unwrapped before
// comparing, since dropping them is a safe formatting change.
// Exit 1 + a report on any semantic drift.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sqlc-dev/marino/ast"
	"github.com/sqlc-dev/marino/format"
	"github.com/sqlc-dev/marino/parser"
)

type deparen struct{}

func (v deparen) Enter(n ast.Node) (ast.Node, bool) { return n, false }
func (v deparen) Leave(n ast.Node) (ast.Node, bool) {
	if p, ok := n.(*ast.ParenthesesExpr); ok {
		return p.Expr, true
	}
	return n, true
}

func canonical(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := parser.New()
	stmts, _, err := p.Parse(string(b), "", "")
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	var out []string
	flags := format.RestoreStringSingleQuotes | format.RestoreKeyWordUppercase | format.RestoreNameBackQuotes
	for _, s := range stmts {
		n, _ := s.Accept(deparen{})
		var sb strings.Builder
		if err := n.Restore(format.NewRestoreCtx(flags, &sb)); err != nil {
			return nil, fmt.Errorf("restore: %w", err)
		}
		out = append(out, sb.String())
	}
	return out, nil
}

func main() {
	a, err := canonical(os.Args[1])
	if err != nil {
		fmt.Printf("UNPARSEABLE %s: %s\n", os.Args[1], err)
		os.Exit(2)
	}
	b, err := canonical(os.Args[2])
	if err != nil {
		fmt.Printf("UNPARSEABLE %s: %s\n", os.Args[2], err)
		os.Exit(2)
	}
	bad := false
	if len(a) != len(b) {
		fmt.Printf("statement count %d != %d\n", len(a), len(b))
		bad = true
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			fmt.Printf("stmt %d differs:\n  orig: %s\n  fmt:  %s\n", i+1, a[i], b[i])
			bad = true
		}
	}
	if bad {
		os.Exit(1)
	}
}
