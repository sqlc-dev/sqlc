package preprocess_test

import (
	"encoding/json"
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
// the input and the expected results, grouped by engine:
//
//	testdata/<engine>/<case>/input.sql        the query file
//	testdata/<engine>/<case>/output.sql       the rewritten SQL
//	testdata/<engine>/<case>/side_table.json  what the preprocessor recorded
//	testdata/<engine>/<case>/stderr.txt       reported errors (only when invalid)
//
// Regenerate the golden files with `go test ./internal/sql/preprocess -update`.
func TestRewrite(t *testing.T) {
	for _, dir := range cases(t) {
		t.Run(filepath.ToSlash(dir), func(t *testing.T) {
			path := filepath.Join("testdata", dir)
			engine := config.Engine(filepath.Dir(dir))

			input := readFile(t, filepath.Join(path, "input.sql"))
			res := preprocess.File(engine, input)

			compare(t, filepath.Join(path, "output.sql"), res.Text)
			compare(t, filepath.Join(path, "side_table.json"), sideTable(res))
			compare(t, filepath.Join(path, "stderr.txt"), errors(input, res))

			// A rewrite only ever replaces text with something shorter or
			// equal in lines, so a line number taken from the rewritten text is
			// never past the end of the original.
			if want, got := strings.Count(input, "\n"), strings.Count(res.Text, "\n"); got > want {
				t.Errorf("rewrite added lines: got %d, want at most %d", got, want)
			}
		})
	}
}

// The JSON shape of the side table the preprocessor records alongside the
// rewritten SQL. Locations are offsets into the rewritten text; origins are the
// offsets they map back to in the original source.
type statement struct {
	Start      int         `json:"start"`
	End        int         `json:"end"`
	Dollar     bool        `json:"dollar"`
	Params     []parameter `json:"params,omitempty"`
	Embeds     []embed     `json:"embeds,omitempty"`
	Error      string      `json:"error,omitempty"`
	ErrorAt    int         `json:"error_at,omitempty"`
	ParamError string      `json:"param_error,omitempty"`
}

type parameter struct {
	Number   int    `json:"number"`
	Location int    `json:"location"`
	Origin   int    `json:"origin"`
	Name     string `json:"name,omitempty"`
	// Named is false for a placeholder the user wrote themselves, which the
	// preprocessor numbers but does not name.
	Named bool `json:"named"`
	// Nullable is true when sqlc.narg() forced the parameter nullable, even
	// though the schema says it is NOT NULL.
	Nullable bool `json:"nullable,omitempty"`
	Slice    bool `json:"slice,omitempty"`
}

type embed struct {
	Table    string `json:"table"`
	Location int    `json:"location"`
	Origin   int    `json:"origin"`
	Orig     string `json:"orig"`
}

// sideTable renders everything the preprocessor recorded, reading it back
// through the same API the compiler uses.
func sideTable(res *preprocess.Result) string {
	out := make([]statement, 0, len(res.Statements()))
	for _, stmt := range res.Statements() {
		s := statement{Start: stmt.Start, End: stmt.End, Dollar: stmt.Dollar}
		if stmt.Err != nil {
			s.Error = stmt.Err.Error()
			if e, ok := stmt.Err.(*sqlerr.Error); ok {
				s.ErrorAt = e.Location
			}
		}
		if stmt.ParamErr != nil {
			s.ParamError = stmt.ParamErr.Error()
		}

		locations := make([]int, 0, len(stmt.Numbers))
		for loc := range stmt.Numbers {
			locations = append(locations, loc)
		}
		sort.Ints(locations)
		for _, loc := range locations {
			number := stmt.Numbers[loc]
			name, _ := stmt.Params.NameFor(number)
			// Merging against a NOT NULL parameter is how the compiler asks
			// whether the user overrode nullability with sqlc.narg().
			p, isNamed := stmt.Params.FetchMerge(number, named.NewInferredParam(name, true))
			_, isSlice := stmt.Slices[loc]
			s.Params = append(s.Params, parameter{
				Number:   number,
				Location: loc,
				Origin:   res.Origin(loc),
				Name:     name,
				Named:    isNamed,
				Nullable: isNamed && !p.NotNull(),
				Slice:    isSlice,
			})
		}

		for _, e := range stmt.Embeds {
			s.Embeds = append(s.Embeds, embed{
				Table:    e.Table.Name,
				Location: e.Location,
				Origin:   res.Origin(e.Location),
				Orig:     e.Orig(),
			})
		}
		out = append(out, s)
	}

	blob, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(blob) + "\n"
}

// errors renders every validation error the preprocessor recorded, one per
// line, as "line:column: message" against the original source.
func errors(input string, res *preprocess.Result) string {
	var out []string
	for _, stmt := range res.Statements() {
		for _, err := range []error{stmt.Err, stmt.ParamErr} {
			if err == nil {
				continue
			}
			loc := stmt.Start
			if e, ok := err.(*sqlerr.Error); ok && e.Location != 0 {
				loc = e.Location
			}
			line, column := source.LineNumber(input, res.Origin(loc))
			out = append(out, fmt.Sprintf("%d:%d: %s", line, column, err))
		}
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, "\n") + "\n"
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
