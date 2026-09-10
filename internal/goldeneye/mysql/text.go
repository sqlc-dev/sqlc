package mysql

import (
	"strings"
)

// MySQL prints a statement back with every name resolved, every table
// alias kept, every SELECT * expanded and every binary operation
// parenthesised — in the optimizer trace's expanded_query, once per query
// block, and in the note EXPLAIN leaves. Such a text is tokenised here and
// read for three things: the tables it names and their aliases, the
// columns each query block projects, and the operand a parameter is
// compared with or assigned to.

type tokKind int

const (
	tkIdent  tokKind = iota // `a`.`b`.`c`
	tkVar                   // @`name`
	tkNumber                // 42, 1.5
	tkString                // 'text'
	tkWord                  // a keyword, a function name or a <marker>
	tkOp                    // =, <=>, +, like is a word
	tkLParen
	tkRParen
	tkComma
)

type token struct {
	kind  tokKind
	text  string   // the token as printed; a word is lower-cased
	parts []string // an identifier's parts, unquoted; a variable's name
	start int      // offsets into the source
	end   int
}

// text is a tokenised statement.
type text struct {
	src   string
	toks  []token
	match []int // the index of each parenthesis's partner, or -1
}

func tokenize(src string) text {
	var toks []token
	i := 0
	emit := func(kind tokKind, start, end int, parts []string) {
		toks = append(toks, token{kind: kind, text: src[start:end], parts: parts, start: start, end: end})
	}
	for i < len(src) {
		c := src[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i:], "*/")
			if end < 0 {
				i = len(src)
			} else {
				i += end + 2
			}
		case c == '`':
			start := i
			var parts []string
			for {
				part, end := readQuoted(src, i)
				parts = append(parts, part)
				i = end
				if i+1 < len(src) && src[i] == '.' && src[i+1] == '`' {
					i++
					continue
				}
				break
			}
			emit(tkIdent, start, i, parts)
		case c == '@':
			start := i
			i++
			for i < len(src) && src[i] == '@' {
				i++
			}
			if i < len(src) && src[i] == '`' {
				name, end := readQuoted(src, i)
				i = end
				emit(tkVar, start, i, []string{name})
			} else {
				for i < len(src) && isWordByte(src[i]) {
					i++
				}
				emit(tkVar, start, i, []string{src[start+1 : i]})
			}
		case c == '\'' || c == '"':
			start := i
			i = skipString(src, i)
			emit(tkString, start, i, nil)
		case c >= '0' && c <= '9' || c == '.' && i+1 < len(src) && src[i+1] >= '0' && src[i+1] <= '9':
			start := i
			for i < len(src) && (src[i] >= '0' && src[i] <= '9' || src[i] == '.') {
				i++
			}
			if i < len(src) && (src[i] == 'e' || src[i] == 'E') {
				i++
				if i < len(src) && (src[i] == '+' || src[i] == '-') {
					i++
				}
				for i < len(src) && src[i] >= '0' && src[i] <= '9' {
					i++
				}
			}
			emit(tkNumber, start, i, nil)
		case c == '<' && i+1 < len(src) && isWordByte(src[i+1]) && strings.IndexByte(src[i:], '>') > 0 && isMarker(src[i:i+strings.IndexByte(src[i:], '>')+1]):
			// <cache>, <in_optimizer>, <exists>: what the optimizer wrapped
			// an expression in. <cache> says nothing about the expression
			// and is dropped; the others are read as calls.
			end := i + strings.IndexByte(src[i:], '>') + 1
			if src[i:end] != "<cache>" {
				emit(tkWord, i, end, nil)
				toks[len(toks)-1].text = src[i+1 : end-1]
			}
			i = end
		case c == '(':
			emit(tkLParen, i, i+1, nil)
			i++
		case c == ')':
			emit(tkRParen, i, i+1, nil)
			i++
		case c == ',':
			emit(tkComma, i, i+1, nil)
			i++
		case isWordByte(c):
			start := i
			for i < len(src) && isWordByte(src[i]) {
				i++
			}
			emit(tkWord, start, i, nil)
			toks[len(toks)-1].text = strings.ToLower(src[start:i])
		default:
			n := 1
			for _, op := range operators {
				if strings.HasPrefix(src[i:], op) {
					n = len(op)
					break
				}
			}
			emit(tkOp, i, i+n, nil)
			i += n
		}
	}
	t := text{src: src, toks: toks, match: make([]int, len(toks))}
	var open []int
	for i, tok := range toks {
		t.match[i] = -1
		switch tok.kind {
		case tkLParen:
			open = append(open, i)
		case tkRParen:
			if len(open) > 0 {
				j := open[len(open)-1]
				open = open[:len(open)-1]
				t.match[i] = j
				t.match[j] = i
			}
		}
	}
	return t
}

// operators is every symbolic operator MySQL prints, longest first.
var operators = []string{
	"<=>", "->>", "<>", "!=", ">=", "<=", "->", "<<", ">>", "||", "&&", ":=",
	"=", "<", ">", "+", "-", "*", "/", "%", "!", "^", "&", "|", "~", ".",
}

func isMarker(s string) bool {
	for _, c := range []byte(s[1 : len(s)-1]) {
		if !isWordByte(c) {
			return false
		}
	}
	return true
}

func isWordByte(c byte) bool {
	return c == '_' || c == '$' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// readQuoted reads the backticked identifier at i, unescaping a doubled
// backtick, and returns it with the index just past its closing quote.
func readQuoted(s string, i int) (string, int) {
	var b strings.Builder
	j := i + 1
	for j < len(s) {
		if s[j] == '`' {
			if j+1 < len(s) && s[j+1] == '`' {
				b.WriteByte('`')
				j += 2
				continue
			}
			return b.String(), j + 1
		}
		b.WriteByte(s[j])
		j++
	}
	return b.String(), j
}

// skipString returns the index just past the string literal starting at i,
// honouring backslash escapes and doubled quotes.
func skipString(s string, i int) int {
	q := s[i]
	j := i + 1
	for j < len(s) {
		switch {
		case s[j] == '\\' && j+1 < len(s):
			j += 2
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

// wordOperators are the operators MySQL spells as words. and/or/not/is
// are left out: MySQL parenthesises each predicate, so a parameter that
// is a whole predicate has nothing to be compared with.
var wordOperators = map[string]bool{
	"like": true, "in": true, "regexp": true, "rlike": true, "div": true, "mod": true,
	"xor": true, "sounds": true, "escape": true, "between": true, "collate": true,
	"member": true,
}

// keywords are the words that can precede a parenthesis without naming a
// function, so that "(a, b)" after one is a list rather than a call.
var keywords = map[string]bool{
	"in": true, "values": true, "value": true, "and": true, "or": true, "not": true,
	"where": true, "on": true, "from": true, "select": true, "having": true,
	"set": true, "when": true, "then": true, "else": true, "case": true, "end": true,
	"like": true, "between": true, "by": true, "group": true, "order": true,
	"limit": true, "offset": true, "using": true, "join": true, "as": true, "is": true,
	"xor": true, "div": true, "mod": true, "interval": true, "distinct": true,
	"all": true, "any": true, "some": true, "escape": true, "regexp": true,
	"rlike": true, "sounds": true, "collate": true, "union": true, "except": true,
	"intersect": true, "with": true, "insert": true, "into": true, "update": true,
	"delete": true, "replace": true, "key": true, "duplicate": true, "do": true,
	"partition": true, "over": true, "window": true, "rows": true, "range": true,
	"straight_join": true, "cross": true, "inner": true, "left": true, "right": true,
	"natural": true, "lateral": true, "of": true, "for": true, "lock": true,
	"share": true, "null": true, "true": true, "false": true, "unknown": true,
	"binary": true, "row": true, "asc": true, "desc": true, "sql_calc_found_rows": true,
	"high_priority": true, "low_priority": true, "ignore": true, "table": true,
	"member": true, "exists": true,
}

// literalWords are the words that are operands on their own.
var literalWords = map[string]bool{
	"null": true, "true": true, "false": true, "unknown": true,
	"current_timestamp": true, "current_date": true, "current_time": true,
	"current_user": true, "localtime": true, "localtimestamp": true, "utc_timestamp": true,
}

func (t text) at(i int) token {
	if i < 0 || i >= len(t.toks) {
		return token{kind: -1}
	}
	return t.toks[i]
}

func (t text) isWord(i int, word string) bool {
	tok := t.at(i)
	return tok.kind == tkWord && tok.text == word
}

// isCall reports whether the word at i names a function, given that a
// parenthesis follows it.
func (t text) isCall(i int) bool {
	tok := t.at(i)
	return tok.kind == tkWord && !keywords[tok.text] || tok.kind == tkWord && tok.text == "exists"
}

// isOperator reports whether the token at i joins two operands.
func (t text) isOperator(i int) bool {
	tok := t.at(i)
	switch tok.kind {
	case tkOp:
		return tok.text != "." && tok.text != "~" && tok.text != "!"
	case tkWord:
		return wordOperators[tok.text]
	}
	return false
}

// operandEnd returns the index just past the operand starting at i, or -1
// when no operand starts there: a name, a literal, a call with its
// arguments, or a parenthesised group.
func (t text) operandEnd(i int) int {
	tok := t.at(i)
	switch tok.kind {
	case tkIdent, tkVar, tkNumber, tkString:
		return i + 1
	case tkLParen:
		if t.match[i] < 0 {
			return -1
		}
		return t.match[i] + 1
	case tkWord:
		if literalWords[tok.text] {
			return i + 1
		}
		if t.at(i+1).kind == tkLParen && t.isCall(i) {
			return t.operandEnd(i + 1)
		}
		if tok.text == "not" || tok.text == "binary" || tok.text == "interval" {
			return t.operandEnd(i + 1)
		}
	case tkOp:
		if tok.text == "-" || tok.text == "+" || tok.text == "~" || tok.text == "!" {
			return t.operandEnd(i + 1)
		}
	}
	return -1
}

// operandStart returns the index of the operand ending just before j, or
// -1 when none ends there.
func (t text) operandStart(j int) int {
	k := j - 1
	tok := t.at(k)
	switch tok.kind {
	case tkIdent, tkVar, tkNumber, tkString:
		return k
	case tkRParen:
		open := t.match[k]
		if open < 0 {
			return -1
		}
		if t.at(open-1).kind == tkWord && t.isCall(open-1) {
			return open - 1
		}
		return open
	case tkWord:
		if literalWords[tok.text] {
			return k
		}
	}
	return -1
}

// group returns the index of the parenthesis enclosing token i, or -1.
func (t text) group(i int) int {
	depth := 0
	for k := i - 1; k >= 0; k-- {
		switch t.toks[k].kind {
		case tkRParen:
			depth++
		case tkLParen:
			if depth == 0 {
				return k
			}
			depth--
		}
	}
	return -1
}

// list splits the tokens strictly inside the group opened at open into its
// comma-separated members, each a token range.
func (t text) list(open int) [][2]int {
	close := t.match[open]
	if close < 0 {
		return nil
	}
	var members [][2]int
	start := open + 1
	for i := open + 1; i < close; i++ {
		switch t.toks[i].kind {
		case tkLParen:
			i = t.match[i]
		case tkComma:
			members = append(members, [2]int{start, i})
			start = i + 1
		}
	}
	if start < close {
		members = append(members, [2]int{start, close})
	}
	return members
}

// item is one entry of a SELECT list: the expression and the name it is
// projected as.
type item struct {
	start, end int
	alias      string
}

// selectItems reads the list of the first SELECT at the top level of the
// text, which is the statement's own for a query block's expanded_query
// and for a statement whose CTEs come first.
func (t text) selectItems() []item {
	depth := 0
	start := -1
	for i, tok := range t.toks {
		switch tok.kind {
		case tkLParen:
			depth++
		case tkRParen:
			depth--
		case tkWord:
			if depth == 0 && tok.text == "select" {
				start = i + 1
			}
		}
		if start >= 0 {
			break
		}
	}
	if start < 0 {
		return nil
	}
	for t.at(start).kind == tkWord && (t.at(start).text == "distinct" || t.at(start).text == "all" || t.at(start).text == "sql_calc_found_rows" || t.at(start).text == "straight_join" || t.at(start).text == "high_priority") {
		start++
	}
	end := len(t.toks)
	for i := start; i < len(t.toks); i++ {
		switch tok := t.toks[i]; tok.kind {
		case tkLParen:
			i = t.match[i]
			if i < 0 {
				return nil
			}
		case tkWord:
			switch tok.text {
			case "from", "where", "group", "having", "order", "limit", "union", "into", "window", "for", "lock", "except", "intersect":
				end = i
			}
		}
		if end != len(t.toks) {
			break
		}
	}
	var items []item
	itemStart := start
	flush := func(itemEnd int) {
		it := item{start: itemStart, end: itemEnd}
		if itemEnd-itemStart >= 3 && t.isWord(itemEnd-2, "as") && t.at(itemEnd-1).kind == tkIdent {
			it.alias = t.toks[itemEnd-1].parts[0]
			it.end = itemEnd - 2
		}
		items = append(items, it)
	}
	for i := start; i < end; i++ {
		switch t.toks[i].kind {
		case tkLParen:
			i = t.match[i]
		case tkComma:
			flush(i)
			itemStart = i + 1
		}
	}
	if itemStart < end {
		flush(end)
	}
	return items
}

// tableRef is a table as a FROM clause names it, with the schema when it
// spells one out.
type tableRef struct {
	schema, name string
}

// aliases reads the table aliases the text declares: a table name followed
// by a bare identifier, as `users` `u` in a FROM clause. A SELECT list
// alias is always introduced by AS, so the shape means nothing else.
func (t text) aliases(into map[string]tableRef) {
	for i := 0; i+1 < len(t.toks); i++ {
		tok, next := t.toks[i], t.toks[i+1]
		if tok.kind != tkIdent || next.kind != tkIdent || len(next.parts) != 1 || len(tok.parts) > 2 {
			continue
		}
		if t.at(i-1).kind == tkOp && t.at(i-1).text == "." {
			continue
		}
		ref := tableRef{name: tok.parts[len(tok.parts)-1]}
		if len(tok.parts) == 2 {
			ref.schema = tok.parts[0]
		}
		into[next.parts[0]] = ref
	}
}

// find returns the index of the variable named name, or -1.
func (t text) findVar(name string) int {
	for i, tok := range t.toks {
		if tok.kind == tkVar && len(tok.parts) == 1 && tok.parts[0] == name {
			return i
		}
	}
	return -1
}

// findNumber returns the index of the number printed as lit, or -1.
func (t text) findNumber(lit string) int {
	for i, tok := range t.toks {
		if tok.kind == tkNumber && tok.text == lit {
			return i
		}
	}
	return -1
}

// A partner is what a parameter stands against in a statement.
type partnerKind int

const (
	partnerNone       partnerKind = iota
	partnerOperand                // the operand it is compared with, or the argument beside it
	partnerProjection             // the SELECT list item it is projected as
	partnerInsert                 // the column of an INSERT ... VALUES row it is written to
	partnerLimit                  // the count of a LIMIT or OFFSET
)

type partner struct {
	kind       partnerKind
	start, end int // partnerOperand: the operand's tokens
	index      int // partnerProjection: the item; partnerInsert: the position in the row
}

// partner finds what the parameter at token index s stands against.
func (t text) partner(s int) partner {
	lo, hi := s, s+1
	if t.at(lo-1).kind == tkLParen && t.at(hi).kind == tkRParen {
		lo, hi = lo-1, hi+1
	}
	if t.at(hi).kind == tkWord && t.at(hi).text == "as" {
		for i, it := range t.selectItems() {
			if it.start <= lo && hi <= it.end {
				return partner{kind: partnerProjection, index: i}
			}
		}
	}
	// The count of a LIMIT or OFFSET, which is a number rather than a
	// variable since MySQL allows nothing else there.
	if t.at(lo-1).kind == tkWord && (t.at(lo-1).text == "limit" || t.at(lo-1).text == "offset") ||
		t.at(lo-1).kind == tkComma && t.at(lo-2).kind == tkNumber && t.isWord(lo-3, "limit") {
		return partner{kind: partnerLimit}
	}
	// Compared: the operand on the other side of the operator.
	if t.isOperator(lo - 1) {
		q := lo - 2
		if t.isWord(q, "not") {
			q--
		}
		if st := t.operandStart(q + 1); st >= 0 {
			// The second bound of BETWEEN ... AND: the operand before BETWEEN.
			if t.isWord(lo-1, "and") {
				if t.isWord(st-1, "between") || t.isWord(st-1, "not") && t.isWord(st-2, "between") {
					q = st - 2
					if t.isWord(q, "not") {
						q--
					}
					if st2 := t.operandStart(q + 1); st2 >= 0 {
						return partner{kind: partnerOperand, start: st2, end: q + 1}
					}
				}
			}
			return partner{kind: partnerOperand, start: st, end: q + 1}
		}
	}
	if t.isWord(lo-1, "and") {
		if st := t.operandStart(lo - 1); st >= 0 && (t.isWord(st-1, "between") || t.isWord(st-1, "not") && t.isWord(st-2, "between")) {
			q := st - 2
			if t.isWord(q, "not") {
				q--
			}
			if st2 := t.operandStart(q + 1); st2 >= 0 {
				return partner{kind: partnerOperand, start: st2, end: q + 1}
			}
		}
	}
	if t.isOperator(hi) {
		n := hi + 1
		if t.isWord(n, "not") {
			n++
		}
		if end := t.operandEnd(n); end >= 0 {
			return partner{kind: partnerOperand, start: n, end: end}
		}
	}
	// A member of a list: the row of an INSERT ... VALUES, the values of an
	// IN, or the arguments of a call, where the argument beside it says
	// what it is, a column above an expression above a constant.
	if t.at(lo-1).kind == tkComma || t.at(lo-1).kind == tkLParen {
		open := t.group(lo)
		if open < 0 {
			return partner{}
		}
		members := t.list(open)
		index := -1
		for i, m := range members {
			if m[0] <= lo && hi <= m[1] {
				index = i
			}
		}
		head := open - 1
		for t.at(head).kind == tkComma && t.at(head-1).kind == tkRParen {
			head = t.match[head-1] - 1
		}
		switch {
		case t.isWord(head, "values") || t.isWord(head, "value"):
			return partner{kind: partnerInsert, index: index}
		case t.isWord(head, "in"):
			q := head - 1
			if t.isWord(q, "not") {
				q--
			}
			if st := t.operandStart(q + 1); st >= 0 {
				return partner{kind: partnerOperand, start: st, end: q + 1}
			}
		case t.at(head).kind == tkWord && t.isCall(head):
			best, bestRank := -1, 3
			for i, m := range members {
				if i == index {
					continue
				}
				if r := t.rank(m[0], m[1]); r < bestRank {
					best, bestRank = i, r
				}
			}
			if best >= 0 {
				return partner{kind: partnerOperand, start: members[best][0], end: members[best][1]}
			}
		}
	}
	return partner{}
}

// rank orders operands by how much they say about a parameter beside them.
func (t text) rank(start, end int) int {
	switch tok := t.at(start); {
	case end-start == 1 && tok.kind == tkIdent:
		return 0
	case tok.kind == tkVar:
		return 3
	case end-start == 1:
		return 2
	}
	return 1
}

// idents lists the qualified identifiers in a token range: the columns an
// expression reads.
func (t text) idents(start, end int) []token {
	var out []token
	for i := start; i < end; i++ {
		if tok := t.toks[i]; tok.kind == tkIdent && len(tok.parts) >= 2 {
			out = append(out, tok)
		}
	}
	return out
}

// slice is the source text of a token range.
func (t text) slice(start, end int) string {
	if start >= end {
		return ""
	}
	return t.src[t.toks[start].start:t.toks[end-1].end]
}

// quote backticks an identifier.
func quote(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
