package preprocess_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/preprocess"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

var update = flag.Bool("update", false, "update the testdata golden files")

// TestRewrite runs every case under testdata. Each case is a directory holding
// an input.sql, the expected output.sql and, when the input is invalid, a
// stderr.txt with the reported error. Cases are grouped by engine:
//
//	testdata/<engine>/<case>/input.sql
//	testdata/<engine>/<case>/output.sql
//	testdata/<engine>/<case>/stderr.txt   (optional)
func TestRewrite(t *testing.T) {
	for _, dir := range cases(t) {
		t.Run(filepath.ToSlash(dir), func(t *testing.T) {
			path := filepath.Join("testdata", dir)
			engine := config.Engine(filepath.Dir(dir))

			input := readFile(t, filepath.Join(path, "input.sql"))
			res, err := preprocess.File(engine, input)
			if err != nil {
				t.Fatalf("preprocess.File: %s", err)
			}

			compare(t, filepath.Join(path, "output.sql"), res.Text)
			compare(t, filepath.Join(path, "stderr.txt"), errors(input, res))
		})
	}
}

// cases returns every testdata directory that holds an input.sql, relative to
// testdata and in a stable order.
func cases(t *testing.T) []string {
	t.Helper()
	var dirs []string
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != "input.sql" {
			return nil
		}
		rel, err := filepath.Rel("testdata", filepath.Dir(path))
		if err != nil {
			return err
		}
		dirs = append(dirs, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		t.Fatal("no test cases found in testdata")
	}
	sort.Strings(dirs)
	return dirs
}

// errors renders every validation error the preprocessor recorded, one per
// line, as "line:column: message" against the original source.
func errors(input string, res *preprocess.Result) string {
	var out []string
	for offset := 0; offset < len(res.Text); {
		stmt := res.Statement(offset)
		if stmt.End <= offset {
			break
		}
		offset = stmt.End
		if stmt.Err == nil {
			continue
		}
		loc := 0
		if e, ok := stmt.Err.(*sqlerr.Error); ok {
			loc = res.Origin(e.Location)
		}
		line, column := source.LineNumber(input, loc)
		out = append(out, fmt.Sprintf("%d:%d: %s", line, column, stmt.Err))
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, "\n") + "\n"
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	blob, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(blob)
}

// compare checks actual against a golden file. An empty expectation means the
// file should not exist.
func compare(t *testing.T, path, actual string) {
	t.Helper()
	if *update {
		if actual == "" {
			os.Remove(path)
			return
		}
		if err := os.WriteFile(path, []byte(actual), 0644); err != nil {
			t.Fatal(err)
		}
		return
	}
	var expected string
	if blob, err := os.ReadFile(path); err == nil {
		expected = string(blob)
	} else if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("%s differed (-want +got):\n%s", filepath.Base(path), diff)
	}
}

func mustRewrite(t *testing.T, engine config.Engine, src string) *preprocess.Result {
	t.Helper()
	res, err := preprocess.File(engine, src)
	if err != nil {
		t.Fatalf("preprocess.File(%s): %s", engine, err)
	}
	return res
}

func TestParamSet(t *testing.T) {
	res := mustRewrite(t, config.EnginePostgreSQL,
		"SELECT * FROM users WHERE a = sqlc.arg(alpha) AND b = sqlc.narg(beta) AND c = ANY(sqlc.slice(gamma));")
	stmt := res.Statement(0)

	for num, want := range map[int]string{1: "alpha", 2: "beta", 3: "gamma"} {
		got, ok := stmt.Params.NameFor(num)
		if !ok {
			t.Fatalf("no name for parameter %d", num)
		}
		if got != want {
			t.Errorf("parameter %d: got %q, want %q", num, got, want)
		}
	}

	// sqlc.narg() marks the parameter nullable even when inference says
	// otherwise.
	beta, _ := stmt.Params.FetchMerge(2, named.NewInferredParam("beta", true))
	if beta.NotNull() {
		t.Error("sqlc.narg parameter should be nullable")
	}
	gamma, _ := stmt.Params.FetchMerge(3, named.NewParam("gamma"))
	if !gamma.IsSqlcSlice() {
		t.Error("sqlc.slice parameter should be marked as a slice")
	}
}

func TestEmbedSpans(t *testing.T) {
	src := "SELECT sqlc.embed(a), b.* FROM a, b;"
	res := mustRewrite(t, config.EnginePostgreSQL, src)
	stmt := res.Statement(0)
	if len(stmt.Embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(stmt.Embeds))
	}
	e := stmt.Embeds[0]
	if e.Table.Name != "a" {
		t.Errorf("embed table: got %q, want %q", e.Table.Name, "a")
	}
	if want := strings.Index(res.Text, "a.*"); e.Location != want {
		t.Errorf("embed location: got %d, want %d", e.Location, want)
	}
	// A star reference the user wrote must not be mistaken for an embed.
	if _, ok := stmt.Embeds.Find(strings.Index(res.Text, "b.*")); ok {
		t.Error("user-written b.* was reported as an embed")
	}
	if got, want := e.Orig(), "sqlc.embed(a)"; got != want {
		t.Errorf("Orig: got %q, want %q", got, want)
	}
}

func TestSliceSpans(t *testing.T) {
	src := "SELECT * FROM t WHERE a IN (sqlc.slice(ids)) AND b = sqlc.arg(x);"
	res := mustRewrite(t, config.EngineMySQL, src)
	stmt := res.Statement(0)
	if len(stmt.Slices) != 1 {
		t.Fatalf("expected 1 slice, got %d", len(stmt.Slices))
	}
	// The recorded offset points at the placeholder itself, past the marker,
	// which is where the engine reports the parameter node.
	want := strings.Index(res.Text, "*/?") + len("*/")
	if got, ok := stmt.Slices[want]; !ok {
		t.Errorf("slice offsets %v do not contain %d", stmt.Slices, want)
	} else if got != "ids" {
		t.Errorf("slice name: got %q, want %q", got, "ids")
	}
}

func TestOriginMapping(t *testing.T) {
	src := "SELECT * FROM t WHERE a = sqlc.arg(alpha) AND b = sqlc.arg(beta);"
	res := mustRewrite(t, config.EnginePostgreSQL, src)

	// Text before the first rewrite maps to itself.
	if got := res.Origin(3); got != 3 {
		t.Errorf("Origin(3): got %d, want 3", got)
	}
	// Text after a rewrite maps back across the length change.
	newB, oldB := strings.Index(res.Text, "b ="), strings.Index(src, "b =")
	if got := res.Origin(newB); got != oldB {
		t.Errorf("Origin(%d): got %d, want %d", newB, got, oldB)
	}
	// An offset inside a placeholder maps to the start of what it replaced.
	newParam := strings.Index(res.Text, "$2")
	oldParam := strings.Index(src, "sqlc.arg(beta)")
	if got := res.Origin(newParam); got != oldParam {
		t.Errorf("Origin(%d): got %d, want %d", newParam, got, oldParam)
	}
}

func TestLinesArePreserved(t *testing.T) {
	// Downstream error reporting counts lines in the rewritten text, so a
	// rewrite must never add or remove one.
	src := "-- name: Get :one\nSELECT *\nFROM users\nWHERE id = sqlc.arg(id)\n  AND name = @name;\n"
	res := mustRewrite(t, config.EnginePostgreSQL, src)
	if want, got := strings.Count(src, "\n"), strings.Count(res.Text, "\n"); want != got {
		t.Errorf("line count changed: got %d, want %d", got, want)
	}
}

func TestStatementLookup(t *testing.T) {
	src := "SELECT sqlc.arg(a);\nSELECT sqlc.arg(b), sqlc.arg(c);\n"
	res := mustRewrite(t, config.EnginePostgreSQL, src)
	first := res.Statement(0)
	second := res.Statement(strings.Index(res.Text, "$1, $2"))
	if first == second {
		t.Fatal("expected two distinct statements")
	}
	if name, _ := first.Params.NameFor(1); name != "a" {
		t.Errorf("first statement parameter 1: got %q, want %q", name, "a")
	}
	if name, _ := second.Params.NameFor(2); name != "c" {
		t.Errorf("second statement parameter 2: got %q, want %q", name, "c")
	}
}

func TestUnknownEngine(t *testing.T) {
	if _, err := preprocess.File(config.Engine("nope"), "SELECT 1;"); err == nil {
		t.Fatal("expected an error for an unknown engine")
	}
}
