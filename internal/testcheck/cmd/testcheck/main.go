// Command testcheck verifies the analyze cases under internal/endtoend/testdata
// against a real database. It generates nothing: each engine package reads a
// case's schema, fixture and queries, asks the database what it makes of
// them, and compares the answer with the output.json the case committed.
//
// Usage, from this directory:
//
//	go run ./cmd/testcheck install clickhouse   # download the pinned clickhouse binary
//	go run ./cmd/testcheck check [engine]       # check every case, or one engine's
//
// `go test ./...` runs the same checks as tests, skipping engines whose
// database is not available.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/sqlc-dev/sqlc/internal/testcheck/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/testcheck/endtoend"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "testcheck:", err)
		os.Exit(1)
	}
}

const usage = `usage:
  testcheck install clickhouse [-version V]
      download the pinned clickhouse binary into the user cache directory
  testcheck check [engine]
      verify the analyze cases against the database, for every engine or one`

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return errors.New("a command is required")
	}
	switch args[0] {
	case "install":
		return install(ctx, args[1:], stdout, stderr)
	case "check":
		return check(ctx, args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprintln(stdout, usage)
		return nil
	}
	fmt.Fprintln(stderr, usage)
	return fmt.Errorf("unknown command %q", args[0])
}

func install(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != clickhouse.Engine {
		return errors.New("install takes the engine to install: clickhouse")
	}
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(stderr)
	version := fs.String("version", clickhouse.DefaultVersion, "ClickHouse release to install")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	path, err := clickhouse.Install(ctx, *version, runtime.GOOS, runtime.GOARCH, stderr)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, path)
	return nil
}

func check(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	engine := ""
	if len(args) > 0 {
		engine = args[0]
	}
	failed := 0
	if engine == "" || engine == clickhouse.Engine {
		binary, err := clickhouse.Locate()
		if err != nil {
			if engine != "" {
				return err
			}
			fmt.Fprintf(stderr, "skipping clickhouse: %v\n", err)
		} else {
			cases, err := endtoend.Cases(clickhouse.Engine)
			if err != nil {
				return err
			}
			for _, c := range cases {
				diff, err := clickhouse.Check(ctx, binary, c)
				switch {
				case err != nil:
					failed++
					fmt.Fprintf(stdout, "ERROR %s: %v\n", c.Name, err)
				case diff != "":
					failed++
					fmt.Fprintf(stdout, "FAIL  %s (-committed +clickhouse)\n%s\n", c.Name, diff)
				default:
					fmt.Fprintf(stdout, "ok    %s\n", c.Name)
				}
			}
		}
	} else {
		return fmt.Errorf("unknown engine %q", engine)
	}
	if failed > 0 {
		return fmt.Errorf("%d case(s) did not match", failed)
	}
	return nil
}
