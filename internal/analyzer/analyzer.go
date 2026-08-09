package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
)

type CachedAnalyzer struct {
	a           Analyzer
	config      config.Config
	configBytes []byte
	db          config.Database
	store       *cache.Cache
}

func Cached(a Analyzer, c config.Config, db config.Database) *CachedAnalyzer {
	return &CachedAnalyzer{
		a:      a,
		config: c,
		db:     db,
	}
}

// Create a new error here

func (c *CachedAnalyzer) Analyze(ctx context.Context, n ast.Node, q string, schema []string, np *named.ParamSet) (*analysis.Analysis, error) {
	result, rerun, err := c.analyze(ctx, n, q, schema, np)
	if rerun {
		if err != nil {
			slog.Warn("first analysis failed with error", "err", err)
		}
		return c.a.Analyze(ctx, n, q, schema, np)
	}
	return result, err
}

// The name of the sole output blob a QueryAnalysis action produces.
const analysisOutput = "analysis.pb"

func (c *CachedAnalyzer) analyze(ctx context.Context, n ast.Node, q string, schema []string, np *named.ParamSet) (*analysis.Analysis, bool, error) {
	// Only cache queries for managed databases. We can't be certain the
	// database is in an unchanged state otherwise
	if !c.db.Managed {
		return nil, true, nil
	}

	if c.store == nil {
		var err error
		c.store, err = cache.Open()
		if err != nil {
			return nil, true, err
		}
	}
	store := c.store

	if c.configBytes == nil {
		var err error
		c.configBytes, err = json.Marshal(c.config)
		if err != nil {
			return nil, true, err
		}
	}

	// Analyzing a query is an action whose inputs are the configuration, the
	// schema migrations, and the query itself. (The sqlc binary is an
	// implicit input of every action.)
	action := store.NewAction("QueryAnalysis").
		AddInput("config", c.configBytes)
	for _, m := range schema {
		action.AddInput("schema", []byte(m))
	}
	actionDigest := action.AddInput("query", []byte(q)).Digest()

	if cached, err := store.Actions.Get(actionDigest); err == nil {
		contents, err := store.CAS.Get(cached.Outputs[analysisOutput])
		if err == nil {
			var a analysis.Analysis
			if err := proto.Unmarshal(contents, &a); err == nil {
				return &a, false, nil
			}
		}
		if !errors.Is(err, cache.ErrNotFound) {
			slog.Warn("reading analysis from cache failed", "err", err)
		}
	}

	result, err := c.a.Analyze(ctx, n, q, schema, np)

	if err == nil {
		contents, err := proto.Marshal(result)
		if err != nil {
			slog.Warn("unable to marshal analysis", "err", err)
			return result, false, nil
		}
		outDigest, err := store.CAS.Put(contents)
		if err != nil {
			slog.Warn("saving analysis to cache failed", "err", err)
			return result, false, nil
		}
		err = store.Actions.Put(actionDigest, &cache.ActionResult{
			Outputs: map[string]cache.Digest{analysisOutput: outDigest},
		})
		if err != nil {
			slog.Warn("saving analysis action result failed", "err", err)
		}
	}

	return result, false, err
}

func (c *CachedAnalyzer) Close(ctx context.Context) error {
	if c.store != nil {
		c.store.Close()
	}
	return c.a.Close(ctx)
}

type Analyzer interface {
	Analyze(context.Context, ast.Node, string, []string, *named.ParamSet) (*analysis.Analysis, error)
	Close(context.Context) error
}
