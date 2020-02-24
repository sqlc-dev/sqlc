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

	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/dinosql/kotlin"
	"github.com/kyleconroy/sqlc/internal/mysql"
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

func printFileErr(stderr io.Writer, dir string, fileErr dinosql.FileErr) {
	filename := strings.TrimPrefix(fileErr.Filename, dir+"/")
	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
}

type outPair struct {
	Gen config.SQLGen
	config.SQL
}

func Generate(dir string, stderr io.Writer) (map[string]string, error) {
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
		var result dinosql.Generateable

		// TODO: This feels like a hack that will bite us later
		sql.Schema = filepath.Join(dir, sql.Schema)
		sql.Queries = filepath.Join(dir, sql.Queries)

		var name string
		parseOpts := dinosql.ParserOpts{}
		if sql.Gen.Go != nil {
			name = combo.Go.Package
		} else if sql.Gen.Kotlin != nil {
			parseOpts.UsePositionalParameters = true
			name = combo.Kotlin.Package
		}

		result, errored = parse(name, dir, sql.SQL, combo, parseOpts, stderr)
		if errored {
			break
		}

		var files map[string]string
		var out string
		if sql.Gen.Go != nil {
			out = combo.Go.Out
			files, err = dinosql.Generate(result, combo)
		} else if sql.Gen.Kotlin != nil {
			out = combo.Kotlin.Out
			ktRes, ok := result.(kotlin.KtGenerateable)
			if !ok {
				err = fmt.Errorf("kotlin not supported for engine %s", combo.Package.Engine)
				break
			}
			files, err = kotlin.KtGenerate(ktRes, combo)
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

func parse(name, dir string, sql config.SQL, combo config.CombinedSettings, parserOpts dinosql.ParserOpts, stderr io.Writer) (dinosql.Generateable, bool) {
	switch sql.Engine {
	case config.EngineMySQL:
		// Experimental MySQL support
		q, err := mysql.GeneratePkg(name, sql.Schema, sql.Queries, combo)
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			if parserErr, ok := err.(*dinosql.ParserErr); ok {
				for _, fileErr := range parserErr.Errs {
					printFileErr(stderr, dir, fileErr)
				}
			} else {
				fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
			}
			return nil, true
		}
		return q, false

	case config.EnginePostgreSQL:
		c, err := dinosql.ParseCatalog(sql.Schema)
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			if parserErr, ok := err.(*dinosql.ParserErr); ok {
				for _, fileErr := range parserErr.Errs {
					printFileErr(stderr, dir, fileErr)
				}
			} else {
				fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
			}
			return nil, true
		}

		q, err := dinosql.ParseQueries(c, sql.Queries, parserOpts)
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			if parserErr, ok := err.(*dinosql.ParserErr); ok {
				for _, fileErr := range parserErr.Errs {
					printFileErr(stderr, dir, fileErr)
				}
			} else {
				fmt.Fprintf(stderr, "error parsing queries: %s\n", err)
			}
			return nil, true
		}
		return &kotlin.Result{Result: q}, false

	case config.EngineXLemon, config.EngineXDolphin, config.EngineXElephant:
		r, err := compiler.Run(sql, combo)
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			fmt.Fprintf(stderr, "error: %s\n", err)
			return nil, true
		}
		return r, false

	default:
		panic("invalid engine")
	}
}
