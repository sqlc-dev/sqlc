package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"
	"runtime/trace"

	"golang.org/x/sync/errgroup"

	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/opts"
)

type outputPair struct {
	Gen    config.SQLGen
	Plugin *config.Codegen

	config.SQL
}

type resultProcessor interface {
	Pairs(context.Context, *config.Config) []outputPair
	ProcessResult(context.Context, config.CombinedSettings, outputPair, *compiler.Result) error
}

func processQuerySets(ctx context.Context, rp resultProcessor, conf *config.Config, baseDir string, stderr io.Writer) error {
	errored := false

	pairs := rp.Pairs(ctx, conf)
	grp, gctx := errgroup.WithContext(ctx)
	grp.SetLimit(runtime.GOMAXPROCS(0))

	stderrs := make([]bytes.Buffer, len(pairs))

	for i, pair := range pairs {
		sql := pair
		errout := &stderrs[i]

		grp.Go(func() error {
			combo := config.Combine(*conf, sql.SQL)
			if sql.Plugin != nil {
				combo.Codegen = *sql.Plugin
			}

			absSchema := make([]string, 0, len(sql.Schema))
			for _, s := range sql.Schema {
				absSchema = append(absSchema, resolvePath(baseDir, s))
			}
			sql.Schema = absSchema

			absQueries := make([]string, 0, len(sql.Queries))
			for _, q := range sql.Queries {
				absQueries = append(absQueries, resolvePath(baseDir, q))
			}
			sql.Queries = absQueries

			var name, lang string
			parseOpts := opts.Parser{
				Debug: debug.Debug,
			}

			switch {
			case sql.Gen.Go != nil:
				name = combo.Go.Package
				lang = "golang"

			case sql.Plugin != nil:
				lang = fmt.Sprintf("process:%s", sql.Plugin.Plugin)
				name = sql.Plugin.Plugin
			}

			packageRegion := trace.StartRegion(gctx, "package")
			trace.Logf(gctx, "", "name=%s plugin=%s", name, lang)

			result, failed := parse(gctx, name, baseDir, sql.SQL, combo, parseOpts, errout)
			if failed {
				packageRegion.End()
				errored = true
				return nil
			}
			if err := rp.ProcessResult(gctx, combo, sql, result); err != nil {
				fmt.Fprintf(errout, "# package %s\n", name)
				fmt.Fprintf(errout, "error generating code: %s\n", err)
				errored = true
			}
			packageRegion.End()
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	if errored {
		for i := range stderrs {
			if _, err := io.Copy(stderr, &stderrs[i]); err != nil {
				return err
			}
		}
		return fmt.Errorf("errored")
	}
	return nil
}
