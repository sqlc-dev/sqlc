package tracer

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"

	"github.com/sqlc-dev/sqlc/internal/sqlcdebug"
)

var debugTrace = sqlcdebug.New("trace")

// Path returns the file to which Go's runtime tracer should write its
// output, derived from SQLCDEBUG=trace=...
func Path() string {
	v := debugTrace.Value()
	if v == "1" {
		return "trace.out"
	}
	return v
}

// Start starts Go's runtime tracing facility. Traces are written to the
// path returned by [Path]. It also starts a new [*trace.Task] that will
// be stopped when the cleanup is called.
func Start(base context.Context) (_ context.Context, cleanup func(), _ error) {
	f, err := os.Create(Path())
	if err != nil {
		return base, cleanup, fmt.Errorf("failed to create trace output file: %v", err)
	}

	if err := trace.Start(f); err != nil {
		return base, cleanup, fmt.Errorf("failed to start trace: %v", err)
	}

	ctx, task := trace.NewTask(base, "sqlc")

	return ctx, func() {
		defer f.Close()
		defer trace.Stop()
		defer task.End()
	}, nil
}
