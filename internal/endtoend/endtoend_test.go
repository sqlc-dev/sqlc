package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyleconroy/sqlc/internal/cmd"
)

func TestExamples(t *testing.T) {
	t.Parallel()

	examples, _ := filepath.Abs(filepath.Join("..", "..", "examples"))
	var stderr bytes.Buffer

	output, err := cmd.Generate(examples, &stderr)
	if err != nil {
		t.Fatalf("%s", stderr.String())
	}

	cmpDirectory(t, examples, output)
}

func TestReplay(t *testing.T) {
	t.Parallel()

	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, replay := range files {
		tc := replay.Name()
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path, _ := filepath.Abs(filepath.Join("testdata", tc))
			var stderr bytes.Buffer
			output, err := cmd.Generate(path, &stderr)
			if err != nil {
				t.Fatalf("%s", stderr.String())
			}
			cmpDirectory(t, path, output)
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

	if len(expected) == 0 {
		t.Fatalf("expected output is empty: %s", expected)
	}

	if !cmp.Equal(expected, actual) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			name := name
			tn := strings.Replace(name, dir+"/", "", -1)
			t.Run(tn, func(t *testing.T) {
				if actual[name] == "" {
					t.Fatalf("%s is empty", name)
				}
				if diff := cmp.Diff(contents, actual[name]); diff != "" {
					t.Errorf("%s differed (-want +got):\n%s", name, diff)
				}
			})
		}
	}
}
