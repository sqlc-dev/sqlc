package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log/slog"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/analysis"
	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
)

type CachedAnalyzer struct {
	a           Analyzer
	config      config.Config
	configBytes []byte
	db          config.Database
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

func (c *CachedAnalyzer) analyze(ctx context.Context, n ast.Node, q string, schema []string, np *named.ParamSet) (*analysis.Analysis, bool, error) {
	// Only cache queries for managed databases. We can't be certain the
	// database is in an unchanged state otherwise
	if !c.db.Managed {
		return nil, true, nil
	}

	dir, err := cache.AnalysisDir()
	if err != nil {
		return nil, true, err
	}

	if c.configBytes == nil {
		c.configBytes, err = json.Marshal(c.config)
		if err != nil {
			return nil, true, err
		}
	}

	// Calculate cache key
	h := fnv.New64()
	h.Write([]byte(info.Version))
	h.Write(c.configBytes)
	for _, m := range schema {
		h.Write([]byte(m))
	}
	h.Write([]byte(q))

	key := fmt.Sprintf("%x", h.Sum(nil))
	path := filepath.Join(dir, key)
	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, true, err
		}
		var a analysis.Analysis
		if err := proto.Unmarshal(contents, &a); err != nil {
			return nil, true, err
		}
		return &a, false, nil
	}

	result, err := c.a.Analyze(ctx, n, q, schema, np)

	if err == nil {
		contents, err := proto.Marshal(result)
		if err != nil {
			slog.Warn("unable to marshal analysis", "err", err)
			return result, false, nil
		}
		if err := os.WriteFile(path, contents, 0644); err != nil {
			slog.Warn("saving analysis to disk failed", "err", err)
			return result, false, nil
		}
	}

	return result, false, err
}

func (c *CachedAnalyzer) Close(ctx context.Context) error {
	return c.a.Close(ctx)
}

type Analyzer interface {
	Analyze(context.Context, ast.Node, string, []string, *named.ParamSet) (*analysis.Analysis, error)
	Close(context.Context) error
}
