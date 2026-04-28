package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"

	"github.com/sqlc-dev/sqlc/internal/api"
	"github.com/sqlc-dev/sqlc/internal/cmd"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/sqltest/docker"
	"github.com/sqlc-dev/sqlc/internal/sqltest/native"
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
			res := api.Generate(ctx, api.GenerateOptions{
				Dir:    path,
				Stderr: &stderr,
			})
			if len(res.Errors) > 0 {
				t.Fatalf("sqlc generate failed: %s", stderr.String())
			}
			cmpDirectory(t, path, res.Files)
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
				api.Generate(ctx, api.GenerateOptions{
					Dir:    path,
					Stderr: &stderr,
				})
			}
		})
	}
}

// textContext describes a TestReplay scenario. Mutate returns the config
// filename (relative to the test directory) that should be passed to the
// command under test. The "base" context returns "" to use the project's
// existing sqlc config; the "managed-db" context writes a mutated copy of the
// config to a temporary file inside the test directory and returns its name.
type textContext struct {
	Mutate  func(*testing.T, string) string
	Enabled func() bool
}

// writeMutatedConfig parses the sqlc config in dir, applies mutate to the
// in-memory Config (which is always v2-shaped, even when the file on disk is
// v1), forces version "2", and writes the result to a temp file alongside the
// original. The temp file is removed when the test ends.
func writeMutatedConfig(t *testing.T, dir string, mutate func(*config.Config)) string {
	t.Helper()
	original, conf, err := readSqlcConfig(dir)
	if err != nil {
		t.Fatalf("read sqlc config from %s: %s", dir, err)
	}

	// Parsing v1 configs converts them to a v2-shaped Config. Force version "2"
	// so the mutated config can be re-parsed as v2 from disk.
	conf.Version = "2"
	mutate(conf)

	f, err := os.CreateTemp(dir, "sqlc.test-*"+filepath.Ext(original))
	if err != nil {
		t.Fatalf("create temp config in %s: %s", dir, err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })

	enc := yaml.NewEncoder(f)
	if err := enc.Encode(conf); err != nil {
		f.Close()
		t.Fatalf("write temp config %s: %s", f.Name(), err)
	}
	if err := enc.Close(); err != nil {
		f.Close()
		t.Fatalf("close yaml encoder for %s: %s", f.Name(), err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close temp config %s: %s", f.Name(), err)
	}
	return filepath.Base(f.Name())
}

func readSqlcConfig(dir string) (string, *config.Config, error) {
	for _, name := range []string{"sqlc.yaml", "sqlc.yml", "sqlc.json"} {
		path := filepath.Join(dir, name)
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return path, nil, err
		}
		defer f.Close()
		conf, err := config.ParseConfig(f)
		if err != nil {
			return path, nil, fmt.Errorf("parse %s: %w", path, err)
		}
		return path, &conf, nil
	}
	return "", nil, fmt.Errorf("no sqlc config found in %s", dir)
}

func TestReplay(t *testing.T) {
	// Ensure that this environment variable is always set to true when running
	// end-to-end tests
	os.Setenv("SQLC_DUMMY_VALUE", "true")

	// t.Parallel()
	ctx := context.Background()

	var mysqlURI, postgresURI string

	// First, check environment variables
	if uri := os.Getenv("POSTGRESQL_SERVER_URI"); uri != "" {
		postgresURI = uri
	}
	if uri := os.Getenv("MYSQL_SERVER_URI"); uri != "" {
		mysqlURI = uri
	}

	// Try Docker for any missing databases
	if postgresURI == "" || mysqlURI == "" {
		if err := docker.Installed(); err == nil {
			if postgresURI == "" {
				host, err := docker.StartPostgreSQLServer(ctx)
				if err != nil {
					t.Logf("docker postgresql startup failed: %s", err)
				} else {
					postgresURI = host
				}
			}
			if mysqlURI == "" {
				host, err := docker.StartMySQLServer(ctx)
				if err != nil {
					t.Logf("docker mysql startup failed: %s", err)
				} else {
					mysqlURI = host
				}
			}
		}
	}

	// Try native installation for any missing databases (Linux only)
	if postgresURI == "" || mysqlURI == "" {
		if err := native.Supported(); err == nil {
			if postgresURI == "" {
				host, err := native.StartPostgreSQLServer(ctx)
				if err != nil {
					t.Logf("native postgresql startup failed: %s", err)
				} else {
					postgresURI = host
				}
			}
			if mysqlURI == "" {
				host, err := native.StartMySQLServer(ctx)
				if err != nil {
					t.Logf("native mysql startup failed: %s", err)
				} else {
					mysqlURI = host
				}
			}
		}
	}

	// Log which databases are available
	t.Logf("PostgreSQL available: %v (URI: %s)", postgresURI != "", postgresURI)
	t.Logf("MySQL available: %v (URI: %s)", mysqlURI != "", mysqlURI)

	contexts := map[string]textContext{
		"base": {
			Mutate:  func(t *testing.T, path string) string { return "" },
			Enabled: func() bool { return true },
		},
		"managed-db": {
			Mutate: func(t *testing.T, path string) string {
				return writeMutatedConfig(t, path, func(c *config.Config) {
					// Add all servers - tests will fail if database isn't available
					c.Servers = []config.Server{
						{Name: "postgres", Engine: config.EnginePostgreSQL, URI: postgresURI},
						{Name: "mysql", Engine: config.EngineMySQL, URI: mysqlURI},
					}
					for i := range c.SQL {
						switch c.SQL[i].Engine {
						case config.EnginePostgreSQL, config.EngineMySQL, config.EngineSQLite:
							c.SQL[i].Database = &config.Database{Managed: true}
						}
					}
				})
			},
			Enabled: func() bool {
				// Enabled if at least one database URI is available
				return postgresURI != "" || mysqlURI != ""
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

				configFile := testctx.Mutate(t, path)
				cmdOpts := cmd.Options{
					Env: cmd.Env{
						Debug:      opts.DebugFromString(args.Env["SQLCDEBUG"]),
						Experiment: opts.ExperimentFromString(args.Env["SQLCEXPERIMENT"]),
					},
					Stderr: &stderr,
				}

				switch args.Command {
				case "diff":
					res := api.Generate(ctx, api.GenerateOptions{
						Dir:    path,
						File:   configFile,
						Stderr: &stderr,
						Diff:   true,
					})
					if len(res.Errors) > 0 {
						err = res.Errors[0]
					}
				case "generate":
					res := api.Generate(ctx, api.GenerateOptions{
						Dir:    path,
						File:   configFile,
						Stderr: &stderr,
					})
					output = res.Files
					if len(res.Errors) > 0 {
						err = res.Errors[0]
					}
					if err == nil {
						cmpDirectory(t, path, output)
					}
				case "vet":
					err = cmd.Vet(ctx, path, configFile, &cmdOpts)
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
		// Mutated configs written by writeMutatedConfig.
		if strings.HasPrefix(filepath.Base(path), "sqlc.test-") {
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
				api.Generate(ctx, api.GenerateOptions{
					Dir:    path,
					Stderr: &stderr,
				})
			}
		})
	}
}
