// Package postgresql contains end-to-end tests for the PostgreSQL engine.
package postgresql

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/sqlc-dev/sqlc/internal/cmd"
	"github.com/sqlc-dev/sqlc/internal/enginetest/testcases"
)

func TestEndToEnd(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testdataDir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}

	// Walk through all test directories
	err = filepath.Walk(testdataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for sqlc config files
		if info.Name() != "sqlc.yaml" && info.Name() != "sqlc.json" {
			return nil
		}

		dir := filepath.Dir(path)
		testName := strings.TrimPrefix(dir, testdataDir+string(filepath.Separator))

		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			runTest(ctx, t, dir)
		})

		return filepath.SkipDir
	})

	if err != nil {
		t.Fatal(err)
	}
}

func runTest(ctx context.Context, t *testing.T, dir string) {
	t.Helper()

	var stderr bytes.Buffer
	opts := &cmd.Options{
		Env:    cmd.Env{},
		Stderr: &stderr,
	}

	// Check for expected stderr
	expectedStderr := readExpectedStderr(t, dir)

	output, err := cmd.Generate(ctx, dir, "", opts)

	// If we expect an error, check stderr matches
	if len(expectedStderr) > 0 {
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		diff := cmp.Diff(
			strings.TrimSpace(expectedStderr),
			strings.TrimSpace(stderr.String()),
			stderrTransformer(),
		)
		if diff != "" {
			t.Fatalf("stderr differed (-want +got):\n%s", diff)
		}
		return
	}

	if err != nil {
		t.Fatalf("sqlc generate failed: %s", stderr.String())
	}

	cmpDirectory(t, dir, output)
}

func readExpectedStderr(t *testing.T, dir string) string {
	t.Helper()

	paths := []string{
		filepath.Join(dir, "stderr.txt"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			blob, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			return string(blob)
		}
	}
	return ""
}

func stderrTransformer() cmp.Option {
	return cmp.Transformer("Stderr", func(in string) string {
		s := strings.Replace(in, "\r", "", -1)
		return strings.Replace(s, "\\", "/", -1)
	})
}

func lineEndings() cmp.Option {
	return cmp.Transformer("LineEndings", func(in string) string {
		return strings.Replace(in, "\r\n", "\n", -1)
	})
}

func cmpDirectory(t *testing.T, dir string, actual map[string]string) {
	t.Helper()

	expected := map[string]string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		blob, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		expected[path] = string(blob)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		lineEndings(),
	}

	if !cmp.Equal(expected, actual, opts...) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			if actual[name] == "" {
				t.Errorf("%s is empty", name)
				continue
			}
			if diff := cmp.Diff(contents, actual[name], opts...); diff != "" {
				t.Errorf("%s differed (-want +got):\n%s", name, diff)
			}
		}
	}
}

// Engine returns the engine type for this package
func Engine() testcases.Engine {
	return testcases.EnginePostgreSQL
}
