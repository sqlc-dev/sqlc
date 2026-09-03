// Package endtoend finds the analyze cases under internal/endtoend/testdata
// and compares an engine's own answer with the output a case committed,
// the way the dialect package compares a generated dialect with the
// committed one.
package endtoend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Case is one analyze case: the files sqlc analyze ran with, the fixture
// loaded before the queries run against a real database, and the output
// sqlc committed.
type Case struct {
	Name    string // analyze_params/clickhouse
	Dir     string
	Schema  string
	Query   string
	Fixture string // empty when the case has no fixture.sql
	Output  string
}

// Testdata returns the end-to-end testdata directory, found relative to this
// source file so the working directory does not matter.
func Testdata() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("cannot locate the testcheck source directory")
	}
	dir := filepath.Join(filepath.Dir(file), "..", "..", "endtoend", "testdata")
	if _, err := os.Stat(dir); err != nil {
		return "", err
	}
	return filepath.Clean(dir), nil
}

// Cases lists the analyze cases for an engine: every analyze_*/<engine>
// directory whose exec.json runs the analyze command. A case that asks for
// the AST is skipped, since only sqlc can print that.
func Cases(engine string) ([]Case, error) {
	root, err := Testdata()
	if err != nil {
		return nil, err
	}
	dirs, err := filepath.Glob(filepath.Join(root, "analyze_*", engine))
	if err != nil {
		return nil, err
	}
	var cases []Case
	for _, dir := range dirs {
		c, ok, err := load(dir)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", dir, err)
		}
		if ok {
			cases = append(cases, c)
		}
	}
	return cases, nil
}

func load(dir string) (Case, bool, error) {
	blob, err := os.ReadFile(filepath.Join(dir, "exec.json"))
	if errors.Is(err, os.ErrNotExist) {
		return Case{}, false, nil
	}
	if err != nil {
		return Case{}, false, err
	}
	var exec struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}
	if err := json.Unmarshal(blob, &exec); err != nil {
		return Case{}, false, fmt.Errorf("exec.json: %w", err)
	}
	if exec.Command != "analyze" {
		return Case{}, false, nil
	}
	var schema, query string
	for i := 0; i < len(exec.Args); i++ {
		switch arg := exec.Args[i]; arg {
		case "--schema", "-s", "--dialect", "-d":
			i++
			if arg == "--schema" || arg == "-s" {
				if i < len(exec.Args) {
					schema = exec.Args[i]
				}
			}
		case "--ast":
			return Case{}, false, nil
		default:
			if !strings.HasPrefix(arg, "-") {
				query = arg
			}
		}
	}
	if schema == "" || query == "" {
		return Case{}, false, errors.New("exec.json: analyze needs --schema and a query file")
	}
	c := Case{
		Name:   filepath.Join(filepath.Base(filepath.Dir(dir)), filepath.Base(dir)),
		Dir:    dir,
		Schema: filepath.Join(dir, schema),
		Query:  filepath.Join(dir, query),
		Output: filepath.Join(dir, "output.json"),
	}
	if fixture := filepath.Join(dir, "fixture.sql"); fileExists(fixture) {
		c.Fixture = fixture
	}
	return c, true, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Compare checks an engine's answer against the case's committed output
// byte for byte and returns a line diff when they differ, or "" when they
// match.
func (c Case) Compare(got []byte) (string, error) {
	want, err := os.ReadFile(c.Output)
	if err != nil {
		return "", err
	}
	if bytes.Equal(want, got) {
		return "", nil
	}
	return dialect.Diff(string(want), string(got)), nil
}
