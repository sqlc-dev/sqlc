package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/kyleconroy/sqlc/internal/cmd"
)

func TestExamples(t *testing.T) {
	t.Parallel()
	examples, err := filepath.Abs(filepath.Join("..", "..", "examples"))
	if err != nil {
		t.Fatal(err)
	}

	files, err := ioutil.ReadDir(examples)
	if err != nil {
		t.Fatal(err)
	}

	for _, replay := range files {
		if !replay.IsDir() {
			continue
		}
		tc := replay.Name()
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(examples, tc)
			var stderr bytes.Buffer
			output, err := cmd.Generate(cmd.Env{}, path, &stderr)
			if err != nil {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, output)
		})
	}
}

type testConfig struct {
	ExperimentalParserOnly bool `json:"experimental_parser_only"`
}

func TestReplay(t *testing.T) {
	t.Parallel()
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, replay := range files {
		if !replay.IsDir() {
			continue
		}
		tc := replay.Name()
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path, _ := filepath.Abs(filepath.Join("testdata", tc))
			var stderr bytes.Buffer
			expected := expectedStderr(t, path)
			output, err := cmd.Generate(cmd.Env{ExperimentalParser: true}, path, &stderr)
			if len(expected) == 0 && err != nil {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, output)
			if diff := cmp.Diff(expected, stderr.String()); diff != "" {
				t.Errorf("stderr differed (-want +got):\n%s", diff)
			}
		})
		t.Run(tc+"/deprecated-parser", func(t *testing.T) {
			t.Parallel()
			path, _ := filepath.Abs(filepath.Join("testdata", tc))
			// TODO: Extract test configuration into a method
			confPath := filepath.Join(path, "endtoend.json")
			if _, err := os.Stat(confPath); !os.IsNotExist(err) {
				blob, err := ioutil.ReadFile(confPath)
				if err != nil {
					t.Fatal(err)
				}
				var conf testConfig
				if err := json.Unmarshal(blob, &conf); err != nil {
					t.Fatal(err)
				}
				if conf.ExperimentalParserOnly {
					t.Skip("experimental parser only")
				}
			}
			var stderr bytes.Buffer
			expected := expectedStderr(t, path)
			output, err := cmd.Generate(cmd.Env{ExperimentalParser: false}, path, &stderr)
			if len(expected) == 0 && err != nil {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, output)
			if diff := cmp.Diff(expected, stderr.String()); diff != "" {
				t.Errorf("stderr differed (-want +got):\n%s", diff)
			}
		})
	}
}

func cmpDirectory(t *testing.T, dir string, actual map[string]string) {
	expected := map[string]string{}
	var ff = func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if file.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".kt") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") || strings.Contains(path, "src/test/") {
			return nil
		}
		// TODO(mightyguava): Remove this after sqlc-kotlin-runtime is published to Maven.
		if strings.HasSuffix(path, "Query.kt") {
			return nil
		}
		blob, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		expected[path] = string(blob)
		return nil
	}
	if err := filepath.Walk(dir, ff); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(expected, actual, cmpopts.EquateEmpty()) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			name := name
			tn := strings.Replace(name, dir+"/", "", -1)
			t.Run(tn, func(t *testing.T) {
				if actual[name] == "" {
					t.Errorf("%s is empty", name)
					return
				}
				if diff := cmp.Diff(contents, actual[name]); diff != "" {
					t.Errorf("%s differed (-want +got):\n%s", name, diff)
				}
			})
		}
	}
}

func expectedStderr(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "stderr.txt")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		blob, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return string(blob)
	}
	return ""
}
