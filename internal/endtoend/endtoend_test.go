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

func TestCodeGeneration(t *testing.T) {
	// Change to the top-level directory of the project
	examples, _ := filepath.Abs(filepath.Join("..", "..", "examples"))
	var stderr bytes.Buffer

	output, err := cmd.Generate(examples, &stderr)
	if err != nil {
		t.Fatalf("%s", stderr.String())
	}

	cmpDirectory(t, examples, output)
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
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
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
			if actual[name] == "" {
				t.Errorf("%s is empty", name)
				continue
			}
			if diff := cmp.Diff(contents, actual[name]); diff != "" {
				t.Errorf("%s differed (-want +got):\n%s", name, diff)
			}
		}
		for name, _ := range actual {
			t.Log(name)
		}
	}
}
