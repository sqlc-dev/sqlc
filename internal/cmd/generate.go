package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/codegen/kotlin"
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/multierr"
	"github.com/kyleconroy/sqlc/internal/mysql"
	"github.com/kyleconroy/sqlc/internal/opts"
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

func Generate(e Env, dir string, stderr io.Writer) (map[string]string, error) {
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
		fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
		return nil, errors.New("config file missing")
	}

	if !yamlMissing && !jsonMissing {
		fmt.Fprintln(stderr, "error parsing sqlc.json: both files present")
		return nil, errors.New("sqlc.json and sqlc.yaml present")
	}

	configPath := yamlPath
	if yamlMissing {
		configPath = jsonPath
	}

	blob, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
		return nil, err
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
		fmt.Fprintf(stderr, "error parsing sqlc.json: %s\n", err)
		return nil, err
	}

	debug, err := opts.DebugFromEnv()
	if err != nil {
		fmt.Fprintf(stderr, "error parsing SQLCDEBUG: %s\n", err)
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
	}

	for _, sql := range pairs {
		combo := config.Combine(conf, sql.SQL)

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

		var name string
		parseOpts := opts.Parser{
			Debug: debug,
		}
		if sql.Gen.Go != nil {
			name = combo.Go.Package
		} else if sql.Gen.Kotlin != nil {
			parseOpts.UsePositionalParameters = true
			name = combo.Kotlin.Package
		}

		var files map[string]string
		var out string

		// TODO: Note about how this will be going away
		if sql.Engine == config.EngineMySQL {
			result, errored := parseMySQL(e, name, dir, sql.SQL, combo, parseOpts, stderr)
			if errored {
				break
			}
			out = combo.Go.Out
			files, err = golang.DeprecatedGenerate(result, combo)
		} else {
			result, errored := parse(e, name, dir, sql.SQL, combo, parseOpts, stderr)
			if errored {
				break
			}
			switch {
			case sql.Gen.Go != nil:
				out = combo.Go.Out
				files, err = golang.Generate(result, combo)
			case sql.Gen.Kotlin != nil:
				out = combo.Kotlin.Out
				files, err = kotlin.Generate(result, combo)
			default:
				panic("missing language backend")
			}
		}

		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			fmt.Fprintf(stderr, "error generating code: %s\n", err)
			errored = true
			continue
		}
		for n, source := range files {
			filename := filepath.Join(dir, out, n)
			output[filename] = source
		}
	}

	if errored {
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}

// Experimental MySQL support
func parseMySQL(e Env, name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, stderr io.Writer) (golang.Generateable, bool) {
	q, err := mysql.GeneratePkg(name, sql.Schema, sql.Queries, combo)
	if err != nil {
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
	return q, false
}

func parse(e Env, name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts opts.Parser, stderr io.Writer) (*compiler.Result, bool) {
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
