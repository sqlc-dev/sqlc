package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	ctx := context.Background()

	examples, err := filepath.Abs(filepath.Join("..", "..", "examples"))
	if err != nil {
		t.Fatal(err)
	}

	files, err := os.ReadDir(examples)
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
			output, err := cmd.Generate(ctx, cmd.Env{ExperimentalFeatures: true}, path, "", &stderr)
			if err != nil {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, output)
		})
	}
}

func BenchmarkExamples(b *testing.B) {
	ctx := context.Background()
	examples, err := filepath.Abs(filepath.Join("..", "..", "examples"))
	if err != nil {
		b.Fatal(err)
	}
	files, err := os.ReadDir(examples)
	if err != nil {
		b.Fatal(err)
	}
	for _, replay := range files {
		if !replay.IsDir() {
			continue
		}
		tc := replay.Name()
		b.Run(tc, func(b *testing.B) {
			path := filepath.Join(examples, tc)
			for i := 0; i < b.N; i++ {
				var stderr bytes.Buffer
				cmd.Generate(ctx, cmd.Env{ExperimentalFeatures: true}, path, "", &stderr)
			}
		})
	}
}

func TestReplay(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	var dirs []string
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.json" || info.Name() == "sqlc.yaml" {
			dirs = append(dirs, filepath.Dir(path))
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, replay := range dirs {
		tc := replay
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			var output map[string]string
			var err error

			path, _ := filepath.Abs(tc)
			args := parseExec(t, path)
			expected := expectedStderr(t, path)

			switch args.Command {
			case "diff":
				err = cmd.Diff(ctx, cmd.Env{ExperimentalFeatures: true}, path, "", &stderr)
			case "generate":
				output, err = cmd.Generate(ctx, cmd.Env{ExperimentalFeatures: true}, path, "", &stderr)
				if err == nil {
					cmpDirectory(t, path, output)
				}
			default:
				t.Fatalf("unknown command")
			}

			if len(expected) == 0 && err != nil {
				t.Fatalf("sqlc %s failed: %s", args.Command, stderr.String())
			}

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
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".kt") && !strings.HasSuffix(path, ".py") && !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".txt") {
			return nil
		}
		// TODO: Figure out a better way to ignore certain files
		if strings.HasSuffix(path, ".txt") && filepath.Base(path) != "hello.txt" {
			return nil
		}
		if filepath.Base(path) == "sqlc.json" {
			return nil
		}
		if strings.Contains(path, "/kotlin/build") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") || strings.Contains(path, "src/test/") {
			return nil
		}
		if strings.Contains(path, "/python/.venv") || strings.Contains(path, "/python/src/tests/") ||
			strings.HasSuffix(path, "__init__.py") || strings.Contains(path, "/python/src/dbtest/") ||
			strings.Contains(path, "/python/.mypy_cache") {
			return nil
		}
		blob, err := os.ReadFile(path)
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
		blob, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return string(blob)
	}
	return ""
}

type exec struct {
	Command string `json:"command"`
}

func parseExec(t *testing.T, dir string) exec {
	t.Helper()
	var e exec
	path := filepath.Join(dir, "exec.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		blob, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(blob, &e); err != nil {
			t.Fatal(err)
		}
	}
	if e.Command == "" {
		e.Command = "generate"
	}
	return e
}

func BenchmarkReplay(b *testing.B) {
	ctx := context.Background()
	var dirs []string
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.json" || info.Name() == "sqlc.yaml" {
			dirs = append(dirs, filepath.Dir(path))
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	for _, replay := range dirs {
		tc := replay
		b.Run(tc, func(b *testing.B) {
			path, _ := filepath.Abs(tc)
			for i := 0; i < b.N; i++ {
				var stderr bytes.Buffer
				cmd.Generate(ctx, cmd.Env{ExperimentalFeatures: true}, path, "", &stderr)
			}
		})
	}
}
