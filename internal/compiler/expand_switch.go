package compiler

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

// sqlcSwitchBranch is a single branch of a sqlc.switch(...) macro: a key that
// names the generated query variant and a SQL fragment that is spliced into the
// query in place of the whole sqlc.switch(...) call.
type sqlcSwitchBranch struct {
	key      string // "name_asc", or "else" for the default branch
	fragment string // author-authored SQL, e.g. "authors.name ASC"
}

// isSqlcFunc reports whether node is a call to sqlc.<name>.
func isSqlcFunc(node ast.Node, name string) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok || call.Func == nil {
		return false
	}
	return call.Func.Schema == "sqlc" && call.Func.Name == name
}

// stringConst extracts a string literal value from an A_Const node.
func stringConst(node ast.Node) (string, bool) {
	c, ok := node.(*ast.A_Const)
	if !ok {
		return "", false
	}
	s, ok := c.Val.(*ast.String)
	if !ok {
		return "", false
	}
	return s.Str, true
}

// camelize turns a branch key like "name_asc" into "NameAsc" so it can be
// appended to a query name and remain a valid Go identifier.
func camelize(s string) string {
	var b strings.Builder
	upper := true
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			upper = true
			continue
		}
		if upper {
			b.WriteRune(unicode.ToUpper(r))
			upper = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// stmtSource pairs a statement to compile with the source text its byte
// locations are relative to.
type stmtSource struct {
	raw *ast.RawStmt
	src string
}

// statementSources returns the statements to compile for a single parsed
// statement. Normally that is just the statement itself; if it contains a
// sqlc.switch(...) macro, it is the re-parsed branch variants it expands into.
func (c *Compiler) statementSources(raw *ast.RawStmt, src string) ([]stmtSource, error) {
	variants, err := c.expandSqlcSwitch(raw, src)
	if err != nil {
		return nil, err
	}
	if variants == nil {
		return []stmtSource{{raw: raw, src: src}}, nil
	}
	var sources []stmtSource
	for _, v := range variants {
		stmts, err := c.parser.Parse(strings.NewReader(v))
		if err != nil {
			return nil, err
		}
		for i := range stmts {
			sources = append(sources, stmtSource{raw: stmts[i].Raw, src: v})
		}
	}
	return sources, nil
}

// expandSqlcSwitch looks for a sqlc.switch(...) macro in a statement and, if
// present, returns one rewritten SQL string per branch. Each rewritten string
// is a normal query: the whole sqlc.switch(...) call is replaced by that
// branch's SQL fragment and the "-- name:" comment is renamed to
// <QueryName><BranchKey>. Each variant is re-parsed and analyzed as an ordinary
// query, so a bad column reference in a fragment is a compile error, and a
// generated name that collides with another query is caught by the normal
// duplicate-query-name check.
//
// The macro is recognized from the AST the same way sqlc.arg/sqlc.slice are
// (a FuncCall with schema "sqlc"), so it works in exactly the clauses where
// those macros parse. Returning (nil, nil) means there is no sqlc.switch and the
// statement should be compiled as-is.
func (c *Compiler) expandSqlcSwitch(raw *ast.RawStmt, src string) ([]string, error) {
	found := astutils.Search(raw, func(n ast.Node) bool { return isSqlcFunc(n, "switch") })
	if len(found.Items) == 0 {
		// Some parsers (e.g. SQLite for ORDER BY) silently discard a function
		// call they cannot place rather than erroring. If the text clearly
		// contains the macro but no node survived, fail loudly instead of
		// emitting the unexpanded call into the generated SQL.
		if stmtSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen); err == nil &&
			strings.Contains(stmtSQL, "sqlc.switch") {
			return nil, fmt.Errorf("sqlc.switch() is not supported in this position for this engine")
		}
		return nil, nil
	}
	if len(found.Items) > 1 {
		return nil, fmt.Errorf("only one sqlc.switch() per query is supported")
	}
	call := found.Items[0].(*ast.FuncCall)

	// sqlc.switch() is only allowed where it does not change the result shape
	// (WHERE, ORDER BY, ...), never in the SELECT projection: different branches
	// there could produce different columns and break the shared row type.
	if sel, ok := raw.Stmt.(*ast.SelectStmt); ok && sel.TargetList != nil {
		inTarget := astutils.Search(sel.TargetList, func(n ast.Node) bool { return isSqlcFunc(n, "switch") })
		if len(inTarget.Items) > 0 {
			return nil, fmt.Errorf("sqlc.switch() is not allowed in the SELECT list; use it in WHERE or ORDER BY")
		}
	}

	if call.Args == nil || len(call.Args.Items) < 2 {
		return nil, fmt.Errorf("sqlc.switch() requires a selector and at least one sqlc.when()/sqlc.else() branch")
	}

	// args[0] is the runtime selector (e.g. @sort). It plays no role in the
	// generated code (one function per branch, named by branch key) but is
	// required so the intent is explicit.
	branches := make([]sqlcSwitchBranch, 0, len(call.Args.Items)-1)
	seenElse := false
	for _, arg := range call.Args.Items[1:] {
		switch {
		case isSqlcFunc(arg, "when"):
			when := arg.(*ast.FuncCall)
			if when.Args == nil || len(when.Args.Items) != 2 {
				return nil, fmt.Errorf("sqlc.when() requires exactly 2 arguments: a key and a SQL fragment")
			}
			key, ok := stringConst(when.Args.Items[0])
			if !ok {
				return nil, fmt.Errorf("sqlc.when() key must be a string literal")
			}
			frag, ok := stringConst(when.Args.Items[1])
			if !ok {
				return nil, fmt.Errorf("sqlc.when() fragment must be a string literal")
			}
			branches = append(branches, sqlcSwitchBranch{key: key, fragment: frag})
		case isSqlcFunc(arg, "else"):
			if seenElse {
				return nil, fmt.Errorf("sqlc.switch() allows at most one sqlc.else()")
			}
			seenElse = true
			els := arg.(*ast.FuncCall)
			if els.Args == nil || len(els.Args.Items) != 1 {
				return nil, fmt.Errorf("sqlc.else() requires exactly 1 argument: a SQL fragment")
			}
			frag, ok := stringConst(els.Args.Items[0])
			if !ok {
				return nil, fmt.Errorf("sqlc.else() fragment must be a string literal")
			}
			branches = append(branches, sqlcSwitchBranch{key: "else", fragment: frag})
		default:
			return nil, fmt.Errorf("sqlc.switch() branches must be sqlc.when() or sqlc.else() calls")
		}
	}

	// Locate the byte span of the whole sqlc.switch(...) call within its
	// statement so it can be replaced textually with each branch fragment.
	stmtSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		return nil, err
	}
	switchStart := call.Location - raw.StmtLocation
	if switchStart < 0 || switchStart >= len(stmtSQL) {
		return nil, fmt.Errorf("could not locate sqlc.switch() in source")
	}
	switchEnd, err := matchParen(stmtSQL, switchStart)
	if err != nil {
		return nil, err
	}

	name, _, err := metadata.ParseQueryNameAndType(stmtSQL, metadata.CommentSyntax(c.parser.CommentSyntax()))
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("sqlc.switch() requires the query to have a -- name: annotation")
	}

	variants := make([]string, 0, len(branches))
	for _, br := range branches {
		spliced := stmtSQL[:switchStart] + br.fragment + stmtSQL[switchEnd+1:]
		// The plucked statement may exclude its trailing ";" (it can fall
		// outside StmtLen), so re-parsing a branch without one could yield an
		// empty statement. Normalize to exactly one terminator.
		spliced = strings.TrimRight(spliced, " \t\r\n;") + ";"
		newName := name + camelize(br.key)
		// Rename only the name comment. "name: <Name>" is shared by all comment
		// styles (-- /* #), so a single replacement is enough.
		spliced = strings.Replace(spliced, "name: "+name, "name: "+newName, 1)
		variants = append(variants, spliced)
	}
	return variants, nil
}

// matchParen returns the index of the ')' that closes the first '(' at or after
// start in s, skipping single-quoted string literals so parentheses inside a
// fragment (e.g. coalesce(x, 0)) do not throw off the depth count.
func matchParen(s string, start int) (int, error) {
	i := start
	for i < len(s) && s[i] != '(' {
		i++
	}
	if i >= len(s) {
		return 0, fmt.Errorf("could not locate opening parenthesis of sqlc.switch()")
	}
	depth := 0
	for ; i < len(s); i++ {
		switch s[i] {
		case '\'':
			// Advance to the closing quote, honoring '' escapes.
			for i++; i < len(s); i++ {
				if s[i] == '\'' {
					if i+1 < len(s) && s[i+1] == '\'' {
						i++
						continue
					}
					break
				}
			}
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return 0, fmt.Errorf("could not locate closing parenthesis of sqlc.switch()")
}
