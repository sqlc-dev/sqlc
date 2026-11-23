package preprocess

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/constants"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// A macro is a statement annotated ":macro". It is not a query: it generates
// no code and is never parsed on its own. Its body is a piece of SQL that is
// pasted, in parentheses, wherever another query calls sqlc.macro.<Name>().
//
//	-- name: ActiveAuthors :macro
//	-- @param org bigint
//	SELECT * FROM authors WHERE org_id = sqlc.arg(org) AND deleted_at IS NULL;
//
//	-- name: ListActiveAuthors :many
//	SELECT a.id, a.name FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a;
type macro struct {
	name     string
	filename string
	params   []macroParam
	// body is the statement without its comments, its terminator or its line
	// breaks, so pasting it never adds a line to the query it lands in.
	body string
	// spansLines is set when the body could not be put on one line.
	spansLines bool
}

// macroParam is a parameter declared with "-- @param name [type]". A typed
// parameter casts whatever it is given to that type.
type macroParam struct {
	name string
	typ  string
}

// Macros is the set of macros visible to a group of query files. Macros are
// expanded lazily — only where a query calls one — so a definition may sit in
// any file of the group, before or after its callers.
type Macros struct {
	d    Dialect
	ok   bool
	defs map[string]*macro
}

// NewMacros returns an empty set of macros for an engine. An engine that is not
// preprocessed has no sqlc syntax, and so no macros: its set stays empty and
// File returns the source unchanged.
func NewMacros(engine config.Engine) *Macros {
	d, ok := DialectFor(engine)
	return newMacros(d, ok)
}

func newMacros(d Dialect, ok bool) *Macros {
	return &Macros{d: d, ok: ok, defs: map[string]*macro{}}
}

// commentSyntax is how the dialect writes the comments annotations live in.
func (d Dialect) commentSyntax() source.CommentSyntax {
	return source.CommentSyntax{Dash: true, SlashStar: true, Hash: d.HashComment}
}

// fold returns the form of a parameter name that references are matched by.
func (d Dialect) fold(name string) string {
	if d.FoldIdentifier {
		return strings.ToLower(name)
	}
	return name
}

// macroName returns the name of the macro a statement defines, and whether it
// defines one at all.
func (d Dialect) macroName(stmt string) (string, bool) {
	name, cmd, err := metadata.ParseQueryNameAndType(stmt, metadata.CommentSyntax(d.commentSyntax()))
	if err != nil || cmd != metadata.CmdMacro {
		return "", false
	}
	return name, true
}

// Collect records every macro defined in src. The errors it returns are
// problems with the definitions themselves, located in src; a definition that
// has one is still recorded, so its callers are not reported a second time for
// calling a macro that does not exist.
func (m *Macros) Collect(filename, src string) []*sqlerr.Error {
	if !m.ok {
		return nil
	}
	var errs []*sqlerr.Error
	l := &lexer{d: m.d, src: src}
	for _, span := range l.statements() {
		start, end := span[0], span[1]
		stmt := src[start:end]
		name, ok := m.d.macroName(stmt)
		if !ok {
			continue
		}
		fail := func(format string, args ...any) {
			errs = append(errs, &sqlerr.Error{
				Message:  fmt.Sprintf("macro %q ", name) + fmt.Sprintf(format, args...),
				Location: start,
			})
		}
		if prev, exists := m.defs[name]; exists {
			if prev.filename != filename {
				fail("is already defined in %s", prev.filename)
			} else {
				fail("is already defined")
			}
			continue
		}

		def := &macro{name: name, filename: filename}
		m.defs[name] = def

		comments, err := source.CleanedComments(stmt, m.d.commentSyntax())
		if err != nil {
			fail("has unreadable comments: %s", err)
			continue
		}
		seen := map[string]bool{}
		for _, line := range comments {
			fields := strings.Fields(line)
			if len(fields) == 0 || fields[0] != constants.QueryFlagParam {
				continue
			}
			if len(fields) == 1 {
				fail("has a @param without a name")
				continue
			}
			param := macroParam{name: m.d.fold(fields[1]), typ: strings.Join(fields[2:], " ")}
			if seen[param.name] {
				fail("declares the parameter %q twice", fields[1])
				continue
			}
			seen[param.name] = true
			def.params = append(def.params, param)
		}

		body, ok := flatten(m.d, stmt)
		def.body = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(body), ";"))
		def.spansLines = !ok
		if def.body == "" {
			fail("has no body")
		}
	}
	return errs
}

// File rewrites src like the package-level File does, after expanding every
// call to a macro in the set and blanking out the definitions src holds.
func (m *Macros) File(src string) *Result {
	if !m.ok {
		return &Result{Text: src}
	}
	text, edits, blanks, errs := m.expandFile(src)
	res := rewrite(m.d, text, errs)
	res.macroEdits = edits
	res.blanks = blanks
	return res
}

// Expand returns src with its macro definitions blanked and its macro calls
// expanded, and the rest of its sqlc syntax as it was written. It is for a
// reader that takes sqlc syntax as it comes but has no way to take a macro: a
// schema parser, when a schema file holds queries as well. A statement whose
// macros cannot be expanded is returned as it was written; File reports why.
func (m *Macros) Expand(src string) string {
	if !m.ok || len(m.defs) == 0 {
		return src
	}
	text, _, _, _ := m.expandFile(src)
	return text
}

// CallsMacro reports whether text holds a call to a macro, whether or not the
// macro exists. A tool that rewrites queries without expanding them has to
// leave such a statement alone: the name of a macro is matched exactly, and a
// parser that reads the call as a function is free to fold it.
func CallsMacro(engine config.Engine, text string) bool {
	d, ok := DialectFor(engine)
	if !ok {
		return false
	}
	l := &lexer{d: d, src: text}
	for i := 0; i < len(text); {
		if stop := l.skip(i); stop >= 0 {
			if stop <= i {
				stop = i + 1
			}
			i = stop
			continue
		}
		if !isIdentStart(text[i]) {
			i++
			continue
		}
		j := l.identEnd(i)
		if strings.EqualFold(text[i:j], "sqlc") {
			// A call that is not well formed is still a call.
			if _, ok, err := l.macroCall(i, j); ok || err != nil {
				return true
			}
		}
		i = j
	}
	return false
}

// expandFile returns src with its macro definitions blanked and its macro
// calls expanded, the spans it replaced, the spans it blanked, and one error —
// located in the returned text — for every statement it had to leave alone.
func (m *Macros) expandFile(src string) (string, []edit, [][2]int, []*sqlerr.Error) {
	var (
		b      strings.Builder
		edits  []edit
		blanks [][2]int
		errs   []*sqlerr.Error
	)
	b.Grow(len(src))

	l := &lexer{d: m.d, src: src}
	for _, span := range l.statements() {
		start, end := span[0], span[1]
		stmt := src[start:end]

		// A definition is not SQL the engine should see. It is blanked rather
		// than cut so every other statement keeps its offset and its line.
		if _, ok := m.d.macroName(stmt); ok {
			b.WriteString(blank(stmt))
			blanks = append(blanks, [2]int{start, end})
			continue
		}

		var local []edit
		out, err := m.expand(stmt, nil, nil, &local)
		if err != nil {
			// Like any other invalid statement, this one is copied through as
			// written and the error is reported against it.
			err.Location += b.Len()
			errs = append(errs, err)
			b.WriteString(stmt)
			continue
		}
		for _, e := range local {
			edits = append(edits, edit{
				oldStart: start + e.oldStart, oldEnd: start + e.oldEnd,
				newStart: b.Len() + e.newStart, newEnd: b.Len() + e.newEnd,
			})
		}
		b.WriteString(out)
	}
	return b.String(), edits, blanks, errs
}

// expand copies text, replacing each macro call with its expansion. Inside a
// macro's body bind is non-nil, and maps each of the macro's parameters to
// the argument it stands for; everywhere else sqlc.arg() is left for the
// placeholder rewrite that follows.
//
// stack names the macros being expanded, outermost first. record, when set,
// collects the span of each call replaced, relative to text and to the result.
// An error is located in text.
func (m *Macros) expand(text string, bind map[string]string, stack []string, record *[]edit) (string, *sqlerr.Error) {
	l := &lexer{d: m.d, src: text}
	inBody := bind != nil

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); {
		if stop := l.skip(i); stop >= 0 {
			if stop <= i {
				stop = i + 1
			}
			b.WriteString(text[i:stop])
			i = stop
			continue
		}
		c := text[i]
		switch {
		case isIdentStart(c):
			j := l.identEnd(i)
			if !strings.EqualFold(text[i:j], "sqlc") {
				b.WriteString(text[i:j])
				i = j
				continue
			}
			if call, ok, err := l.macroCall(i, j); err != nil {
				return "", err
			} else if ok {
				expansion, err := m.call(l, call, bind, stack)
				if err != nil {
					return "", err
				}
				newStart := b.Len()
				b.WriteString(expansion)
				if record != nil {
					*record = append(*record, edit{
						oldStart: call.start, oldEnd: call.end,
						newStart: newStart, newEnd: b.Len(),
					})
				}
				i = call.end
				continue
			}
			if inBody {
				// A reference to one of the macro's own parameters. Any other
				// sqlc call is copied whole, so its argument is never taken
				// for a reference itself.
				if occ, next, ok := l.sqlcCall(i, j); ok {
					arg, bound := bind[m.d.fold(occ.name)]
					if occ.err == nil && bound && occ.kind != kindEmbed {
						if occ.kind != kindArg {
							return "", &sqlerr.Error{
								Message:  fmt.Sprintf("the macro parameter %q can only be referenced with sqlc.arg()", occ.name),
								Location: i,
							}
						}
						b.WriteString(arg)
					} else {
						b.WriteString(text[i:next])
					}
					i = next
					continue
				}
			}
			b.WriteString(text[i:j])
			i = j

		case inBody && c == '@' && m.d.AtSign && isIdentStart(l.at(i+1)):
			j := l.identEnd(i + 1)
			if arg, bound := bind[m.d.fold(text[i+1:j])]; bound {
				b.WriteString(arg)
			} else {
				b.WriteString(text[i:j])
			}
			i = j

		case inBody && ((c == '$' && m.d.DollarNumber && isDigit(l.at(i+1))) || (c == '?' && m.d.Question)):
			// A positional placeholder would be numbered against whichever
			// query the macro lands in, which is never what was meant.
			j := i + 1
			for j < len(text) && isDigit(text[j]) {
				j++
			}
			return "", &sqlerr.Error{
				Message:  fmt.Sprintf("positional placeholder %q is not allowed in a macro; use sqlc.arg()", text[i:j]),
				Location: i,
			}

		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String(), nil
}

// macroCall is a "sqlc.macro.<name>(<args>)" call found in a lexer's source.
type macroCall struct {
	name       string
	start, end int
	// open and closing are the offsets of the call's parentheses.
	open, closing int
}

// macroCall parses the macro call that starts at start, where afterSqlc is the
// offset just past the "sqlc" identifier. It reports false when the text is
// not a macro call at all.
func (l *lexer) macroCall(start, afterSqlc int) (macroCall, bool, *sqlerr.Error) {
	i := l.skipSpace(afterSqlc)
	if l.at(i) != '.' {
		return macroCall{}, false, nil
	}
	i = l.skipSpace(i + 1)
	if !isIdentStart(l.at(i)) {
		return macroCall{}, false, nil
	}
	j := l.identEnd(i)
	if !strings.EqualFold(l.src[i:j], "macro") {
		return macroCall{}, false, nil
	}
	i = l.skipSpace(j)
	if l.at(i) != '.' {
		return macroCall{}, false, nil
	}
	i = l.skipSpace(i + 1)
	if !isIdentStart(l.at(i)) {
		return macroCall{}, false, nil
	}
	j = l.identEnd(i)
	call := macroCall{name: l.src[i:j], start: start}

	call.open = l.skipSpace(j)
	if l.at(call.open) != '(' {
		return call, false, &sqlerr.Error{
			Message:  fmt.Sprintf("sqlc.macro.%s must be called with an argument list", call.name),
			Location: start,
		}
	}
	call.closing = l.matchParen(call.open)
	if call.closing < 0 {
		return call, false, &sqlerr.Error{
			Message:  fmt.Sprintf("sqlc.macro.%s is missing its closing parenthesis", call.name),
			Location: start,
		}
	}
	call.end = call.closing + 1
	return call, true, nil
}

// call returns the expansion of one macro call: the macro's body, with each
// parameter replaced by its argument, in parentheses.
func (m *Macros) call(l *lexer, call macroCall, bind map[string]string, stack []string) (string, *sqlerr.Error) {
	fail := func(format string, args ...any) (string, *sqlerr.Error) {
		return "", &sqlerr.Error{Message: fmt.Sprintf(format, args...), Location: call.start}
	}

	def, ok := m.defs[call.name]
	if !ok {
		return fail("macro %q is not defined", call.name)
	}
	if slices.Contains(stack, call.name) {
		return fail("macro %q expands itself: %s", call.name, strings.Join(append(stack, call.name), " -> "))
	}

	if def.spansLines {
		return fail("in macro %q: %s", call.name, errSpansLines)
	}

	args := l.splitArgs(call.open+1, call.closing)
	if len(args) != len(def.params) {
		return fail("macro %q takes %s; got %d", call.name, plural(len(def.params), "argument"), len(args))
	}

	// Arguments belong to the caller: they are expanded where they were
	// written, and are not looked at again once they are inside the body.
	callee := make(map[string]string, len(args))
	offset := call.open + 1
	for i, arg := range args {
		if strings.TrimSpace(arg) == "" {
			return fail("macro %q is missing argument %d", call.name, i+1)
		}
		expanded, err := m.expand(arg, bind, stack, nil)
		if err != nil {
			err.Location += offset
			return "", err
		}
		flat, ok := flatten(m.d, expanded)
		if !ok {
			return "", &sqlerr.Error{Message: errSpansLines, Location: offset}
		}
		offset += len(arg) + 1
		param := def.params[i]
		callee[param.name] = m.d.macroArg(param, strings.TrimSpace(flat))
	}

	body, err := m.expand(def.body, callee, append(stack, call.name), nil)
	if err != nil {
		// The body is not in the file the error is reported against, so the
		// error points at the call that reached it.
		return fail("in macro %q: %s", call.name, err.Message)
	}
	return "(" + body + ")", nil
}

func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

// macroArg returns the text that stands in for a parameter. A typed parameter
// casts its argument. An untyped one takes it as written, in parentheses
// unless it is already a single term, so that the operators around the
// reference cannot regroup an expression.
func (d Dialect) macroArg(p macroParam, arg string) string {
	if p.typ != "" {
		return "CAST(" + arg + " AS " + p.typ + ")"
	}
	if d.atomic(arg) {
		return arg
	}
	return "(" + arg + ")"
}

// atomic reports whether s is a single term: a literal, a placeholder, a
// possibly qualified name, a function call or something already in
// parentheses.
func (d Dialect) atomic(s string) bool {
	l := &lexer{d: d, src: s}
	n := len(s)
	if n == 0 {
		return false
	}
	if l.skip(0) == n {
		return true // a string literal or a quoted identifier
	}
	switch c := s[0]; {
	case c == '(':
		return l.matchParen(0) == n-1
	case isDigit(c):
		for i := 1; i < n; i++ {
			if !isDigit(s[i]) && s[i] != '.' {
				return false
			}
		}
		return true
	case c == '$' || c == '?':
		for i := 1; i < n; i++ {
			if !isDigit(s[i]) {
				return false
			}
		}
		return true
	case c == '@':
		return isIdentStart(l.at(1)) && l.identEnd(1) == n
	}
	for i := 0; ; {
		switch c := s[i]; {
		case isIdentStart(c):
			i = l.identEnd(i)
		case c == '"' || (c == '`' && d.Backtick):
			i = l.skip(i)
		default:
			return false
		}
		if i >= n {
			return true
		}
		switch s[i] {
		case '.':
			if i++; i >= n {
				return false
			}
		case '(':
			return l.matchParen(i) == n-1
		default:
			return false
		}
	}
}

// errSpansLines is the error for text flatten could not put on one line.
const errSpansLines = "a string or quoted name that spans lines cannot be expanded by a macro"

// flatten returns s on a single line and without its comments, which is the
// form text is pasted into another query in. A comment that reached the query
// would be read as one of its own annotations, and a line break would move
// every line after it. Optimizer hints are comments in name only, and stay.
//
// It reports false when a line break is part of a literal, where it cannot be
// taken out without changing what the literal says.
func flatten(d Dialect, s string) (string, bool) {
	l := &lexer{d: d, src: s}
	var b strings.Builder
	b.Grow(len(s))
	space := func() {
		if b.Len() > 0 && !strings.HasSuffix(b.String(), " ") {
			b.WriteByte(' ')
		}
	}
	ok := true
	for i := 0; i < len(s); {
		if stop := l.skip(i); stop >= 0 {
			if stop <= i {
				stop = i + 1
			}
			if l.isComment(i) && !strings.HasPrefix(s[i:], "/*+") && !strings.HasPrefix(s[i:], "/*!") {
				space()
			} else {
				if strings.Contains(s[i:stop], "\n") {
					ok = false
				}
				b.WriteString(s[i:stop])
			}
			i = stop
			continue
		}
		if isSpace(s[i]) {
			space()
		} else {
			b.WriteByte(s[i])
		}
		i++
	}
	return b.String(), ok
}

// isComment reports whether the token skip found at i is a comment rather than
// a literal.
func (l *lexer) isComment(i int) bool {
	switch l.at(i) {
	case '-':
		return l.at(i+1) == '-'
	case '/':
		return l.at(i+1) == '*'
	case '#':
		return l.d.HashComment
	}
	return false
}

// blank returns s with everything but its line breaks replaced by spaces.
func blank(s string) string {
	out := []byte(s)
	for i, c := range out {
		if c != '\n' {
			out[i] = ' '
		}
	}
	return string(out)
}
