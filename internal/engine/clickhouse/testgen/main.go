// Command testgen records what ClickHouse itself says about a set of sqlc
// queries, so the answer can be committed as a golden file and compared with
// what sqlc's own analysis produces.
//
// It loads a schema and a fixture into an ephemeral `clickhouse local`
// process, runs each query found in a sqlc query file against that data, and
// prints the result column types, nullability and source tables along with
// the parameters each query binds, in the same JSON shape `sqlc analyze`
// prints.
//
// The clickhouse binary is downloaded once per pinned version with
// `testgen install`, or supplied through the CLICKHOUSE environment variable.
//
// Usage:
//
//	go run . install
//	go run . analyze --schema schema.sql --fixture fixture.sql query.sql
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "testgen:", err)
		os.Exit(1)
	}
}

const usage = `usage:
  testgen install [-version V]
      download the pinned clickhouse binary into the user cache directory
  testgen analyze [-clickhouse PATH] --schema FILE [--fixture FILE] QUERY_FILE
      analyze every query in QUERY_FILE and print the result as JSON

The binary is looked up in the CLICKHOUSE environment variable first, then in
the cache directory populated by "testgen install".`

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return errors.New("a command is required")
	}
	switch args[0] {
	case "install":
		return runInstall(ctx, args[1:], stdout, stderr)
	case "analyze":
		return runAnalyze(ctx, args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprintln(stdout, usage)
		return nil
	default:
		fmt.Fprintln(stderr, usage)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runInstall(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(stderr)
	version := fs.String("version", DefaultVersion, "ClickHouse release to install")
	if err := fs.Parse(args); err != nil {
		return err
	}
	path, err := Install(ctx, *version, runtime.GOOS, runtime.GOARCH, stderr)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, path)
	return nil
}

func runAnalyze(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(stderr)
	binary := fs.String("clickhouse", "", "path to the clickhouse binary (defaults to $CLICKHOUSE, then the cache)")
	schemaPath := fs.String("schema", "", "path to the schema file")
	fixturePath := fs.String("fixture", "", "path to the fixture file loaded after the schema")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, usage)
		return errors.New("analyze takes exactly one query file")
	}
	if *schemaPath == "" {
		return errors.New("--schema is required")
	}
	if *binary == "" {
		path, err := Locate()
		if err != nil {
			return err
		}
		*binary = path
	}

	schema, err := os.ReadFile(*schemaPath)
	if err != nil {
		return err
	}
	var fixture []byte
	if *fixturePath != "" {
		fixture, err = os.ReadFile(*fixturePath)
		if err != nil {
			return err
		}
	}
	querySrc, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	queries, err := parseQueries(string(querySrc))
	if err != nil {
		return fmt.Errorf("%s: %w", fs.Arg(0), err)
	}

	out, err := analyze(ctx, local{binary: *binary}, string(schema), string(fixture), queries)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
