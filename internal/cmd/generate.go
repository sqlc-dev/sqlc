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

	"github.com/sqlc-dev/sqlc/internal/codegen/golang"
	genjson "github.com/sqlc-dev/sqlc/internal/codegen/json"
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/config/convert"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/ext"
	"github.com/sqlc-dev/sqlc/internal/ext/process"
	"github.com/sqlc-dev/sqlc/internal/ext/wasm"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/sqlcdebug"
)

var debugDumpCatalog = sqlcdebug.New("dumpcatalog")

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

// sourceFiles holds in-memory config and optional file contents so generate can run
// without reading from disk (e.g. in tests). For production, Generate reads from FS and
// fills FileContents before calling generate.
type sourceFiles struct {
	Config       *config.Config
	ConfigPath   string
	Dir          string
	FileContents map[string][]byte // path -> content; keys match paths used when reading (e.g. filepath.Join(dir, "schema.sql"))
}

// Generate runs codegen for the given directory and config file, reading all input from disk.
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

	return generate(ctx, &sourceFiles{
		Config:       conf,
		ConfigPath:   configPath,
		Dir:          dir,
		FileContents: nil,
	}, o)
}

// generate runs codegen from in-memory or on-disk inputs (see sourceFiles).
func generate(ctx context.Context, inputs *sourceFiles, o *Options) (map[string]string, error) {
	g := &generator{
		dir:                    inputs.Dir,
		output:                 map[string]string{},
		codegenHandlerOverride: o.CodegenHandlerOverride,
	}
	if err := processQuerySets(ctx, g, inputs, o); err != nil {
		return nil, err
	}
	return g.output, nil
}

type generator struct {
	m                      sync.Mutex
	dir                    string
	output                 map[string]string
	codegenHandlerOverride grpc.ClientConnInterface
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
	out, resp, err := codegen(ctx, combo, sql, result, g.codegenHandlerOverride)
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
		modelsAbsout = filepath.Join(g.dir, sql.Gen.Go.OutputModelsPath)
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

func parse(ctx context.Context, name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, stderr io.Writer) (*compiler.Result, bool) {
	defer trace.StartRegion(ctx, "parse").End()
	c, err := compiler.NewCompiler(sql, combo, parserOpts)
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
	if debugDumpCatalog.Value() == "1" {
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

func codegen(ctx context.Context, combo config.CombinedSettings, sql OutputPair, result *compiler.Result, codegenOverride grpc.ClientConnInterface) (string, *plugin.GenerateResponse, error) {
	defer trace.StartRegion(ctx, "codegen").End()
	req := codeGenRequest(result, combo)
	var handler grpc.ClientConnInterface
	var out string
	switch {
	case sql.Plugin != nil:
		out = sql.Plugin.Out
		if codegenOverride != nil {
			handler = codegenOverride
		} else {
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
			global, found := combo.Global.Options[plug.Name]
			if found {
				opts, err := convert.YAMLtoJSON(global)
				if err != nil {
					return "", nil, fmt.Errorf("invalid global options: %w", err)
				}
				req.GlobalOptions = opts
			}
		}
		opts, err := convert.YAMLtoJSON(sql.Plugin.Options)
		if err != nil {
			return "", nil, fmt.Errorf("invalid plugin options: %w", err)
		}
		if err := validateExternalPluginOptions(opts, sql.Plugin.Plugin); err != nil {
			return "", nil, err
		}
		req.PluginOptions = opts

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

// driverOnlyOptions are options that apply only to built-in Go codegen, not to external plugins.
var driverOnlyOptions = []string{"sql_package", "sql_driver"}

// validateExternalPluginOptions returns an error if plugin options contain sql_package or sql_driver.
// External codegen plugins define their own database driver; these options are not supported.
func validateExternalPluginOptions(opts []byte, pluginName string) error {
	if len(opts) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(opts, &m); err != nil {
		return nil
	}
	var invalid []string
	for _, key := range driverOnlyOptions {
		if _, ok := m[key]; ok {
			invalid = append(invalid, key)
		}
	}
	if len(invalid) == 0 {
		return nil
	}
	return fmt.Errorf("plugin %q: options %q are not supported for external codegen plugins; the plugin defines its own database driver (these options only apply to built-in Go codegen)", pluginName, invalid)
}
