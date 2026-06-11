package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
)

// GenerateOptions controls a single Generate invocation. Paths declared in the
// configuration are resolved relative to the current working directory, so
// callers wanting a different base directory should either pass absolute
// paths in the config or os.Chdir before calling.
type GenerateOptions struct {
	// Config is the sqlc configuration as a YAML or JSON document. Required.
	Config io.Reader

	// Stderr receives diagnostic output. If nil, output is discarded.
	Stderr io.Writer

	// Write writes the generated files to disk after a successful generate.
	// Failures are reported via GenerateResult.Errors.
	Write bool

	// Diff compares each generated file against any existing file on disk and
	// writes a unified diff for differences to Stderr. If any differences are
	// found, an error is appended to GenerateResult.Errors.
	Diff bool

	// BaseDir is the directory relative paths in Config are resolved against,
	// and the prefix stripped from file paths shown in parse errors and diff
	// labels. When empty, BaseDir defaults to the current working directory.
	BaseDir string

	// EnableProcessPlugins controls whether the configuration may invoke
	// process-based plugins. When false (the zero value), Generate fails
	// before parsing or codegen runs if the configuration declares any
	// process plugin. Process plugins execute arbitrary local commands, so
	// callers must opt in explicitly. The sqlc CLI populates this from
	// SQLCDEBUG: it defaults to true, and SQLCDEBUG=processplugins=0 turns
	// it off.
	EnableProcessPlugins bool
}

// GenerateResult is the outcome of a Generate call.
type GenerateResult struct {
	// Files maps absolute output paths to generated file contents.
	Files map[string]string

	// Errors collects any errors encountered. A non-empty Errors slice means
	// generation did not fully succeed.
	Errors []error
}

// Generate parses the sqlc configuration referenced by opts and runs every
// configured codegen target.
func Generate(ctx context.Context, opts GenerateOptions) GenerateResult {
	stderr := opts.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	res := GenerateResult{Files: map[string]string{}}

	if opts.Config == nil {
		err := errors.New("GenerateOptions.Config is required")
		fmt.Fprintln(stderr, err)
		res.Errors = append(res.Errors, err)
		return res
	}

	conf, err := config.ParseConfig(opts.Config)
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprint(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprint(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprint(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing config: %s\n", err)
		res.Errors = append(res.Errors, err)
		return res
	}

	if err := config.Validate(&conf); err != nil {
		fmt.Fprintf(stderr, "error validating config: %s\n", err)
		res.Errors = append(res.Errors, err)
		return res
	}

	if !opts.EnableProcessPlugins {
		for _, plug := range conf.Plugins {
			if plug.Process != nil {
				err := fmt.Errorf("process plugin %q declared but EnableProcessPlugins is false", plug.Name)
				fmt.Fprintf(stderr, "error validating config: %s\n", err)
				res.Errors = append(res.Errors, err)
				return res
			}
		}
	}

	g := &generator{output: map[string]string{}, baseDir: opts.BaseDir}

	if err := processQuerySets(ctx, g, &conf, opts.BaseDir, stderr); err != nil {
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
		if err := diffFiles(ctx, opts.BaseDir, res.Files, stderr); err != nil {
			res.Errors = append(res.Errors, err)
		}
	}

	return res
}

const errMessageNoVersion = `The configuration must have a version number.
Set the version to 1 or 2 at the top of the config:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration has an invalid version number.
The supported version can only be "1" or "2".
`

const errMessageNoPackages = `No packages are configured`

type generator struct {
	m       sync.Mutex
	baseDir string
	output  map[string]string
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
	defer g.m.Unlock()

	absout := resolvePath(g.baseDir, out)

	// When the Go codegen is configured to emit the models file into a
	// separate package directory, route that file to its own absolute path.
	// This is the only file allowed to live outside of `out`.
	var (
		modelsFileName string
		modelsAbsout   string
		modelsAbsfile  string
	)
	if sql.Gen.Go != nil && sql.Gen.Go.OutputModelsPath != "" && sql.Gen.Go.ModelsEmitEnabled() {
		modelsFileName = sql.Gen.Go.OutputModelsFileName
		if modelsFileName == "" {
			modelsFileName = "models.go"
		}
		modelsAbsout = resolvePath(g.baseDir, sql.Gen.Go.OutputModelsPath)
		modelsAbsfile = filepath.Join(modelsAbsout, modelsFileName)
	}

	for n, source := range files {
		if modelsFileName != "" && n == modelsFileName {
			// Models file routed to a separate package directory.
			if strings.Contains(modelsAbsfile, "..") {
				return fmt.Errorf("invalid file output path: %s", modelsAbsfile)
			}
			if !strings.HasPrefix(modelsAbsfile, modelsAbsout) {
				return fmt.Errorf("invalid file output path: %s", modelsAbsfile)
			}
			g.output[modelsAbsfile] = source
			continue
		}
		filename := resolvePath(g.baseDir, filepath.Join(out, n))
		if strings.Contains(filename, "..") {
			return fmt.Errorf("invalid file output path: %s", filename)
		}
		if !strings.HasPrefix(filename, absout) {
			return fmt.Errorf("invalid file output path: %s", filename)
		}
		g.output[filename] = source
	}
	return nil
}

// resolvePath joins p with baseDir when p is relative. baseDir is treated as
// the current working directory when empty.
func resolvePath(baseDir, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if baseDir == "" {
		abs, err := filepath.Abs(p)
		if err == nil {
			return abs
		}
		return p
	}
	return filepath.Join(baseDir, p)
}
