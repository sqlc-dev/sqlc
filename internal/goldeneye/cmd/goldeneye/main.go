// Command goldeneye generates the dialect seeds under
// internal/engine/<engine>/dialect from a live database, checks the
// committed ones against it, and checks the analyze cases under
// internal/endtoend/testdata against what the database itself reports.
//
// Usage, from internal/goldeneye:
//
//	go run ./cmd/goldeneye install clickhouse   # download the pinned clickhouse binary
//	go run ./cmd/goldeneye install sqlite       # build the pinned sqlite3 shells from source
//	go run ./cmd/goldeneye generate [engine]    # rewrite the generated files from the database
//	go run ./cmd/goldeneye check [engine]       # compare the committed files and analyze cases with the database
//
// Without an engine, generate and check cover every engine whose database
// is available and say which ones they skipped. `go test ./...` runs the
// same checks as tests.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/duckdb"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/mysql"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/postgresql"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/sqlite"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "goldeneye:", err)
		os.Exit(1)
	}
}

const usage = `usage:
  goldeneye install clickhouse|sqlite [-version V]
      put the pinned release of an engine into the user cache directory: clickhouse is
      downloaded, sqlite is built from the downloaded amalgamation with cc or $CC
  goldeneye generate [engine]
      rewrite the generated dialect files from the database, for every available engine or one
  goldeneye check [engine]
      compare the committed dialect files and analyze cases with the database, for every available engine or one
      postgresql and mysql read the server POSTGRESQL_SERVER_URI and MYSQL_SERVER_URI name

engines: clickhouse, duckdb, mysql, postgresql, sqlite`

// engine is one database goldeneye knows how to read a dialect from.
type engine struct {
	name string
	// dir is the engine directory the dialect lives under, when it is not
	// named after the engine: MySQL's is dolphin, after its parser.
	dir string
	// locate finds the database — a binary or a connection URL — or says
	// why it is not available.
	locate func() (string, error)
	// version describes the release the database is, for the log.
	version func(context.Context, string) (string, error)
	// generate reads the dialect from the database.
	generate func(context.Context, string) (dialect.Files, error)
	// analyze asks the database what it reports for an analyze case, in
	// the shape sqlc analyze prints. Nil for an engine without one yet.
	analyze func(context.Context, string, endtoend.Case) ([]byte, error)
}

var engines = []engine{
	{clickhouse.Engine, "", clickhouse.Locate, clickhouse.Version, clickhouse.Generate, clickhouse.Analyze},
	{duckdb.Engine, "", duckdb.Locate, duckdb.Version, duckdb.Generate, nil},
	{mysql.Engine, mysql.Dir, mysql.Locate, mysql.Version, mysql.Generate, mysql.Analyze},
	{postgresql.Engine, "", postgresql.Locate, postgresql.Version, postgresql.Generate, nil},
	{sqlite.Engine, "", sqlite.Locate, sqlite.Version, sqlite.Generate, sqlite.Analyze},
}

// dialectDir returns the engine's dialect directory.
func (e engine) dialectDir() (string, error) {
	if e.dir != "" {
		return dialect.Dir(e.dir)
	}
	return dialect.Dir(e.name)
}

// installer puts the binary an engine is read through in place, for the
// engines that need no server and whose release is pinned in their package.
type installer struct {
	defaultVersion string
	install        func(ctx context.Context, version, goos, goarch string, progress io.Writer) (string, error)
}

var installers = map[string]installer{
	clickhouse.Engine: {clickhouse.DefaultVersion, clickhouse.Install},
	sqlite.Engine:     {sqlite.DefaultVersion, sqlite.Install},
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return errors.New("a command is required")
	}
	switch args[0] {
	case "install":
		return install(ctx, args[1:], stdout, stderr)
	case "generate":
		return forEach(ctx, args[1:], stderr, generate)
	case "check":
		return forEach(ctx, args[1:], stderr, check)
	case "help", "-h", "--help":
		fmt.Fprintln(stdout, usage)
		return nil
	}
	fmt.Fprintln(stderr, usage)
	return fmt.Errorf("unknown command %q", args[0])
}

func install(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("install takes the engine to install: clickhouse or sqlite")
	}
	inst, ok := installers[args[0]]
	if !ok {
		return fmt.Errorf("install takes the engine to install, clickhouse or sqlite, not %q", args[0])
	}
	fs := flag.NewFlagSet("install "+args[0], flag.ContinueOnError)
	fs.SetOutput(stderr)
	version := fs.String("version", inst.defaultVersion, args[0]+" release to install")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	path, err := inst.install(ctx, *version, runtime.GOOS, runtime.GOARCH, stderr)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, path)
	return nil
}

// forEach runs a command over the named engine, or over every engine that
// is available when none is named. An engine that is named has to be
// available; one that is not named is skipped with a note when it is not.
func forEach(ctx context.Context, args []string, stderr io.Writer, cmd func(context.Context, engine, string, io.Writer) error) error {
	if len(args) > 1 {
		return errors.New("at most one engine may be named")
	}
	name := ""
	if len(args) == 1 {
		name = args[0]
	}
	found := false
	var failed []string
	for _, e := range engines {
		if name != "" && e.name != name {
			continue
		}
		found = true
		handle, err := e.locate()
		if err != nil {
			if name != "" {
				return err
			}
			fmt.Fprintf(stderr, "skipping %s: %v\n", e.name, err)
			continue
		}
		if err := cmd(ctx, e, handle, stderr); err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", e.name, err)
			failed = append(failed, e.name)
		}
	}
	if !found {
		return fmt.Errorf("unknown engine %q", name)
	}
	if len(failed) > 0 {
		return fmt.Errorf("%d engine(s) failed: %v", len(failed), failed)
	}
	return nil
}

func generate(ctx context.Context, e engine, handle string, stderr io.Writer) error {
	version, err := e.version(ctx, handle)
	if err != nil {
		return err
	}
	fmt.Fprintf(stderr, "%s: generating from %s\n", e.name, version)
	files, err := e.generate(ctx, handle)
	if err != nil {
		return err
	}
	dir, err := e.dialectDir()
	if err != nil {
		return err
	}
	if err := dialect.Write(dir, files); err != nil {
		return err
	}
	fmt.Fprintf(stderr, "%s: wrote %d file(s) to %s\n", e.name, len(files), dir)
	return nil
}

func check(ctx context.Context, e engine, handle string, stderr io.Writer) error {
	version, err := e.version(ctx, handle)
	if err != nil {
		return err
	}
	fmt.Fprintf(stderr, "%s: checking against %s\n", e.name, version)
	files, err := e.generate(ctx, handle)
	if err != nil {
		return err
	}
	dir, err := e.dialectDir()
	if err != nil {
		return err
	}
	report, err := dialect.Check(dir, files)
	if err != nil {
		return err
	}
	if report != "" {
		return fmt.Errorf("%s does not match the database\n%s", dir, report)
	}
	fmt.Fprintf(stderr, "%s: ok, %d file(s) match\n", e.name, len(files))
	return checkAnalyzeCases(ctx, e, handle, stderr)
}

// checkAnalyzeCases compares what the database reports for each of the
// engine's analyze cases with the case's committed output.
func checkAnalyzeCases(ctx context.Context, e engine, handle string, stderr io.Writer) error {
	if e.analyze == nil {
		return nil
	}
	cases, err := endtoend.Cases(e.name)
	if err != nil {
		return err
	}
	var report strings.Builder
	for _, c := range cases {
		got, err := e.analyze(ctx, handle, c)
		if err != nil {
			fmt.Fprintf(&report, "%s: %v\n", c.Name, err)
			continue
		}
		diff, err := c.Compare(got)
		if err != nil {
			return err
		}
		if diff != "" {
			fmt.Fprintf(&report, "%s (-committed +database)\n%s", c.Name, diff)
		}
	}
	if report.Len() > 0 {
		return fmt.Errorf("analyze cases do not match the database\n%s", report.String())
	}
	fmt.Fprintf(stderr, "%s: ok, %d analyze case(s) match\n", e.name, len(cases))
	return nil
}
