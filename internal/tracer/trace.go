package tracer

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"

	"github.com/kyleconroy/sqlc/internal/debug"
)

var noop Ender

func init() {
	noop = NoopRegion{}
}

type NoopRegion struct {
}

func (n NoopRegion) End() {
}

type Ender interface {
	End()
}

func Start(base context.Context) (context.Context, func(), error) {
	if !debug.Traced {
		return base, func() {}, nil
	}

	f, err := os.Create(debug.Debug.Trace)
	if err != nil {
		return base, func() {}, fmt.Errorf("failed to create trace output file: %v", err)
	}

	if err := trace.Start(f); err != nil {
		return base, func() {}, fmt.Errorf("failed to start trace: %v", err)
	}

	ctx, task := trace.NewTask(base, "sqlc")

	return ctx, func() {
		defer f.Close()
		defer trace.Stop()
		defer task.End()
	}, nil
}
