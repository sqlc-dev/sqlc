package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/bundler"
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
)

func init() {
	pushCmd.Flags().StringSliceP("tag", "t", nil, "tag this push with a value")
}

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"upload"},
	Short:   "Push the schema, queries, and configuration for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		tags, err := cmd.Flags().GetStringSlice("tag")
		if err != nil {
			return err
		}
		opts := &Options{
			Env:    ParseEnv(cmd),
			Stderr: stderr,
			Tags:   tags,
		}
		if err := Push(cmd.Context(), dir, name, opts); err != nil {
			fmt.Fprintf(stderr, "error pushing: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
}

type pusher struct {
	m       sync.Mutex
	results []*bundler.QuerySetArchive
}

func (g *pusher) Pairs(ctx context.Context, conf *config.Config) []OutputPair {
	var pairs []OutputPair
	for _, sql := range conf.SQL {
		pairs = append(pairs, OutputPair{
			SQL: sql,
		})
	}
	return pairs
}

func (g *pusher) ProcessResult(ctx context.Context, combo config.CombinedSettings, sql OutputPair, result *compiler.Result) error {
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
		return err
	}
	if e.DryRun {
		return up.DumpRequestOut(ctx, p.results)
	} else {
		return up.Upload(ctx, p.results, opts.Tags)
	}
}
