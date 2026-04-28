package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
				Config: openConfigReader(t, path),
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
			cfg := openConfigBytes(b, path)
			for i := 0; i < b.N; i++ {
				var stderr bytes.Buffer
				api.Generate(ctx, api.GenerateOptions{
					Config: bytes.NewReader(cfg),
					Stderr: &stderr,
				})
			}
		})
	}
}

// openConfigReader reads the sqlc config in dir, rewrites every relative
// schema/queries/output path to an absolute one (so api.Generate doesn't have
// to know the config's directory), and returns the result as an io.Reader.
func openConfigReader(t testing.TB, dir string) io.Reader {
	return bytes.NewReader(openConfigBytes(t, dir))
}

func openConfigBytes(t testing.TB, dir string) []byte {
	t.Helper()
	data, _ := mutatedConfigBytes(t, dir, nil)
	return data
}

// mutatedConfigBytes parses the sqlc config in dir, applies mutate (when
// non-nil) to the in-memory Config, makes every path absolute relative to dir,
// and re-encodes as YAML. Parsing v1 configs converts them to a v2-shaped
// Config; we force version "2" so the result can be parsed back by api.Generate.
//
// When mutate is non-nil, the encoded bytes are also written to a temp file
// alongside the original (cleaned up at test end) and the filename is returned
// so callers like cmd.Vet that still take a config-file path can use it.
func mutatedConfigBytes(t testing.TB, dir string, mutate func(*config.Config)) ([]byte, string) {
	t.Helper()
	original, conf, err := readSqlcConfig(dir)
	if err != nil {
		t.Fatalf("read sqlc config from %s: %s", dir, err)
	}

	conf.Version = "2"
	absolutizePaths(conf, dir)
	if mutate != nil {
		mutate(conf)
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(conf); err != nil {
		t.Fatalf("encode config: %s", err)
	}
	if err := enc.Close(); err != nil {
		t.Fatalf("close yaml encoder: %s", err)
	}
	data := buf.Bytes()

	if mutate == nil {
		return data, ""
	}

	f, err := os.CreateTemp(dir, "sqlc.test-*"+filepath.Ext(original))
	if err != nil {
		t.Fatalf("create temp config in %s: %s", dir, err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.Write(data); err != nil {
		f.Close()
		t.Fatalf("write temp config %s: %s", f.Name(), err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close temp config %s: %s", f.Name(), err)
	}
	return data, filepath.Base(f.Name())
}

func absolutizePaths(conf *config.Config, dir string) {
	abs := func(p string) string {
		if p == "" || filepath.IsAbs(p) {
			return p
		}
		return filepath.Join(dir, p)
	}
	for i := range conf.SQL {
		s := &conf.SQL[i]
		for j, p := range s.Schema {
			s.Schema[j] = abs(p)
		}
		for j, p := range s.Queries {
			s.Queries[j] = abs(p)
		}
		if s.Gen.Go != nil {
			s.Gen.Go.Out = abs(s.Gen.Go.Out)
		}
		if s.Gen.JSON != nil {
			s.Gen.JSON.Out = abs(s.Gen.JSON.Out)
		}
		for j := range s.Codegen {
			s.Codegen[j].Out = abs(s.Codegen[j].Out)
		}
	}
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

// configRef is the result of preparing a config for a single TestReplay case.
// reader is for api.Generate; file is for cmd.Vet which still takes a path.
type configRef struct {
	reader io.Reader
	file   string
}

type textContext struct {
	Config  func(*testing.T, string) configRef
	Enabled func() bool
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
			Config: func(t *testing.T, dir string) configRef {
				data, _ := mutatedConfigBytes(t, dir, nil)
				return configRef{reader: bytes.NewReader(data)}
			},
			Enabled: func() bool { return true },
		},
		"managed-db": {
			Config: func(t *testing.T, dir string) configRef {
				data, file := mutatedConfigBytes(t, dir, func(c *config.Config) {
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
				return configRef{reader: bytes.NewReader(data), file: file}
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

				cfg := testctx.Config(t, path)
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
						Config: cfg.reader,
						Stderr: &stderr,
						Diff:   true,
					})
					if len(res.Errors) > 0 {
						err = res.Errors[0]
					}
				case "generate":
					res := api.Generate(ctx, api.GenerateOptions{
						Config: cfg.reader,
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
					err = cmd.Vet(ctx, path, cfg.file, &cmdOpts)
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
		// Mutated configs written by mutatedConfigBytes.
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
			cfg := openConfigBytes(b, path)
			for i := 0; i < b.N; i++ {
				var stderr bytes.Buffer
				api.Generate(ctx, api.GenerateOptions{
					Config: bytes.NewReader(cfg),
					Stderr: &stderr,
				})
			}
		})
	}
}
