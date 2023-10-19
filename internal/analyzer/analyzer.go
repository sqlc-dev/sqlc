package analyzer

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/analyzer/pb"
	"github.com/sqlc-dev/sqlc/internal/cache"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
)

type CachedAnalyzer struct {
	a          Analyzer
	configPath string
	config     []byte
	db         config.Database
}

func Cached(a Analyzer, cp string, db config.Database) *CachedAnalyzer {
	return &CachedAnalyzer{
		a:          a,
		configPath: cp,
		db:         db,
	}
}

func (c *CachedAnalyzer) Analyze(ctx context.Context, n ast.Node, q string, schema []string, np *named.ParamSet) (*pb.Analysis, error) {
	// Only cache queries for managed databases. We can't be certain the the
	// database is in an unchanged state otherwise
	if !c.db.Managed {
		return c.a.Analyze(ctx, n, q, schema, np)
	}

	dir, err := cache.AnalysisDir()
	if err != nil {
		return nil, err
	}

	if c.config == nil {
		c.config, err = os.ReadFile(c.configPath)
		if err != nil {
			return nil, err
		}
	}

	// Calculate cache key
	h := fnv.New64()
	h.Write([]byte(info.Version))
	h.Write(c.config)
	for _, m := range schema {
		h.Write([]byte(m))
	}
	h.Write([]byte(q))

	key := fmt.Sprintf("%x", h.Sum(nil))
	path := filepath.Join(dir, key)
	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var a pb.Analysis
		return &a, proto.Unmarshal(contents, &a)
	}

	result, err := c.a.Analyze(ctx, n, q, schema, np)

	if err == nil {
		contents, err := proto.Marshal(result)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, contents, 0644); err != nil {
			return nil, err
		}
	}

	return result, err
}

func (c *CachedAnalyzer) Close(ctx context.Context) error {
	return c.a.Close(ctx)
}

type Analyzer interface {
	Analyze(context.Context, ast.Node, string, []string, *named.ParamSet) (*pb.Analysis, error)
	Close(context.Context) error
}
