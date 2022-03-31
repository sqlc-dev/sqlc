package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"strings"

	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/codegen/kotlin"
	"github.com/kyleconroy/sqlc/internal/codegen/python"
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/multierr"
	"github.com/kyleconroy/sqlc/internal/opts"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The only supported version is "1".
`

const errMessageNoPackages = `No packages are configured`

func printFileErr(stderr io.Writer, dir string, fileErr *multierr.FileError) {
	filename := strings.TrimPrefix(fileErr.Filename, dir+"/")
	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
}

type outPair struct {
	Gen config.SQLGen
	config.SQL
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
		if sql.Gen.Kotlin != nil {
			pairs = append(pairs, outPair{
				SQL: sql,
				Gen: config.SQLGen{Kotlin: sql.Gen.Kotlin},
			})
		}
		if sql.Gen.Python != nil {
			pairs = append(pairs, outPair{
				SQL: sql,
				Gen: config.SQLGen{Python: sql.Gen.Python},
			})
		}
	}

	for _, sql := range pairs {
		combo := config.Combine(*conf, sql.SQL)

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
		if sql.Gen.Go != nil {
			name = combo.Go.Package
			lang = "golang"
		} else if sql.Gen.Kotlin != nil {
			if sql.Engine == config.EnginePostgreSQL {
				parseOpts.UsePositionalParameters = true
			}
			lang = "kotlin"
			name = combo.Kotlin.Package
		} else if sql.Gen.Python != nil {
			lang = "python"
			name = combo.Python.Package
		}

		var packageRegion *trace.Region
		if debug.Traced {
			packageRegion = trace.StartRegion(ctx, "package")
			trace.Logf(ctx, "", "name=%s dir=%s language=%s", name, dir, lang)
		}

		result, failed := parse(ctx, e, name, dir, sql.SQL, combo, parseOpts, stderr)
		if failed {
			if packageRegion != nil {
				packageRegion.End()
			}
			errored = true
			break
		}

		var region *trace.Region
		if debug.Traced {
			region = trace.StartRegion(ctx, "codegen")
		}
		var genfunc func(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error)
		var out string
		switch {
		case sql.Gen.Go != nil:
			out = combo.Go.Out
			genfunc = golang.Generate
		case sql.Gen.Kotlin != nil:
			out = combo.Kotlin.Out
			genfunc = kotlin.Generate
		case sql.Gen.Python != nil:
			out = combo.Python.Out
			genfunc = python.Generate
		default:
			panic("missing language backend")
		}
		resp, err := genfunc(codeGenRequest(result, combo))
		if region != nil {
			region.End()
		}
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			fmt.Fprintf(stderr, "error generating code: %s\n", err)
			errored = true
			if packageRegion != nil {
				packageRegion.End()
			}
			continue
		}
		files := map[string]string{}
		for _, file := range resp.Files {
			files[file.Name] = string(file.Contents)
		}
		for n, source := range files {
			filename := filepath.Join(dir, out, n)
			output[filename] = source
		}
		if packageRegion != nil {
			packageRegion.End()
		}
	}

	if errored {
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}

func parse(ctx context.Context, e Env, name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, stderr io.Writer) (*compiler.Result, bool) {
	if debug.Traced {
		defer trace.StartRegion(ctx, "parse").End()
	}
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
