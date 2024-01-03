package main

import (
	"bytes"
	"context"
	"os"
	osexec "os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/sqlc-dev/sqlc/internal/cmd"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func lineEndings() cmp.Option {
	return cmp.Transformer("LineEndings", func(in string) string {
		// Replace Windows new lines with Unix newlines
		return strings.Replace(in, "\r\n", "\n", -1)
	})
}

func stderrTransformer() cmp.Option {
	return cmp.Transformer("Stderr", func(in string) string {
		s := strings.Replace(in, "\r", "", -1)
		return strings.Replace(s, "\\", "/", -1)
	})
}

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
			opts := &cmd.Options{
				Env:    cmd.Env{},
				Stderr: &stderr,
			}
			output, err := cmd.Generate(ctx, path, "", opts)
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
				opts := &cmd.Options{
					Env:    cmd.Env{},
					Stderr: &stderr,
				}
				cmd.Generate(ctx, path, "", opts)
			}
		})
	}
}

type textContext struct {
	Mutate  func(*testing.T, string) func(*config.Config)
	Enabled func() bool
}

func TestReplay(t *testing.T) {
	// Ensure that this environment variable is always set to true when running
	// end-to-end tests
	os.Setenv("SQLC_DUMMY_VALUE", "true")

	// t.Parallel()
	ctx := context.Background()

	contexts := map[string]textContext{
		"base": {
			Mutate:  func(t *testing.T, path string) func(*config.Config) { return func(c *config.Config) {} },
			Enabled: func() bool { return true },
		},
		"managed-db": {
			Mutate: func(t *testing.T, path string) func(*config.Config) {
				return func(c *config.Config) {
					c.Cloud.Project = "01HAQMMECEYQYKFJN8MP16QC41" // TODO: Read from environment
					for i := range c.SQL {
						files := []string{}
						for _, s := range c.SQL[i].Schema {
							files = append(files, filepath.Join(path, s))
						}
						switch c.SQL[i].Engine {
						case config.EnginePostgreSQL:
							uri := local.PostgreSQL(t, files)
							c.SQL[i].Database = &config.Database{
								URI: uri,
							}
						// case config.EngineMySQL:
						// 	uri := local.MySQL(t, files)
						// 	c.SQL[i].Database = &config.Database{
						// 		URI: uri,
						// 	}
						default:
							c.SQL[i].Database = &config.Database{
								Managed: true,
							}
						}
					}
				}
			},
			Enabled: func() bool {
				// Return false if no auth token exists
				if len(os.Getenv("SQLC_AUTH_TOKEN")) == 0 {
					return false
				}
				if len(os.Getenv("POSTGRESQL_SERVER_URI")) == 0 {
					return false
				}
				// if len(os.Getenv("MYSQL_SERVER_URI")) == 0 {
				// 	return false
				// }
				return true
			},
		},
	}

	for name, testctx := range contexts {
		name := name
		testctx := testctx

		if !testctx.Enabled() {
			continue
		}

		for _, replay := range FindTests(t, "testdata", name) {
			tc := replay
			t.Run(filepath.Join(name, tc.Name), func(t *testing.T) {
				t.Parallel()

				var stderr bytes.Buffer
				var output map[string]string
				var err error

				path, _ := filepath.Abs(tc.Path)
				args := tc.Exec
				if args == nil {
					args = &Exec{Command: "generate"}
				}
				expected := string(tc.Stderr)

				if args.Process != "" {
					_, err := osexec.LookPath(args.Process)
					if err != nil {
						t.Skipf("executable not found: %s %s", args.Process, err)
					}
				}

				if len(args.Contexts) > 0 {
					if !slices.Contains(args.Contexts, name) {
						t.Skipf("unsupported context: %s", name)
					}
				}

				if len(args.OS) > 0 {
					if !slices.Contains(args.OS, runtime.GOOS) {
						t.Skipf("unsupported os: %s", runtime.GOOS)
					}
				}

				opts := cmd.Options{
					Env: cmd.Env{
						Debug:    opts.DebugFromString(args.Env["SQLCDEBUG"]),
						NoRemote: true,
					},
					Stderr:       &stderr,
					MutateConfig: testctx.Mutate(t, path),
				}

				switch args.Command {
				case "diff":
					err = cmd.Diff(ctx, path, "", &opts)
				case "generate":
					output, err = cmd.Generate(ctx, path, "", &opts)
					if err == nil {
						cmpDirectory(t, path, output)
					}
				case "vet":
					err = cmd.Vet(ctx, path, "", &opts)
				default:
					t.Fatalf("unknown command")
				}

				if len(expected) == 0 && err != nil {
					t.Fatalf("sqlc %s failed: %s", args.Command, stderr.String())
				}

				diff := cmp.Diff(
					strings.TrimSpace(expected),
					strings.TrimSpace(stderr.String()),
					stderrTransformer(),
				)
				if diff != "" {
					t.Fatalf("stderr differed (-want +got):\n%s", diff)
				}
			})
		}
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
		if filepath.Base(path) == "exec.json" {
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

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		lineEndings(),
	}

	if !cmp.Equal(expected, actual, opts...) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			name := name
			if actual[name] == "" {
				t.Errorf("%s is empty", name)
				return
			}
			if diff := cmp.Diff(contents, actual[name], opts...); diff != "" {
				t.Errorf("%s differed (-want +got):\n%s", name, diff)
			}
		}
	}
}

func BenchmarkReplay(b *testing.B) {
	ctx := context.Background()
	var dirs []string
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.json" || info.Name() == "sqlc.yaml" || info.Name() == "sqlc.yml" {
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
				opts := &cmd.Options{
					Env:    cmd.Env{},
					Stderr: &stderr,
				}
				cmd.Generate(ctx, path, "", opts)
			}
		})
	}
}
