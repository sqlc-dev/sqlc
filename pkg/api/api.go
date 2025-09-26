package api

import (
	"context"

	"github.com/sqlc-dev/sqlc/internal/cmd"
)

type Options struct {
	Dir      string // working directory for relative paths
	Filename string
	Options  *cmd.Options
}

func Generate(ctx context.Context, opt Options) (map[string]string, error) {
	return cmd.Generate(ctx, opt.Dir, opt.Filename, opt.Options)
}

func Verify(ctx context.Context, opt Options) error {
	return cmd.Verify(ctx, opt.Dir, opt.Filename, opt.Options)
}

func Vet(ctx context.Context, opt Options) error {
	return cmd.Vet(ctx, opt.Dir, opt.Filename, opt.Options)
}
