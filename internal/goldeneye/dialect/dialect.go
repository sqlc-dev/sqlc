// Package dialect describes the files a dialect seed is made of — the JSONL
// records under internal/engine/<engine>/dialect that give an engine its
// type system and standard library — and provides what a generator needs to
// write them and a check needs to compare them with what is committed.
//
// The record types mirror the ones in internal/core/seed, which reads the
// files. They are repeated here rather than imported so that this module
// never shares code with the analysis it checks: the files are the contract,
// and a change to their shape has to be made on both sides.
package dialect

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// The files a dialect directory is made of.
const (
	SettingsFile  = "dialect.json"
	TypesFile     = "types.jsonl"
	OperatorsFile = "operators.jsonl"
	CastsFile     = "casts.jsonl"
	FunctionsFile = "functions.jsonl"
	RelationsFile = "relations.jsonl"
)

// ExtensionsDir is the directory under a dialect holding one directory per
// extension the dialect knows, each a smaller bundle of the same files.
const ExtensionsDir = "extensions"

// Type is a type the dialect defines. Aliases are spellings of the same type
// that a schema may use in a column definition.
type Type struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Aliases  []string `json:"aliases,omitempty"`
}

// Operator is a single operator overload.
type Operator struct {
	Name   string `json:"name"`
	Left   string `json:"left"`
	Right  string `json:"right"`
	Result string `json:"result"`
}

// Cast is a single cast between two types. Context is 'i'mplicit,
// 'a'ssignment or 'e'xplicit.
type Cast struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Context string `json:"context,omitempty"`
}

// Function is a function the dialect ships with. Kind is 'f'unction,
// 'a'ggregate, 'w'indow or 'p'rocedure.
type Function struct {
	Name     string `json:"name"`
	Kind     string `json:"kind,omitempty"`
	Args     []Arg  `json:"args,omitempty"`
	Returns  string `json:"returns"`
	Nullable bool   `json:"nullable,omitempty"`
}

// Arg is one of a function's parameters. Mode is 'i'n, 'o'ut, 'b'oth,
// 't'able or 'v'ariadic, and defaults to in.
type Arg struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type"`
	Mode       string `json:"mode,omitempty"`
	HasDefault bool   `json:"has_default,omitempty"`
}

// Relation is a table or view the dialect ships with, such as one of
// PostgreSQL's system catalogs. Kind is 'r' for a table or 'v' for a view,
// and defaults to a table.
type Relation struct {
	Catalog string   `json:"catalog,omitempty"`
	Schema  string   `json:"schema"`
	Name    string   `json:"name"`
	Kind    string   `json:"kind,omitempty"`
	Columns []Column `json:"columns"`
}

// Column is one of a relation's columns.
type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	NotNull bool   `json:"not_null,omitempty"`
	Array   bool   `json:"array,omitempty"`
	Length  int    `json:"length,omitempty"`
}

// Files is what a generator produces: the content of each file it writes,
// keyed by slash-separated path relative to the dialect directory, such as
// "functions.jsonl" or "extensions/hstore/types.jsonl". A generator owns
// only the files it produces; the rest of a dialect, dialect.json above all,
// is written by hand and is neither generated nor checked.
type Files map[string][]byte

// JSONL encodes records one per line, the way the committed files are
// written: encoding/json's compact form, with a newline after each record.
func JSONL[T any](records []T) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, record := range records {
		if err := enc.Encode(record); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// Dir returns the dialect directory of an engine,
// internal/engine/<engine>/dialect, found relative to this source file so
// the working directory does not matter.
func Dir(engine string) (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("cannot locate the goldeneye source directory")
	}
	dir := filepath.Join(filepath.Dir(file), "..", "..", "engine", engine, "dialect")
	if _, err := os.Stat(dir); err != nil {
		return "", err
	}
	return filepath.Clean(dir), nil
}

// Names lists the files in name order, so that a run is reported and
// written the same way every time.
func (f Files) Names() []string {
	names := make([]string, 0, len(f))
	for name := range f {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Write writes the files into dir, creating extension directories as
// needed. Files the generator did not produce are left alone.
func Write(dir string, files Files) error {
	for _, name := range files.Names() {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, files[name], 0o644); err != nil {
			return err
		}
	}
	return nil
}

// Check compares the files with the ones committed in dir and returns a
// report of every file that differs, or "" when each matches byte for byte.
func Check(dir string, files Files) (string, error) {
	var report strings.Builder
	for _, name := range files.Names() {
		want, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(&report, "%s: not committed\n", name)
			continue
		}
		if err != nil {
			return "", err
		}
		if bytes.Equal(want, files[name]) {
			continue
		}
		fmt.Fprintf(&report, "%s (-committed +database)\n%s", name, Diff(string(want), string(files[name])))
	}
	return report.String(), nil
}

// Diff is a line diff of two texts, marking lines only in want with "-" and
// lines only in got with "+", showing two lines of context around each
// change and eliding the rest. The committed files run to thousands of
// lines and usually differ in a handful, so the common head and tail are
// stripped before the longest common subsequence of the middle is found.
func Diff(want, got string) string {
	a := lines(want)
	b := lines(got)
	head := 0
	for head < len(a) && head < len(b) && a[head] == b[head] {
		head++
	}
	tail := 0
	for tail < len(a)-head && tail < len(b)-head && a[len(a)-1-tail] == b[len(b)-1-tail] {
		tail++
	}
	var edits []edit
	for _, line := range a[:head] {
		edits = append(edits, edit{' ', line})
	}
	edits = append(edits, lcsEdits(a[head:len(a)-tail], b[head:len(b)-tail])...)
	for _, line := range a[len(a)-tail:] {
		edits = append(edits, edit{' ', line})
	}

	const context = 2
	show := make([]bool, len(edits))
	for i, e := range edits {
		if e.kind == ' ' {
			continue
		}
		for j := max(0, i-context); j <= min(len(edits)-1, i+context); j++ {
			show[j] = true
		}
	}
	var out strings.Builder
	elided := false
	for i, e := range edits {
		if !show[i] {
			elided = true
			continue
		}
		if elided && out.Len() > 0 {
			out.WriteString("...\n")
		}
		elided = false
		fmt.Fprintf(&out, "%c %s\n", e.kind, e.line)
	}
	return out.String()
}

type edit struct {
	kind byte // ' ', '-' or '+'
	line string
}

func lines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

// lcsEdits is the edit script between a and b by longest common subsequence.
func lcsEdits(a, b []string) []edit {
	lcs := make([][]int, len(a)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	var edits []edit
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			edits = append(edits, edit{' ', a[i]})
			i++
			j++
		case lcs[i+1][j] >= lcs[i][j+1]:
			edits = append(edits, edit{'-', a[i]})
			i++
		default:
			edits = append(edits, edit{'+', b[j]})
			j++
		}
	}
	for ; i < len(a); i++ {
		edits = append(edits, edit{'-', a[i]})
	}
	for ; j < len(b); j++ {
		edits = append(edits, edit{'+', b[j]})
	}
	return edits
}
