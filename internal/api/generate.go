package api

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
)

// GenerateOptions controls a single Generate invocation.
type GenerateOptions struct {
	// Dir is the working directory used to resolve the config file and any
	// relative schema/query paths within it.
	Dir string

	// File is the configuration filename to use, relative to Dir. When empty,
	// Generate looks for sqlc.yaml, sqlc.yml, or sqlc.json in Dir.
	File string

	// Stderr receives diagnostic output. If nil, output is discarded.
	Stderr io.Writer

	// Write, when true, writes the generated files to disk after a successful
	// generate. Failures are reported via GenerateResult.Errors.
	Write bool

	// Diff, when true, compares each generated file against any existing file
	// on disk and writes a unified diff for differences to Stderr. If any
	// differences are found, an error is appended to GenerateResult.Errors.
	Diff bool
}

// GenerateResult is the outcome of a Generate call. Files maps absolute output
// paths to file contents; callers are responsible for writing them to disk if
// desired. Errors collects any errors encountered during code generation.
type GenerateResult struct {
	// Files maps absolute output paths to generated file contents.
	Files map[string]string

	// Errors collects any errors encountered. A non-empty Errors slice means
	// generation did not fully succeed.
	Errors []error
}

// Generate parses the sqlc configuration referenced by opts and runs every
// configured codegen target. The returned GenerateResult always has a non-nil
// Files map; the map is empty when generation fails before any files are
// produced.
func Generate(ctx context.Context, opts GenerateOptions) GenerateResult {
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	res := GenerateResult{Files: map[string]string{}}

	configPath, conf, err := readConfig(stderr, opts.Dir, opts.File)
	if err != nil {
		res.Errors = append(res.Errors, err)
		return res
	}

	base := filepath.Base(configPath)
	if err := config.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		res.Errors = append(res.Errors, err)
		return res
	}

	g := &generator{
		dir:    opts.Dir,
		output: map[string]string{},
	}

	if err := processQuerySets(ctx, g, conf, opts.Dir, stderr); err != nil {
		res.Errors = append(res.Errors, err)
		return res
	}

	res.Files = g.output

	if opts.Write {
		if err := writeFiles(ctx, res.Files, stderr); err != nil {
			res.Errors = append(res.Errors, err)
		}
	}

	if opts.Diff {
		if err := diffFiles(ctx, opts.Dir, res.Files, stderr); err != nil {
			res.Errors = append(res.Errors, err)
		}
	}

	return res
}

type generator struct {
	m      sync.Mutex
	dir    string
	output map[string]string
}

func (g *generator) Pairs(ctx context.Context, conf *config.Config) []outputPair {
	var pairs []outputPair
	for _, sql := range conf.SQL {
		if sql.Gen.Go != nil {
			pairs = append(pairs, outputPair{
				SQL: sql,
				Gen: config.SQLGen{Go: sql.Gen.Go},
			})
		}
		if sql.Gen.JSON != nil {
			pairs = append(pairs, outputPair{
				SQL: sql,
				Gen: config.SQLGen{JSON: sql.Gen.JSON},
			})
		}
		for i := range sql.Codegen {
			pairs = append(pairs, outputPair{
				SQL:    sql,
				Plugin: &sql.Codegen[i],
			})
		}
	}
	return pairs
}

func (g *generator) ProcessResult(ctx context.Context, combo config.CombinedSettings, sql outputPair, result *compiler.Result) error {
	out, resp, err := codegen(ctx, combo, sql, result)
	if err != nil {
		return err
	}
	files := map[string]string{}
	for _, file := range resp.Files {
		files[file.Name] = string(file.Contents)
	}
	g.m.Lock()

	// out is specified by the user, not a plugin
	absout := filepath.Join(g.dir, out)

	for n, source := range files {
		filename := filepath.Join(g.dir, out, n)
		// filepath.Join calls filepath.Clean which should remove all "..", but
		// double check to make sure
		if strings.Contains(filename, "..") {
			return fmt.Errorf("invalid file output path: %s", filename)
		}
		// The output file must be contained inside the output directory
		if !strings.HasPrefix(filename, absout) {
			return fmt.Errorf("invalid file output path: %s", filename)
		}
		g.output[filename] = source
	}
	g.m.Unlock()
	return nil
}
