package spanner

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner/apiv1/spannerpb"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Spanner keeps no catalog of its types, functions or operators, so the
// library is read by asking: each candidate is compiled in PLAN mode, which
// runs nothing, against a probe database, and what Spanner accepts is kept
// with the type it reports.

// libraryDatabase is the probe database the library is read from.
const libraryDatabase = "goldeneye_library"

// librarySchema is the probe database: a table with a NOT NULL column of
// each type a candidate's argument can be, the search indexes SCORE and
// SEARCH need, the sequence GET_NEXT_SEQUENCE_VALUE draws from, and a
// commit timestamp column for PENDING_COMMIT_TIMESTAMP.
var librarySchema = []string{
	`CREATE TABLE probe (
  id INT64 NOT NULL,
  i INT64 NOT NULL,
  s STRING(MAX) NOT NULL,
  f FLOAT64 NOT NULL,
  f32 FLOAT32 NOT NULL,
  n NUMERIC NOT NULL,
  b BOOL NOT NULL,
  byt BYTES(MAX) NOT NULL,
  d DATE NOT NULL,
  ts TIMESTAMP NOT NULL,
  j JSON NOT NULL,
  ai ARRAY<INT64> NOT NULL,
  ast ARRAY<STRING(MAX)> NOT NULL,
  af ARRAY<FLOAT64> NOT NULL,
  af32 ARRAY<FLOAT32> NOT NULL,
  aj ARRAY<JSON> NOT NULL,
  ab ARRAY<BYTES(MAX)> NOT NULL,
  tl TOKENLIST AS (TOKENIZE_FULLTEXT(s)) HIDDEN,
  tl_substring TOKENLIST AS (TOKENIZE_SUBSTRING(s)) HIDDEN,
  tl_ngrams TOKENLIST AS (TOKENIZE_NGRAMS(s)) HIDDEN,
) PRIMARY KEY (id)`,
	`CREATE SEARCH INDEX probe_search ON probe(tl, tl_substring, tl_ngrams)`,
	`CREATE SEQUENCE probe_sequence OPTIONS (sequence_kind = 'bit_reversed_positive')`,
	`CREATE TABLE probe_commits (
  id INT64 NOT NULL,
  committed_at TIMESTAMP OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY (id)`,
}

// probeValues is the expression a candidate's argument of each type is
// called with: the probe table's column of that type, so that the call is
// made the way a query makes it, on a value rather than a constant.
var probeValues = map[string]string{
	"int64":          "i",
	"string":         "s",
	"float64":        "f",
	"float32":        "f32",
	"numeric":        "n",
	"bool":           "b",
	"bytes":          "byt",
	"date":           "d",
	"timestamp":      "ts",
	"json":           "j",
	"interval":       "INTERVAL 1 DAY",
	"tokenlist":      "tl",
	"array(int64)":   "ai",
	"array(string)":  "ast",
	"array(float64)": "af",
	"array(float32)": "af32",
	"array(json)":    "aj",
	"array(bytes)":   "ab",
}

// libraryTypes are the types the dialect lists, in the order types.jsonl
// does, with the category sqlc groups each in and the spelling Spanner is
// asked to cast to. A type Spanner does not have, BigQuery's BIGNUMERIC or
// DATETIME, is GoogleSQL's all the same, and the engine is GoogleSQL's,
// shared with BigQuery; those have no probe and are listed as they are.
var libraryTypes = []struct {
	dialect.Type
	probe string
}{
	{dialect.Type{Name: "bool", Category: "B"}, "BOOL"},
	{dialect.Type{Name: "int64", Category: "N"}, "INT64"},
	{dialect.Type{Name: "float32", Category: "N"}, "FLOAT32"},
	{dialect.Type{Name: "float64", Category: "N"}, "FLOAT64"},
	{dialect.Type{Name: "numeric", Category: "N"}, "NUMERIC"},
	{dialect.Type{Name: "bignumeric", Category: "N"}, ""},
	{dialect.Type{Name: "string", Category: "S"}, "STRING"},
	{dialect.Type{Name: "bytes", Category: "S"}, "BYTES"},
	{dialect.Type{Name: "uuid", Category: "S"}, "UUID"},
	{dialect.Type{Name: "date", Category: "D"}, "DATE"},
	{dialect.Type{Name: "datetime", Category: "D"}, ""},
	{dialect.Type{Name: "time", Category: "D"}, ""},
	{dialect.Type{Name: "timestamp", Category: "D"}, "TIMESTAMP"},
	{dialect.Type{Name: "interval", Category: "T"}, "INTERVAL"},
	{dialect.Type{Name: "json", Category: "U"}, "JSON"},
	{dialect.Type{Name: "geography", Category: "U"}, ""},
	{dialect.Type{Name: "struct", Category: "U"}, "STRUCT<a INT64>"},
	{dialect.Type{Name: "enum", Category: "U"}, ""},
	{dialect.Type{Name: "proto", Category: "U"}, ""},
	{dialect.Type{Name: "tokenlist", Category: "U"}, ""},
	{dialect.Type{Name: "array", Category: "A"}, "ARRAY<INT64>"},
	{dialect.Type{Name: "range", Category: "U"}, ""},
	{dialect.Type{Name: "map", Category: "U"}, ""},
}

//go:embed library.jsonl
var libraryCandidates []byte

// candidate is one line of library.jsonl.
type candidate struct {
	Name      string   `json:"name"`
	Kind      string   `json:"kind,omitempty"`
	Args      []string `json:"args,omitempty"`
	Returns   string   `json:"returns,omitempty"`
	Nullable  bool     `json:"nullable,omitempty"`
	NeverNull bool     `json:"never_null,omitempty"`
	// Expr is the expression to compile in place of name(args), and Query
	// the whole statement, with $1, $2... standing for the arguments.
	Expr  string `json:"expr,omitempty"`
	Query string `json:"query,omitempty"`
}

// polymorphic reports whether a candidate takes a T or an AT.
func (c candidate) polymorphic() bool {
	for _, a := range c.Args {
		if a == "T" || a == "AT" {
			return true
		}
	}
	return false
}

// sql is the statement that compiles a call of the candidate with the
// arguments given, as types.
func (c candidate) sql(args []string) (string, error) {
	vals := make([]string, len(args))
	for i, a := range args {
		if typ, lit, ok := strings.Cut(a, "="); ok {
			_ = typ
			vals[i] = lit
			continue
		}
		v, ok := probeValues[a]
		if !ok {
			if c.Expr == "" && c.Query == "" {
				return "", fmt.Errorf("%s: no probe value for %s", c.Name, a)
			}
			v = "NULL"
		}
		vals[i] = v
	}
	fill := func(tmpl string) string {
		for i := len(vals); i >= 1; i-- {
			tmpl = strings.ReplaceAll(tmpl, "$"+strconv.Itoa(i), vals[i-1])
		}
		return tmpl
	}
	if c.Query != "" {
		return fill(c.Query), nil
	}
	call := c.Name + "(" + strings.Join(vals, ", ") + ")"
	if c.Expr != "" {
		call = fill(c.Expr)
	}
	return "SELECT " + call + " AS r FROM probe", nil
}

// seedType is the type an argument is written as in functions.jsonl.
func seedType(arg string) string {
	switch arg {
	case "T":
		return "anyelement"
	case "AT":
		return "anyarray"
	}
	typ, _, _ := strings.Cut(arg, "=")
	return typ
}

// readLibrary reads types.jsonl, functions.jsonl and operators.jsonl from
// the probe database.
func readLibrary(ctx context.Context, s *server, session string) (dialect.Files, error) {
	compile := func(sql string) (*spannerpb.Type, error) {
		req := &spannerpb.ExecuteSqlRequest{Session: session, Sql: sql, QueryMode: spannerpb.ExecuteSqlRequest_PLAN}
		if isDML(sql) {
			req.Transaction = &spannerpb.TransactionSelector{Selector: &spannerpb.TransactionSelector_Begin{
				Begin: &spannerpb.TransactionOptions{Mode: &spannerpb.TransactionOptions_ReadWrite_{ReadWrite: &spannerpb.TransactionOptions_ReadWrite{}}},
			}}
			req.Seqno = 1
		}
		rs, err := s.data.ExecuteSql(ctx, req)
		if err != nil {
			return nil, err
		}
		if rs.Metadata != nil && rs.Metadata.Transaction != nil && len(rs.Metadata.Transaction.Id) > 0 {
			s.data.Rollback(context.WithoutCancel(ctx), &spannerpb.RollbackRequest{Session: session, TransactionId: rs.Metadata.Transaction.Id})
		}
		if rs.Metadata == nil || rs.Metadata.RowType == nil || len(rs.Metadata.RowType.Fields) == 0 {
			return nil, nil
		}
		return rs.Metadata.RowType.Fields[0].Type, nil
	}

	types, err := readLibraryTypes(compile)
	if err != nil {
		return nil, err
	}
	funcs, err := readLibraryFunctions(compile)
	if err != nil {
		return nil, err
	}
	operators, err := readLibraryOperators(compile)
	if err != nil {
		return nil, err
	}
	files := dialect.Files{}
	for name, records := range map[string]any{
		dialect.TypesFile:     types,
		dialect.FunctionsFile: funcs,
		dialect.OperatorsFile: operators,
	} {
		var blob []byte
		var err error
		switch r := records.(type) {
		case []dialect.Type:
			blob, err = dialect.JSONL(r)
		case []dialect.Function:
			blob, err = dialect.JSONL(r)
		case []dialect.Operator:
			blob, err = dialect.JSONL(r)
		}
		if err != nil {
			return nil, err
		}
		files[name] = blob
	}
	return files, nil
}

type compiler func(sql string) (*spannerpb.Type, error)

// readLibraryTypes checks that Spanner has each type the dialect lists with
// a probe, failing when one is gone, since sqlc would then accept a schema
// Spanner rejects.
func readLibraryTypes(compile compiler) ([]dialect.Type, error) {
	types := make([]dialect.Type, 0, len(libraryTypes))
	for _, t := range libraryTypes {
		if t.probe != "" {
			// IS NULL, since Spanner returns no STRUCT as a column.
			if _, err := compile("SELECT CAST(NULL AS " + t.probe + ") IS NULL AS r"); err != nil {
				return nil, fmt.Errorf("type %s: %w", t.Name, err)
			}
		}
		types = append(types, t.Type)
	}
	return types, nil
}

// readLibraryFunctions compiles every candidate in library.jsonl and lists
// the ones Spanner accepts, sorted by name and then by arity, the order the
// candidates give otherwise kept.
func readLibraryFunctions(compile compiler) ([]dialect.Function, error) {
	var funcs []dialect.Function
	scanner := bufio.NewScanner(bytes.NewReader(libraryCandidates))
	for line := 1; scanner.Scan(); line++ {
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		var c candidate
		if err := json.Unmarshal([]byte(text), &c); err != nil {
			return nil, fmt.Errorf("library.jsonl:%d: %w", line, err)
		}
		returns, ok, err := probeCandidate(compile, c)
		if err != nil {
			return nil, fmt.Errorf("library.jsonl:%d: %w", line, err)
		}
		if !ok {
			continue
		}
		fn := dialect.Function{Name: c.Name, Kind: c.Kind, Returns: returns, Nullable: c.Nullable, NeverNull: c.NeverNull}
		for _, a := range c.Args {
			fn.Args = append(fn.Args, dialect.Arg{Type: seedType(a)})
		}
		funcs = append(funcs, fn)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(funcs, func(i, j int) bool {
		if funcs[i].Name != funcs[j].Name {
			return funcs[i].Name < funcs[j].Name
		}
		return len(funcs[i].Args) < len(funcs[j].Args)
	})
	return funcs, nil
}

// probeCandidate compiles a candidate and returns the type it is seeded
// with, or false when Spanner rejects it. A candidate that is not there is
// left out rather than failing the run: a release of Spanner that drops a
// function shows in the diff of functions.jsonl.
func probeCandidate(compile compiler, c candidate) (string, bool, error) {
	if !c.polymorphic() {
		sql, err := c.sql(c.Args)
		if err != nil {
			return "", false, err
		}
		t, err := compile(sql)
		if err != nil {
			return "", false, nil
		}
		got := typeSpelling(t)
		if c.Query != "" && c.Returns != "" {
			// The statement's result, if it has one, is not the call's.
			return c.Returns, true, nil
		}
		if c.Returns != "" && c.Returns != got {
			return "", false, fmt.Errorf("%s returns %s, not %s", c.Name, got, c.Returns)
		}
		return got, true, nil
	}
	for _, elem := range []string{"int64", "string"} {
		args := make([]string, len(c.Args))
		for i, a := range c.Args {
			switch a {
			case "T":
				args[i] = elem
			case "AT":
				args[i] = "array(" + elem + ")"
			default:
				args[i] = a
			}
		}
		sql, err := c.sql(args)
		if err != nil {
			return "", false, err
		}
		t, err := compile(sql)
		if err != nil {
			return "", false, nil
		}
		want := c.Returns
		for i := len(args); i >= 1; i-- {
			want = strings.ReplaceAll(want, "$"+strconv.Itoa(i), args[i-1])
		}
		switch want {
		case "anyelement":
			want = elem
		case "anyarray":
			want = "array(" + elem + ")"
		}
		if got := typeSpelling(t); got != want {
			return "", false, fmt.Errorf("%s(%s) returns %s, not %s", c.Name, strings.Join(args, ", "), got, want)
		}
	}
	return c.Returns, true, nil
}

// typeSpelling spells a type the way functions.jsonl does.
func typeSpelling(t *spannerpb.Type) string {
	if t == nil {
		return ""
	}
	if t.Code == spannerpb.TypeCode_ARRAY {
		return "array(" + typeSpelling(t.ArrayElementType) + ")"
	}
	return strings.ToLower(t.Code.String())
}

// libraryOperandTypes are the operands each operator is tried with.
var (
	numericOperands = []string{"int64", "float64", "float32", "numeric"}
	operatorPairs   = func() [][3]string {
		var pairs [][3]string
		for _, op := range []string{"+", "-", "*", "/"} {
			for _, l := range numericOperands {
				for _, r := range numericOperands {
					pairs = append(pairs, [3]string{op, l, r})
				}
			}
		}
		for _, t := range []string{"string", "bytes"} {
			pairs = append(pairs, [3]string{"||", t, t})
		}
		for _, op := range []string{"+", "-"} {
			pairs = append(pairs, [3]string{op, "timestamp", "interval"}, [3]string{op, "interval", "interval"})
		}
		return append(pairs,
			[3]string{"+", "interval", "timestamp"},
			[3]string{"-", "date", "date"}, [3]string{"-", "timestamp", "timestamp"},
			[3]string{"*", "interval", "int64"}, [3]string{"/", "interval", "int64"},
			[3]string{"&", "int64", "int64"}, [3]string{"|", "int64", "int64"}, [3]string{"^", "int64", "int64"},
			[3]string{"<<", "int64", "int64"}, [3]string{">>", "int64", "int64"},
			[3]string{"&", "bytes", "bytes"}, [3]string{"|", "bytes", "bytes"}, [3]string{"^", "bytes", "bytes"},
		)
	}()
)

// readLibraryOperators compiles each binary operator over the operands it
// is tried with and lists the ones Spanner accepts with their result. The
// comparisons are left to dialect.json's rules, which make every one of
// them a BOOL.
func readLibraryOperators(compile compiler) ([]dialect.Operator, error) {
	var ops []dialect.Operator
	for _, p := range operatorPairs {
		t, err := compile("SELECT " + probeValues[p[1]] + " " + p[0] + " " + probeValues[p[2]] + " AS r FROM probe")
		if err != nil {
			continue
		}
		ops = append(ops, dialect.Operator{Name: p[0], Left: p[1], Right: p[2], Result: typeSpelling(t)})
	}
	return ops, nil
}
