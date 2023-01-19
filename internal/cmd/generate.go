package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/trace"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/codegen/json"
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/config/convert"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/ext"
	"github.com/kyleconroy/sqlc/internal/ext/process"
	"github.com/kyleconroy/sqlc/internal/ext/wasm"
	"github.com/kyleconroy/sqlc/internal/multierr"
	"github.com/kyleconroy/sqlc/internal/opts"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 or 2 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The supported version can only be "1" or "2".
`

const errMessageNoPackages = `No packages are configured`

func printFileErr(stderr io.Writer, dir string, fileErr *multierr.FileError) {
	filename := strings.TrimPrefix(fileErr.Filename, dir+"/")
	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
}

type outPair struct {
	Gen    config.SQLGen
	Plugin *config.Codegen

	config.SQL
}

func findPlugin(conf config.Config, name string) (*config.Plugin, error) {
	for _, plug := range conf.Plugins {
		if plug.Name == name {
			return &plug, nil
		}
	}
	return nil, fmt.Errorf("plugin not found")
}

func readConfig(stderr io.Writer, dir, filename string) (string, *config.Config, error) {
	configPath := ""
	if filename != "" {
		configPath = filepath.Join(dir, filename)
	} else {
		var yamlMissing, jsonMissing bool
		yamlPath := filepath.Join(dir, "sqlc.yaml")
		jsonPath := filepath.Join(dir, "sqlc.json")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			yamlMissing = true
		}
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			jsonMissing = true
		}

		if yamlMissing && jsonMissing {
			fmt.Fprintln(stderr, "error parsing configuration files. sqlc.yaml or sqlc.json: file does not exist")
			return "", nil, errors.New("config file missing")
		}

		if !yamlMissing && !jsonMissing {
			fmt.Fprintln(stderr, "error: both sqlc.json and sqlc.yaml files present")
			return "", nil, errors.New("sqlc.json and sqlc.yaml present")
		}

		configPath = yamlPath
		if yamlMissing {
			configPath = jsonPath
		}
	}

	base := filepath.Base(configPath)
	blob, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}

	conf, err := config.ParseConfig(bytes.NewReader(blob))
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprintf(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprintf(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprintf(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing %s: %s\n", base, err)
		return "", nil, err
	}

	return configPath, &conf, nil
}

func Generate(ctx context.Context, e Env, dir, filename string, stderr io.Writer) (map[string]string, error) {
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(configPath)
	if err := config.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return nil, err
	}

	if err := e.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return nil, err
	}

	output := map[string]string{}
	errored := false

	var pairs []outPair
	for _, sql := range conf.SQL {
		if sql.Gen.Go != nil {
			pairs = append(pairs, outPair{
				SQL: sql,
				Gen: config.SQLGen{Go: sql.Gen.Go},
			})
		}
		if sql.Gen.JSON != nil {
			pairs = append(pairs, outPair{
				SQL: sql,
				Gen: config.SQLGen{JSON: sql.Gen.JSON},
			})
		}
		for i, _ := range sql.Codegen {
			pairs = append(pairs, outPair{
				SQL:    sql,
				Plugin: &sql.Codegen[i],
			})
		}
	}

	var m sync.Mutex
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

			// TODO: This feels like a hack that will bite us later
			joined := make([]string, 0, len(sql.Schema))
			for _, s := range sql.Schema {
				joined = append(joined, filepath.Join(dir, s))
			}
			sql.Schema = joined

			joined = make([]string, 0, len(sql.Queries))
			for _, q := range sql.Queries {
				joined = append(joined, filepath.Join(dir, q))
			}
			sql.Queries = joined

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
			trace.Logf(gctx, "", "name=%s dir=%s plugin=%s", name, dir, lang)

			result, failed := parse(gctx, name, dir, sql.SQL, combo, parseOpts, errout)
			if failed {
				packageRegion.End()
				errored = true
				return nil
			}

			out, resp, err := codegen(gctx, combo, sql, result)
			if err != nil {
				fmt.Fprintf(errout, "# package %s\n", name)
				fmt.Fprintf(errout, "error generating code: %s\n", err)
				errored = true
				packageRegion.End()
				return nil
			}

			files := map[string]string{}
			for _, file := range resp.Files {
				files[file.Name] = string(file.Contents)
			}

			m.Lock()
			for n, source := range files {
				filename := filepath.Join(dir, out, n)
				output[filename] = source
			}
			m.Unlock()

			packageRegion.End()
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return nil, err
	}
	if errored {
		for i, _ := range stderrs {
			if _, err := io.Copy(stderr, &stderrs[i]); err != nil {
				return nil, err
			}
		}
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}

func parse(ctx context.Context, name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, stderr io.Writer) (*compiler.Result, bool) {
	defer trace.StartRegion(ctx, "parse").End()
	c := compiler.NewCompiler(sql, combo)
	if err := c.ParseCatalog(sql.Schema); err != nil {
		fmt.Fprintf(stderr, "# package %s\n", name)
		if parserErr, ok := err.(*multierr.Error); ok {
			for _, fileErr := range parserErr.Errs() {
				printFileErr(stderr, dir, fileErr)
			}
		} else {
			fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
		}
		return nil, true
	}
	if parserOpts.Debug.DumpCatalog {
		debug.Dump(c.Catalog())
	}
	if err := c.ParseQueries(sql.Queries, parserOpts); err != nil {
		fmt.Fprintf(stderr, "# package %s\n", name)
		if parserErr, ok := err.(*multierr.Error); ok {
			for _, fileErr := range parserErr.Errs() {
				printFileErr(stderr, dir, fileErr)
			}
		} else {
			fmt.Fprintf(stderr, "error parsing queries: %s\n", err)
		}
		return nil, true
	}
	return c.Result(), false
}

func codegen(ctx context.Context, combo config.CombinedSettings, sql outPair, result *compiler.Result) (string, *plugin.CodeGenResponse, error) {
	defer trace.StartRegion(ctx, "codegen").End()
	req := codeGenRequest(result, combo)
	var handler ext.Handler
	var out string
	switch {
	case sql.Gen.Go != nil:
		out = combo.Go.Out
		handler = ext.HandleFunc(golang.Generate)

	case sql.Gen.JSON != nil:
		out = combo.JSON.Out
		handler = ext.HandleFunc(json.Generate)

	case sql.Plugin != nil:
		out = sql.Plugin.Out
		plug, err := findPlugin(combo.Global, sql.Plugin.Plugin)
		if err != nil {
			return "", nil, fmt.Errorf("plugin not found: %s", err)
		}

		switch {
		case plug.Process != nil:
			handler = &process.Runner{
				Cmd: plug.Process.Cmd,
			}
		case plug.WASM != nil:
			handler = &wasm.Runner{
				URL:    plug.WASM.URL,
				SHA256: plug.WASM.SHA256,
			}
		default:
			return "", nil, fmt.Errorf("unsupported plugin type")
		}

		opts, err := convert.YAMLtoJSON(sql.Plugin.Options)
		if err != nil {
			return "", nil, fmt.Errorf("invalid plugin options")
		}
		req.PluginOptions = opts

	default:
		return "", nil, fmt.Errorf("missing language backend")
	}
	resp, err := handler.Generate(ctx, req)
	return out, resp, err
}
