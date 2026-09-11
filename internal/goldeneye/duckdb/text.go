package duckdb

import (
	"strings"
)

// DuckDB prints a plan with every column by its bare name and every
// aliased expression by its alias, and describes a query no further than
// its result columns' names and types, so which table a result column is
// read from and which column a parameter is compared with are read from
// the query text: the select list's items, the FROM clause's tables and
// aliases, and the operand beside each parameter. The text is tokenized
// and read at the top level of the statement, outside the parentheses a
// subquery or a CTE body sits in.

// token is one lexical element of a query.
type token struct {
	kind byte // 'w' word, 'q' quoted identifier, 's' string, 'n' number, 'p' parameter, 'o' operator or punctuation
	text string
}

var operators = []string{"::", "<>", "!=", "<=", ">=", "!~~", "~~", "->>", "->", "||", "=", "<", ">", "(", ")", ",", ".", "*", "+", "-", "/", "%", "[", "]", ";", "^", "&", "|", "~", "!", ":"}

// tokenize splits a query into tokens, dropping comments. A parameter is
// $ followed by digits, which is how bind spells every parameter.
func tokenize(src string) []token {
	var out []token
	i := 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case strings.HasPrefix(src[i:], "--"):
			end := strings.IndexByte(src[i:], '\n')
			if end < 0 {
				i = len(src)
			} else {
				i += end
			}
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i:], "*/")
			if end < 0 {
				i = len(src)
			} else {
				i += end + 2
			}
		case c == '\'':
			end := quotedEnd(src, i)
			out = append(out, token{'s', src[i:end]})
			i = end
		case c == '"':
			end := quotedEnd(src, i)
			out = append(out, token{'q', strings.ReplaceAll(src[i+1:end-1], `""`, `"`)})
			i = end
		case c == '$' && i+1 < len(src) && isDigit(src[i+1]):
			end := i + 1
			for end < len(src) && isDigit(src[end]) {
				end++
			}
			out = append(out, token{'p', src[i+1 : end]})
			i = end
		case isDigit(c) || (c == '.' && i+1 < len(src) && isDigit(src[i+1])):
			end := i
			for end < len(src) && (isDigit(src[end]) || src[end] == '.' || src[end] == 'e' || src[end] == 'E' || src[end] == '_') {
				end++
			}
			out = append(out, token{'n', src[i:end]})
			i = end
		case isWordStart(c):
			end := i
			for end < len(src) && isWordByte(src[end]) {
				end++
			}
			out = append(out, token{'w', src[i:end]})
			i = end
		default:
			matched := false
			for _, op := range operators {
				if strings.HasPrefix(src[i:], op) {
					out = append(out, token{'o', op})
					i += len(op)
					matched = true
					break
				}
			}
			if !matched {
				out = append(out, token{'o', string(c)})
				i++
			}
		}
	}
	return out
}

// quotedEnd returns the index just past the quoted token starting at i,
// whose delimiter is escaped by doubling it.
func quotedEnd(s string, i int) int {
	q := s[i]
	j := i + 1
	for j < len(s) {
		switch {
		case s[j] == q && j+1 < len(s) && s[j+1] == q:
			j += 2
		case s[j] == q:
			return j + 1
		default:
			j++
		}
	}
	return len(s)
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }
func isWordStart(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}
func isWordByte(c byte) bool { return isWordStart(c) || isDigit(c) }

// text is a tokenized query.
type text []token

func (t text) at(i int) token {
	if i < 0 || i >= len(t) {
		return token{}
	}
	return t[i]
}

// isWord reports whether token i is the keyword, in any case.
func (t text) isWord(i int, word string) bool {
	tok := t.at(i)
	return tok.kind == 'w' && strings.EqualFold(tok.text, word)
}

// isOp reports whether token i is the operator.
func (t text) isOp(i int, op string) bool {
	tok := t.at(i)
	return tok.kind == 'o' && tok.text == op
}

// isName reports whether token i can name something: a word or a quoted
// identifier.
func (t text) isName(i int) bool {
	k := t.at(i).kind
	return k == 'w' || k == 'q'
}

// group returns the index of the parenthesis closing the one at i.
func (t text) group(i int) int {
	depth := 0
	for j := i; j < len(t); j++ {
		switch {
		case t.isOp(j, "("):
			depth++
		case t.isOp(j, ")"):
			depth--
			if depth == 0 {
				return j
			}
		}
	}
	return len(t) - 1
}

// ref is a column reference: a qualifier, empty for none, and a column,
// or "*" for every column of the qualifier.
type ref struct {
	qualifier, column string
}

// readRef reads a dotted name starting at i, returning the parts and the
// index after them.
func (t text) readRef(i int) ([]string, int) {
	var parts []string
	for t.isName(i) || t.isOp(i, "*") {
		parts = append(parts, t.at(i).text)
		if !t.isOp(i+1, ".") {
			return parts, i + 1
		}
		i += 2
	}
	return parts, i
}

// keywords that end a list of tables or select items at the top level.
var clauseKeywords = map[string]bool{
	"where": true, "group": true, "order": true, "limit": true, "offset": true, "having": true,
	"qualify": true, "window": true, "returning": true, "set": true, "values": true, "select": true,
	"union": true, "except": true, "intersect": true, "on": true, "from": true, "join": true,
	"left": true, "right": true, "full": true, "inner": true, "cross": true, "outer": true,
	"natural": true, "asof": true, "semi": true, "anti": true, "positional": true, "lateral": true,
	"using": true, "with": true, "into": true, "default": true, "fetch": true, "for": true,
}

// tableRef is one table of a statement's scope, by name and alias.
type tableRef struct {
	name, alias string
}

// scope is what the statement reads and writes: the tables of its FROM,
// JOIN and USING clauses and the table a DML statement targets, with the
// names of its CTEs, which are not tables.
type scope struct {
	kind   string // select, insert, update, delete
	target tableRef
	tables []tableRef
	ctes   map[string]bool
}

// readScope reads a statement's scope from its top-level tokens.
func (t text) readScope() scope {
	sc := scope{ctes: map[string]bool{}}
	i := 0
	// WITH name [(cols)] AS (...), ...
	if t.isWord(0, "with") {
		i = 1
		if t.isWord(i, "recursive") {
			i++
		}
		for t.isName(i) {
			sc.ctes[strings.ToLower(t.at(i).text)] = true
			i++
			if t.isOp(i, "(") {
				i = t.group(i) + 1
			}
			if t.isWord(i, "as") {
				i++
			}
			if t.isWord(i, "not") || t.isWord(i, "materialized") {
				for !t.isOp(i, "(") && i < len(t) {
					i++
				}
			}
			if t.isOp(i, "(") {
				i = t.group(i) + 1
			}
			if t.isOp(i, ",") {
				i++
				continue
			}
			break
		}
	}
	sc.kind = "select"
	switch {
	case t.isWord(i, "insert"):
		sc.kind = "insert"
		for i < len(t) && !t.isWord(i, "into") {
			i++
		}
		sc.target, i = t.readTable(i + 1)
	case t.isWord(i, "update"):
		sc.kind = "update"
		sc.target, i = t.readTable(i + 1)
	case t.isWord(i, "delete"):
		sc.kind = "delete"
		if t.isWord(i+1, "from") {
			i++
		}
		sc.target, i = t.readTable(i + 1)
	}
	for ; i < len(t); i++ {
		switch {
		case t.isOp(i, "("):
			i = t.group(i)
		case t.isWord(i, "from") || t.isWord(i, "join"):
			i = t.readTables(i+1, &sc.tables) - 1
		case t.isWord(i, "using") && !t.isOp(i+1, "("):
			i = t.readTables(i+1, &sc.tables) - 1
		}
	}
	return sc
}

// readTable reads a table name with an optional alias.
func (t text) readTable(i int) (tableRef, int) {
	parts, j := t.readRef(i)
	if len(parts) == 0 {
		return tableRef{}, i
	}
	ref := tableRef{name: strings.ToLower(parts[len(parts)-1])}
	if t.isWord(j, "as") {
		j++
	}
	if t.isName(j) && !clauseKeywords[strings.ToLower(t.at(j).text)] {
		ref.alias = strings.ToLower(t.at(j).text)
		j++
	}
	return ref, j
}

// readTables reads a comma-separated list of table references, skipping
// subqueries and table functions, and returns the index after it.
func (t text) readTables(i int, into *[]tableRef) int {
	for i < len(t) {
		if t.isOp(i, "(") {
			// A derived table or a table function's arguments.
			i = t.group(i) + 1
			if t.isWord(i, "as") {
				i++
			}
			if t.isName(i) && !clauseKeywords[strings.ToLower(t.at(i).text)] {
				i++
			}
		} else if t.isName(i) && !clauseKeywords[strings.ToLower(t.at(i).text)] {
			ref, j := t.readTable(i)
			if t.isOp(j, "(") {
				// A table function: its result is not a table.
				j = t.group(j) + 1
			} else {
				*into = append(*into, ref)
			}
			i = j
		} else {
			return i
		}
		if t.isOp(i, ",") {
			i++
			continue
		}
		return i
	}
	return i
}

// item is one entry of a select list or a RETURNING list.
type item struct {
	star bool   // every column of qualifier, or of every table
	ref  *ref   // a column reference, when the item is one
	name string // the alias, or the column name of a reference
}

// items splits a select or RETURNING list into its items, given the index
// of the first token after the keyword.
func (t text) items(start int) []item {
	var out []item
	i := start
	for i < len(t) {
		end := i
		for end < len(t) && !t.isOp(end, ",") && !(t.at(end).kind == 'w' && clauseKeywords[strings.ToLower(t.at(end).text)] && !t.isWord(end, "on")) && !t.isOp(end, ";") {
			if t.isOp(end, "(") {
				end = t.group(end)
			}
			end++
		}
		if end > i {
			out = append(out, t.item(i, end))
		}
		if !t.isOp(end, ",") {
			break
		}
		i = end + 1
	}
	return out
}

// item classifies the tokens of one select item.
func (t text) item(start, end int) item {
	parts, j := t.readRef(start)
	if len(parts) > 0 && parts[len(parts)-1] == "*" {
		it := item{star: true}
		if len(parts) > 1 {
			it.ref = &ref{qualifier: strings.ToLower(parts[len(parts)-2]), column: "*"}
		}
		return it
	}
	var it item
	if len(parts) > 0 && j == end {
		it.ref = &ref{column: parts[len(parts)-1]}
		if len(parts) > 1 {
			it.ref.qualifier = strings.ToLower(parts[len(parts)-2])
		}
		it.name = parts[len(parts)-1]
		return it
	}
	if len(parts) > 0 && j == end-1 && t.isName(end-1) {
		// A reference followed by a bare alias.
		it.ref = &ref{column: parts[len(parts)-1]}
		if len(parts) > 1 {
			it.ref.qualifier = strings.ToLower(parts[len(parts)-2])
		}
		it.name = t.at(end - 1).text
		return it
	}
	if len(parts) > 0 && j == end-2 && t.isWord(end-2, "as") && t.isName(end-1) {
		it.ref = &ref{column: parts[len(parts)-1]}
		if len(parts) > 1 {
			it.ref.qualifier = strings.ToLower(parts[len(parts)-2])
		}
		it.name = t.at(end - 1).text
		return it
	}
	if t.isWord(end-2, "as") && t.isName(end-1) {
		it.name = t.at(end - 1).text
	}
	return it
}

// selectItems finds the statement's select list: the items after its
// top-level SELECT, or every column when a FROM-first query has none.
func (t text) selectItems() []item {
	for i := 0; i < len(t); i++ {
		switch {
		case t.isOp(i, "("):
			i = t.group(i)
		case t.isWord(i, "select"):
			j := i + 1
			if t.isWord(j, "distinct") || t.isWord(j, "all") {
				j++
				if t.isWord(j, "on") && t.isOp(j+1, "(") {
					j = t.group(j+1) + 1
				}
			}
			return t.items(j)
		}
	}
	if t.isWord(0, "from") {
		return []item{{star: true}}
	}
	return nil
}

// returningItems finds the items of a top-level RETURNING clause.
func (t text) returningItems() []item {
	for i := 0; i < len(t); i++ {
		switch {
		case t.isOp(i, "("):
			i = t.group(i)
		case t.isWord(i, "returning"):
			return t.items(i + 1)
		}
	}
	return nil
}

// find returns the index of parameter k.
func (t text) find(k string) int {
	for i, tok := range t {
		if tok.kind == 'p' && tok.text == k {
			return i
		}
	}
	return -1
}

var comparisons = map[string]bool{"=": true, "<>": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true, "~~": true, "!~~": true}

// partner finds the column a parameter is compared with or assigned to:
// the column reference on the other side of the comparison, LIKE or IN
// it is an operand of, or the column a SET assigns it to.
func (t text) partner(k string) (ref, bool) {
	i := t.find(k)
	if i < 0 {
		return ref{}, false
	}
	// col op $k, col [NOT] LIKE $k, col IN ($k, ...)
	j := i - 1
	if t.isOp(j, "(") || t.isOp(j, ",") {
		for j >= 0 && (t.isOp(j, ",") || t.at(j).kind == 'p' || t.at(j).kind == 's' || t.at(j).kind == 'n') {
			j--
		}
		if t.isOp(j, "(") {
			j--
		}
	}
	if t.at(j).kind == 'o' && comparisons[t.at(j).text] || t.isWord(j, "like") || t.isWord(j, "ilike") || t.isWord(j, "glob") || t.isWord(j, "in") {
		if t.isWord(j-1, "not") {
			j--
		}
		if r, ok := t.refEnding(j - 1); ok {
			return r, true
		}
	}
	// $k op col
	j = i + 1
	if t.at(j).kind == 'o' && comparisons[t.at(j).text] {
		if r, ok := t.refStarting(j + 1); ok {
			return r, true
		}
	}
	return ref{}, false
}

// refEnding reads the column reference whose last token is at i.
func (t text) refEnding(i int) (ref, bool) {
	if !t.isName(i) {
		return ref{}, false
	}
	r := ref{column: t.at(i).text}
	if t.isOp(i-1, ".") && t.isName(i-2) {
		r.qualifier = strings.ToLower(t.at(i - 2).text)
	}
	return r, true
}

// refStarting reads the column reference starting at i, which must not be
// a function call.
func (t text) refStarting(i int) (ref, bool) {
	parts, j := t.readRef(i)
	if len(parts) == 0 || t.isOp(j, "(") {
		return ref{}, false
	}
	r := ref{column: parts[len(parts)-1]}
	if len(parts) > 1 {
		r.qualifier = strings.ToLower(parts[len(parts)-2])
	}
	return r, true
}

// castOf returns the type a parameter is cast to, as $k::T or CAST($k AS
// T), spelled as the query spells it.
func (t text) castOf(k string) (string, bool) {
	i := t.find(k)
	if i < 0 {
		return "", false
	}
	if t.isOp(i+1, "::") {
		return t.typeText(i + 2), true
	}
	if t.isWord(i-2, "cast") && t.isOp(i-1, "(") && t.isWord(i+1, "as") {
		return t.typeText(i + 2), true
	}
	return "", false
}

// typeText reads a type starting at i: a dotted name, optional arguments
// in parentheses and any number of [] or [N] suffixes.
func (t text) typeText(i int) string {
	start := i
	_, i = t.readRef(i)
	for t.isName(i) && !clauseKeywords[strings.ToLower(t.at(i).text)] && !t.isWord(i, "as") {
		// Multi-word names: TIMESTAMP WITH TIME ZONE.
		i++
	}
	if t.isOp(i, "(") {
		i = t.group(i) + 1
	}
	for t.isOp(i, "[") {
		for i < len(t) && !t.isOp(i, "]") {
			i++
		}
		i++
	}
	var b strings.Builder
	for j := start; j < i && j < len(t); j++ {
		tok := t.at(j)
		if j > start && (tok.kind == 'w' || tok.kind == 'q' || tok.kind == 'n' || tok.kind == 's') && (t.at(j-1).kind == 'w' || t.at(j-1).kind == 'q' || t.isOp(j-1, ",")) {
			b.WriteByte(' ')
		}
		b.WriteString(tok.text)
	}
	return b.String()
}

// valuesPosition reports, for a parameter inside an INSERT's VALUES, the
// index of the column it is a value for, counting positions within its
// row.
func (t text) valuesPosition(k string) (int, bool) {
	i := t.find(k)
	if i < 0 {
		return 0, false
	}
	for j := 0; j < len(t); j++ {
		if t.isOp(j, "(") {
			j = t.group(j)
			continue
		}
		if !t.isWord(j, "values") {
			continue
		}
		// Rows follow, each a parenthesised list.
		for r := j + 1; t.isOp(r, "("); r = t.group(r) + 2 {
			end := t.group(r)
			if i < r || i > end {
				if !t.isOp(end+1, ",") {
					break
				}
				continue
			}
			pos, depth := 0, 0
			for q := r + 1; q < i; q++ {
				switch {
				case t.isOp(q, "("):
					depth++
				case t.isOp(q, ")"):
					depth--
				case t.isOp(q, ",") && depth == 0:
					pos++
				}
			}
			return pos, true
		}
		return 0, false
	}
	return 0, false
}

// insertColumns lists the columns an INSERT names, or nil for all.
func (t text) insertColumns() []string {
	for i := 0; i < len(t); i++ {
		if t.isWord(i, "into") {
			_, j := t.readRef(i + 1)
			if t.isWord(j, "as") {
				j += 2
			} else if t.isName(j) && !clauseKeywords[strings.ToLower(t.at(j).text)] {
				j++
			}
			if !t.isOp(j, "(") {
				return nil
			}
			var cols []string
			for q := j + 1; q < t.group(j); q++ {
				if t.isName(q) {
					cols = append(cols, t.at(q).text)
				}
			}
			return cols
		}
	}
	return nil
}
