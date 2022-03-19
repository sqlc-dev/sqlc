package cmd

import (
	"context"
	"io"
	"os"

	"github.com/kyleconroy/sqlc/internal/bundler"
)

func createPkg(ctx context.Context, e Env, dir, filename string, stderr io.Writer) error {
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}
	up := bundler.NewUploader(configPath, dir, conf)
	if err := up.Validate(); err != nil {
		return err
	}
	output, err := Generate(ctx, e, dir, filename, stderr)
	if err != nil {
		os.Exit(1)
	}
	if e.DryRun {
		return up.DumpRequestOut(ctx, output)
	} else {
		return up.Upload(ctx, output)
	}
}
