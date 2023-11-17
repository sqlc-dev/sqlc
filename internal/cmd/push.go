package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/bundler"
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
)

type pusher struct {
	m       sync.Mutex
	results []*bundler.QuerySetArchive
}

func (g *pusher) ProcessResult(ctx context.Context, combo config.CombinedSettings, sql outPair, result *compiler.Result) error {
	req := codeGenRequest(result, combo)
	g.m.Lock()
	g.results = append(g.results, &bundler.QuerySetArchive{
		Name:    sql.Name,
		Schema:  sql.Schema,
		Queries: sql.Queries,
		Request: req,
	})
	g.m.Unlock()
	return nil
}

func Push(ctx context.Context, dir, filename string, opts *Options) error {
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
	p := &pusher{}
	if err := Process(ctx, p, dir, filename, opts); err != nil {
		os.Exit(1)
	}
	if e.DryRun {
		return up.DumpRequestOut(ctx, p.results)
	} else {
		return up.Upload(ctx, p.results)
	}
}
