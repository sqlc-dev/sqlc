package sqlite

import "github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"

// signature is what a function returns and what its arguments hold, read
// from the amalgamation by source.signature. Args types the leading
// arguments in order; a position it does not cover is "any", which the
// analyzer resolves to the argument's own type. Variadic types the
// arguments an overload of variable arity repeats, "any" when empty. How
// many arguments an overload takes, and how many of them it requires,
// comes from the shell. Nullable is decided afterwards: for an aggregate by
// running it over no rows, for a scalar by the nullable list below.
type signature struct {
	Args     []string
	Variadic string
	Returns  string
	Nullable bool
}

// args builds an overload's parameters from the shell's argument count. A
// fixed count is that many parameters; a negative count is the required
// leading parameters, any further ones the signature types marked as
// having a default, and a variadic tail.
func (s signature) args(narg int) []dialect.Arg {
	if narg >= 0 {
		args := make([]dialect.Arg, narg)
		for i := range args {
			args[i] = dialect.Arg{Type: s.argType(i)}
		}
		return args
	}
	required := minArgs(narg)
	n := max(required, len(s.Args))
	args := make([]dialect.Arg, 0, n+1)
	for i := 0; i < n; i++ {
		args = append(args, dialect.Arg{Type: s.argType(i), HasDefault: i >= required})
	}
	tail := s.Variadic
	if tail == "" {
		tail = "any"
	}
	return append(args, dialect.Arg{Type: tail, Mode: "v"})
}

func (s signature) argType(i int) string {
	if i < len(s.Args) {
		return s.Args[i]
	}
	return "any"
}

// inlineReturns is what the functions the VDBE implements in bytecode
// return, by the INLINEFUNC_* constant their registration carries. All but
// one hand back one of their arguments, which is what the default of "any"
// says; sqlite_offset is a byte offset.
var inlineReturns = map[string]string{
	"INLINEFUNC_sqlite_offset": "integer",
}

// omitted are functions the dialect leaves out: ones that exist for their
// side effect and return nothing a query can use, and ones an extension
// uses to pass pointers to itself.
var omitted = map[string]bool{
	// Loads a shared library and returns NULL.
	"load_extension": true,
	// Writes to the error log and returns NULL.
	"sqlite_log": true,
	// FTS3's and FTS5's ways of passing pointers to their virtual tables,
	// not functions a query calls.
	"fts3_tokenizer": true,
	"fts5":           true,
	// A debugging aid of GEOPOLY's.
	"geopoly_debug": true,
}

// nullable are the scalar functions that return NULL for arguments that
// are not: a lookup that finds nothing, an input that does not parse, a
// value with no sign. The source cannot say this — a SQLite function
// returns NULL as often by setting no result as by calling
// sqlite3_result_null — so it is listed by the documentation of each.
// Aggregates are not listed: whether one returns NULL over no rows is
// found by running it over none.
var nullable = map[string]bool{
	// Core functions.
	"nullif":                   true,
	"sign":                     true,
	"sqlite_compileoption_get": true,
	"sqlite_offset":            true,
	"unhex":                    true,
	// Date and time functions, for a time value they cannot parse.
	"date":      true,
	"datetime":  true,
	"julianday": true,
	"strftime":  true,
	"time":      true,
	"timediff":  true,
	"unixepoch": true,
	// JSON functions, for a path that leads nowhere.
	"->":                true,
	"->>":               true,
	"json_array_length": true,
	"json_extract":      true,
	"json_type":         true,
	"jsonb_extract":     true,
	// Window functions, for a row outside the frame.
	"first_value": true,
	"lag":         true,
	"last_value":  true,
	"lead":        true,
	"nth_value":   true,
	// FTS5, for a table with no locale.
	"fts5_get_locale": true,
	// GEOPOLY functions, for an argument that is not a polygon.
	"geopoly_area":           true,
	"geopoly_bbox":           true,
	"geopoly_blob":           true,
	"geopoly_ccw":            true,
	"geopoly_contains_point": true,
	"geopoly_json":           true,
	"geopoly_overlap":        true,
	"geopoly_svg":            true,
	"geopoly_within":         true,
	"geopoly_xform":          true,
}
