package sqlite

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// source is what the amalgamation says about how SQLite's functions answer.
// A function is registered in one of a few shapes — the FuncDef macros of
// the built-in tables, a sqlite3_create_function call, or a struct table an
// extension walks — each naming the C functions that implement it. Those
// implementations set their result through sqlite3_result_* and read their
// arguments through sqlite3_value_*, which is as close as SQLite comes to
// declaring a signature.
type source struct {
	funcs   map[string]cfunc
	aliases map[string]string
	regs    map[string]*registration
}

// cfunc is one C function definition: its parameter list and its body.
type cfunc struct {
	params string
	body   string
}

// registration is the union of everything the source registers under one
// SQL function name.
type registration struct {
	// scalar implements a scalar overload: its results and its arguments.
	scalar []string
	// step and final implement an aggregate or window overload: step reads
	// the arguments, final and value set the result.
	step  []string
	final []string
	// inline names the INLINEFUNC_* constant of a function the VDBE
	// implements in bytecode, which has no C body to read.
	inline string
	// json records what a JSON function's registration says it returns,
	// "text" or "blob", since the json_ and jsonb_ forms share an
	// implementation; jsonAlways says the registration promises JSON text
	// whatever the implementation might otherwise return, as -> does.
	json       string
	jsonAlways bool
	// table marks a registration found only in a struct table, which is
	// consulted when nothing more direct registered the name.
	table bool
}

// The FuncDef macros, with the positions of the C functions in each.
var macroEntry = regexp.MustCompile(`\b(FUNCTION2|FUNCTION|VFUNCTION|SFUNCTION|MFUNCTION|JFUNCTION|INLINE_FUNC|DFUNCTION|PURE_DATE|STR_FUNCTION|LIKEFUNC|WAGGREGATE|WINDOWFUNCX|WINDOWFUNCALL|WINDOWFUNCNOOP)\(`)

var (
	createFunction = regexp.MustCompile(`\bsqlite3_create_(window_)?function\(`)
	defineAlias    = regexp.MustCompile(`(?m)^#define\s+(\w+)\s+(\w+)\s*$`)
	definition     = regexp.MustCompile(`(?m)^(?:static\s+|SQLITE_PRIVATE\s+)?(?:const\s+)?(?:unsigned\s+)?[A-Za-z_]\w*(?:\s*\*+\s*|\s+)([A-Za-z_]\w*)\(`)
	initializer    = regexp.MustCompile(`\{[^{}]*"([a-z_0-9>-]+)"[^{}]*\}`)
	identifier     = regexp.MustCompile(`\b([A-Za-z_]\w*)\b`)
	resultCall     = regexp.MustCompile(`\bsqlite3_result_(\w+)\(`)
	valueCall      = regexp.MustCompile(`\bsqlite3_value_(\w+)\(\s*(?:argv|apVal|apArg)\[(\w+)\]`)
	callee         = regexp.MustCompile(`\b([A-Za-z_]\w*)\s*\(`)
)

// readSource reads the amalgamation.
func readSource(path string) (*source, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(data)
	s := &source{
		funcs:   map[string]cfunc{},
		aliases: map[string]string{},
		regs:    map[string]*registration{},
	}
	for _, m := range defineAlias.FindAllStringSubmatch(text, -1) {
		s.aliases[m[1]] = m[2]
	}
	s.readDefinitions(text)
	s.readMacros(text)
	s.readCreateFunctions(text)
	s.readTables(text)
	return s, nil
}

// readDefinitions indexes every function defined at the left margin.
func (s *source) readDefinitions(text string) {
	for _, m := range definition.FindAllStringSubmatchIndex(text, -1) {
		name := text[m[2]:m[3]]
		params, end, ok := balanced(text, m[1]-1, '(', ')')
		if !ok {
			continue
		}
		i := end
		for i < len(text) && (text[i] == ' ' || text[i] == '\n' || text[i] == '\r' || text[i] == '\t') {
			i++
		}
		if i >= len(text) || text[i] != '{' {
			continue
		}
		body, _, ok := balanced(text, i, '{', '}')
		if !ok {
			continue
		}
		if _, dup := s.funcs[name]; !dup {
			s.funcs[name] = cfunc{params: params, body: body}
		}
	}
}

// balanced returns the text between the bracket at open and its match,
// skipping string and character literals and comments, and the index after
// the closing bracket.
func balanced(text string, open int, lb, rb byte) (string, int, bool) {
	depth := 0
	for i := open; i < len(text); i++ {
		switch c := text[i]; {
		case c == '"' || c == '\'':
			for i++; i < len(text) && text[i] != c; i++ {
				if text[i] == '\\' {
					i++
				}
			}
		case c == '/' && i+1 < len(text) && text[i+1] == '*':
			end := strings.Index(text[i+2:], "*/")
			if end < 0 {
				return "", 0, false
			}
			i += end + 3
		case c == '/' && i+1 < len(text) && text[i+1] == '/':
			end := strings.IndexByte(text[i:], '\n')
			if end < 0 {
				return "", 0, false
			}
			i += end
		case c == lb:
			depth++
		case c == rb:
			depth--
			if depth == 0 {
				return text[open+1 : i], i + 1, true
			}
		}
	}
	return "", 0, false
}

// arguments splits a bracketed argument list at the commas of its own
// level, and reports the index after the closing bracket.
func arguments(text string, open int) ([]string, int, bool) {
	inner, end, ok := balanced(text, open, '(', ')')
	if !ok {
		return nil, 0, false
	}
	var args []string
	depth, start := 0, 0
	for i := 0; i < len(inner); i++ {
		switch inner[i] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case '"':
			for i++; i < len(inner) && inner[i] != '"'; i++ {
				if inner[i] == '\\' {
					i++
				}
			}
		case ',':
			if depth == 0 {
				args = append(args, strings.TrimSpace(inner[start:i]))
				start = i + 1
			}
		}
	}
	args = append(args, strings.TrimSpace(inner[start:]))
	return args, end, true
}

func (s *source) reg(name string) *registration {
	name = strings.ToLower(name)
	r, ok := s.regs[name]
	if !ok {
		r = &registration{}
		s.regs[name] = r
	}
	return r
}

// readMacros reads the FuncDef tables. The macro definitions themselves
// match too, and are told apart by their formal parameter names.
func (s *source) readMacros(text string) {
	for _, m := range macroEntry.FindAllStringSubmatchIndex(text, -1) {
		macro := text[m[2]:m[3]]
		args, _, ok := arguments(text, m[1]-1)
		if !ok || len(args) < 2 || args[0] == "zName" || args[0] == "name" {
			continue
		}
		name := args[0]
		switch macro {
		case "FUNCTION", "FUNCTION2", "VFUNCTION", "SFUNCTION", "DFUNCTION", "PURE_DATE", "STR_FUNCTION":
			if len(args) > 4 {
				r := s.reg(name)
				r.scalar = append(r.scalar, args[4])
			}
		case "MFUNCTION":
			if len(args) > 3 {
				r := s.reg(name)
				r.scalar = append(r.scalar, args[3])
			}
		case "JFUNCTION":
			// JFUNCTION(zName, nArg, bUseCache, bWS, bRS, bJsonB, iArg, xFunc)
			if len(args) > 7 {
				r := s.reg(name)
				r.scalar = append(r.scalar, args[7])
				r.json = "text"
				if args[5] == "1" {
					r.json = "blob"
				}
				r.jsonAlways = strings.Contains(args[6], "JSON_JSON")
			}
		case "INLINE_FUNC":
			if len(args) > 2 {
				s.reg(name).inline = args[2]
			}
		case "LIKEFUNC":
			r := s.reg(name)
			r.scalar = append(r.scalar, "likeFunc")
		case "WAGGREGATE":
			// WAGGREGATE(zName, nArg, arg, nc, xStep, xFinal, xValue, xInverse, f)
			// The JSON aggregates use it too, told by their implementation,
			// with JSON_BLOB as the user data of the jsonb_ forms.
			if len(args) > 6 {
				r := s.reg(name)
				r.step = append(r.step, args[4])
				r.final = append(r.final, args[5], args[6])
				if strings.HasPrefix(args[4], "json") {
					r.json = "text"
					if strings.Contains(args[2], "JSON_BLOB") {
						r.json = "blob"
					}
				}
			}
		case "WINDOWFUNCX", "WINDOWFUNCALL":
			r := s.reg(name)
			r.step = append(r.step, name+"StepFunc")
			r.final = append(r.final, name+"ValueFunc")
		case "WINDOWFUNCNOOP":
			s.reg(name).inline = "bytecode"
		}
	}
}

// readCreateFunctions reads the functions extensions register by calling
// sqlite3_create_function or sqlite3_create_window_function with a literal
// name.
func (s *source) readCreateFunctions(text string) {
	for _, m := range createFunction.FindAllStringSubmatchIndex(text, -1) {
		window := m[2] != -1
		args, _, ok := arguments(text, m[1]-1)
		if !ok || len(args) < 8 || !strings.HasPrefix(args[1], `"`) {
			continue
		}
		r := s.reg(strings.Trim(args[1], `"`))
		if window {
			// (db, zName, nArg, eTextRep, pApp, xStep, xFinal, xValue, xInverse, xDestroy)
			r.step = given(r.step, args[5])
			r.final = given(r.final, args[6], args[7])
		} else {
			// (db, zName, nArg, eTextRep, pApp, xFunc, xStep, xFinal)
			r.scalar = given(r.scalar, args[5])
			r.step = given(r.step, args[6])
			r.final = given(r.final, args[7])
		}
	}
}

// given appends the implementations a registration names, leaving out the
// methods it passes as 0 or NULL.
func given(impls []string, names ...string) []string {
	for _, name := range names {
		if name != "0" && name != "NULL" {
			impls = append(impls, name)
		}
	}
	return impls
}

// readTables reads the struct tables extensions walk to register their
// functions — geopoly's aFunc, FTS5's aBuiltin, FTS3's aOverload — each row
// an initializer holding the function's name and the C functions that
// implement it. Any C function in the row counts as an implementation;
// which reads arguments and which sets the result comes out in the
// reading. A row is consulted only for a name nothing else registered.
func (s *source) readTables(text string) {
	for _, m := range initializer.FindAllStringSubmatch(text, -1) {
		var impls []string
		for _, id := range identifier.FindAllStringSubmatch(m[0], -1) {
			if _, ok := s.funcs[s.resolve(id[1])]; ok {
				impls = append(impls, id[1])
			}
		}
		if len(impls) == 0 {
			continue
		}
		name := strings.ToLower(m[1])
		r, ok := s.regs[name]
		if ok && !r.table {
			continue
		}
		r = s.reg(name)
		r.table = true
		r.scalar = append(r.scalar, impls...)
	}
}

// resolve follows #define aliases from a C function name to the one that
// is defined.
func (s *source) resolve(name string) string {
	for i := 0; i < 8; i++ {
		alias, ok := s.aliases[name]
		if !ok {
			return name
		}
		name = alias
	}
	return name
}

// noop reports a C function that does nothing: the stand-in for a function
// the VDBE implements in bytecode, whose result is one of its arguments.
func noop(name string) bool {
	return strings.HasPrefix(name, "noop")
}

// The kinds a result or argument call names, by the sqlite3_result_* or
// sqlite3_value_* suffix with any 64 dropped. Suffixes not listed — error,
// subtype, type, bytes, dup — say nothing about a type.
var (
	resultKinds = map[string]string{
		"int": "integer", "double": "real",
		"text": "text", "text16": "text", "text16le": "text", "text16be": "text",
		"blob": "blob", "zeroblob": "blob",
		"value": "any", "pointer": "any",
	}
	valueKinds = map[string]string{
		"int": "integer", "double": "real",
		"text": "text", "text16": "text", "text16le": "text", "text16be": "text",
		"blob": "blob",
	}
)

func kindOf(kinds map[string]string, suffix string) (string, bool) {
	k, ok := kinds[strings.TrimSuffix(suffix, "64")]
	return k, ok
}

// results collects the kinds a C function sets its result to, following
// the helpers it calls a few levels down, since many functions hand their
// result to one.
func (s *source) results(name string, depth int, seen map[string]bool) map[string]bool {
	kinds := map[string]bool{}
	name = s.resolve(name)
	if noop(name) {
		kinds["any"] = true
		return kinds
	}
	fn, ok := s.funcs[name]
	if !ok || seen[name] || depth > 3 {
		return kinds
	}
	seen[name] = true
	for _, m := range resultCall.FindAllStringSubmatch(fn.body, -1) {
		if k, ok := kindOf(resultKinds, m[1]); ok {
			kinds[k] = true
		}
	}
	for _, m := range callee.FindAllStringSubmatch(fn.body, -1) {
		c := m[1]
		if c == name || strings.HasPrefix(c, "sqlite3_") || strings.HasPrefix(c, "sqlite3Vdbe") {
			continue
		}
		if _, ok := s.funcs[c]; ok {
			for k := range s.results(c, depth+1, seen) {
				kinds[k] = true
			}
		}
	}
	return kinds
}

// args collects the kinds a C function reads each of its arguments as, by
// position, and the kind it reads a run of arguments as under a loop
// index. An implementation called through the FTS5 extension API is handed
// the table as an implicit first argument, so its positions shift by one.
func (s *source) args(name string, positions map[int]map[string]bool, variadic map[string]bool) {
	fn, ok := s.funcs[s.resolve(name)]
	if !ok {
		return
	}
	shift := 0
	if strings.Contains(fn.params, "Fts5ExtensionApi") {
		shift = 1
		if positions[0] == nil {
			positions[0] = map[string]bool{}
		}
		positions[0]["any"] = true
	}
	for _, m := range valueCall.FindAllStringSubmatch(fn.body, -1) {
		k, ok := kindOf(valueKinds, m[1])
		if !ok {
			continue
		}
		var pos int
		if _, err := fmt.Sscanf(m[2], "%d", &pos); err != nil {
			variadic[k] = true
			continue
		}
		pos += shift
		if positions[pos] == nil {
			positions[pos] = map[string]bool{}
		}
		positions[pos][k] = true
	}
}

// single reduces the kinds seen at one position to a type: the one kind
// seen, or "any" for a mixture or nothing.
func single(kinds map[string]bool) string {
	if len(kinds) == 1 {
		for k := range kinds {
			return k
		}
	}
	return "any"
}

// signature derives what the source says a SQL function returns and takes.
// A result of one kind is that type. A function that returns one of its
// arguments, or a mixture of kinds, takes the type of its first argument,
// which the seed spells "any" — except that integer and real together widen
// to real, as SQLite's own arithmetic does, and text and blob together to
// text, since a function that returns either is handing back the bytes it
// was given, and the legacy compiler cannot follow "any" to an argument.
func (s *source) signature(name string) (signature, error) {
	r, ok := s.regs[strings.ToLower(name)]
	if !ok {
		return signature{}, fmt.Errorf("the amalgamation registers no function named %s", name)
	}
	var sig signature
	switch {
	case r.inline != "":
		sig.Returns = inlineReturns[r.inline]
		if sig.Returns == "" {
			sig.Returns = "any"
		}
	default:
		kinds := map[string]bool{}
		for _, fn := range append(append([]string{}, r.scalar...), r.final...) {
			for k := range s.results(fn, 0, map[string]bool{}) {
				kinds[k] = true
			}
		}
		jsonOnly := r.json != "" && len(kinds) <= 2 && !kinds["integer"] && !kinds["real"] && !kinds["any"]
		switch {
		case len(kinds) == 0:
			return signature{}, fmt.Errorf("cannot tell what %s returns: no sqlite3_result call in %s", name, strings.Join(append(append([]string{}, r.scalar...), r.final...), ", "))
		case r.jsonAlways:
			sig.Returns = "text"
		case kinds["any"]:
			sig.Returns = "any"
		case len(kinds) == 1:
			for k := range kinds {
				sig.Returns = k
			}
		case jsonOnly:
			// The shared implementation writes JSON text or JSONB; the
			// registration says which this form gets.
			sig.Returns = r.json
		case len(kinds) == 2 && kinds["integer"] && kinds["real"]:
			sig.Returns = "real"
		case len(kinds) == 2 && kinds["text"] && kinds["blob"]:
			sig.Returns = "text"
		default:
			sig.Returns = "any"
		}
	}

	positions := map[int]map[string]bool{}
	variadic := map[string]bool{}
	for _, fn := range append(append([]string{}, r.scalar...), r.step...) {
		s.args(fn, positions, variadic)
	}
	if len(positions) > 0 {
		var order []int
		for pos := range positions {
			order = append(order, pos)
		}
		sort.Ints(order)
		sig.Args = make([]string, order[len(order)-1]+1)
		for i := range sig.Args {
			sig.Args[i] = single(positions[i])
		}
		// Trailing positions nothing typed say nothing, unless a variadic
		// tail follows them, when they keep it from starting early.
		for len(variadic) == 0 && len(sig.Args) > 0 && sig.Args[len(sig.Args)-1] == "any" {
			sig.Args = sig.Args[:len(sig.Args)-1]
		}
	}
	if len(variadic) == 1 {
		sig.Variadic = single(variadic)
	}
	return sig, nil
}
