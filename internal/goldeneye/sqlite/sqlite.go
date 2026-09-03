// Package sqlite generates the SQLite dialect seed under
// internal/engine/sqlite/dialect from the sqlite3 shell sqlite.org
// publishes, run against an in-memory database that needs no server.
//
// SQLite describes its functions as far as their names, their kinds and the
// number of arguments each takes — pragma_function_list — and no further:
// it types values rather than columns or functions, so nothing in the
// database says what a function returns or what it expects. functions.jsonl
// is therefore built from both sides. The shell says which functions exist,
// how many arguments each overload takes and whether it aggregates, and the
// signatures table in this package says what each returns and what its
// arguments are meant to hold. A built-in function the shell reports that
// the table does not know fails generation rather than being guessed at,
// and so does a table entry the shell does not report. SQLite has no
// catalog of types or operators, so types.jsonl and operators.jsonl are
// hand-written.
//
// The shell is downloaded once per pinned version by Install, or supplied
// through the SQLITE3 environment variable.
package sqlite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name of the engine directory the dialect lives under.
const Engine = "sqlite"

// functionList is every function the connection knows, one row per
// overload: name, whether the library builds it in or the shell or an
// extension registered it, 's', 'a' or 'w' for how it is called, and how
// many arguments it takes. Fixed arities come before the variable one so
// that an exact match is found first by anything reading the file in
// order.
const functionList = `
SELECT name, builtin, type, narg
FROM pragma_function_list
ORDER BY name, narg < 0, narg`

type functionRow struct {
	Name    string `json:"name"`
	Builtin int    `json:"builtin"`
	Type    string `json:"type"`
	NArg    int    `json:"narg"`
}

// Version reports the release a shell is.
func Version(ctx context.Context, binary string) (string, error) {
	out, err := exec.CommandContext(ctx, binary, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("sqlite3 --version: %w", err)
	}
	return "SQLite " + strings.TrimSpace(string(out)), nil
}

// query runs a SQL statement against an in-memory database and decodes the
// shell's JSON output into rows.
func query(ctx context.Context, binary, sql string, rows any) error {
	cmd := exec.CommandContext(ctx, binary, "-json", "-bail", ":memory:", sql)
	out, err := cmd.Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return fmt.Errorf("sqlite3: %s: %s", err, bytes.TrimSpace(exit.Stderr))
		}
		return fmt.Errorf("sqlite3: %w", err)
	}
	// An empty result set prints nothing rather than [].
	if len(bytes.TrimSpace(out)) == 0 {
		return nil
	}
	return json.Unmarshal(out, rows)
}

// minArgs decodes a negative argument count from pragma_function_list,
// which prints the count SQLite keeps for the function as it is: -1 for any
// number of arguments, and the two values reserved for built-ins, -3 for
// one or more and -4 for two or more — see matchQuality in the SQLite
// source. No function is counted as -2.
func minArgs(narg int) int {
	if narg < -2 {
		return -2 - narg
	}
	return 0
}

// functionKind maps the type the shell reports onto the seed's letters.
// The shell says 's' for a scalar function, 'a' for an aggregate and 'w'
// for one with a window implementation — but every one of SQLite's
// aggregates has one, so 'w' covers avg and count as well as row_number.
// The two are told apart by calling the function without an OVER clause,
// which only a window function refuses.
func functionKind(ctx context.Context, binary string, row functionRow) (string, error) {
	switch row.Type {
	case "s":
		return "", nil
	case "a":
		return "a", nil
	case "w":
		window, err := isWindowFunction(ctx, binary, row)
		if err != nil {
			return "", err
		}
		if window {
			return "w", nil
		}
		return "a", nil
	}
	return "", fmt.Errorf("sqlite: function %s has unknown type %q", row.Name, row.Type)
}

// isWindowFunction reports whether the shell refuses to call a function
// without an OVER clause. Any other complaint — an aggregate objecting to
// the NULLs it is handed — means the call was accepted as an aggregate.
func isWindowFunction(ctx context.Context, binary string, row functionRow) (bool, error) {
	n := row.NArg
	if n < 0 {
		n = minArgs(n)
	}
	args := strings.TrimSuffix(strings.Repeat("NULL, ", n), ", ")
	cmd := exec.CommandContext(ctx, binary, ":memory:", fmt.Sprintf("SELECT %s(%s)", row.Name, args))
	var stderr bytes.Buffer
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return false, nil
	}
	if strings.Contains(stderr.String(), "misuse of window function") {
		return true, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return false, nil
	}
	return false, fmt.Errorf("sqlite3: %w", err)
}

// readFunctions turns the shell's list into the dialect's, one record per
// overload of every function the signatures table knows. A built-in
// function without a signature is an error; a function without one that
// the shell or a bundled extension registered is not the dialect's
// business.
func readFunctions(ctx context.Context, binary string, rows []functionRow) ([]dialect.Function, error) {
	var funcs []dialect.Function
	seen := map[string]bool{}
	reported := map[string]bool{}
	var missing []string
	for _, row := range rows {
		reported[row.Name] = true
		if omitted[row.Name] {
			continue
		}
		sig, ok := signatures[row.Name]
		if !ok {
			if row.Builtin != 0 && !seen[row.Name] {
				seen[row.Name] = true
				missing = append(missing, row.Name)
			}
			continue
		}
		key := row.Name + "\x00" + strconv.Itoa(row.NArg)
		if seen[key] {
			continue
		}
		seen[key] = true
		kind, err := functionKind(ctx, binary, row)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, dialect.Function{
			Name:     row.Name,
			Kind:     kind,
			Args:     sig.args(row.NArg),
			Returns:  sig.Returns,
			Nullable: sig.Nullable,
		})
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("sqlite: no signature for built-in function(s) %s: add them to signatures.go", strings.Join(missing, ", "))
	}
	var stale []string
	for name := range signatures {
		if !reported[name] {
			stale = append(stale, name)
		}
	}
	if len(stale) > 0 {
		sort.Strings(stale)
		return nil, fmt.Errorf("sqlite: signatures.go lists function(s) the shell does not report: %s", strings.Join(stale, ", "))
	}
	return funcs, nil
}

// Generate reads the dialect from the shell.
func Generate(ctx context.Context, binary string) (dialect.Files, error) {
	var rows []functionRow
	if err := query(ctx, binary, functionList, &rows); err != nil {
		return nil, err
	}
	funcs, err := readFunctions(ctx, binary, rows)
	if err != nil {
		return nil, err
	}
	functions, err := dialect.JSONL(funcs)
	if err != nil {
		return nil, err
	}
	return dialect.Files{dialect.FunctionsFile: functions}, nil
}
