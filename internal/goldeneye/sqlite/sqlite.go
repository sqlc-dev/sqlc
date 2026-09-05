// Package sqlite generates the SQLite dialect seed under
// internal/engine/sqlite/dialect from sqlite3 shells built from the
// amalgamation sqlite.org publishes, run against in-memory databases that
// need no server, and from the amalgamation itself.
//
// SQLite describes its functions as far as their names, their kinds and the
// number of arguments each takes — pragma_function_list — and no further:
// it types values rather than columns or functions, so nothing in the
// database says what a function returns or what it expects. The source
// does, in its way. Every function is registered with the C functions that
// implement it, and those set their result through sqlite3_result_* and
// read their arguments through sqlite3_value_*, which is as close as SQLite
// comes to declaring a signature. So functions.jsonl is built from both:
// the shell says which functions exist, how many arguments each overload
// takes and whether it aggregates, and the amalgamation says what each
// returns and what its arguments hold. A function the shell reports that
// the source does not register, or whose implementation sets no result,
// fails generation rather than being guessed at. What neither can say —
// which scalar functions return NULL for arguments that are not — is a
// short list in signatures.go; for aggregates it is found by running each
// over no rows.
//
// Which functions a SQLite has is decided when it is compiled, so the
// dialect treats compile options the way the PostgreSQL dialect treats
// contrib extensions. functions.jsonl is what a build with the default
// options has, and each further option gets a directory under extensions/
// holding the functions a build with that option adds — found the way
// CREATE EXTENSION's additions are, by comparing the catalog with and
// without. SQLite has no catalog of types or operators, so types.jsonl and
// operators.jsonl are hand-written.
//
// The shells are built once per pinned version by Install.
package sqlite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
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

// key tells one overload from another.
func (r functionRow) key() string {
	return r.Name + "/" + strconv.Itoa(r.NArg)
}

// Version reports the release the default shell is.
func Version(ctx context.Context, dir string) (string, error) {
	out, err := exec.CommandContext(ctx, builds()[0].binary(dir), "--version").Output()
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

// shell is one built sqlite3 and the functions it reports.
type shell struct {
	build  build
	binary string
	rows   []functionRow
}

// readShell lists a build's functions, after checking that the shell was
// built with the options the build says — a cached shell from before the
// option lists changed would otherwise describe the wrong dialect.
func readShell(ctx context.Context, dir string, b build) (*shell, error) {
	s := &shell{build: b, binary: b.binary(dir)}
	var used []struct {
		Option string `json:"option"`
		Used   int    `json:"used"`
	}
	known := map[string]bool{}
	for _, b := range builds() {
		for _, opt := range b.flags() {
			known[opt] = true
		}
	}
	var clauses []string
	for _, opt := range slices.Sorted(maps.Keys(known)) {
		clauses = append(clauses, fmt.Sprintf("SELECT '%s' AS option, sqlite_compileoption_used('%s') AS used", opt, opt))
	}
	if err := query(ctx, s.binary, strings.Join(clauses, " UNION ALL "), &used); err != nil {
		return nil, err
	}
	for _, u := range used {
		want := 0
		if slices.Contains(b.flags(), u.Option) {
			want = 1
		}
		if u.Used != want {
			return nil, fmt.Errorf("sqlite: the %s shell was not built with the options it should have been (%s is %d): remove %s and run `go run ./cmd/goldeneye install sqlite` again", b.name, u.Option, u.Used, dir)
		}
	}
	if err := query(ctx, s.binary, functionList, &s.rows); err != nil {
		return nil, err
	}
	return s, nil
}

// generator accumulates the functions of every build, reading their
// signatures from the amalgamation, and remembers which names were
// reported so that the lists in signatures.go can be checked against the
// shells at the end.
type generator struct {
	ctx      context.Context
	src      *source
	reported map[string]bool
}

// overload is one row of a shell's list with what was found out about it.
type overload struct {
	row  functionRow
	kind string
	sig  signature
}

// functions turns rows into records, one per overload. Every row has to be
// a function the amalgamation registers unless omitted. In the default
// build only what the library builds in counts: the rest is what the shell
// registers on top — edit, sha3, the extensions it bundles — which is not
// the dialect's business, and which a comparison with the default build has
// already removed from an option's rows.
func (g *generator) functions(s *shell, base bool) ([]dialect.Function, error) {
	var overloads []overload
	seen := map[string]bool{}
	var missing []string
	for _, row := range s.rows {
		g.reported[row.Name] = true
		if omitted[row.Name] || seen[row.key()] || base && row.Builtin == 0 {
			continue
		}
		seen[row.key()] = true
		sig, err := g.src.signature(row.Name)
		if err != nil {
			if !seen[row.Name] {
				seen[row.Name] = true
				missing = append(missing, err.Error())
			}
			continue
		}
		kind, err := functionKind(g.ctx, s.binary, row)
		if err != nil {
			return nil, err
		}
		overloads = append(overloads, overload{row: row, kind: kind, sig: sig})
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("sqlite: %s build: %s", s.build.name, strings.Join(missing, "; "))
	}
	empty, err := overNoRows(g.ctx, s.binary, overloads)
	if err != nil {
		return nil, err
	}
	funcs := make([]dialect.Function, 0, len(overloads))
	for _, o := range overloads {
		isNullable := nullable[o.row.Name]
		if o.kind == "a" {
			isNullable = empty[o.row.key()] == "null"
		}
		funcs = append(funcs, dialect.Function{
			Name:     o.row.Name,
			Kind:     o.kind,
			Args:     o.sig.args(o.row.NArg),
			Returns:  o.sig.Returns,
			Nullable: isNullable,
		})
	}
	return funcs, nil
}

// overNoRows runs every aggregate over no rows and reports the type of
// what each returns, "null" for the ones — avg, max, group_concat — that
// return NULL when there is nothing to aggregate, and not count or total.
// The probes go through one shell process, each statement labelled with
// its overload so that the answers can be told apart.
func overNoRows(ctx context.Context, binary string, overloads []overload) (map[string]string, error) {
	var script strings.Builder
	for _, o := range overloads {
		if o.kind != "a" {
			continue
		}
		n := o.row.NArg
		if n < 0 {
			n = minArgs(n)
		}
		args := strings.TrimSuffix(strings.Repeat("x, ", n), ", ")
		fmt.Fprintf(&script, "SELECT '%s', typeof(%s(%s)) FROM (SELECT NULL AS x) WHERE 0;\n", o.row.key(), o.row.Name, args)
	}
	results := map[string]string{}
	if script.Len() == 0 {
		return results, nil
	}
	cmd := exec.CommandContext(ctx, binary, "-list", ":memory:", script.String())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return nil, fmt.Errorf("sqlite3: %w", err)
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		key, typ, ok := strings.Cut(line, "|")
		if ok {
			results[key] = typ
		}
	}
	for _, o := range overloads {
		if o.kind == "a" {
			if _, ok := results[o.row.key()]; !ok {
				return nil, fmt.Errorf("sqlite: cannot run %s over no rows: %s", o.row.Name, strings.TrimSpace(stderr.String()))
			}
		}
	}
	return results, nil
}

// added returns the rows of an option's shell that the default shell does
// not have: what the option adds.
func added(opt, base *shell) *shell {
	have := map[string]bool{}
	for _, row := range base.rows {
		have[row.key()] = true
	}
	diff := &shell{build: opt.build, binary: opt.binary}
	for _, row := range opt.rows {
		if !have[row.key()] {
			diff.rows = append(diff.rows, row)
		}
	}
	return diff
}

// Generate reads the dialect from the shells under dir and the
// amalgamation they were built from.
func Generate(ctx context.Context, dir string) (dialect.Files, error) {
	src, err := readSource(filepath.Join(dir, "src", "sqlite3.c"))
	if err != nil {
		return nil, err
	}
	g := &generator{ctx: ctx, src: src, reported: map[string]bool{}}
	all := builds()
	base, err := readShell(ctx, dir, all[0])
	if err != nil {
		return nil, err
	}
	funcs, err := g.functions(base, true)
	if err != nil {
		return nil, err
	}
	files := dialect.Files{}
	if files[dialect.FunctionsFile], err = dialect.JSONL(funcs); err != nil {
		return nil, err
	}
	for _, b := range all[1:] {
		s, err := readShell(ctx, dir, b)
		if err != nil {
			return nil, err
		}
		funcs, err := g.functions(added(s, base), false)
		if err != nil {
			return nil, err
		}
		if len(funcs) == 0 {
			return nil, fmt.Errorf("sqlite: %s adds no functions over the default build; drop it from the extensions in install.go", strings.Join(b.options, " "))
		}
		if files[path.Join(dialect.ExtensionsDir, b.name, dialect.FunctionsFile)], err = dialect.JSONL(funcs); err != nil {
			return nil, err
		}
	}
	var stale []string
	for _, list := range []map[string]bool{omitted, nullable} {
		for name := range list {
			if !g.reported[name] {
				stale = append(stale, name)
			}
		}
	}
	if len(stale) > 0 {
		sort.Strings(stale)
		return nil, fmt.Errorf("sqlite: signatures.go lists function(s) no shell reports: %s", strings.Join(stale, ", "))
	}
	return files, nil
}
