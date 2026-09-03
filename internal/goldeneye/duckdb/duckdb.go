// Package duckdb generates the DuckDB dialect seed under
// internal/engine/duckdb/dialect — types.jsonl, functions.jsonl and
// operators.jsonl — from a live DuckDB CLI, the same way the postgresql
// package generates PostgreSQL's from a live server. The CLI must be the
// DuckDB 2.0 build darkwing is pinned against; it is located through the
// DUCKDB environment variable, falling back to "duckdb" on PATH.
package duckdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name of the engine directory the dialect lives under.
const Engine = "duckdb"

// Locate finds the DuckDB CLI: the DUCKDB environment variable wins, then
// "duckdb" on PATH.
func Locate() (string, error) {
	if path := os.Getenv("DUCKDB"); path != "" {
		return path, nil
	}
	path, err := exec.LookPath("duckdb")
	if err != nil {
		return "", errors.New("no duckdb CLI found: set DUCKDB to the DuckDB 2.0 binary darkwing is pinned against, or put duckdb on PATH")
	}
	return path, nil
}

// Version reports the release a CLI is.
func Version(ctx context.Context, binary string) (string, error) {
	out, err := exec.CommandContext(ctx, binary, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("duckdb --version: %w", err)
	}
	return "DuckDB " + strings.TrimSpace(string(out)), nil
}

// query runs a SQL statement against an in-memory database and decodes the
// CLI's JSON output into rows.
func query(ctx context.Context, binary, sql string, rows any) error {
	cmd := exec.CommandContext(ctx, binary, "-json", ":memory:", "-c", sql)
	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("duckdb: %s: %s", err, exit.Stderr)
		}
		return fmt.Errorf("duckdb: %w", err)
	}
	// An empty result set prints nothing rather than [].
	if len(out) == 0 {
		return nil
	}
	return json.Unmarshal(out, rows)
}

type typeRow struct {
	TypeName    string  `json:"type_name"`
	LogicalType string  `json:"logical_type"`
	Category    *string `json:"type_category"`
}

type functionRow struct {
	Name           string   `json:"function_name"`
	FunctionType   string   `json:"function_type"`
	ParameterTypes []string `json:"parameter_types"`
	Varargs        *string  `json:"varargs"`
	ReturnType     *string  `json:"return_type"`
}

// metaTypes are type ids that never describe a column's value: sentinels and
// binder-internal types the dump lists alongside the real ones.
var metaTypes = map[string]bool{
	"null":    true,
	"unknown": true,
	"invalid": true,
	"any":     true,
	"type":    true,
	"lambda":  true,
	"table":   true,
	"pointer": true,
}

// categoryLetter maps duckdb_types() categories onto the PostgreSQL category
// letters the seed package uses.
func categoryLetter(category *string) string {
	if category == nil {
		return "U"
	}
	switch *category {
	case "NUMERIC":
		return "N"
	case "STRING":
		return "S"
	case "DATETIME":
		return "D"
	case "BOOLEAN":
		return "B"
	case "COMPOSITE":
		return "C"
	default:
		return "U"
	}
}

func readTypes(ctx context.Context, binary string) ([]dialect.Type, error) {
	var rows []typeRow
	err := query(ctx, binary, `
SELECT type_name, logical_type, type_category
FROM duckdb_types()
WHERE database_name = 'system'
ORDER BY type_name`, &rows)
	if err != nil {
		return nil, err
	}

	// Group the dump's one-row-per-spelling by logical type: the spelling
	// matching the logical type id is the canonical name, the rest are
	// aliases.
	grouped := map[string]*dialect.Type{}
	var order []string
	for _, row := range rows {
		logical := strings.ToLower(row.LogicalType)
		if metaTypes[logical] {
			continue
		}
		t, ok := grouped[logical]
		if !ok {
			t = &dialect.Type{Name: logical, Category: categoryLetter(row.Category)}
			grouped[logical] = t
			order = append(order, logical)
		}
		if t.Category == "U" {
			t.Category = categoryLetter(row.Category)
		}
		if name := strings.ToLower(row.TypeName); name != logical {
			t.Aliases = append(t.Aliases, name)
		}
	}

	sort.Strings(order)
	types := make([]dialect.Type, 0, len(order)+1)
	// "any" stands in for the generic parameters of DuckDB's polymorphic
	// functions (ANY, T, K, V); the analyzer resolves a call returning it
	// to the type of the call's first argument.
	types = append(types, dialect.Type{Name: "any", Category: "U"})
	for _, name := range order {
		types = append(types, *grouped[name])
	}
	return types, nil
}

// typeNames is the set of names the seed declares, for filtering out function
// overloads over generic or binder-internal types.
func typeNames(types []dialect.Type) map[string]bool {
	names := map[string]bool{}
	for _, t := range types {
		names[t.Name] = true
		for _, alias := range t.Aliases {
			names[alias] = true
		}
	}
	return names
}

// seedTypeName maps a duckdb_functions() type spelling to a seeded type name:
// lowercased, with modifiers dropped (DECIMAL(18,3) is a decimal), array
// markers kept (fixed sizes and generic bounds included: DOUBLE[3] and
// ANY[] are arrays), and the generic parameters of polymorphic functions
// mapped to "any". It reports false for the internal types a seed cannot
// describe.
func seedTypeName(name string, known map[string]bool) (string, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	if close := strings.LastIndexByte(name, ']'); close == len(name)-1 {
		if open := strings.LastIndexByte(name, '['); open != -1 {
			elementName, ok := seedTypeName(name[:open], known)
			return elementName + "[]", ok
		}
	}
	if open := strings.IndexByte(name, '('); open != -1 {
		name = name[:open]
	}
	switch name {
	case "any", "t", "k", "v":
		return "any", true
	}
	if !known[name] {
		return "", false
	}
	return name, true
}

// isOperatorName reports a function named by symbols rather than an
// identifier — the spelling of a binary or prefix operator.
func isOperatorName(name string) bool {
	return strings.IndexFunc(name, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' && r != '$'
	}) != -1
}

func functionKind(functionType string) string {
	switch functionType {
	case "aggregate":
		return "a"
	case "window":
		return "w"
	default:
		return ""
	}
}

func readFunctions(ctx context.Context, binary string, known map[string]bool) ([]dialect.Function, []dialect.Operator, error) {
	var rows []functionRow
	err := query(ctx, binary, `
SELECT DISTINCT function_name, function_type, parameter_types, varargs, return_type
FROM duckdb_functions()
WHERE database_name = 'system'
  AND schema_name = 'main'
  AND function_type IN ('scalar', 'aggregate', 'window', 'macro')
ORDER BY function_name, parameter_types::VARCHAR, return_type`, &rows)
	if err != nil {
		return nil, nil, err
	}

	var funcs []dialect.Function
	var operators []dialect.Operator
	seenFunc := map[string]bool{}
	seenOp := map[string]bool{}
	for _, row := range rows {
		if row.ReturnType == nil {
			continue
		}
		returns, ok := seedTypeName(*row.ReturnType, known)
		if !ok {
			continue
		}
		args := make([]string, 0, len(row.ParameterTypes))
		resolved := true
		for _, param := range row.ParameterTypes {
			arg, ok := seedTypeName(param, known)
			if !ok {
				resolved = false
				break
			}
			args = append(args, arg)
		}
		if !resolved {
			continue
		}

		if isOperatorName(row.Name) {
			// Operator-named functions become operator entries; sqlc's
			// engine reports them as operator applications. Only the
			// binary form over concrete scalar types is representable —
			// the seed's operator table, unlike its function table, does
			// not register generic or array types on demand.
			if row.FunctionType != "scalar" || len(args) != 2 || row.Varargs != nil {
				continue
			}
			if concrete := !strings.Contains(returns, "any") && !strings.Contains(returns, "[]") &&
				!strings.Contains(args[0], "any") && !strings.Contains(args[0], "[]") &&
				!strings.Contains(args[1], "any") && !strings.Contains(args[1], "[]"); !concrete {
				continue
			}
			key := row.Name + "\x00" + strings.Join(args, "\x00")
			if seenOp[key] {
				continue
			}
			seenOp[key] = true
			operators = append(operators, dialect.Operator{
				Name:   row.Name,
				Left:   args[0],
				Right:  args[1],
				Result: returns,
			})
			continue
		}

		fn := dialect.Function{
			Name:    row.Name,
			Kind:    functionKind(row.FunctionType),
			Returns: returns,
			// An aggregate over no rows returns NULL — except count,
			// which returns 0.
			Nullable: row.FunctionType == "aggregate" && !strings.HasPrefix(row.Name, "count"),
		}
		for _, arg := range args {
			fn.Args = append(fn.Args, dialect.Arg{Type: arg})
		}
		if row.Varargs != nil {
			vararg, ok := seedTypeName(*row.Varargs, known)
			if !ok {
				continue
			}
			fn.Args = append(fn.Args, dialect.Arg{Type: vararg, Mode: "v"})
		}

		key := fn.Name + "\x00" + fn.Kind + "\x00" + strings.Join(args, "\x00")
		if seenFunc[key] {
			continue
		}
		seenFunc[key] = true
		funcs = append(funcs, fn)
	}
	return funcs, operators, nil
}

// Generate reads the dialect from the CLI.
func Generate(ctx context.Context, binary string) (dialect.Files, error) {
	types, err := readTypes(ctx, binary)
	if err != nil {
		return nil, err
	}
	funcs, operators, err := readFunctions(ctx, binary, typeNames(types))
	if err != nil {
		return nil, err
	}
	files := dialect.Files{}
	if files[dialect.TypesFile], err = dialect.JSONL(types); err != nil {
		return nil, err
	}
	if files[dialect.FunctionsFile], err = dialect.JSONL(funcs); err != nil {
		return nil, err
	}
	if files[dialect.OperatorsFile], err = dialect.JSONL(operators); err != nil {
		return nil, err
	}
	return files, nil
}
