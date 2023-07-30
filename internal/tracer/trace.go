package tracer

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"

	"github.com/sqlc-dev/sqlc/internal/debug"
)

// Start starts Go's runtime tracing facility.
// Traces will be written to the file named by [debug.Debug.Trace].
// It also starts a new [*trace.Task] that will be stopped when the cleanup is called.
func Start(base context.Context) (_ context.Context, cleanup func(), _ error) {
	f, err := os.Create(debug.Debug.Trace)
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
