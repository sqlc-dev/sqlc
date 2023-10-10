package cmd

import (
	"context"
	"os"

	"github.com/sqlc-dev/sqlc/internal/bundler"
)

func createPkg(ctx context.Context, dir, filename string, opts *Options) error {
	e := opts.Env
	stderr := opts.Stderr
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}
	up := bundler.NewUploader(configPath, dir, conf)
	if err := up.Validate(); err != nil {
		return err
	}
	output, err := Generate(ctx, dir, filename, opts)
	if err != nil {
		os.Exit(1)
	}
	if e.DryRun {
		return up.DumpRequestOut(ctx, output)
	} else {
		return up.Upload(ctx, output)
	}
}
