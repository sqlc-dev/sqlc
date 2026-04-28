package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"

	"github.com/sqlc-dev/sqlc/internal/api"
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
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

// Generate is a thin wrapper around api.Generate that translates between the
// CLI's Options struct and api.GenerateOptions. New callers should prefer
// api.Generate directly.
func Generate(ctx context.Context, dir, filename string, o *Options) (map[string]string, error) {
	res := api.Generate(ctx, api.GenerateOptions{
		Dir:                   dir,
		File:                  filename,
		Stderr:                o.Stderr,
		DisableProcessPlugins: !o.Env.Debug.ProcessPlugins,
		MutateConfig:          o.MutateConfig,
	})
	if len(res.Errors) > 0 {
		return res.Files, res.Errors[0]
	}
	return res.Files, nil
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
