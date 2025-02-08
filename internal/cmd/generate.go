package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang"
	genjson "github.com/sqlc-dev/sqlc/internal/codegen/json"
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/config/convert"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/ext"
	"github.com/sqlc-dev/sqlc/internal/ext/process"
	"github.com/sqlc-dev/sqlc/internal/ext/wasm"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/remote"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
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
	filename, err := filepath.Rel(dir, fileErr.Filename)
	if err != nil {
		filename = fileErr.Filename
	}
	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
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
		var yamlMissing, jsonMissing, ymlMissing bool
		yamlPath := filepath.Join(dir, "sqlc.yaml")
		ymlPath := filepath.Join(dir, "sqlc.yml")
		jsonPath := filepath.Join(dir, "sqlc.json")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			yamlMissing = true
		}
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			jsonMissing = true
		}

		if _, err := os.Stat(ymlPath); os.IsNotExist(err) {
			ymlMissing = true
		}

		if yamlMissing && ymlMissing && jsonMissing {
			fmt.Fprintln(stderr, "error parsing configuration files. sqlc.(yaml|yml) or sqlc.json: file does not exist")
			return "", nil, errors.New("config file missing")
		}

		if (!yamlMissing || !ymlMissing) && !jsonMissing {
			fmt.Fprintln(stderr, "error: both sqlc.json and sqlc.(yaml|yml) files present")
			return "", nil, errors.New("sqlc.json and sqlc.(yaml|yml) present")
		}

		if jsonMissing {
			if yamlMissing {
				configPath = ymlPath
			} else {
				configPath = yamlPath
			}
		} else {
			configPath = jsonPath
		}
	}

	base := filepath.Base(configPath)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}
	defer file.Close()

	conf, err := config.ParseConfig(file)
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprint(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprint(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprint(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing %s: %s\n", base, err)
		return "", nil, err
	}

	return configPath, &conf, nil
}

func Generate(ctx context.Context, dir, filename string, o *Options) (map[string]string, error) {
	e := o.Env
	stderr := o.Stderr

	configPath, conf, err := o.ReadConfig(dir, filename)
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

	// Comment on why these two methods exist
	if conf.Cloud.Project != "" && e.Remote && !e.NoRemote {
		return remoteGenerate(ctx, configPath, conf, dir, stderr)
	}

	g := &generator{
		dir:    dir,
		output: map[string]string{},
	}

	if err := processQuerySets(ctx, g, conf, dir, o); err != nil {
		return nil, err
	}

	return g.output, nil
}

type generator struct {
	m      sync.Mutex
	dir    string
	output map[string]string
}

func (g *generator) Pairs(ctx context.Context, conf *config.Config) []OutputPair {
	var pairs []OutputPair
	for _, sql := range conf.SQL {
		if sql.Gen.Go != nil {
			pairs = append(pairs, OutputPair{
				SQL: sql,
				Gen: config.SQLGen{Go: sql.Gen.Go},
			})
		}
		if sql.Gen.JSON != nil {
			pairs = append(pairs, OutputPair{
				SQL: sql,
				Gen: config.SQLGen{JSON: sql.Gen.JSON},
			})
		}
		for i := range sql.Codegen {
			pairs = append(pairs, OutputPair{
				SQL:    sql,
				Plugin: &sql.Codegen[i],
			})
		}
	}
	return pairs
}

func (g *generator) ProcessResult(ctx context.Context, combo config.CombinedSettings, sql OutputPair, result *compiler.Result) error {
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

func remoteGenerate(ctx context.Context, configPath string, conf *config.Config, dir string, stderr io.Writer) (map[string]string, error) {
	rpcClient, err := remote.NewClient(conf.Cloud)
	if err != nil {
		fmt.Fprintf(stderr, "error creating rpc client: %s\n", err)
		return nil, err
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error reading config file %s: %s\n", configPath, err)
		return nil, err
	}

	rpcReq := remote.GenerateRequest{
		Version: info.Version,
		Inputs:  []*remote.File{{Path: filepath.Base(configPath), Bytes: configBytes}},
	}

	for _, pkg := range conf.SQL {
		for _, paths := range []config.Paths{pkg.Schema, pkg.Queries} {
			for i, relFilePath := range paths {
				paths[i] = filepath.Join(dir, relFilePath)
			}
			files, err := sqlpath.Glob(paths)
			if err != nil {
				fmt.Fprintf(stderr, "error globbing paths: %s\n", err)
				return nil, err
			}
			for _, filePath := range files {
				fileBytes, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Fprintf(stderr, "error reading file %s: %s\n", filePath, err)
					return nil, err
				}
				fileRelPath, _ := filepath.Rel(dir, filePath)
				rpcReq.Inputs = append(rpcReq.Inputs, &remote.File{Path: fileRelPath, Bytes: fileBytes})
			}
		}
	}

	rpcResp, err := rpcClient.Generate(ctx, &rpcReq)
	if err != nil {
		rpcStatus, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		fmt.Fprintf(stderr, "rpc error: %s", rpcStatus.Message())
		return nil, rpcStatus.Err()
	}

	if rpcResp.ExitCode != 0 {
		fmt.Fprintf(stderr, "%s", rpcResp.Stderr)
		return nil, errors.New("remote execution returned with non-zero exit code")
	}

	output := map[string]string{}
	for _, file := range rpcResp.Outputs {
		output[filepath.Join(dir, file.Path)] = string(file.Bytes)
	}

	return output, nil
}

func parse(ctx context.Context, name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, stderr io.Writer) (*compiler.Result, bool) {
	defer trace.StartRegion(ctx, "parse").End()
	c, err := compiler.NewCompiler(sql, combo)
	defer func() {
		if c != nil {
			c.Close(ctx)
		}
	}()
	if err != nil {
		fmt.Fprintf(stderr, "error creating compiler: %s\n", err)
		return nil, true
	}
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

func codegen(ctx context.Context, combo config.CombinedSettings, sql OutputPair, result *compiler.Result) (string, *plugin.GenerateResponse, error) {
	defer trace.StartRegion(ctx, "codegen").End()
	req := codeGenRequest(result, combo)
	var handler grpc.ClientConnInterface
	var out string
	switch {
	case sql.Plugin != nil:
		out = sql.Plugin.Out
		plug, err := findPlugin(combo.Global, sql.Plugin.Plugin)
		if err != nil {
			return "", nil, fmt.Errorf("plugin not found: %s", err)
		}

		switch {
		case plug.Process != nil:
			handler = &process.Runner{
				Cmd:    plug.Process.Cmd,
				Env:    plug.Env,
				Format: plug.Process.Format,
			}
		case plug.WASM != nil:
			handler = &wasm.Runner{
				URL:    plug.WASM.URL,
				SHA256: plug.WASM.SHA256,
				Env:    plug.Env,
			}
		default:
			return "", nil, fmt.Errorf("unsupported plugin type")
		}

		opts, err := convert.YAMLtoJSON(sql.Plugin.Options)
		if err != nil {
			return "", nil, fmt.Errorf("invalid plugin options: %w", err)
		}
		req.PluginOptions = opts

		global, found := combo.Global.Options[plug.Name]
		if found {
			opts, err := convert.YAMLtoJSON(global)
			if err != nil {
				return "", nil, fmt.Errorf("invalid global options: %w", err)
			}
			req.GlobalOptions = opts
		}

	case sql.Gen.Go != nil:
		out = combo.Go.Out
		handler = ext.HandleFunc(golang.Generate)
		opts, err := json.Marshal(sql.Gen.Go)
		if err != nil {
			return "", nil, fmt.Errorf("opts marshal failed: %w", err)
		}
		req.PluginOptions = opts

		if combo.Global.Overrides.Go != nil {
			opts, err := json.Marshal(combo.Global.Overrides.Go)
			if err != nil {
				return "", nil, fmt.Errorf("opts marshal failed: %w", err)
			}
			req.GlobalOptions = opts
		}

	case sql.Gen.JSON != nil:
		out = combo.JSON.Out
		handler = ext.HandleFunc(genjson.Generate)
		opts, err := json.Marshal(sql.Gen.JSON)
		if err != nil {
			return "", nil, fmt.Errorf("opts marshal failed: %w", err)
		}
		req.PluginOptions = opts

	default:
		return "", nil, fmt.Errorf("missing language backend")
	}
	client := plugin.NewCodegenServiceClient(handler)
	resp, err := client.Generate(ctx, req)
	return out, resp, err
}
